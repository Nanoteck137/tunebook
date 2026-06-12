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

type History struct {
	Id     string `json:"id"`

	UserId string `json:"userId"`
	TrackId string `json:"trackId"`

	ListenedAt   int64  `json:"listenedAt"`
	PlaybackType string `json:"playbackType"`
	Status       string `json:"status"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBHistory(c pyrin.Context, history database.UserTrackHistory) History {
	return History{
		Id:     history.Id,

		UserId: history.UserId,
		TrackId: history.TrackId,

		ListenedAt:   history.ListenedAt,
		PlaybackType: history.PlaybackType,
		Status:       history.Status,

		Created: formatTime(history.Created),
		Updated: formatTime(history.Updated),
	}
}

type GetHistory struct {
	Page    types.Page `json:"page"`
	History []History  `json:"history"`
}

type GetHistoryById struct {
	History History `json:"history"`
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
			Name:         "GetHistory",
			Method:       http.MethodGet,
			Path:         "/history",
			ResponseType: GetHistory{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				items, page, err := app.HistoryService().GetHistory(
					ctx,
					service.GetHistoryParams{
						Page:   pageParams,
						Filter: filterParams,
					},
				)
				if err != nil {
					return nil, handleHistoryServiceErrors(err)
				}

				res := GetHistory{
					Page:    page,
					History: make([]History, len(items)),
				}

				for i, item := range items {
					res.History[i] = ConvertDBHistory(c, item)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetHistoryById",
			Method:       http.MethodGet,
			Path:         "/history/:id",
			ResponseType: GetHistoryById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				history, err := app.HistoryService().GetHistoryById(
					ctx,
					service.GetHistoryByIdParams{
						HistoryId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, handleHistoryServiceErrors(err)
				}

				return GetHistoryById{
					History: ConvertDBHistory(c, history),
				}, nil
			},
		},
	)
}
