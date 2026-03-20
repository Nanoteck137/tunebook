package service

import (
	"log/slog"

	"github.com/nanoteck137/dwebble/database"
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
