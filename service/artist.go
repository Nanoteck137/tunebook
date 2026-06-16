package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var (
	ErrArtistServiceArtistNotFound = errors.New("artist service: artist not found")
)

type ArtistService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewArtistService(
	logger *slog.Logger,
	db *database.Database,
) *ArtistService {
	return &ArtistService{
		logger: logger,
		db:     db,
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

		return nil, types.Page{}, err
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

		return database.Artist{}, err
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
		return nil, err
	}

	return albums, nil
}
