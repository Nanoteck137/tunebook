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
	ErrTrackServiceTrackNotFound        = trackErr.New("track not found")
	ErrTrackServiceFilterNotFound       = trackErr.New("filter not found")
	ErrTrackServiceAnchorFilterNotFound = trackErr.New("anchor filter not found")
	ErrTrackServiceUnauthorized         = trackErr.New("unauthorized")
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
	Page  types.PageParams
	Query types.QueryParams

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

		params.Query.Filter = dbFilter.Filter
	}

	tracks, page, err := s.db.GetTracks(ctx, database.GetTracksParams{
		Page:  params.Page,
		Query: params.Query,
	})
	if err != nil {
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
	Query    types.QueryParams
	FilterId string
}

func (s *TrackService) GetFavoriteTracks(
	ctx context.Context,
	params GetFavoriteTracksParams,
) ([]database.UserFavoriteTrack, types.Page, error) {
	if params.FilterId != "" {
		dbFilter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err == nil {
			params.Query.Filter = dbFilter.Filter
		}
	}

	tracks, page, err := s.db.GetUserFavoriteTracks(
		ctx,
		database.GetUserFavoriteTracksParams{
			UserId: params.UserId,
			Page:   params.Page,
			Query:  params.Query,
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
	// if err := database.ValidateTrackFilter(params.Filter); err != nil {
	// 	var filterErr *database.FilterError
	// 	if errors.As(err, &filterErr) {
	// 		return "", &database.QueryError{Filter: filterErr}
	// 	}
	//
	// 	return "", trackErr.Wrap("create track filter: validate", err)
	// }

	position, err := s.db.GetNextTrackFilterPosition(ctx, params.UserId)
	if err != nil {
		return "", trackErr.Wrap("create track filter: db get next position", err)
	}

	filterId, err := s.db.CreateTrackFilter(
		ctx,
		database.CreateTrackFilterParams{
			UserId:   params.UserId,
			Name:     params.Name,
			Filter:   params.Filter,
			Position: position,
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
		// if err := database.ValidateTrackFilter(*params.Filter); err != nil {
		// 	var filterErr *database.FilterError
		// 	if errors.As(err, &filterErr) {
		// 		return &database.QueryError{Filter: filterErr}
		// 	}
		// 	return trackErr.Wrap("update track filter: validate", err)
		// }

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

	tx, err := s.db.Begin()
	if err != nil {
		return trackErr.Wrap("delete track filter: db begin", err)
	}
	defer tx.Rollback()

	err = tx.DeleteTrackFilter(ctx, filter.Id)
	if err != nil {
		return trackErr.Wrap("delete track filter: db delete", err)
	}

	err = tx.ReorderTrackFiltersAfterDelete(ctx, filter.UserId, filter.Position)
	if err != nil {
		return trackErr.Wrap("delete track filter: db reorder filters", err)
	}

	err = tx.Commit()
	if err != nil {
		return trackErr.Wrap("delete track filter: db commit", err)
	}

	return nil
}

type ReorderTrackFiltersParams struct {
	UserId string

	Before         bool
	AnchorFilterId string
	FilterIds      []string
}

func (s *TrackService) ReorderTrackFilters(
	ctx context.Context,
	params ReorderTrackFiltersParams,
) error {
	current, err := s.db.GetTrackFiltersByUserId(ctx, params.UserId)
	if err != nil {
		return trackErr.Wrap("reorder track filters: db get filters", err)
	}

	index := make(map[string]database.TrackFilter, len(current))
	for _, filter := range current {
		index[filter.Id] = filter
	}

	filters := make([]database.TrackFilter, 0, len(params.FilterIds))
	for _, id := range params.FilterIds {
		filter, ok := index[id]
		if !ok {
			continue
		}

		filters = append(filters, filter)
	}

	if len(filters) == 0 {
		return nil
	}

	if params.AnchorFilterId != "" {
		if _, ok := index[params.AnchorFilterId]; !ok {
			return ErrTrackServiceAnchorFilterNotFound
		}
	}

	moveSet := make(map[string]bool, len(filters))
	for _, filter := range filters {
		moveSet[filter.Id] = true
	}

	stationary := make([]database.TrackFilter, 0, len(current))
	for _, filter := range current {
		if !moveSet[filter.Id] {
			stationary = append(stationary, filter)
		}
	}

	insertAt := 0
	if params.AnchorFilterId != "" {
		for i, filter := range stationary {
			if filter.Id == params.AnchorFilterId {
				insertAt = i + 1
				break
			}
		}
	}

	spliced := make([]database.TrackFilter, 0, len(current))
	spliced = append(spliced, stationary[:insertAt]...)
	spliced = append(spliced, filters...)
	spliced = append(spliced, stationary[insertAt:]...)

	tx, err := s.db.Begin()
	if err != nil {
		return trackErr.Wrap("reorder track filters: db begin", err)
	}
	defer tx.Rollback()

	for i, filter := range spliced {
		err := tx.UpdateTrackFilter(
			ctx,
			filter.Id,
			database.TrackFilterChanges{
				Position: database.Change[int]{
					Value:   i,
					Changed: i != filter.Position,
				},
			},
		)
		if err != nil {
			return trackErr.Wrap("reorder track filters: db update filter", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return trackErr.Wrap("reorder track filters: db commit", err)
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
