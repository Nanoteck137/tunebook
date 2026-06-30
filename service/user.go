package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
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

	db      *database.Database
	dataDir types.DataDir

	imageService *ImageService
}

func NewUserService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
	imageService *ImageService,
) *UserService {
	return &UserService{
		logger:       logger,
		db:           db,
		dataDir:      dataDir,
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

type GetUserStatsParams struct {
	UserId string
}

func (s *UserService) GetUserStats(
	ctx context.Context,
	params GetUserStatsParams,
) (database.UserStats, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.UserStats{}, ErrUserServiceUserNotFound
		}

		return database.UserStats{}, userErr.Wrap("get user stats", err)
	}

	stats, err := s.db.GetUserStats(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.UserStats{
				UserId: params.UserId,
			}, nil
		}

		return database.UserStats{}, userErr.Wrap("get user stats", err)
	}

	return stats, nil
}

type GetFavoriteTracksParams struct {
	UserId   string
	Page     types.PageParams
	Filter   types.FilterParams
	FilterId string
}

func (s *UserService) GetFavoriteTracks(
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
	PictureUrl  *string
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

	if params.PictureUrl != nil {
		picture, err := s.imageService.DownloadPictureForUser(
			ctx,
			DownloadPictureForUserParams{
				UserId: user.Id,
				Url:    *params.PictureUrl,
			},
		)
		if err != nil {
			return userErr.Wrap("update me: download picture", err)
		}

		changes.Picture = database.Change[sql.NullString]{
			Value: sql.NullString{
				String: picture,
				Valid:  picture != "",
			},
			Changed: picture != user.Picture.String,
		}
	}

	err = s.db.UpdateUser(ctx, user.Id, changes)
	if err != nil {
		return userErr.Wrap("update me: db update", err)
	}

	if params.PictureUrl != nil {
		err = os.RemoveAll(s.dataDir.CacheImages().User(user.Id))
		if err != nil {
			return userErr.Wrap("update me: remove cache", err)
		}
	}

	return nil
}

type UploadUserImageParams struct {
	UserId string

	File *multipart.FileHeader
}

func (s *UserService) UploadUserImage(
	ctx context.Context,
	params UploadUserImageParams,
) error {
	user, err := s.GetUserById(ctx, GetUserByIdParams{
		UserId: params.UserId,
	})
	if err != nil {
		return userErr.Wrap("upload user image: get user", err)
	}

	picture, err := s.imageService.UploadImageForUser(
		ctx,
		UploadImageForUserParams{
			UserId: user.Id,
			File:   params.File,
		},
	)
	if err != nil {
		return userErr.Wrap("upload user image: upload", err)
	}

	err = s.db.UpdateUser(ctx, user.Id, database.UserChanges{
		Picture: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: picture,
				Valid:  picture != "",
			},
			Changed: picture != user.Picture.String,
		},
	})
	if err != nil {
		return userErr.Wrap("upload user image: db update", err)
	}

	err = os.RemoveAll(s.dataDir.CacheImages().User(user.Id))
	if err != nil {
		return userErr.Wrap("upload user image: remove cache", err)
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

func (s *UserService) RecalculateUserStats(ctx context.Context, userId string) error {
	numTracksPlayed, err := s.db.GetCompletedPlayCount(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get completed play count", err)
	}

	numTracksSkipped, err := s.db.GetSkippedPlayCount(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get skipped play count", err)
	}

	listeningTimeMs, err := s.db.GetCompletedListeningTime(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get listening time", err)
	}

	lastListenedAt, err := s.db.GetLastListenedAt(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get last listened at", err)
	}

	numPlaylistsCreated, err := s.db.GetUserPlaylistCount(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get playlist count", err)
	}

	numFavoriteTracks, err := s.db.GetUserFavoriteCount(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get favorite count", err)
	}

	err = s.db.SetUserStats(ctx, database.SetUserStatsParams{
		UserId: userId,

		NumTracksPlayed:     numTracksPlayed,
		NumTracksSkipped:    numTracksSkipped,
		NumPlaylistsCreated: numPlaylistsCreated,
		NumFavoriteTracks:   numFavoriteTracks,
		ListeningTime:       listeningTimeMs,
		LastListenedAt:      lastListenedAt,
	})
	if err != nil {
		return userErr.Wrap("recalculate stats: set", err)
	}

	return nil
}

func (s *UserService) RecalculateAllUserStats(ctx context.Context) error {
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return userErr.Wrap("recalculate all stats: get users", err)
	}

	for _, user := range users {
		err := s.RecalculateUserStats(ctx, user.Id)
		if err != nil {
			return err
		}
	}

	return nil
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
