package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
)

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	CoverArt types.Images `json:"coverArt"`

	OwnerId          string         `json:"ownerId"`
	OwnerDisplayName string         `json:"ownerDisplayName"`
	OwnerPicture     types.Images `json:"ownerPicture"`

	TrackCount int64 `json:"trackCount"`
}

func ConvertDBPlaylist(c pyrin.Context, playlist database.Playlist) Playlist {
	return Playlist{
		Id:               playlist.Id,
		Name:             playlist.Name,
		CoverArt:         ConvertPlaylistCoverURL(c, playlist.Id, playlist.CoverArt),
		OwnerId:          playlist.OwnerId,
		OwnerDisplayName: playlist.OwnerDisplayName,
		OwnerPicture:     ConvertUserPictureURL(c, playlist.OwnerId, playlist.OwnerPicture),
		TrackCount:       playlist.TrackCount.Int64,
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
		validate.Field(&b.Filter, validate.Required, validateFilter),
	)
}

type EditPlaylistFilterBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b *EditPlaylistFilterBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	b.Filter = anvil.StringPtr(b.Filter)
}

func testPlaylistItemFilter(val string) error {
	// TODO(patrik): Replace with PlaylistResolverAdapter when it exists
	return database.TestFilter(val, &adapter.TrackResolverAdapter{})
}

// TODO(patrik): Move
var validateFilter = validate.By(func(value any) error {
	switch value := value.(type) {
	case *string:
		return testPlaylistItemFilter(*value)
	case string:
		return testPlaylistItemFilter(value)
	default:
		panic(fmt.Sprintf("validateFilter: Unknown type: %T", value))
	}

})

func (b EditPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil), validateFilter),
	)
}

type EditPlaylistBody struct {
	Name *string `json:"name,omitempty"`

	CoverUrl *string `json:"coverUrl,omitempty"`
}

func (b *EditPlaylistBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)

	b.CoverUrl = anvil.StringPtr(b.CoverUrl)
}

func (b EditPlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.CoverUrl, validate.Required.When(b.CoverUrl != nil)),
	)
}

type ReorderPlaylistItemsBody struct {
	Before        bool     `json:"before"`
	AnchorTrackId string   `json:"anchorTrackId"`
	TrackIds      []string `json:"trackIds"`
}

