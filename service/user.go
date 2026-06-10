package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/nanoteck137/tunebook/types"
)

var userErr = NewServiceErrCreator("user")

var (
	ErrUserServiceUserNotFound        = userErr.New("user not found")
	ErrUserServicePlaylistNotFound    = userErr.New("playlist not found")
	ErrUserServiceTrackNotFound       = userErr.New("track not found")
	ErrUserServiceTrackFilterNotFound = userErr.New("track filter not found")
	ErrUserServiceApiTokenNotFound    = userErr.New("api token not found")
	ErrUserServiceUnauthorized        = userErr.New("unauthorized")
)

type UserService struct {
	logger *slog.Logger

	db *database.Database

	imageService *ImageService
}

func NewUserService(
	logger *slog.Logger,
	db *database.Database,
	imageService *ImageService,
) *UserService {
	return &UserService{
		logger:       logger,
		db:           db,
		imageService: imageService,
	}
}

type GetUserByIdParams struct {
	UserId string
}

func (s *UserService) GetUserById(
	ctx context.Context,
	params GetUserByIdParams,
) (database.User, error) {
	user, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.User{}, ErrUserServiceUserNotFound
		}

		return database.User{}, userErr.Wrap("get user by id", err)
	}

	return user, nil
}

type GetFavoriteTracksParams struct {
	UserId string
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *UserService) GetFavoriteTracks(
	ctx context.Context,
	params GetFavoriteTracksParams,
) ([]database.UserFavoriteTrack, types.Page, error) {
	tracks, page, err := s.db.GetUserFavoriteTracks(
		ctx,
		database.GetUserFavoriteTracksParams{
			UserId: params.UserId,
			Page:   params.Page,
			Filter: params.Filter,
		},
	)
	if err != nil {
		return nil, types.Page{}, userErr.Wrap("get favorite tracks: db get", err)
	}

	for i := range tracks {
		tracks[i].Track.Order = utils.Pointer((i + 1) + (page.Page * page.PerPage))
	}

	return tracks, page, nil
}

type GetTrackFiltersParams struct {
	UserId string
}

func (s *UserService) GetTrackFilters(
	ctx context.Context,
	params GetTrackFiltersParams,
) ([]database.TrackFilter, error) {
	filters, err := s.db.GetTrackFiltersByUserId(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get track filters: db get", err)
	}

	return filters, nil
}

type UpdateMeParams struct {
	UserId string

	DisplayName *string
}

func (s *UserService) UpdateMe(
	ctx context.Context,
	params UpdateMeParams,
) error {
	user, err := s.GetUserById(ctx, GetUserByIdParams{
		UserId: params.UserId,
	})
	if err != nil {
		return userErr.Wrap("update me: get user", err)
	}

	changes := database.UserChanges{}

	if params.DisplayName != nil {
		changes.DisplayName = database.Change[string]{
			Value:   *params.DisplayName,
			Changed: *params.DisplayName != user.DisplayName,
		}
	}

	err = s.db.UpdateUser(ctx, user.Id, changes)
	if err != nil {
		return userErr.Wrap("update me: db update", err)
	}

	return nil
}

type SetQuickPlaylistParams struct {
	UserId     string
	PlaylistId string
}

func (s *UserService) SetQuickPlaylist(
	ctx context.Context,
	params SetQuickPlaylistParams,
) error {
	if params.PlaylistId != "" {
		_, err := s.db.GetPlaylistById(ctx, params.PlaylistId)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return ErrUserServicePlaylistNotFound
			}

			return userErr.Wrap("set quick playlist: db get playlist", err)
		}
	}

	err := s.db.UpdateUserSettings(ctx, database.UserSettings{
		Id: params.UserId,
		QuickPlaylist: sql.NullString{
			String: params.PlaylistId,
			Valid:  params.PlaylistId != "",
		},
	})
	if err != nil {
		return userErr.Wrap("set quick playlist: db update settings", err)
	}

	return nil
}

