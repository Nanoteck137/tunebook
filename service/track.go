package service

import (
	"log/slog"

	"github.com/nanoteck137/dwebble/database"
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
