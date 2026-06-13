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
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *HistoryService) GetTrackHistory(
	ctx context.Context,
	params GetTrackHistoryParams,
) ([]database.TrackHistory, types.Page, error) {
	items, page, err := s.db.GetTrackHistory(ctx, database.GetTrackHistoryParams{
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

	return history, nil
}
