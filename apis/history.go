package apis

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/app"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/anvil"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/validate"
)

type TrackHistory struct {
	Id string `json:"id"`

	UserId string `json:"userId"`
	Track  Track  `json:"track"`

	ClientEventId string `json:"clientEventId,omitempty"`
	Source        string `json:"source"`

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

		ClientEventId: history.ClientEventId.String,
		Source:        history.Source,

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

	ClientEventId string `json:"clientEventId,omitempty"`
	ListenedAt    *int64 `json:"listenedAt,omitempty"`
	Source        string `json:"source"`
}

func (b *PushTrackHistoryBody) Transform() {
	b.TrackId = anvil.String(b.TrackId)
	b.PlaybackType = anvil.String(b.PlaybackType)
	b.ClientEventId = anvil.String(b.ClientEventId)
	b.Source = anvil.String(b.Source)
}

func (b PushTrackHistoryBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.TrackId, validate.Required),
		validate.Field(&b.PlaybackType, validate.In("normal", "radio")),
		validate.Field(&b.Source, validate.Required, validate.In("web", "android")),
		validate.Field(&b.ClientEventId, validate.RuneLength(0, 128)),
		validate.Field(&b.ListenedAt, validate.Min(int64(0))),
	)
}

type PushTrackHistory struct {
	Id            string `json:"id"`
	Duplicate     bool   `json:"duplicate"`
	Status        string `json:"status,omitempty"`
	ListenedAt    int64  `json:"listenedAt,omitempty"`
	PercentPlayed int    `json:"percentPlayed,omitempty"`
	Source        string `json:"source"`
}

type PushTrackHistoryBatchBody struct {
	Events []PushTrackHistoryBody `json:"events"`
}

func (b *PushTrackHistoryBatchBody) Transform() {
	for i := range b.Events {
		b.Events[i].Transform()
	}
}

func (b PushTrackHistoryBatchBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Events, validate.Required),
	)
}

type PushTrackHistoryBatchResult struct {
	ClientEventId string `json:"clientEventId,omitempty"`
	Id            string `json:"id,omitempty"`
	Result        string `json:"result"`
	Status        string `json:"status,omitempty"`
	ListenedAt    int64  `json:"listenedAt,omitempty"`
	PercentPlayed int    `json:"percentPlayed,omitempty"`
	Duplicate     bool   `json:"duplicate"`
	Reason        string `json:"reason,omitempty"`
	Error         string `json:"error,omitempty"`
	Source        string `json:"source,omitempty"`
}

type PushTrackHistoryBatch struct {
	Results []PushTrackHistoryBatchResult `json:"results"`
}

func handleHistoryServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrHistoryServiceHistoryNotFound):
		return HistoryNotFound()
	case errors.Is(err, service.ErrHistoryServiceInvalidHistoryEvent):
		return InvalidHistoryEvent()
	case errors.Is(err, service.ErrHistoryServiceInvalidListenedAt):
		return InvalidListenedAt()
	case errors.Is(err, service.ErrHistoryServiceInvalidPlaybackType):
		return InvalidPlaybackType()
	case errors.Is(err, service.ErrHistoryServiceInvalidSource):
		return InvalidSource()
	}

	var queryErr *database.QueryError
	if errors.As(err, &queryErr) {
		var filterErr, sortErr error
		if queryErr.Filter != nil {
			filterErr = queryErr.Filter
		}
		if queryErr.Sort != nil {
			sortErr = queryErr.Sort
		}
		return QueryError(filterErr, sortErr)
	}

	return err
}

func historyBatchErrorResult(
	clientEventId string,
	source string,
	err error,
) PushTrackHistoryBatchResult {
	result := PushTrackHistoryBatchResult{
		ClientEventId: clientEventId,
		Result:        "error",
		Source:        source,
	}

	switch {
	case errors.Is(err, service.ErrHistoryServiceInvalidHistoryEvent):
		result.Reason = "invalid_history_event"
		result.Error = "invalid history event"
	case errors.Is(err, service.ErrHistoryServiceInvalidListenedAt):
		result.Reason = "invalid_listened_at"
		result.Error = "invalid listenedAt"
	case errors.Is(err, service.ErrHistoryServiceInvalidPlaybackType):
		result.Reason = "invalid_playback_type"
		result.Error = "invalid playback type"
	case errors.Is(err, service.ErrHistoryServiceInvalidSource):
		result.Reason = "invalid_source"
		result.Error = "invalid source"
	default:
		slog.Error(
			"history batch event failed",
			slog.String("clientEventId", clientEventId),
			slog.String("error", err.Error()),
		)
		result.Reason = "internal_error"
		result.Error = "internal error"
	}

	return result
}

