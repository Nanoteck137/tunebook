package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

type QueueItem struct {
	QueueItemId string `json:"queueItemId"`
	Track       Track  `json:"track"`
}

func ConvertDBQueueItem(
	c pyrin.Context,
	item database.QueueItemTrack,
) QueueItem {
	return QueueItem{
		QueueItemId: item.Id,
		Track:       ConvertDBTrack(c, item.Track),
	}
}

func ConvertDBQueueIdItem(item database.QueueItemEntry) QueueIdItem {
	return QueueIdItem{
		QueueItemId: item.QueueItemId,
		TrackId:     item.TrackId,
	}
}

type GetQueue struct {
	CurrentIndex int         `json:"currentIndex"`
	Items        []QueueItem `json:"items"`
	Page         types.Page  `json:"page"`
}

type GetQueueIds struct {
	CurrentIndex int           `json:"currentIndex"`
	Items        []QueueIdItem `json:"items"`
}

type QueueIdItem struct {
	QueueItemId string `json:"queueItemId"`
	TrackId     string `json:"trackId"`
}

type GetQueueItem struct {
	Item QueueItem `json:"item"`
}

type ReplaceQueueBody struct {
	TrackIds     []string `json:"trackIds"`
	CurrentIndex *int     `json:"currentIndex,omitempty"`
	Shuffle      bool     `json:"shuffle,omitempty"`
}

type AddQueueItemsBody struct {
	TrackIds []string `json:"trackIds"`
	Position string   `json:"position"`
}

type AddToQueueBody struct {
	Source              string   `json:"source"`
	SourceId            string   `json:"sourceId"`
	TrackIds            []string `json:"trackIds,omitempty"`
	Position            string   `json:"position"`
	Shuffle             bool     `json:"shuffle,omitempty"`
	CurrentIndex        *int     `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string  `json:"queueIndexToTrackId,omitempty"`
}

type AddAlbumToQueueBody struct {
	FilterId            string  `json:"filterId,omitempty"`
	Position            string  `json:"position"`
	Shuffle             bool    `json:"shuffle,omitempty"`
	CurrentIndex        *int    `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string `json:"queueIndexToTrackId,omitempty"`
}

type AddArtistToQueueBody struct {
	FilterId            string  `json:"filterId,omitempty"`
	Position            string  `json:"position"`
	Shuffle             bool    `json:"shuffle,omitempty"`
	CurrentIndex        *int    `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string `json:"queueIndexToTrackId,omitempty"`
}

type AddPlaylistToQueueBody struct {
	FilterId            string  `json:"filterId,omitempty"`
	Position            string  `json:"position"`
	Shuffle             bool    `json:"shuffle,omitempty"`
	CurrentIndex        *int    `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string `json:"queueIndexToTrackId,omitempty"`
}

