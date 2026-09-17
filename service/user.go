package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"path"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/tools/pretty"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var userErr = NewServiceErrCreator("user")

var (
	ErrUserServiceUserNotFound     = userErr.New("user not found")
	ErrUserServicePlaylistNotFound = userErr.New("playlist not found")
	ErrUserServiceApiTokenNotFound = userErr.New("api token not found")
	ErrUserServiceUnauthorized     = userErr.New("unauthorized")
)

type UserService struct {
	logger *slog.Logger

	db *database.Database

	filesystem   *FilesystemService
	imageService *ImageService

	emitter broker.EventEmitter
}

func NewUserService(
	logger *slog.Logger,
	db *database.Database,
	filesystem *FilesystemService,
	imageService *ImageService,
	emitter broker.EventEmitter,
) *UserService {
	return &UserService{
		logger:       logger,
		db:           db,
		filesystem:   filesystem,
		imageService: imageService,
		emitter:      emitter,
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
			stats = database.UserStats{
				UserId: params.UserId,
			}
		} else {
			return database.UserStats{}, userErr.Wrap("get user stats", err)
		}
	}

	numFavoriteTracks, err := s.db.GetUserFavoriteCount(ctx, params.UserId)
	if err != nil {
		if !errors.Is(err, database.ErrItemNotFound) {
			return database.UserStats{}, userErr.Wrap(
				"get user stats: favorite tracks", err)
		}
	} else {
		stats.NumFavoriteTracks = numFavoriteTracks
	}

	numPlaylistsCreated, err := s.db.GetUserPlaylistCount(ctx, params.UserId)
	if err != nil {
		if !errors.Is(err, database.ErrItemNotFound) {
			return database.UserStats{}, userErr.Wrap(
				"get user stats: playlists created", err)
		}
	} else {
		stats.NumPlaylistsCreated = numPlaylistsCreated
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

type GetAllUserYearReviewsParams struct {
	UserId string
}

func (s *UserService) GetAllUserYearReviewsParams(
	ctx context.Context,
	params GetAllUserYearReviewsParams,
) ([]database.UserYearReview, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, ErrUserServiceUserNotFound
		}

		return nil, userErr.Wrap("get all user year reviews: user", err)
	}

	res, err := s.db.GetUserYearReviews(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get all user year reviews: reviews", err)
	}

	pretty.Println(res)

	return res, nil
}

type GetUserYearReviewParams struct {
	UserId string
	Year   int
}

func (s *UserService) GetUserYearReview(
	ctx context.Context,
	params GetUserYearReviewParams,
) (database.UserYearReview, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.UserYearReview{}, ErrUserServiceUserNotFound
		}

		return database.UserYearReview{},
			userErr.Wrap("get user year review: user", err)
	}

	review, err := s.db.GetUserYearReview(ctx, params.UserId, params.Year)
	if err != nil {
		return database.UserYearReview{},
			userErr.Wrap("get user year review: reviews", err)
	}

	pretty.Println(review)

	return review, nil
}

type GetUserYearReviewTracksParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (s *UserService) GetUserYearReviewTracks(
	ctx context.Context,
	params GetUserYearReviewTracksParams,
) ([]database.UserYearReviewTrack, types.Page, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, types.Page{}, ErrUserServiceUserNotFound
		}

		return nil, types.Page{}, userErr.Wrap("get user year review tracks: user", err)
	}

	tracks, page, err := s.db.GetUserYearReviewTracks(
		ctx,
		database.GetUserYearReviewTracksParams{
			UserId: params.UserId,
			Year:   params.Year,
			Page:   params.Page,
			Query:  params.Query,
		},
	)
	if err != nil {
		return nil, types.Page{},
			userErr.Wrap("get user year review tracks: reviews", err)
	}

	for i, track := range tracks {
		tracks[i].Track.Order = utils.Pointer(track.Rank)
	}

	pretty.Println(tracks)
	pretty.Println(page)

	return tracks, page, nil
}

type GetUserYearReviewAlbumsParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (s *UserService) GetUserYearReviewAlbums(
	ctx context.Context,
	params GetUserYearReviewAlbumsParams,
) ([]database.UserYearReviewAlbum, types.Page, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, types.Page{}, ErrUserServiceUserNotFound
		}

		return nil, types.Page{}, 
			userErr.Wrap("get user year review albums: user", err)
	}

	albums, page, err := s.db.GetUserYearReviewAlbums(
		ctx,
		database.GetUserYearReviewAlbumsParams{
			UserId: params.UserId,
			Year:   params.Year,
			Page:   params.Page,
			Query:  params.Query,
		},
	)
	if err != nil {
		return nil, types.Page{},
			userErr.Wrap("get user year review albums: reviews", err)
	}

	pretty.Println(albums)
	pretty.Println(page)

	return albums, page, nil
}

type GetUserYearReviewArtistsParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (s *UserService) GetUserYearReviewArtists(
	ctx context.Context,
	params GetUserYearReviewArtistsParams,
) ([]database.UserYearReviewArtist, types.Page, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, types.Page{}, ErrUserServiceUserNotFound
		}

		return nil, types.Page{}, 
			userErr.Wrap("get user year review artists: user", err)
	}

	artists, page, err := s.db.GetUserYearReviewArtists(
		ctx,
		database.GetUserYearReviewArtistsParams{
			UserId: params.UserId,
			Year:   params.Year,
			Page:   params.Page,
			Query:  params.Query,
		},
	)
	if err != nil {
		return nil, types.Page{},
			userErr.Wrap("get user year review artists: reviews", err)
	}

	pretty.Println(artists)
	pretty.Println(page)

	return artists, page, nil
}

type GetUserYearStatsParams struct {
	UserId string
}

func (s *UserService) GetUserYearStats(
	ctx context.Context,
	params GetUserYearStatsParams,
) ([]database.UserYearReview, error) {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, ErrUserServiceUserNotFound
		}

		return nil, userErr.Wrap("get user year stats: user", err)
	}

	years, err := s.db.GetUserStatsYears(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get user year stats: years", err)
	}

	reviews, err := s.db.GetUserYearReviews(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get user year stats: reviews", err)
	}

	generated := make(map[int]bool, len(reviews))
	for _, review := range reviews {
		generated[review.Year] = true
	}

	for _, year := range years {
		if generated[year.Year] {
			continue
		}

		err = s.db.GenerateUserReview(ctx, params.UserId, year.Year)
		if err != nil {
			return nil, userErr.Wrap("get user year stats: generate review", err)
		}
	}

	reviews, err = s.db.GetUserYearReviews(ctx, params.UserId)
	if err != nil {
		return nil, userErr.Wrap("get user year stats: reviews", err)
	}

	return reviews, nil
}

type GenerateUserReviewParams struct {
	UserId string
	Year   int
}

func (s *UserService) GenerateUserReview(
	ctx context.Context,
	params GenerateUserReviewParams,
) error {
	_, err := s.db.GetUserById(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrUserServiceUserNotFound
		}

		return userErr.Wrap("generate user review: user", err)
	}

	err = s.db.GenerateUserReview(ctx, params.UserId, params.Year)
	if err != nil {
		return userErr.Wrap("generate user review", err)
	}

	return nil
}

type UserYearReviewTrack struct {
	Rank      int
	PlayCount int

	Track database.Track
}

type UserYearReviewAlbum struct {
	Rank      int
	PlayCount int

	Album database.Album
}

type UserYearReviewArtist struct {
	Rank      int
	PlayCount int

	Artist database.Artist
}

type UserYearReviewInnerTrack struct {
	Rank      int
	PlayCount int
	PlayTime  int64

	Track database.Track
}

type UserYearReviewMonthDetail struct {
	Month     int
	PlayCount int
	PlayTime  int64

	AvgCompletion float64
	SkipCount     int
	UniqueTracks  int
	FavoritePlays int

	TrackCount  int
	AlbumCount  int
	ArtistCount int

	Tracks  []UserYearReviewTrack
	Albums  []UserYearReviewAlbum
	Artists []UserYearReviewArtist

	Tags    []database.UserYearReviewMonthTag
	Decades []database.UserYearReviewMonthDecade
}

// type GetUserYearReviewParams struct {
// 	UserId string
// 	Year   int
// }

