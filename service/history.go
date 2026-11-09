package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var historyErr = NewServiceErrCreator("history")

var (
	ErrHistoryServiceHistoryNotFound = historyErr.New("history not found")

	ErrHistoryServiceInvalidHistoryEvent = historyErr.New("invalid history event")
	ErrHistoryServiceInvalidListenedAt   = historyErr.New("invalid listened at")
	ErrHistoryServiceInvalidPlaybackType = historyErr.New("invalid playback type")
	ErrHistoryServiceInvalidSource       = historyErr.New("invalid source")
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
	Query  types.QueryParams

	ListenedAfter  int64
	ListenedBefore int64
}

func (s *HistoryService) GetTrackHistory(
	ctx context.Context,
	params GetTrackHistoryParams,
) ([]database.TrackHistory, types.Page, error) {
	items, page, err := s.db.GetTrackHistory(
		ctx,
		database.GetTrackHistoryParams{
			UserId:         params.UserId,
			Page:           params.Page,
			Query:          params.Query,
			ListenedAfter:  params.ListenedAfter,
			ListenedBefore: params.ListenedBefore,
		},
	)
	if err != nil {
		return nil, types.Page{}, historyErr.Wrap("get track history", err)
	}

	return items, page, nil
}

type PushTrackHistoryParams struct {
	UserId        string
	TrackId       string
	PlaybackType  string
	PercentPlayed int

	ClientEventId string
	ListenedAt    int64
	Source        string
}

type PushTrackHistoryResult struct {
	Id            string
	Duplicate     bool
	Status        string
	ListenedAt    int64
	PercentPlayed int
	Source        string
	Year          int
	Reason        string
}

func resultFromExisting(
	history database.TrackHistory,
	duplicate bool,
) PushTrackHistoryResult {
	return PushTrackHistoryResult{
		Id:            history.Id,
		Duplicate:     duplicate,
		Status:        history.Status,
		ListenedAt:    history.ListenedAt,
		PercentPlayed: history.PercentPlayed,
		Source:        history.Source,
		Year:          time.UnixMilli(history.ListenedAt).In(time.Local).Year(),
	}
}

func (s *HistoryService) PushTrackHistory(
	ctx context.Context,
	params PushTrackHistoryParams,
) (PushTrackHistoryResult, error) {
	if params.TrackId == "" {
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidHistoryEvent
	}

	percentPlayed := params.PercentPlayed
	if percentPlayed < 0 || percentPlayed > 100 {
		s.logger.Warn(
			"clamping percent played",
			slog.Int("percentPlayed", percentPlayed),
		)

		if percentPlayed < 0 {
			percentPlayed = 0
		}
		if percentPlayed > 100 {
			percentPlayed = 100
		}
	}

	playbackType := params.PlaybackType
	if playbackType == "" {
		playbackType = "normal"
	}
	switch playbackType {
	case "normal", "radio":
	default:
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidPlaybackType
	}

	source := params.Source
	switch source {
	case "web", "android":
	default:
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidSource
	}

	if utf8.RuneCountInString(params.ClientEventId) > 128 {
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidHistoryEvent
	}

	now := time.Now()

	listenedAt := params.ListenedAt
	if listenedAt == 0 {
		listenedAt = now.UnixMilli()
	}
	if listenedAt > now.Add(5*time.Minute).UnixMilli() {
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidListenedAt
	}
	if listenedAt < now.AddDate(-10, 0, 0).UnixMilli() {
		return PushTrackHistoryResult{}, ErrHistoryServiceInvalidListenedAt
	}

	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return PushTrackHistoryResult{
				Status: "ignored",
				Reason: "track_not_found",
				Source: source,
			}, nil
		}

		return PushTrackHistoryResult{}, historyErr.Wrap(
			"push track history: get track", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return PushTrackHistoryResult{}, historyErr.Wrap(
			"push track history: begin transaction", err)
	}
	defer tx.Rollback()

	// Check for an existing event before applying the threshold so a retry
	// still returns the original row even if the retried payload is odd.
	if params.ClientEventId != "" {
		existing, err := tx.GetTrackHistoryByClientEventId(
			ctx, params.UserId, params.ClientEventId,
		)
		if err == nil {
			return resultFromExisting(existing, true), nil
		}
		if !errors.Is(err, database.ErrItemNotFound) {
			return PushTrackHistoryResult{}, historyErr.Wrap(
				"push track history: get existing event", err)
		}
	}

	if percentPlayed < 10 {
		return PushTrackHistoryResult{
			Status: "ignored",
			Reason: "below_threshold",
			Source: source,
		}, nil
	}

	status := "skipped"
	if percentPlayed >= 80 {
		status = "completed"
	}

	id, err := tx.CreateTrackHistory(ctx, database.CreateTrackHistoryParams{
		UserId:        params.UserId,
		TrackId:       params.TrackId,
		ClientEventId: params.ClientEventId,
		Source:        source,
		ListenedAt:    listenedAt,
		PlaybackType:  playbackType,
		Status:        status,
		PercentPlayed: percentPlayed,
		Created:       now.UnixMilli(),
		Updated:       now.UnixMilli(),
	})
	if err != nil {
		// Concurrent duplicate with the same clientEventId.
		if params.ClientEventId != "" &&
			errors.Is(err, database.ErrItemAlreadyExists) {
			tx.Rollback()

			existing, getErr := s.db.GetTrackHistoryByClientEventId(
				ctx, params.UserId, params.ClientEventId,
			)
			if getErr == nil {
				return resultFromExisting(existing, true), nil
			}
		}

		return PushTrackHistoryResult{}, historyErr.Wrap(
			"push track history: create history", err)
	}

	// Bucket stats by the event timestamp, not the row creation time.
	eventTime := time.UnixMilli(listenedAt).In(time.Local)
	year := eventTime.Year()
	month := int(eventTime.Month())
	quarter := (month-1)/3 + 1

	playTimeDelta := int64(
		float64(track.Duration) * float64(percentPlayed) / 100.0,
	)

	skipDelta := 0
	if status == "skipped" {
		skipDelta = 1
	}

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
		err := tx.UpsertUserTrackStats(ctx, database.UpsertUserTrackStatsParams{
			UserId:          params.UserId,
			TrackId:         params.TrackId,
			PeriodType:      p.periodType,
			Year:            p.year,
			PeriodValue:     p.periodValue,
			SkipDelta:       skipDelta,
			PlayTimeDelta:   playTimeDelta,
			CompletionDelta: percentPlayed,
		})
		if err != nil {
			return PushTrackHistoryResult{}, historyErr.Wrap(
				"push track history: upsert stats", err)
		}
	}

	err = tx.IncrementUserStats(ctx, database.IncrementUserStatsParams{
		UserId:             params.UserId,
		SkipDelta:          skipDelta,
		ListeningTimeDelta: playTimeDelta,
		LastListenedAt:     listenedAt,
	})
	if err != nil {
		return PushTrackHistoryResult{}, historyErr.Wrap(
			"push track history: increment user stats", err)
	}

	if err := tx.Commit(); err != nil {
		return PushTrackHistoryResult{}, historyErr.Wrap(
			"push track history: commit", err)
	}

	return PushTrackHistoryResult{
		Id:            id,
		Status:        status,
		ListenedAt:    listenedAt,
		PercentPlayed: percentPlayed,
		Source:        source,
		Year:          year,
	}, nil
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
