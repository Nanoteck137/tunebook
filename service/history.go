package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var (
	ErrHistoryServiceHistoryNotFound = errors.New("user history service: history not found")
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
	items, page, err := s.db.GetTrackHistory(ctx, database.GetTrackHistoryParams{
		UserId: params.UserId,
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "history service",
				Message: err.Error(),
			}
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "history service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, err
	}

	return items, page, nil
}

type GetTrackHistoryByIdParams struct {
	HistoryId string
	UserId    string
}

type PushTrackHistoryParams struct {
	UserId       string
	TrackId      string
	ListenedAt   int64
	PlaybackType string
	Status       string
}

func (s *HistoryService) PushTrackHistory(
	ctx context.Context,
	params PushTrackHistoryParams,
) (string, error) {
	id, err := s.db.CreateTrackHistory(ctx, database.CreateTrackHistoryParams{
		UserId:       params.UserId,
		TrackId:      params.TrackId,
		ListenedAt:   params.ListenedAt,
		PlaybackType: params.PlaybackType,
		Status:       params.Status,
	})
	if err != nil {
		return "", err
	}

	return id, nil
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

		return database.TrackHistory{}, err
	}

	if history.UserId != params.UserId {
		return database.TrackHistory{}, ErrHistoryServiceHistoryNotFound
	}

	return history, nil
}
