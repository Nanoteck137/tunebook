package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var artistErr = NewServiceErrCreator("artist")

var (
	ErrArtistServiceArtistNotFound = artistErr.New("artist not found")
)

type ArtistService struct {
	logger      *slog.Logger
	db          *database.Database
	imageService *ImageService
	dataDir     types.DataDir
}

func NewArtistService(
	logger *slog.Logger,
	db *database.Database,
	imageService *ImageService,
	dataDir types.DataDir,
) *ArtistService {
	return &ArtistService{
		logger:       logger,
		db:           db,
		imageService: imageService,
		dataDir:      dataDir,
	}
}

type GetArtistsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *ArtistService) GetArtists(
	ctx context.Context,
	params GetArtistsParams,
) ([]database.Artist, types.Page, error) {
	artists, page, err := s.db.GetArtists(ctx, database.GetArtistsParams{
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "artist service",
				Message: err.Error(),
			}
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "artist service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, artistErr.Wrap("get artists", err)
	}

	return artists, page, nil
}

type GetArtistByIdParams struct {
	ArtistId string
}

func (s *ArtistService) GetArtistById(
	ctx context.Context,
	params GetArtistByIdParams,
) (database.Artist, error) {
	artist, err := s.db.GetArtistById(ctx, params.ArtistId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.Artist{}, ErrArtistServiceArtistNotFound
		}

		return database.Artist{}, artistErr.Wrap("get artist by id", err)
	}

	return artist, nil
}

type GetArtistAlbumsParams struct {
	ArtistId string
}

func (s *ArtistService) GetArtistAlbums(
	ctx context.Context,
	params GetArtistAlbumsParams,
) ([]database.Album, error) {
	artist, err := s.GetArtistById(ctx, GetArtistByIdParams{
		ArtistId: params.ArtistId,
	})
	if err != nil {
		return nil, err
	}

	albums, err := s.db.GetAlbumsByArtist(ctx, artist.Id)
	if err != nil {
		return nil, artistErr.Wrap("get artist albums", err)
	}

	return albums, nil
}

type GetArtistImageParams struct {
	ArtistId    string
	Size        int
	ImageFormat types.ImageFormat
}

func (s *ArtistService) GetArtistImage(
	ctx context.Context,
	params GetArtistImageParams,
) (string, error) {
	artist, err := s.GetArtistById(ctx, GetArtistByIdParams{ArtistId: params.ArtistId})
	if err != nil {
		return "", err
	}

	cacheDir := s.dataDir.CacheImages()

	if err := utils.CreateDirectories([]string{
		cacheDir.String(),
		cacheDir.Artists(),
		cacheDir.Artist(artist.Id),
	}); err != nil {
		return "", artistErr.Wrap("get artist image", err)
	}

	input := ""
	if artist.CoverArt.Valid {
		input = artist.CoverArt.String
	}

	p, err := s.imageService.ProcessImage(ProcessImageParams{
		Input:       input,
		Default:     "default_artist.png",
		OutputDir:   cacheDir.Artist(artist.Id),
		Size:        params.Size,
		ImageFormat: params.ImageFormat,
	})
	if err != nil {
		return "", artistErr.Wrap("get artist image", err)
	}

	return p, nil
}
