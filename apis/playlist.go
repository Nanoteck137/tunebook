package apis

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"path"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
)

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	CoverArt types.Images `json:"coverArt"`
}

func ConvertDBPlaylist(c pyrin.Context, playlist database.Playlist) Playlist {
	return Playlist{
		Id:       playlist.Id,
		Name:     playlist.Name,
		CoverArt: ConvertPlaylistCoverURL(c, playlist.Id, playlist.CoverArt),
	}
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type CreatePlaylist struct {
	Id string `json:"id"`
}

type CreatePlaylistBody struct {
	Name string `json:"name"`
}

func (b *CreatePlaylistBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b CreatePlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type PostPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *PostPlaylistFilterBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Filter = anvil.String(b.Filter)
}

func (b PostPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type GetPlaylistById struct {
	Playlist
}

type GetPlaylistItems struct {
	Page  types.Page `json:"page"`
	Items []Track    `json:"items"`
}

type AddItemToPlaylistBody struct {
	TrackId string `json:"trackId"`
}

type RemovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
}

type PlaylistFilter struct {
	FilterId   string `json:"filterId"`
	PlaylistId string `json:"playlistId"`

	Name   string `json:"name"`
	Filter string `json:"filter"`

	// TODO(patrik): Created, Updated
}

type GetPlaylistFilters struct {
	Filters []PlaylistFilter `json:"filters"`
}

type AddPlaylistFilter struct {
	FilterId string `json:"filterId"`
}

type AddPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *AddPlaylistFilterBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Filter = anvil.String(b.Filter)
}

func (b AddPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.Filter, validate.Required),
	)
}

func InstallPlaylistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetPlaylists",
			Path:         "/playlists",
			Method:       http.MethodGet,
			ResponseType: GetPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlists, err := app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
				if err != nil {
					return nil, err
				}

				res := GetPlaylists{
					Playlists: make([]Playlist, len(playlists)),
				}

				for i, playlist := range playlists {
					res.Playlists[i] = ConvertDBPlaylist(c, playlist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistById",
			Path:         "/playlists/:id",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistById{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				return GetPlaylistById{
					Playlist: ConvertDBPlaylist(c, playlist),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreatePlaylist",
			Path:         "/playlists",
			Method:       http.MethodPost,
			ResponseType: CreatePlaylist{},
			BodyType:     CreatePlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreatePlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				return CreatePlaylist{
					Id: playlist.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreatePlaylistFromFilter",
			Path:         "/playlists/filter",
			Method:       http.MethodPost,
			ResponseType: CreatePlaylist{},
			BodyType:     PostPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[PostPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				playlist, err := tx.CreatePlaylist(ctx, database.CreatePlaylistParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				tracks, err := tx.GetAllTracks(ctx, body.Filter, "")
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					return nil, err
				}

				for _, track := range tracks {
					err = tx.AddItemToPlaylist(ctx, playlist.Id, track.Id, 0)
					if err != nil {
						return nil, err
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return CreatePlaylist{
					Id: playlist.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeletePlaylist",
			Path:   "/playlists/:id",
			Method: http.MethodDelete,
			Errors: []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = app.DB().DeletePlaylist(ctx, playlist.Id)
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.WorkDir().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.WorkDir().Cache().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistItems",
			Path:         "/playlists/:id/items",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistItems{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				playlistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				opts := getPageOptions(q)

				filterId := q.Get("filterId")
				if filterId != "" {
					filter, err := app.DB().GetPlaylistFilterById(ctx, filterId, playlist.Id)
					if err != nil {
						// TODO(patrik): Handle error
						return nil, err
					}

					opts.Filter = filter.Filter
				}

				tracks, pageInfo, err := app.DB().GetPlaylistTracksPaged(ctx, playlist.Id, opts)
				if err != nil {
					return nil, err
				}

				res := GetPlaylistItems{
					Page:  pageInfo,
					Items: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					// TODO(patrik): Replace with track order
					track.Track.Number = sql.NullInt64{
						Int64: int64(track.Order) + 1,
						Valid: true,
					}
					res.Items[i] = ConvertDBTrack(c, track.Track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddItemToPlaylist",
			Path:     "/playlists/:id/items",
			Method:   http.MethodPost,
			BodyType: AddItemToPlaylistBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound, ErrTypePlaylistAlreadyHasTrack},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddItemToPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				// TODO(patrik): Check for trackId exists?
				err = app.DB().AddItemToPlaylist(c.Request().Context(), playlist.Id, body.TrackId, rand.Int())
				if err != nil {
					if errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, PlaylistAlreadyHasTrack()
					}

					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemovePlaylistItem",
			Path:     "/playlists/:id/items",
			Method:   http.MethodDelete,
			BodyType: RemovePlaylistItemBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[RemovePlaylistItemBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				// TODO(patrik): Check for trackId exists?
				err = app.DB().RemovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearPlaylist",
			Path:   "/playlists/:id/items/all",
			Method: http.MethodDelete,
			Errors: []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = app.DB().RemoveAllPlaylistItem(ctx, playlist.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "GeneratePlaylistImage",
			Method: http.MethodPost,
			Path:   "/playlists/:id/images/generate",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				images, err := app.DB().GetPlaylistTrackImages(ctx, playlist.Id, 4)
				if err != nil {
					return nil, err
				}

				imgs := [4]string{
					"",
					"",
					"",
					"",
				}

				for i, img := range images {
					if !img.Valid {
						continue
					}

					imgs[i] = img.String
				}

				workDir := app.WorkDir()
				playlistDir := workDir.Playlist(playlist.Id)

				dirs := []string{
					playlistDir,
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !errors.Is(err, os.ErrExist) {
						return nil, err
					}
				}

				coverArt := "generated.png"
				out := path.Join(playlistDir, coverArt)
				err = utils.GeneratePlaylistCover(imgs, out, 512)
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
					CoverArt: types.Change[sql.NullString]{
						Value: sql.NullString{
							String: coverArt,
							Valid:  coverArt != "",
						},
						Changed: coverArt != playlist.CoverArt.String,
					},
				})
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.WorkDir().Cache().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistFilters",
			Method:       http.MethodGet,
			Path:         "/playlists/:playlistId/filters",
			ResponseType: GetPlaylistFilters{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				ctx := context.Background()

				// TODO(patrik): Get playlist?

				filters, err := app.DB().GetPlaylistFiltersByPlaylistId(ctx, playlistId)
				if err != nil {
					return nil, err
				}

				res := GetPlaylistFilters{
					Filters: make([]PlaylistFilter, len(filters)),
				}

				for i, filter := range filters {
					res.Filters[i] = PlaylistFilter{
						FilterId:   filter.Id,
						PlaylistId: filter.PlaylistId,
						Name:       filter.Name,
						Filter:     filter.Filter,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddPlaylistFilter",
			Method:       http.MethodPost,
			Path:         "/playlists/:playlistId/filters",
			ResponseType: AddPlaylistFilter{},
			BodyType:     AddPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				body, err := pyrin.Body[AddPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = TestFilter(body.Filter)
				if err != nil {
					// TODO(patrik): Better error
					return nil, err
				}

				filterId, err := app.DB().CreatePlaylistFilter(ctx, database.CreatePlaylistFilterParams{
					PlaylistId: playlistId,
					Name:       body.Name,
					Filter:     body.Filter,
				})
				if err != nil {
					return nil, err
				}

				return AddPlaylistFilter{
					FilterId: filterId,
				}, nil
			},
		},
	)
}