type AddFavoritesToQueueBody struct {
	FilterId            string  `json:"filterId,omitempty"`
	Position            string  `json:"position"`
	Shuffle             bool    `json:"shuffle,omitempty"`
	CurrentIndex        *int    `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string `json:"queueIndexToTrackId,omitempty"`
}

type AddTracksToQueueBody struct {
	TrackIds            []string `json:"trackIds"`
	FilterId            string   `json:"filterId,omitempty"`
	Position            string   `json:"position"`
	Shuffle             bool     `json:"shuffle,omitempty"`
	CurrentIndex        *int     `json:"currentIndex,omitempty"`
	QueueIndexToTrackId *string  `json:"queueIndexToTrackId,omitempty"`
}

type SetQueuePositionBody struct {
	Index int `json:"index"`
}

func handleQueueServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrQueueServiceQueueNotFound):
		return QueueNotFound()
	case errors.Is(err, service.ErrQueueServiceItemNotFound):
		return QueueItemNotFound()
	}

	return err
}

func InstallQueueHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetQueue",
			Path:         "/queues/:queueId",
			Method:       http.MethodGet,
			ResponseType: GetQueue{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				result, err := app.QueueService().GetQueue(
					ctx,
					service.GetQueueParams{
						Page:    getPageParams(q, 50),
						QueueId: c.Param("queueId"),
						UserId:  user.Id,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				items := make([]QueueItem, len(result.Items))
				for i, item := range result.Items {
					items[i] = ConvertDBQueueItem(c, item)
				}

				return GetQueue{
					CurrentIndex: result.CurrentIndex,
					Items:        items,
					Page:         result.Page,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetQueueIds",
			Path:         "/queues/:queueId/ids",
			Method:       http.MethodGet,
			ResponseType: GetQueueIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				result, err := app.QueueService().GetQueueIds(
					ctx, c.Param("queueId"), user.Id)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				items := make([]QueueIdItem, len(result.Entries))
				for i, item := range result.Entries {
					items[i] = ConvertDBQueueIdItem(item)
				}

				return GetQueueIds{
					CurrentIndex: result.CurrentIndex,
					Items:        items,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetQueueItemAtIndex",
			Path:         "/queues/:queueId/items/:position",
			Method:       http.MethodGet,
			ResponseType: GetQueueItem{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				position, err := parseIntParam(c, "position")
				if err != nil {
					return nil, err
				}

				item, err := app.QueueService().GetQueueItemAtIndex(
					ctx,
					service.GetQueueItemAtIndexParams{
						QueueId: c.Param("queueId"),
						UserId:  user.Id,
						Index:   position,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return GetQueueItem{
					Item: ConvertDBQueueItem(c, item),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ReplaceQueue",
			Path:     "/queues/:queueId",
			Method:   http.MethodPut,
			BodyType: ReplaceQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[ReplaceQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				err = app.QueueService().ReplaceQueue(
					ctx,
					service.ReplaceQueueParams{
						QueueId:      c.Param("queueId"),
						UserId:       user.Id,
						TrackIds:     body.TrackIds,
						CurrentIndex: currentIndex,
						Shuffle:      body.Shuffle,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddQueueItems",
			Path:     "/queues/:queueId/items",
			Method:   http.MethodPost,
			BodyType: AddQueueItemsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddQueueItemsBody](c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().AddItems(
					ctx,
					service.AddItemsParams{
						QueueId:  c.Param("queueId"),
						UserId:   user.Id,
						TrackIds: body.TrackIds,
						Position: body.Position,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RemoveQueueItem",
			Path:   "/queues/:queueId/items/:itemId",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().RemoveItem(
					ctx,
					service.RemoveItemParams{
						QueueId: c.Param("queueId"),
						UserId:  user.Id,
						ItemId:  c.Param("itemId"),
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddAlbumToQueue",
			Path:     "/queues/:queueId/add/albums/:albumId",
			Method:   http.MethodPost,
			BodyType: AddAlbumToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddAlbumToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddAlbumToQueue(
					ctx,
					service.AddAlbumToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						AlbumId:             c.Param("albumId"),
						FilterId:            body.FilterId,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddArtistToQueue",
			Path:     "/queues/:queueId/add/artists/:artistId",
			Method:   http.MethodPost,
			BodyType: AddArtistToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddArtistToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddArtistToQueue(
					ctx,
					service.AddArtistToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						ArtistId:            c.Param("artistId"),
						FilterId:            body.FilterId,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddPlaylistToQueue",
			Path:     "/queues/:queueId/add/playlists/:playlistId",
			Method:   http.MethodPost,
			BodyType: AddPlaylistToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddPlaylistToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddPlaylistToQueue(
					ctx,
					service.AddPlaylistToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						PlaylistId:          c.Param("playlistId"),
						FilterId:            body.FilterId,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddFavoritesToQueue",
			Path:     "/queues/:queueId/add/favorites/:userId",
			Method:   http.MethodPost,
			BodyType: AddFavoritesToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddFavoritesToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddFavoritesToQueue(
					ctx,
					service.AddFavoritesToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						FavoriteUserId:      c.Param("userId"),
						FilterId:            body.FilterId,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddTracksToQueue",
			Path:     "/queues/:queueId/add/tracks",
			Method:   http.MethodPost,
			BodyType: AddTracksToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddTracksToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddTracksToQueue(
					ctx,
					service.AddTracksToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						TrackIds:            body.TrackIds,
						FilterId:            body.FilterId,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddToQueue",
			Path:     "/queues/:queueId/add",
			Method:   http.MethodPost,
			BodyType: AddToQueueBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddToQueueBody](c)
				if err != nil {
					return nil, err
				}

				currentIndex := 0
				if body.CurrentIndex != nil {
					currentIndex = *body.CurrentIndex
				}

				queueIndexToTrackId := ""
				if body.QueueIndexToTrackId != nil {
					queueIndexToTrackId = *body.QueueIndexToTrackId
				}

				err = app.QueueService().AddToQueue(
					ctx,
					service.AddToQueueParams{
						QueueId:             c.Param("queueId"),
						UserId:              user.Id,
						Source:              body.Source,
						SourceId:            body.SourceId,
						TrackIds:            body.TrackIds,
						Position:            body.Position,
						Shuffle:             body.Shuffle,
						CurrentIndex:        currentIndex,
						QueueIndexToTrackId: queueIndexToTrackId,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SetQueuePosition",
			Path:     "/queues/:queueId/position",
			Method:   http.MethodPatch,
			BodyType: SetQueuePositionBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[SetQueuePositionBody](c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().SetPosition(
					ctx,
					service.SetPositionParams{
						QueueId: c.Param("queueId"),
						UserId:  user.Id,
						Index:   body.Index,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearQueue",
			Path:   "/queues/:queueId",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().ClearQueue(
					ctx,
					service.ClearQueueParams{
						QueueId: c.Param("queueId"),
						UserId:  user.Id,
					},
				)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