type GetUserYearReviewResult struct {
	Review database.UserYearReview

	Tracks  []UserYearReviewTrack
	Albums  []UserYearReviewAlbum
	Artists []UserYearReviewArtist

	TrackCount  int
	AlbumCount  int
	ArtistCount int

	ArtistTracks map[string][]UserYearReviewInnerTrack
	AlbumTracks  map[string][]UserYearReviewInnerTrack

	Tags    []database.UserYearReviewTag
	Decades []database.UserYearReviewDecade

	Months       []database.UserYearReviewMonth
	MonthDetails []UserYearReviewMonthDetail
}

func (s *UserService) ensureUserYearReview(
	ctx context.Context,
	userId string,
	year int,
) (database.UserYearReview, error) {
	_, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.UserYearReview{}, ErrUserServiceUserNotFound
		}

		return database.UserYearReview{}, userErr.Wrap(
			"get user year review: user", err)
	}

	review, err := s.db.GetUserYearReview(ctx, userId, year)
	if errors.Is(err, database.ErrItemNotFound) {
		err = s.db.GenerateUserReview(ctx, userId, year)
		if err != nil {
			return database.UserYearReview{}, userErr.Wrap(
				"get user year review: generate review", err)
		}

		review, err = s.db.GetUserYearReview(ctx, userId, year)
		if err != nil {
			return database.UserYearReview{}, userErr.Wrap(
				"get user year review: review after generate", err)
		}
	} else if err != nil {
		return database.UserYearReview{}, userErr.Wrap(
			"get user year review: review", err)
	}

	return review, nil
}