type GetFavoriteTrackIdsParams struct {
	UserId string
}

func (s *UserService) GetFavoriteTrackIds(
	ctx context.Context,
	params GetFavoriteTrackIdsParams,
) ([]string, error) {
	items, err := s.db.GetUserFavoritesIds(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get favorite track ids: db get", err)
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

func (s *UserService) FavoriteTrack(
	ctx context.Context,
	params FavoriteTrackParams,
) error {

	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceTrackNotFound
		}

		return userErr.Wrap("favorite track: db get track", err)
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

		return userErr.Wrap("favorite track: db create", err)
	}

	return nil
}

type UnfavoriteTrackParams struct {
	UserId  string
	TrackId string
}

func (s *UserService) UnfavoriteTrack(
	ctx context.Context,
	params UnfavoriteTrackParams,
) error {
	err := s.db.DeleteUserFavorite(ctx, params.UserId, params.TrackId)
	if err != nil {
		return userErr.Wrap("unfavorite track: db delete", err)
	}

	return nil
}

type CreateTrackFilterParams struct {
	UserId string

	Name   string
	Filter string
}

func (s *UserService) CreateTrackFilter(
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
		return "", userErr.Wrap("create track filter: db create", err)
	}

	return filterId, nil
}

type UpdateTrackFilterParams struct {
	FilterId string
	UserId   string

	Name   *string
	Filter *string
}

func (s *UserService) UpdateTrackFilter(
	ctx context.Context,
	params UpdateTrackFilterParams,
) error {
	filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceTrackFilterNotFound
		}

		return userErr.Wrap("update track filter: db get track filter", err)
	}

	if filter.UserId != params.UserId {
		return ErrUserServiceUnauthorized
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
		return userErr.Wrap("update track filter: db update", err)
	}

	return nil
}

type DeleteTrackFilterParams struct {
	FilterId string
	UserId   string
}

func (s *UserService) DeleteTrackFilter(
	ctx context.Context,
	params DeleteTrackFilterParams,
) error {
	filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceTrackFilterNotFound
		}

		return userErr.Wrap("delete track filter: db get track filter", err)
	}

	if filter.UserId != params.UserId {
		return ErrUserServiceUnauthorized
	}

	err = s.db.DeleteTrackFilter(ctx, filter.Id)
	if err != nil {
		return userErr.Wrap("delete track filter: db delete", err)
	}

	return nil
}

type GetApiTokensParams struct {
	UserId string
}

func (s *UserService) GetApiTokens(
	ctx context.Context,
	params GetApiTokensParams,
) ([]database.ApiToken, error) {
	tokens, err := s.db.GetAllApiTokensForUser(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get api tokens: db get", err)
	}

	return tokens, nil
}

type CreateApiTokenParams struct {
	UserId string

	Name string
}

func (s *UserService) CreateApiToken(
	ctx context.Context,
	params CreateApiTokenParams,
) (string, error) {
	id, err := s.db.CreateApiToken(ctx, database.CreateApiTokenParams{
		UserId: params.UserId,
		Name:   params.Name,
	})
	if err != nil {
		return "", userErr.Wrap("create api token: db create", err)
	}

	return id, nil
}

type DeleteApiTokenParams struct {
	TokenId string
	UserId  string
}

func (s *UserService) DeleteApiToken(
	ctx context.Context,
	params DeleteApiTokenParams,
) error {
	token, err := s.db.GetApiTokenById(ctx, params.TokenId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceApiTokenNotFound
		}

		return userErr.Wrap("delete api token: db get api token", err)
	}

	if token.UserId != params.UserId {
		return ErrUserServiceUnauthorized
	}

	err = s.db.DeleteApiToken(ctx, token.Id)
	if err != nil {
		return userErr.Wrap("delete api token: db delete", err)
	}

	return nil
}
