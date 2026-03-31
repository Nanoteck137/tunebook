package job

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibrarySyncJob)(nil)

type LibrarySyncJob struct {
	libraryService *service.LibraryService
}

func NewLibrarySyncJob(libraryService *service.LibraryService) *LibrarySyncJob {
	return &LibrarySyncJob{
		libraryService: libraryService,
	}
}

func (j *LibrarySyncJob) Name() string {
	return LibrarySync
}

func (j *LibrarySyncJob) Schedule() string {
	return ""
}

func (j *LibrarySyncJob) Run(ctx context.Context) error {
	return j.libraryService.Sync()
}