// func (s *UserService) GetUserYearReview(
// 	ctx context.Context,
// 	params GetUserYearReviewParams,
// ) (GetUserYearReviewResult, error) {
// 	review, err := s.ensureUserYearReview(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, err
// 	}
//
// 	tracks, err := s.db.GetUserYearReviewTracks(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: tracks", err)
// 	}
//
// 	albums, err := s.db.GetUserYearReviewAlbums(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: albums", err)
// 	}
//
// 	artists, err := s.db.GetUserYearReviewArtists(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: artists", err)
// 	}
//
// 	trackCount := len(tracks)
// 	albumCount := len(albums)
// 	artistCount := len(artists)
//
// 	if len(tracks) > database.TopYearReviewItems {
// 		tracks = tracks[:database.TopYearReviewItems]
// 	}
// 	if len(albums) > database.TopYearReviewItems {
// 		albums = albums[:database.TopYearReviewItems]
// 	}
// 	if len(artists) > database.TopYearReviewItems {
// 		artists = artists[:database.TopYearReviewItems]
// 	}
//
// 	months, err := s.db.GetUserYearReviewMonths(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: months", err)
// 	}
//
// 	artistTracks := make(map[string][]database.UserYearArtistTrack, len(artists))
// 	for _, a := range artists {
// 		rows, err := s.db.GetUserYearArtistTracks(
// 			ctx, params.UserId, params.Year, a.ArtistId)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: artist tracks", err)
// 		}
// 		artistTracks[a.ArtistId] = rows
// 	}
//
// 	albumTracks := make(map[string][]database.UserYearAlbumTrack, len(albums))
// 	for _, a := range albums {
// 		rows, err := s.db.GetUserYearAlbumTracks(
// 			ctx, params.UserId, params.Year, a.AlbumId)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: album tracks", err)
// 		}
// 		albumTracks[a.AlbumId] = rows
// 	}
//
// 	tags, err := s.db.GetUserYearReviewTags(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: tags", err)
// 	}
//
// 	decades, err := s.db.GetUserYearReviewDecades(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: decades", err)
// 	}
//
// 	monthDetails := make([]UserYearReviewMonthDetail, 12)
// 	for i := range monthDetails {
// 		monthDetails[i].Month = i + 1
// 	}
//
// 	for _, m := range months {
// 		if m.Month >= 1 && m.Month <= 12 {
// 			detail := &monthDetails[m.Month-1]
// 			detail.PlayCount = m.PlayCount
// 			detail.PlayTime = m.PlayTime
// 			detail.AvgCompletion = m.AvgCompletion
// 			detail.SkipCount = m.SkipCount
// 			detail.UniqueTracks = m.UniqueTracks
// 			detail.FavoritePlays = m.FavoritePlays
// 		}
// 	}
//
// 	monthDetailTracks := make([][]database.UserYearReviewMonthTrack, 12)
// 	monthDetailAlbums := make([][]database.UserYearReviewMonthAlbum, 12)
// 	monthDetailArtists := make([][]database.UserYearReviewMonthArtist, 12)
// 	monthDetailTags := make([][]database.UserYearReviewMonthTag, 12)
// 	monthDetailDecades := make([][]database.UserYearReviewMonthDecade, 12)
//
// 	monthTrackIds := make([]string, 0)
// 	monthAlbumIds := make([]string, 0)
// 	monthArtistIds := make([]string, 0)
// 	detailIdSeen := make(map[string]bool)
// 	recordDetailId := func(id string, ids *[]string) {
// 		if id != "" && !detailIdSeen[id] {
// 			detailIdSeen[id] = true
// 			*ids = append(*ids, id)
// 		}
// 	}
//
// 	for m := 1; m <= 12; m++ {
// 		tracks, err := s.db.GetUserYearReviewMonthTracks(
// 			ctx, params.UserId, params.Year, m)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: month tracks", err)
// 		}
// 		monthDetailTracks[m-1] = tracks
// 		for _, t := range tracks {
// 			recordDetailId(t.TrackId, &monthTrackIds)
// 		}
//
// 		albums, err := s.db.GetUserYearReviewMonthAlbums(
// 			ctx, params.UserId, params.Year, m)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: month albums", err)
// 		}
// 		monthDetailAlbums[m-1] = albums
// 		for _, a := range albums {
// 			recordDetailId(a.AlbumId, &monthAlbumIds)
// 		}
//
// 		artists, err := s.db.GetUserYearReviewMonthArtists(
// 			ctx, params.UserId, params.Year, m)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: month artists", err)
// 		}
// 		monthDetailArtists[m-1] = artists
// 		for _, a := range artists {
// 			recordDetailId(a.ArtistId, &monthArtistIds)
// 		}
//
// 		tags, err := s.db.GetUserYearReviewMonthTags(
// 			ctx, params.UserId, params.Year, m)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: month tags", err)
// 		}
// 		monthDetailTags[m-1] = tags
//
// 		decades, err := s.db.GetUserYearReviewMonthDecades(
// 			ctx, params.UserId, params.Year, m)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: month decades", err)
// 		}
// 		monthDetailDecades[m-1] = decades
// 	}
//
// 	detailLoadedTracks, err := s.db.GetTracksByIds(ctx, monthTrackIds)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: month tracks load", err)
// 	}
// 	detailLoadedTracksById := make(map[string]database.Track, len(detailLoadedTracks))
// 	for _, t := range detailLoadedTracks {
// 		detailLoadedTracksById[t.Id] = t
// 	}
//
// 	detailLoadedAlbums, err := s.db.GetAlbumsByIds(ctx, monthAlbumIds)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: month albums load", err)
// 	}
// 	detailLoadedAlbumsById := make(map[string]database.Album, len(detailLoadedAlbums))
// 	for _, a := range detailLoadedAlbums {
// 		detailLoadedAlbumsById[a.Id] = a
// 	}
//
// 	detailLoadedArtists, err := s.db.GetArtistsByIds(ctx, monthArtistIds)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: month artists load", err)
// 	}
// 	detailLoadedArtistsById := make(map[string]database.Artist, len(detailLoadedArtists))
// 	for _, a := range detailLoadedArtists {
// 		detailLoadedArtistsById[a.Id] = a
// 	}
//
// 	for m := 1; m <= 12; m++ {
// 		detail := &monthDetails[m-1]
//
// 		for _, t := range monthDetailTracks[m-1] {
// 			track, ok := detailLoadedTracksById[t.TrackId]
// 			if !ok {
// 				continue
// 			}
//
// 			detail.Tracks = append(detail.Tracks, UserYearReviewTrack{
// 				Rank:      t.Rank,
// 				PlayCount: t.PlayCount,
// 				Track:     track,
// 			})
// 		}
//
// 		for _, a := range monthDetailAlbums[m-1] {
// 			album, ok := detailLoadedAlbumsById[a.AlbumId]
// 			if !ok {
// 				continue
// 			}
//
// 			detail.Albums = append(detail.Albums, UserYearReviewAlbum{
// 				Rank:      a.Rank,
// 				PlayCount: a.PlayCount,
// 				Album:     album,
// 			})
// 		}
//
// 		for _, a := range monthDetailArtists[m-1] {
// 			artist, ok := detailLoadedArtistsById[a.ArtistId]
// 			if !ok {
// 				continue
// 			}
//
// 			detail.Artists = append(detail.Artists, UserYearReviewArtist{
// 				Rank:      a.Rank,
// 				PlayCount: a.PlayCount,
// 				Artist:    artist,
// 			})
// 		}
//
// 		detail.Tags = monthDetailTags[m-1]
// 		detail.Decades = monthDetailDecades[m-1]
//
// 		detail.TrackCount = len(monthDetailTracks[m-1])
// 		detail.AlbumCount = len(monthDetailAlbums[m-1])
// 		detail.ArtistCount = len(monthDetailArtists[m-1])
//
// 		if len(detail.Tracks) > database.TopYearReviewItems {
// 			detail.Tracks = detail.Tracks[:database.TopYearReviewItems]
// 		}
// 		if len(detail.Albums) > database.TopYearReviewItems {
// 			detail.Albums = detail.Albums[:database.TopYearReviewItems]
// 		}
// 		if len(detail.Artists) > database.TopYearReviewItems {
// 			detail.Artists = detail.Artists[:database.TopYearReviewItems]
// 		}
// 	}
//
// 	res := GetUserYearReviewResult{
// 		Review: review,
//
// 		Tracks:  make([]UserYearReviewTrack, 0, len(tracks)),
// 		Albums:  make([]UserYearReviewAlbum, 0, len(albums)),
// 		Artists: make([]UserYearReviewArtist, 0, len(artists)),
//
// 		TrackCount:  trackCount,
// 		AlbumCount:  albumCount,
// 		ArtistCount: artistCount,
//
// 		ArtistTracks: make(map[string][]UserYearReviewInnerTrack),
// 		AlbumTracks:  make(map[string][]UserYearReviewInnerTrack),
//
// 		Tags:    tags,
// 		Decades: decades,
//
// 		Months:       months,
// 		MonthDetails: monthDetails,
// 	}
//
// 	for _, t := range tracks {
// 		track, err := s.db.GetTrackById(ctx, t.TrackId)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: track", err)
// 		}
//
// 		res.Tracks = append(res.Tracks, UserYearReviewTrack{
// 			Rank:      t.Rank,
// 			PlayCount: t.PlayCount,
// 			Track:     track,
// 		})
// 	}
//
// 	for _, a := range albums {
// 		album, err := s.db.GetAlbumById(ctx, a.AlbumId)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: album", err)
// 		}
//
// 		res.Albums = append(res.Albums, UserYearReviewAlbum{
// 			Rank:      a.Rank,
// 			PlayCount: a.PlayCount,
// 			Album:     album,
// 		})
// 	}
//
// 	for _, a := range artists {
// 		artist, err := s.db.GetArtistById(ctx, a.ArtistId)
// 		if err != nil {
// 			return GetUserYearReviewResult{}, userErr.Wrap(
// 				"get user year review: artist", err)
// 		}
//
// 		res.Artists = append(res.Artists, UserYearReviewArtist{
// 			Rank:      a.Rank,
// 			PlayCount: a.PlayCount,
// 			Artist:    artist,
// 		})
// 	}
//
// 	trackIds := make([]string, 0)
// 	seen := make(map[string]bool)
// 	recordId := func(id string) {
// 		if !seen[id] {
// 			seen[id] = true
// 			trackIds = append(trackIds, id)
// 		}
// 	}
// 	for _, rows := range artistTracks {
// 		for _, t := range rows {
// 			recordId(t.TrackId)
// 		}
// 	}
// 	for _, rows := range albumTracks {
// 		for _, t := range rows {
// 			recordId(t.TrackId)
// 		}
// 	}
//
// 	loadedTracks, err := s.db.GetTracksByIds(ctx, trackIds)
// 	if err != nil {
// 		return GetUserYearReviewResult{}, userErr.Wrap(
// 			"get user year review: inner tracks", err)
// 	}
// 	loadedTracksById := make(map[string]database.Track, len(loadedTracks))
// 	for _, t := range loadedTracks {
// 		loadedTracksById[t.Id] = t
// 	}
//
// 	for artistId, rows := range artistTracks {
// 		for i, t := range rows {
// 			track, ok := loadedTracksById[t.TrackId]
// 			if !ok {
// 				continue
// 			}
// 			res.ArtistTracks[artistId] = append(
// 				res.ArtistTracks[artistId],
// 				UserYearReviewInnerTrack{
// 					Rank:      i + 1,
// 					PlayCount: t.PlayCount,
// 					PlayTime:  t.PlayTime,
// 					Track:     track,
// 				},
// 			)
// 		}
// 	}
//
// 	for albumId, rows := range albumTracks {
// 		for i, t := range rows {
// 			track, ok := loadedTracksById[t.TrackId]
// 			if !ok {
// 				continue
// 			}
// 			res.AlbumTracks[albumId] = append(
// 				res.AlbumTracks[albumId],
// 				UserYearReviewInnerTrack{
// 					Rank:      i + 1,
// 					PlayCount: t.PlayCount,
// 					PlayTime:  t.PlayTime,
// 					Track:     track,
// 				},
// 			)
// 		}
// 	}
//
// 	return res, nil
// }

