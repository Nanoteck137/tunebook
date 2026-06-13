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

	UserId  string `json:"userId"`
	TrackId string `json:"trackId"`

	ListenedAt   int64  `json:"listenedAt"`
	PlaybackType string `json:"playbackType"`
	Status       string `json:"status"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBTrackHistory(c pyrin.Context, history database.TrackHistory) TrackHistory {
	return TrackHistory{
		Id: history.Id,

		UserId:  history.UserId,
		TrackId: history.TrackId,

		ListenedAt:   history.ListenedAt,
		PlaybackType: history.PlaybackType,
		Status:       history.Status,

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

				ctx := c.Request().Context()

				items, page, err := app.HistoryService().GetTrackHistory(
					ctx,
					service.GetTrackHistoryParams{
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
			Name:         "GetTrackHistoryById",
			Method:       http.MethodGet,
			Path:         "/history/tracks/:id",
			ResponseType: GetHistoryById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				history, err := app.HistoryService().GetTrackHistoryById(
					ctx,
					service.GetTrackHistoryByIdParams{
						HistoryId: c.Param("id"),
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
