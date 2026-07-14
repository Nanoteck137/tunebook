package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var trackErr = NewServiceErrCreator("track")

var (
	ErrTrackServiceTrackNotFound  = trackErr.New("track not found")
	ErrTrackServiceFilterNotFound = trackErr.New("filter not found")
	ErrTrackServiceUnauthorized   = trackErr.New("unauthorized")
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

	FilterId string
}

func (s *TrackService) GetTracks(
	ctx context.Context,
	params GetTracksParams,
) ([]database.Track, types.Page, error) {
	if params.FilterId != "" {
		dbFilter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, types.Page{}, ErrTrackServiceFilterNotFound
			}

			return nil, types.Page{}, trackErr.Wrap("get filter", err)
		}

		params.Filter.Filter = dbFilter.Filter
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

		return nil, types.Page{}, trackErr.Wrap("get tracks", err)
	}

	for i := range tracks {
		tracks[i].Order = utils.Pointer((i + 1) + (page.Page * page.PerPage))
	}

	return tracks, page, nil
}

type GetTrackByIdParams struct {
	TrackId string
}

type GetFavoriteTracksParams struct {
	UserId   string
	Page     types.PageParams
	Filter   types.FilterParams
	FilterId string
}

func (s *TrackService) GetFavoriteTracks(
	ctx context.Context,
	params GetFavoriteTracksParams,
) ([]database.UserFavoriteTrack, types.Page, error) {
	if params.FilterId != "" {
		dbFilter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err == nil {
			params.Filter.Filter = dbFilter.Filter
		}
	}

	tracks, page, err := s.db.GetUserFavoriteTracks(
		ctx,
		database.GetUserFavoriteTracksParams{
			UserId: params.UserId,
			Page:   params.Page,
			Filter: params.Filter,
		},
	)
	if err != nil {
		return nil, types.Page{}, trackErr.Wrap(
			"get favorite tracks: db get", err)
	}

	for i := range tracks {
		tracks[i].Track.Order =
			utils.Pointer((i + 1) + (page.Page * page.PerPage))
	}

	return tracks, page, nil
}

type GetFavoriteTrackIdsParams struct {
	UserId string
}

func (s *TrackService) GetFavoriteTrackIds(
	ctx context.Context,
	params GetFavoriteTrackIdsParams,
) ([]string, error) {
	items, err := s.db.GetUserFavoritesIds(ctx, params.UserId)
	if err != nil {
		return nil, trackErr.Wrap("get favorite track ids: db get", err)
	}

	if items == nil {
		return []string{}, nil
	}

	return items, nil
}

type FavoriteTrackParams struct {
	UserId  string
	TrackId string
}

func (s *TrackService) FavoriteTrack(
	ctx context.Context,
	params FavoriteTrackParams,
) error {

	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrTrackServiceTrackNotFound
		}

		return trackErr.Wrap("favorite track: db get track", err)
	}

	err = s.db.CreateUserFavorite(
		ctx,
		database.CreateUserFavoriteParams{
			UserId:  params.UserId,
			TrackId: track.Id,
		},
	)
	if err != nil {
		if errors.Is(err, database.ErrItemAlreadyExists) {
			return nil
		}

		return trackErr.Wrap("favorite track: db create", err)
	}

	return nil
}

type UnfavoriteTrackParams struct {
	UserId  string
	TrackId string
}

func (s *TrackService) UnfavoriteTrack(
	ctx context.Context,
	params UnfavoriteTrackParams,
) error {
	err := s.db.DeleteUserFavorite(ctx, params.UserId, params.TrackId)
	if err != nil {
		return trackErr.Wrap("unfavorite track: db delete", err)
	}

	return nil
}

type GetTrackFiltersParams struct {
	UserId string
}

func (s *TrackService) GetTrackFilters(
	ctx context.Context,
	params GetTrackFiltersParams,
) ([]database.TrackFilter, error) {
	filters, err := s.db.GetTrackFiltersByUserId(ctx, params.UserId)
	if err != nil {
		return nil, trackErr.Wrap("get track filters: db get", err)
	}

	return filters, nil
}

type CreateTrackFilterParams struct {
	UserId string

	Name   string
	Filter string
}

func (s *TrackService) CreateTrackFilter(
	ctx context.Context,
	params CreateTrackFilterParams,
) (string, error) {
	// TODO(patrik): Test filter

	filterId, err := s.db.CreateTrackFilter(
		ctx,
		database.CreateTrackFilterParams{
			UserId: params.UserId,
			Name:   params.Name,
			Filter: params.Filter,
		},
	)
	if err != nil {
		return "", trackErr.Wrap("create track filter: db create", err)
	}

	return filterId, nil
}

type UpdateTrackFilterParams struct {
	FilterId string
	UserId   string

	Name   *string
	Filter *string
}

func (s *TrackService) UpdateTrackFilter(
	ctx context.Context,
	params UpdateTrackFilterParams,
) error {
	filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrTrackServiceFilterNotFound
		}

		return trackErr.Wrap("update track filter: db get track filter", err)
	}

	if filter.UserId != params.UserId {
		return ErrTrackServiceUnauthorized
	}

	changes := database.TrackFilterChanges{}

	if params.Name != nil {
		changes.Name = database.Change[string]{
			Value:   *params.Name,
			Changed: *params.Name != filter.Name,
		}
	}

	if params.Filter != nil {
		// TODO(patrik): Test filter

		changes.Filter = database.Change[string]{
			Value:   *params.Filter,
			Changed: *params.Filter != filter.Filter,
		}
	}

	err = s.db.UpdateTrackFilter(ctx, filter.Id, changes)
	if err != nil {
		return trackErr.Wrap("update track filter: db update", err)
	}

	return nil
}

type DeleteTrackFilterParams struct {
	FilterId string
	UserId   string
}

func (s *TrackService) DeleteTrackFilter(
	ctx context.Context,
	params DeleteTrackFilterParams,
) error {
	filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrTrackServiceFilterNotFound
		}

		return trackErr.Wrap("delete track filter: db get track filter", err)
	}

	if filter.UserId != params.UserId {
		return ErrTrackServiceUnauthorized
	}

	err = s.db.DeleteTrackFilter(ctx, filter.Id)
	if err != nil {
		return trackErr.Wrap("delete track filter: db delete", err)
	}

	return nil
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

		return database.Track{}, trackErr.Wrap("get track by id", err)
	}

	return track, nil
}