type GetUserYearReviewTopTracksParams struct {
	UserId string
	Year   int
}

type GetUserYearReviewTopTracksResult struct {
	Tracks []UserYearReviewTrack
}

// func (s *UserService) GetUserYearReviewTopTracks(
// 	ctx context.Context,
// 	params GetUserYearReviewTopTracksParams,
// ) (GetUserYearReviewTopTracksResult, error) {
// 	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
// 		return GetUserYearReviewTopTracksResult{}, err
// 	}
//
// 	rows, err := s.db.GetUserYearReviewTracks(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewTopTracksResult{}, userErr.Wrap(
// 			"get user year review top tracks", err)
// 	}
//
// 	ids := make([]string, len(rows))
// 	for i, t := range rows {
// 		ids[i] = t.TrackId
// 	}
//
// 	loaded, err := s.db.GetTracksByIds(ctx, ids)
// 	if err != nil {
// 		return GetUserYearReviewTopTracksResult{}, userErr.Wrap(
// 			"get user year review top tracks: load", err)
// 	}
// 	loadedById := make(map[string]database.Track, len(loaded))
// 	for _, t := range loaded {
// 		loadedById[t.Id] = t
// 	}
//
// 	res := GetUserYearReviewTopTracksResult{
// 		Tracks: make([]UserYearReviewTrack, 0, len(rows)),
// 	}
// 	for _, t := range rows {
// 		track, ok := loadedById[t.TrackId]
// 		if !ok {
// 			continue
// 		}
//
// 		res.Tracks = append(res.Tracks, UserYearReviewTrack{
// 			Rank:      t.Rank,
// 			PlayCount: t.PlayCount,
// 			Track:     track,
// 		})
// 	}
//
// 	return res, nil
// }

