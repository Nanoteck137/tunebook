package apis

import (
	"context"
	"errors"
	"net/http"
	"strconv"

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

func ConvertDBQueueItem(c pyrin.Context, item database.QueueItemTrack) QueueItem {
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
	CurrentIndex int             `json:"currentIndex"`
	Items        []QueueIdItem   `json:"items"`
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

type SetQueuePositionBody struct {
	Index int `json:"index"`
}

func handleQueueServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrQueueServiceQueueNotFound):
		return QueueNotFound()
	case errors.Is(err, service.ErrQueueServiceItemNotFound):
		// TODO(patrik): Replace with its own error
		return TrackNotFound()
	}

	return err
}

func InstallQueueHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetQueue",
			Path:         "/queue",
			Method:       http.MethodGet,
			ResponseType: GetQueue{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 50)

				result, err := app.QueueService().GetQueue(ctx, user.Id, pageParams.Page, pageParams.PerPage)
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
			Path:         "/queue/ids",
			Method:       http.MethodGet,
			ResponseType: GetQueueIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				result, err := app.QueueService().GetQueueIds(ctx, user.Id)
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
			Path:         "/queue/items/:position",
			Method:       http.MethodGet,
			ResponseType: GetQueueItem{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				position, err := strconv.Atoi(c.Param("position"))
				if err != nil {
					return nil, err
				}

				item, err := app.QueueService().GetQueueItemAtIndex(ctx, user.Id, position)
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
			Path:     "/queue",
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

				err = app.QueueService().ReplaceQueue(ctx, service.ReplaceQueueParams{
					UserId:       user.Id,
					TrackIds:     body.TrackIds,
					CurrentIndex: currentIndex,
					Shuffle:      body.Shuffle,
				})
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddQueueItems",
			Path:     "/queue/items",
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

				err = app.QueueService().AddItems(ctx, service.AddItemsParams{
					UserId:   user.Id,
					TrackIds: body.TrackIds,
					Position: body.Position,
				})
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RemoveQueueItem",
			Path:   "/queue/items/:itemId",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().RemoveItem(ctx, service.RemoveItemParams{
					UserId: user.Id,
					ItemId: c.Param("itemId"),
				})
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SetQueuePosition",
			Path:     "/queue/position",
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

				err = app.QueueService().SetPosition(ctx, service.SetPositionParams{
					UserId: user.Id,
					Index:  body.Index,
				})
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearQueue",
			Path:   "/queue",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				err = app.QueueService().ClearQueue(ctx, user.Id)
				if err != nil {
					return nil, handleQueueServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
