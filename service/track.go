package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

var (
	ErrTrackServiceTrackNotFound = errors.New("track service: track not found")
)

type TrackService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewTrackService(
	logger *slog.Logger,
	db *database.Database,
) *TrackService {
	return &TrackService{
		logger: logger,
		db:     db,
	}
}

type GetTracksParams struct {
	Page   types.PageParams
	Filter types.FilterParams

	UserId   string
	FilterId string
}

func (s *TrackService) GetTracks(
	ctx context.Context,
	params GetTracksParams,
) ([]database.Track, types.Page, error) {
	if params.UserId != "" && params.FilterId != "" {
		dbFilter, err := s.db.GetTrackFilterById(ctx, params.FilterId, params.UserId)
		if err == nil {
			params.Filter.Filter = dbFilter.Filter
		}
	}

	tracks, page, err := s.db.GetTracks(ctx, database.GetTracksParams{
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "track service",
				Message: err.Error(),
			}
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "track service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, err
	}

	for i := range tracks {
		tracks[i].Order = utils.IntPtr((i + 1) + (page.Page * page.PerPage))
	}

	return tracks, page, nil
}

type GetTrackByIdParams struct {
	TrackId string
}

func (s *TrackService) GetTrackById(
	ctx context.Context,
	params GetTrackByIdParams,
) (database.Track, error) {
	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.Track{}, ErrTrackServiceTrackNotFound
		}

		return database.Track{}, err
	}

	return track, nil
}
