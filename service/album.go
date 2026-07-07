package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var albumErr = NewServiceErrCreator("album")

var (
	ErrAlbumServiceAlbumNotFound = albumErr.New("album not found")
)

type AlbumService struct {
	logger       *slog.Logger
	db           *database.Database
	imageService *ImageService
	dataDir      types.DataDir
}

func NewAlbumService(
	logger *slog.Logger,
	db *database.Database,
	imageService *ImageService,
	dataDir types.DataDir,
) *AlbumService {
	return &AlbumService{
		logger:       logger,
		db:           db,
		imageService: imageService,
		dataDir:      dataDir,
	}
}

type GetAlbumsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *AlbumService) GetAlbums(
	ctx context.Context,
	params GetAlbumsParams,
) ([]database.Album, types.Page, error) {
	albums, page, err := s.db.GetAlbums(ctx, database.GetAlbumsParams{
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, types.Page{}, &InvalidFilterError{
				Service: "album service",
				Message: err.Error(),
			}
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, types.Page{}, &InvalidSortError{
				Service: "album service",
				Message: err.Error(),
			}
		}

		return nil, types.Page{}, albumErr.Wrap("get albums", err)
	}

	return albums, page, nil
}

type GetAlbumByIdParams struct {
	AlbumId string
}

func (s *AlbumService) GetAlbumById(
	ctx context.Context,
	params GetAlbumByIdParams,
) (database.Album, error) {
	album, err := s.db.GetAlbumById(ctx, params.AlbumId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return database.Album{}, ErrAlbumServiceAlbumNotFound
		}

		return database.Album{}, albumErr.Wrap("get album by id", err)
	}

	return album, nil
}

type GetAlbumTracksParams struct {
	AlbumId string
}

func (s *AlbumService) GetAlbumTracks(
	ctx context.Context,
	params GetAlbumTracksParams,
) ([]database.Track, error) {
	album, err := s.GetAlbumById(ctx, GetAlbumByIdParams{
		AlbumId: params.AlbumId,
	})
	if err != nil {
		return nil, err
	}

	tracks, err := s.db.GetTracksByAlbum(ctx, album.Id)
	if err != nil {
		return nil, albumErr.Wrap("get album tracks", err)
	}

	for i, track := range tracks {
		if track.Number.Valid {
			tracks[i].Order = utils.Pointer(int(track.Number.Int64))
		}
	}

	return tracks, nil
}

type GetAlbumImageParams struct {
	AlbumId     string
	Size        int
	ImageFormat types.ImageFormat
}

func (s *AlbumService) GetAlbumImage(
	ctx context.Context,
	params GetAlbumImageParams,
) (string, error) {
	album, err := s.GetAlbumById(ctx, GetAlbumByIdParams{
		AlbumId: params.AlbumId,
	})
	if err != nil {
		return "", albumErr.Wrap("get album image: get album", err)
	}

	cacheDir := s.dataDir.CacheImages()

	input := ""
	if album.CoverArt.Valid {
		input = album.CoverArt.String
	}

	if err := utils.CreateDirectories([]string{
		cacheDir.String(),
		cacheDir.Albums(),
		cacheDir.Album(album.Id),
	}); err != nil {
		return "", albumErr.Wrap("get album image: mkdir", err)
	}

	p, err := s.imageService.ProcessImage(ProcessImageParams{
		Input:       input,
		Default:     "default_album.png",
		OutputDir:   cacheDir.Album(album.Id),
		Size:        params.Size,
		ImageFormat: params.ImageFormat,
	})
	if err != nil {
		return "", albumErr.Wrap("get album image: process image", err)
	}

	return p, nil
}
