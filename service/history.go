package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var historyErr = NewServiceErrCreator("history")

var (
	ErrHistoryServiceHistoryNotFound = historyErr.New("history not found")
)

type HistoryService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewHistoryService(
	logger *slog.Logger,
	db *database.Database,
) *HistoryService {
	return &HistoryService{
		logger: logger,
		db:     db,
	}
}

type GetTrackHistoryParams struct {
	UserId string
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *HistoryService) GetTrackHistory(
	ctx context.Context,
	params GetTrackHistoryParams,
) ([]database.TrackHistory, types.Page, error) {
	items, page, err := s.db.GetTrackHistory(
		ctx,
		database.GetTrackHistoryParams{
			UserId: params.UserId,
			Page:   params.Page,
			Filter: params.Filter,
		},
	)
	if err != nil {
		// TODO(patrik): Cleanup, when new filter system is used
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "history service",
				Message: err.Error(),
			}
		}

		// TODO(patrik): Cleanup, when new filter system is used
		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "history service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, historyErr.Wrap("get track history", err)
	}

	return items, page, nil
}

type PushTrackHistoryParams struct {
	UserId        string
	TrackId       string
	PlaybackType  string
	PercentPlayed int
}

func (s *HistoryService) PushTrackHistory(
	ctx context.Context,
	params PushTrackHistoryParams,
) (string, error) {
	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", nil
		}

		return "", historyErr.Wrap("push track history: get track", err)
	}

	if params.PercentPlayed < 10 {
		return "", nil
	}

	status := "skipped"
	if params.PercentPlayed >= 80 {
		status = "completed"
	}

	now := time.Now()

	id, err := s.db.CreateTrackHistory(ctx, database.CreateTrackHistoryParams{
		UserId:        params.UserId,
		TrackId:       params.TrackId,
		ListenedAt:    now.UnixMilli(),
		PlaybackType:  params.PlaybackType,
		Status:        status,
		PercentPlayed: params.PercentPlayed,
	})
	if err != nil {
		return "", historyErr.Wrap("push track history", err)
	}

	playTimeDelta := int64(
		float64(track.Duration) * float64(params.PercentPlayed) / 100.0,
	)

	skipDelta := 0
	if status == "skipped" {
		skipDelta = 1
	}

	year := now.Year()
	month := int(now.Month())
	quarter := (month-1)/3 + 1

	periods := []struct {
		periodType  string
		year        int
		periodValue int
	}{
		{"all", 0, 0},
		{"year", year, 0},
		{"quarter", year, quarter},
		{"month", year, month},
	}

	for _, p := range periods {
		err := s.db.UpsertUserTrackStats(
			ctx,
			database.UpsertUserTrackStatsParams{
				UserId:  params.UserId,
				TrackId: params.TrackId,

				PeriodType:  p.periodType,
				Year:        p.year,
				PeriodValue: p.periodValue,

				SkipDelta:     skipDelta,
				PlayTimeDelta: playTimeDelta,
			},
		)
		if err != nil {
			return "", historyErr.Wrap("push track history: upsert stats", err)
		}
	}

	err = s.db.IncrementUserStats(ctx, database.IncrementUserStatsParams{
		UserId:             params.UserId,
		SkipDelta:          skipDelta,
		ListeningTimeDelta: playTimeDelta,
		LastListenedAt:     now.UnixMilli(),
	})
	if err != nil {
		return "", historyErr.Wrap(
			"push track history: increment user stats", err)
	}

	return id, nil
}

type GetTrackHistoryByIdParams struct {
	HistoryId string
	UserId    string
}

func (s *HistoryService) GetTrackHistoryById(
	ctx context.Context,
	params GetTrackHistoryByIdParams,
) (database.TrackHistory, error) {
	history, err := s.db.GetTrackHistoryById(ctx, params.HistoryId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.TrackHistory{}, ErrHistoryServiceHistoryNotFound
		}

		return database.TrackHistory{}, historyErr.Wrap(
			"get track history by id", err)
	}

	if history.UserId != params.UserId {
		return database.TrackHistory{}, ErrHistoryServiceHistoryNotFound
	}

	return history, nil
}