// TODO(patrik): Handle filter errors
func handlePlaylistServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrPlaylistServicePlaylistNotFound):
		return PlaylistNotFound()
	case errors.Is(err, service.ErrPlaylistServiceTrackNotFound):
		return TrackNotFound()
	case errors.Is(err, service.ErrPlaylistServiceTrackAlreadyAdded):
		return PlaylistAlreadyHasTrack()
	case errors.Is(err, service.ErrPlaylistServiceFilterNotFound):
		return PlaylistFilterNotFound()
	case errors.Is(err, service.ErrPlaylistServiceAnchorTrackNotFound):
		// TODO(patrik): Replace with its own error
		return TrackNotFound()
	case errors.Is(err, service.ErrPlaylistServiceNotAuthorized):
		// TODO(patrik): Replace with its own error
		return PlaylistNotFound()
	}

	return err
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

				ctx := c.Request().Context()

				playlists, err := app.PlaylistService().GetPlaylistsByUser(
					ctx,
					user.Id,
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
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
			Path:         "/playlists/:playlistId",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				playlist, err := app.PlaylistService().GetPlaylistById(
					ctx,
					service.GetPlaylistByIdParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
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
				body, err := pyrin.Body[CreatePlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				playlistId, err := app.PlaylistService().CreatePlaylist(
					ctx,
					service.CreatePlaylistParams{
						Name:    body.Name,
						OwnerId: user.Id,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return CreatePlaylist{
					Id: playlistId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditPlaylist",
			Method:   http.MethodPatch,
			Path:     "/playlists/:playlistId",
			BodyType: EditPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[EditPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.PlaylistService().EditPlaylist(
					ctx,
					service.EditPlaylistParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,

						Name:     body.Name,
						CoverUrl: body.CoverUrl,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeletePlaylist",
			Path:   "/playlists/:playlistId",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.PlaylistService().DeletePlaylist(
					ctx,
					service.DeletePlaylistParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:   "UploadPlaylistImage",
			Method: http.MethodPost,
			Path:   "/playlists/:playlistId/image/upload",
			Spec: pyrin.FormSpec{
				Files: map[string]pyrin.FormFileSpec{
					"image": {
						NumExpected: 1,
					},
				},
			},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				files, err := pyrin.FormFiles(c, "image")
				if err != nil {
					return nil, err
				}

				err = app.PlaylistService().UploadPlaylistImage(
					ctx,
					service.UploadPlaylistImageParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
						File:       files[0],
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "GeneratePlaylistImage",
			Method: http.MethodPost,
			Path:   "/playlists/:playlistId/images/generate",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.PlaylistService().GeneratePlaylistImage(
					ctx,
					service.GeneratePlaylistImageParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistItems",
			Path:         "/playlists/:playlistId/items",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistItems{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				tracks, page, err := app.PlaylistService().GetPlaylistItems(
					ctx,
					service.GetPlaylistItemsParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
						Page:       pageParams,
						Filter:     filterParams,
						FilterId:   q.Get("filterId"),
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				res := GetPlaylistItems{
					Page:  page,
					Items: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Items[i] = ConvertDBTrack(c, track.Track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddItemToPlaylist",
			Path:     "/playlists/:playlistId/items",
			Method:   http.MethodPost,
			BodyType: AddItemToPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				body, err := pyrin.Body[AddItemToPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				err = app.PlaylistService().AddItemToPlaylist(
					ctx,
					service.AddItemToPlaylistParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
						TrackId:    body.TrackId,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemovePlaylistItem",
			Path:     "/playlists/:playlistId/items",
			Method:   http.MethodDelete,
			BodyType: RemovePlaylistItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				body, err := pyrin.Body[RemovePlaylistItemBody](c)
				if err != nil {
					return nil, err
				}

				err = app.PlaylistService().RemovePlaylistItem(
					ctx,
					service.RemovePlaylistItemParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
						TrackId:    body.TrackId,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ReorderPlaylistItems",
			Path:     "/playlists/:playlistId/items/reorder",
			Method:   http.MethodPost,
			BodyType: ReorderPlaylistItemsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[ReorderPlaylistItemsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.PlaylistService().ReorderPlaylistItems(
					ctx,
					service.ReorderPlaylistItemsParams{
						PlaylistId:    c.Param("playlistId"),
						UserId:        user.Id,
						Before:        body.Before,
						AnchorTrackId: body.AnchorTrackId,
						TrackIds:      body.TrackIds,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
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
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				filters, err := app.PlaylistService().GetPlaylistFilters(
					ctx,
					service.GetPlaylistFiltersParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
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
			Name:         "CreatePlaylistFilter",
			Method:       http.MethodPost,
			Path:         "/playlists/:playlistId/filters",
			ResponseType: AddPlaylistFilter{},
			BodyType:     AddPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AddPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				filterId, err := app.PlaylistService().CreatePlaylistFilter(
					ctx,
					service.CreatePlaylistFilterParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
						Name:       body.Name,
						Filter:     body.Filter,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return AddPlaylistFilter{
					FilterId: filterId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditPlaylistFilter",
			Method:   http.MethodPatch,
			Path:     "/playlists/:playlistId/filters/:filterId",
			BodyType: EditPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[EditPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.PlaylistService().EditPlaylistFilter(ctx, service.EditPlaylistFilterParams{
					PlaylistId: c.Param("playlistId"),
					UserId:     user.Id,
					FilterId:   c.Param("filterId"),
					Name:       body.Name,
					Filter:     body.Filter,
				})
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
