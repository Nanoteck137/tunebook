package service

import (
	"log/slog"
	"sync/atomic"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
)

type LibraryService struct {
	db     *database.Database
	config *config.Config

	syncRunning atomic.Bool
}

func NewLibraryService(db *database.Database, config *config.Config) *LibraryService {
	return &LibraryService{
		db:          db,
		config:      config,
		syncRunning: atomic.Bool{},
	}
}

func (s *LibraryService) runSync() {
}

func (s *LibraryService) Sync() {
	if s.syncRunning.Load() {
		slog.Error("library syncing already running")
		return
	}
}
