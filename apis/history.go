package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

type TrackHistory struct {
	Id string `json:"id"`

	UserId string `json:"userId"`
	Track  Track  `json:"track"`

	ListenedAt   int64  `json:"listenedAt"`
	PlaybackType string `json:"playbackType"`
	Status       string `json:"status"`

	PercentPlayed int `json:"percentPlayed"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBTrackHistory(
	c pyrin.Context, 
	history database.TrackHistory,
) TrackHistory {
	return TrackHistory{
		Id: history.Id,

		UserId: history.UserId,
		Track:  ConvertDBTrack(c, history.Track),

		ListenedAt:   history.ListenedAt,
		PlaybackType: history.PlaybackType,
		Status:       history.Status,

		PercentPlayed: history.PercentPlayed,

		Created: formatTime(history.Created),
		Updated: formatTime(history.Updated),
	}
}

type GetTrackHistory struct {
	Page    types.Page     `json:"page"`
	History []TrackHistory `json:"history"`
}

type GetHistoryById struct {
	History TrackHistory `json:"history"`
}

type PushTrackHistoryBody struct {
	TrackId       string `json:"trackId"`
	PlaybackType  string `json:"playbackType"`
	PercentPlayed int    `json:"percentPlayed"`
}

type PushTrackHistory struct {
	Id string `json:"id"`
}

func handleHistoryServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrHistoryServiceHistoryNotFound):
		return HistoryNotFound()
	}

	var invalidFilter *service.InvalidFilterError
	if errors.As(err, &invalidFilter) {
		return InvalidFilter(errors.New(invalidFilter.Message))
	}

	var invalidSort *service.InvalidSortError
	if errors.As(err, &invalidSort) {
		return InvalidSort(errors.New(invalidSort.Message))
	}

	return err
}

func InstallHistoryHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTrackHistory",
			Method:       http.MethodGet,
			Path:         "/history/tracks",
			ResponseType: GetTrackHistory{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				items, page, err := app.HistoryService().GetTrackHistory(
					ctx,
					service.GetTrackHistoryParams{
						UserId: user.Id,
						Page:   getPageParams(q, 100),
						Filter: getFilterParams(q),
					},
				)
				if err != nil {
					return nil, handleHistoryServiceErrors(err)
				}

				res := GetTrackHistory{
					Page:    page,
					History: make([]TrackHistory, len(items)),
				}

				for i, item := range items {
					res.History[i] = ConvertDBTrackHistory(c, item)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "PushTrackHistory",
			Method:       http.MethodPost,
			Path:         "/history/tracks",
			ResponseType: PushTrackHistory{},
			BodyType:     PushTrackHistoryBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[PushTrackHistoryBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				id, err := app.HistoryService().PushTrackHistory(
					ctx,
					service.PushTrackHistoryParams{
						UserId:        user.Id,
						TrackId:       body.TrackId,
						PlaybackType:  body.PlaybackType,
						PercentPlayed: body.PercentPlayed,
					},
				)
				if err != nil {
					return nil, err
				}

				return PushTrackHistory{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackHistoryById",
			Method:       http.MethodGet,
			Path:         "/history/tracks/:id",
			ResponseType: GetHistoryById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				history, err := app.HistoryService().GetTrackHistoryById(
					ctx,
					service.GetTrackHistoryByIdParams{
						HistoryId: c.Param("id"),
						UserId:    user.Id,
					},
				)
				if err != nil {
					return nil, handleHistoryServiceErrors(err)
				}

				return GetHistoryById{
					History: ConvertDBTrackHistory(c, history),
				}, nil
			},
		},
	)
}