type GetUserYearReviewTopAlbumsParams struct {
	UserId string
	Year   int
}

type GetUserYearReviewTopAlbumsResult struct {
	Albums []UserYearReviewAlbum
}

// func (s *UserService) GetUserYearReviewTopAlbums(
// 	ctx context.Context,
// 	params GetUserYearReviewTopAlbumsParams,
// ) (GetUserYearReviewTopAlbumsResult, error) {
// 	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
// 		return GetUserYearReviewTopAlbumsResult{}, err
// 	}
//
// 	rows, err := s.db.GetUserYearReviewAlbums(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewTopAlbumsResult{}, userErr.Wrap(
// 			"get user year review top albums", err)
// 	}
//
// 	ids := make([]string, len(rows))
// 	for i, a := range rows {
// 		ids[i] = a.AlbumId
// 	}
//
// 	loaded, err := s.db.GetAlbumsByIds(ctx, ids)
// 	if err != nil {
// 		return GetUserYearReviewTopAlbumsResult{}, userErr.Wrap(
// 			"get user year review top albums: load", err)
// 	}
// 	loadedById := make(map[string]database.Album, len(loaded))
// 	for _, a := range loaded {
// 		loadedById[a.Id] = a
// 	}
//
// 	res := GetUserYearReviewTopAlbumsResult{
// 		Albums: make([]UserYearReviewAlbum, 0, len(rows)),
// 	}
// 	for _, a := range rows {
// 		album, ok := loadedById[a.AlbumId]
// 		if !ok {
// 			continue
// 		}
//
// 		res.Albums = append(res.Albums, UserYearReviewAlbum{
// 			Rank:      a.Rank,
// 			PlayCount: a.PlayCount,
// 			Album:     album,
// 		})
// 	}
//
// 	return res, nil
// }

