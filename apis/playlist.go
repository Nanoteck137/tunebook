package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/anvil"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/validate"
)

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	CoverArt types.Images `json:"coverArt"`

	OwnerId          string       `json:"ownerId"`
	OwnerDisplayName string       `json:"ownerDisplayName"`
	OwnerPicture     types.Images `json:"ownerPicture"`

	TrackCount int64 `json:"trackCount"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBPlaylist(c pyrin.Context, playlist database.Playlist) Playlist {
	return Playlist{
		Id:               playlist.Id,
		Name:             playlist.Name,
		CoverArt:         ConvertPlaylistCoverURL(c, playlist.Id),
		OwnerId:          playlist.OwnerId,
		OwnerDisplayName: playlist.OwnerDisplayName,
		OwnerPicture:     ConvertUserPictureURL(c, playlist.OwnerId),
		TrackCount:       playlist.TrackCount,
		Created:          formatTime(playlist.Created),
		Updated:          formatTime(playlist.Updated),
	}
}

type GetPlaylists struct {
	Page      types.Page `json:"page"`
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

type GetPlaylistById struct {
	Playlist Playlist `json:"playlist"`
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

var trackSchema = database.TrackSchema()

// TODO(patrik): Move
var validateTrackFilter = validate.By(func(value any) error {
	var filterStr string
	switch value := value.(type) {
	case *string:
		filterStr = *value
	case string:
		filterStr = value
	default:
		panic(fmt.Sprintf("validateTrackFilter: Unknown type: %T", value))
	}

	err := database.ValidateFilter(trackSchema, filterStr)
	if err != nil {
		var filterErr *database.FilterError
		if errors.As(err, &filterErr) {
			return &database.QueryError{Filter: filterErr}
		}
		return err
	}

	return nil
})

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

type GetPlaylistItemIds struct {
	Ids []string `json:"ids"`
}

func handlePlaylistServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrPlaylistServicePlaylistNotFound):
		return PlaylistNotFound()
	case errors.Is(err, service.ErrPlaylistServiceTrackNotFound):
		return TrackNotFound()
	case errors.Is(err, service.ErrPlaylistServiceTrackAlreadyAdded):
		return PlaylistAlreadyHasTrack()
	case errors.Is(err, service.ErrPlaylistServiceItemNotFound):
		return PlaylistItemNotFound()
	case errors.Is(err, service.ErrPlaylistServiceFilterNotFound):
		return FilterNotFound()
	case errors.Is(err, service.ErrPlaylistServiceAnchorTrackNotFound):
		return PlaylistAnchorTrackNotFound()
	case errors.Is(err, service.ErrPlaylistServiceNotAuthorized):
		return NotAuthorized()
	case errors.Is(err, service.ErrImageServiceUnsupportedImageFormat):
		return UnsupportedImageType()
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
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				playlists, page, err := app.PlaylistService().GetPlaylists(
					ctx,
					service.GetPlaylistsParams{
						Page:   pageParams,
						Filter: filterParams,
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				res := GetPlaylists{
					Page:      page,
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
				ctx := c.Request().Context()

				playlist, err := app.PlaylistService().GetPlaylistById(
					ctx,
					service.GetPlaylistByIdParams{
						PlaylistId: c.Param("playlistId"),
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

				ctx := context.Background()

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
				ctx := context.Background()

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

				err = app.JobService().PushJob(
					ctx,
					jobs.GeneratePlaylistImage,
					service.GeneratePlaylistImageParams{
						PlaylistId: c.Param("playlistId"),
						UserId:     user.Id,
					},
				)
				if err != nil {
					return nil, err
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

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				tracks, page, err := app.PlaylistService().GetPlaylistItems(
					ctx,
					service.GetPlaylistItemsParams{
						PlaylistId: c.Param("playlistId"),
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
			Name:         "GetPlaylistItemIds",
			Path:         "/playlists/:playlistId/ids",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistItemIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				ids, err := app.PlaylistService().GetPlaylistItemIds(
					ctx,
					service.GetPlaylistItemIdsParams{
						PlaylistId: c.Param("playlistId"),
					},
				)
				if err != nil {
					return nil, handlePlaylistServiceErrors(err)
				}

				res := GetPlaylistItemIds{
					Ids: []string{},
				}

				if ids != nil {
					res.Ids = ids
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
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

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
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

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
				ctx := context.Background()

				body, err := pyrin.Body[ReorderPlaylistItemsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

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
	)
}
