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

type GetHistoryParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *HistoryService) GetHistory(
	ctx context.Context,
	params GetHistoryParams,
) ([]database.UserTrackHistory, types.Page, error) {
	items, page, err := 	s.db.GetUserTrackHistory(ctx, database.GetUserTrackHistoryParams{
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "user history service",
				Message: err.Error(),
			}
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "user history service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, err
	}

	return items, page, nil
}

type GetHistoryByIdParams struct {
	HistoryId string
}

func (s *HistoryService) GetHistoryById(
	ctx context.Context,
	params GetHistoryByIdParams,
) (database.UserTrackHistory, error) {
	history, err := 	s.db.GetUserTrackHistoryById(ctx, params.HistoryId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.UserTrackHistory{}, ErrHistoryServiceHistoryNotFound
		}

		return database.UserTrackHistory{}, err
	}

	return history, nil
}
