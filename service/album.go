package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/nanoteck137/tunebook/types"
)

var (
	ErrAlbumServiceAlbumNotFound = errors.New("album service: album not found")
)

type AlbumService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewAlbumService(
	logger *slog.Logger,
	db *database.Database,
) *AlbumService {
	return &AlbumService{
		logger: logger,
		db:     db,
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

		return nil, types.Page{}, err
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

		return database.Album{}, err
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
		return nil, err
	}

	for i, track := range tracks {
		if track.Number.Valid {
			tracks[i].Order = utils.Pointer(int(track.Number.Int64))
			// TODO(patrik): Should this not be using i?
			// tracks[i].Order = utils.IntPtr(i)
		}
	}

	return tracks, nil
}