type GetUserYearReviewTopArtistsParams struct {
	UserId string
	Year   int
}

type GetUserYearReviewTopArtistsResult struct {
	Artists []UserYearReviewArtist
}

// func (s *UserService) GetUserYearReviewTopArtists(
// 	ctx context.Context,
// 	params GetUserYearReviewTopArtistsParams,
// ) (GetUserYearReviewTopArtistsResult, error) {
// 	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
// 		return GetUserYearReviewTopArtistsResult{}, err
// 	}
//
// 	rows, err := s.db.GetUserYearReviewArtists(ctx, params.UserId, params.Year)
// 	if err != nil {
// 		return GetUserYearReviewTopArtistsResult{}, userErr.Wrap(
// 			"get user year review top artists", err)
// 	}
//
// 	ids := make([]string, len(rows))
// 	for i, a := range rows {
// 		ids[i] = a.ArtistId
// 	}
//
// 	loaded, err := s.db.GetArtistsByIds(ctx, ids)
// 	if err != nil {
// 		return GetUserYearReviewTopArtistsResult{}, userErr.Wrap(
// 			"get user year review top artists: load", err)
// 	}
// 	loadedById := make(map[string]database.Artist, len(loaded))
// 	for _, a := range loaded {
// 		loadedById[a.Id] = a
// 	}
//
// 	res := GetUserYearReviewTopArtistsResult{
// 		Artists: make([]UserYearReviewArtist, 0, len(rows)),
// 	}
// 	for _, a := range rows {
// 		artist, ok := loadedById[a.ArtistId]
// 		if !ok {
// 			continue
// 		}
//
// 		res.Artists = append(res.Artists, UserYearReviewArtist{
// 			Rank:      a.Rank,
// 			PlayCount: a.PlayCount,
// 			Artist:    artist,
// 		})
// 	}
//
// 	return res, nil
// }

type GetUserYearReviewMonthTopTracksParams struct {
	UserId string
	Year   int
	Month  int
}

type GetUserYearReviewMonthTopTracksResult struct {
	Tracks []UserYearReviewTrack
}

func (s *UserService) GetUserYearReviewMonthTopTracks(
	ctx context.Context,
	params GetUserYearReviewMonthTopTracksParams,
) (GetUserYearReviewMonthTopTracksResult, error) {
	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
		return GetUserYearReviewMonthTopTracksResult{}, err
	}

	rows, err := s.db.GetUserYearReviewMonthTracks(
		ctx, params.UserId, params.Year, params.Month)
	if err != nil {
		return GetUserYearReviewMonthTopTracksResult{}, userErr.Wrap(
			"get user year review month top tracks", err)
	}

	ids := make([]string, len(rows))
	for i, t := range rows {
		ids[i] = t.TrackId
	}

	loaded, err := s.db.GetTracksByIds(ctx, ids)
	if err != nil {
		return GetUserYearReviewMonthTopTracksResult{}, userErr.Wrap(
			"get user year review month top tracks: load", err)
	}
	loadedById := make(map[string]database.Track, len(loaded))
	for _, t := range loaded {
		loadedById[t.Id] = t
	}

	res := GetUserYearReviewMonthTopTracksResult{
		Tracks: make([]UserYearReviewTrack, 0, len(rows)),
	}
	for _, t := range rows {
		track, ok := loadedById[t.TrackId]
		if !ok {
			continue
		}

		res.Tracks = append(res.Tracks, UserYearReviewTrack{
			Rank:      t.Rank,
			PlayCount: t.PlayCount,
			Track:     track,
		})
	}

	return res, nil
}

