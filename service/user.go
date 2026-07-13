package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"path"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var userErr = NewServiceErrCreator("user")

var (
	ErrUserServiceUserNotFound        = userErr.New("user not found")
	ErrUserServicePlaylistNotFound    = userErr.New("playlist not found")
	ErrUserServiceApiTokenNotFound    = userErr.New("api token not found")
	ErrUserServiceUnauthorized        = userErr.New("unauthorized")
)

type UserService struct {
	logger *slog.Logger

	db *database.Database

	filesystem   *FilesystemService
	imageService *ImageService
}

func NewUserService(
	logger *slog.Logger,
	db *database.Database,
	filesystem *FilesystemService,
	imageService *ImageService,
) *UserService {
	return &UserService{
		logger:       logger,
		db:           db,
		filesystem:   filesystem,
		imageService: imageService,
	}
}

// TODO(patrik): This is part of the task/job
type UpdateUserStatsParams struct {
	UserId string
}

func (s *UserService) GetAllUsers(
	ctx context.Context,
) ([]database.User, error) {
	return s.db.GetAllUsers(ctx)
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

type GetUserImageParams struct {
	UserId      string
	Size        int
	ImageFormat types.ImageFormat
}

func (s *UserService) GetUserImage(
	ctx context.Context,
	params GetUserImageParams,
) (string, error) {
	user, err := s.GetUserById(ctx, GetUserByIdParams{UserId: params.UserId})
	if err != nil {
		return "", err
	}

	err = s.filesystem.EnsureUserImageCacheDirs(user.Id)
	if err != nil {
		return "", userErr.Wrap("get user image", err)
	}

	input := ""
	if user.Picture.Valid {
		dir := s.filesystem.UserDir(user.Id)
		input = path.Join(dir, user.Picture.String)
	}

	p, err := s.imageService.ProcessImage(ProcessImageParams{
		Input:       input,
		Default:     "default_album.png",
		OutputDir:   s.filesystem.UserImagePath(user.Id),
		Size:        params.Size,
		ImageFormat: params.ImageFormat,
	})
	if err != nil {
		return "", userErr.Wrap("get user image", err)
	}

	return p, nil
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

type GetUserTopTracksParams struct {
	UserId     string
	PeriodType string
	Year       int
	Limit      int
}

func (s *UserService) GetUserTopTracks(
	ctx context.Context,
	params GetUserTopTracksParams,
) ([]database.UserTopTrack, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, ErrUserServiceUserNotFound
		}

		return nil, userErr.Wrap("get user top tracks: user", err)
	}

	tracks, err := s.db.GetUserTopTracks(ctx, database.GetUserTopTracksParams{
		UserId:     params.UserId,
		PeriodType: params.PeriodType,
		Year:       params.Year,
		Limit:      params.Limit,
	})
	if err != nil {
		return nil, userErr.Wrap("get user top tracks: tracks", err)
	}

	return tracks, nil
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
		err = s.filesystem.ClearUserImageCache(user.Id)
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

	err = s.filesystem.ClearUserImageCache(user.Id)
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

func (s *UserService) RecalculateUserStats(
	ctx context.Context, 
	userId string,
) error {
	agg, err := s.db.GetUserTrackStatsAgg(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get track stats agg", err)
	}

	lastListenedAt, err := s.db.GetLastListenedAt(ctx, userId)
	if err != nil {
		return userErr.Wrap("recalculate stats: get last listened at", err)
	}

	numPlaylistsCreated, err := s.db.GetUserPlaylistCount(ctx, userId)
	if err != nil && !errors.Is(err, database.ErrItemNotFound) {
		return userErr.Wrap("recalculate stats: get playlist count", err)
	}

	numFavoriteTracks, err := s.db.GetUserFavoriteCount(ctx, userId)
	if err != nil && !errors.Is(err, database.ErrItemNotFound) {
		return userErr.Wrap("recalculate stats: get favorite count", err)
	}

	err = s.db.SetUserStats(ctx, database.SetUserStatsParams{
		UserId: userId,

		NumTracksPlayed:     agg.NumTracksPlayed,
		NumTracksSkipped:    agg.NumTracksSkipped,
		NumPlaylistsCreated: numPlaylistsCreated,
		NumFavoriteTracks:   numFavoriteTracks,
		ListeningTime:       agg.PlayTime,
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