func processHistoryBatch(
	ctx context.Context,
	historyService *service.HistoryService,
	userId string,
	events []PushTrackHistoryBody,
) ([]PushTrackHistoryBatchResult, []int) {
	results := make([]PushTrackHistoryBatchResult, 0, len(events))
	years := map[int]struct{}{}

	for _, event := range events {
		listenedAt := int64(0)
		if event.ListenedAt != nil {
			listenedAt = *event.ListenedAt
		}

		res, err := historyService.PushTrackHistory(
			ctx,
			service.PushTrackHistoryParams{
				UserId:        userId,
				TrackId:       event.TrackId,
				PlaybackType:  event.PlaybackType,
				PercentPlayed: event.PercentPlayed,
				ClientEventId: event.ClientEventId,
				ListenedAt:    listenedAt,
				Source:        event.Source,
			},
		)
		if err != nil {
			results = append(results, historyBatchErrorResult(
				event.ClientEventId, event.Source, err))
			continue
		}

		if res.Status == "ignored" {
			results = append(results, PushTrackHistoryBatchResult{
				ClientEventId: event.ClientEventId,
				Result:        "ignored",
				Reason:        res.Reason,
				Source:        res.Source,
			})
			continue
		}

		resultType := "accepted"
		if res.Duplicate {
			resultType = "duplicate"
		} else {
			years[res.Year] = struct{}{}
		}

		results = append(results, PushTrackHistoryBatchResult{
			ClientEventId: event.ClientEventId,
			Id:            res.Id,
			Result:        resultType,
			Status:        res.Status,
			ListenedAt:    res.ListenedAt,
			PercentPlayed: res.PercentPlayed,
			Duplicate:     res.Duplicate,
			Source:        res.Source,
		})
	}

	affectedYears := make([]int, 0, len(years))
	for year := range years {
		affectedYears = append(affectedYears, year)
	}
	sort.Ints(affectedYears)

	return results, affectedYears
}

func InstallHistoryHandlers(app *app.App, group pyrin.Group) {
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

				since, err := parseEpochMsQuery(q, "since")
				if err != nil {
					return nil, err
				}

				before, err := parseEpochMsQuery(q, "before")
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				items, page, err := app.HistoryService().GetTrackHistory(
					ctx,
					service.GetTrackHistoryParams{
						UserId:         user.Id,
						Page:           getPageParams(q, 100),
						Query:          getQueryParams(q),
						ListenedAfter:  since,
						ListenedBefore: before,
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

				ctx, cancel := context.WithTimeout(
					c.Request().Context(), 10*time.Second)
				defer cancel()

				listenedAt := int64(0)
				if body.ListenedAt != nil {
					listenedAt = *body.ListenedAt
				}

				result, err := app.HistoryService().PushTrackHistory(
					ctx,
					service.PushTrackHistoryParams{
						UserId:        user.Id,
						TrackId:       body.TrackId,
						PlaybackType:  body.PlaybackType,
						PercentPlayed: body.PercentPlayed,
						ClientEventId: body.ClientEventId,
						ListenedAt:    listenedAt,
						Source:        body.Source,
					},
				)
				if err != nil {
					return nil, handleHistoryServiceErrors(err)
				}

				// An empty id means the play was not recorded (below the
				// minimum percent, or the track is gone), so nothing
				// changed and there is nothing to refresh.
				if result.Id != "" && !result.Duplicate {
					dispatchReviewForYear(
						context.WithoutCancel(ctx),
						app,
						user.Id,
						result.Year,
					)
				}

				return PushTrackHistory{
					Id:            result.Id,
					Duplicate:     result.Duplicate,
					Status:        result.Status,
					ListenedAt:    result.ListenedAt,
					PercentPlayed: result.PercentPlayed,
					Source:        result.Source,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "PushTrackHistoryBatch",
			Method:       http.MethodPost,
			Path:         "/history/tracks/batch",
			ResponseType: PushTrackHistoryBatch{},
			BodyType:     PushTrackHistoryBatchBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[PushTrackHistoryBatchBody](c)
				if err != nil {
					return nil, err
				}

				if len(body.Events) == 0 {
					return nil, InvalidHistoryEvent()
				}
				if len(body.Events) > 100 {
					return nil, HistoryBatchTooLarge()
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx, cancel := context.WithTimeout(
					c.Request().Context(), 30*time.Second)
				defer cancel()

				results, years := processHistoryBatch(
					ctx, app.HistoryService(), user.Id, body.Events)

				reviewCtx := context.WithoutCancel(ctx)
				for _, year := range years {
					dispatchReviewYear(reviewCtx, app, user.Id, year)
				}
				if len(years) > 0 {
					dispatchReviewTotal(reviewCtx, app, user.Id)
				}

				return PushTrackHistoryBatch{
					Results: results,
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