type GetUserYearReviewMonthTopAlbumsParams struct {
	UserId string
	Year   int
	Month  int
}

type GetUserYearReviewMonthTopAlbumsResult struct {
	Albums []UserYearReviewAlbum
}

func (s *UserService) GetUserYearReviewMonthTopAlbums(
	ctx context.Context,
	params GetUserYearReviewMonthTopAlbumsParams,
) (GetUserYearReviewMonthTopAlbumsResult, error) {
	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
		return GetUserYearReviewMonthTopAlbumsResult{}, err
	}

	rows, err := s.db.GetUserYearReviewMonthAlbums(
		ctx, params.UserId, params.Year, params.Month)
	if err != nil {
		return GetUserYearReviewMonthTopAlbumsResult{}, userErr.Wrap(
			"get user year review month top albums", err)
	}

	ids := make([]string, len(rows))
	for i, a := range rows {
		ids[i] = a.AlbumId
	}

	loaded, err := s.db.GetAlbumsByIds(ctx, ids)
	if err != nil {
		return GetUserYearReviewMonthTopAlbumsResult{}, userErr.Wrap(
			"get user year review month top albums: load", err)
	}
	loadedById := make(map[string]database.Album, len(loaded))
	for _, a := range loaded {
		loadedById[a.Id] = a
	}

	res := GetUserYearReviewMonthTopAlbumsResult{
		Albums: make([]UserYearReviewAlbum, 0, len(rows)),
	}
	for _, a := range rows {
		album, ok := loadedById[a.AlbumId]
		if !ok {
			continue
		}

		res.Albums = append(res.Albums, UserYearReviewAlbum{
			Rank:      a.Rank,
			PlayCount: a.PlayCount,
			Album:     album,
		})
	}

	return res, nil
}

type GetUserYearReviewMonthTopArtistsParams struct {
	UserId string
	Year   int
	Month  int
}

type GetUserYearReviewMonthTopArtistsResult struct {
	Artists []UserYearReviewArtist
}

func (s *UserService) GetUserYearReviewMonthTopArtists(
	ctx context.Context,
	params GetUserYearReviewMonthTopArtistsParams,
) (GetUserYearReviewMonthTopArtistsResult, error) {
	if _, err := s.ensureUserYearReview(ctx, params.UserId, params.Year); err != nil {
		return GetUserYearReviewMonthTopArtistsResult{}, err
	}

	rows, err := s.db.GetUserYearReviewMonthArtists(
		ctx, params.UserId, params.Year, params.Month)
	if err != nil {
		return GetUserYearReviewMonthTopArtistsResult{}, userErr.Wrap(
			"get user year review month top artists", err)
	}

	ids := make([]string, len(rows))
	for i, a := range rows {
		ids[i] = a.ArtistId
	}

	loaded, err := s.db.GetArtistsByIds(ctx, ids)
	if err != nil {
		return GetUserYearReviewMonthTopArtistsResult{}, userErr.Wrap(
			"get user year review month top artists: load", err)
	}
	loadedById := make(map[string]database.Artist, len(loaded))
	for _, a := range loaded {
		loadedById[a.Id] = a
	}

	res := GetUserYearReviewMonthTopArtistsResult{
		Artists: make([]UserYearReviewArtist, 0, len(rows)),
	}
	for _, a := range rows {
		artist, ok := loadedById[a.ArtistId]
		if !ok {
			continue
		}

		res.Artists = append(res.Artists, UserYearReviewArtist{
			Rank:      a.Rank,
			PlayCount: a.PlayCount,
			Artist:    artist,
		})
	}

	return res, nil
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

	s.emitter.EmitEventToUser(
		params.UserId,
		broker.QuickPlaylistChangedEvent{
			PlaylistId: params.PlaylistId,
		},
	)

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
