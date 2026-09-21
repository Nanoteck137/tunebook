package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibrarySyncJob)(nil)

const LibrarySync = "library-sync"

type LibrarySyncJob struct {
	libraryService *service.LibraryService
}

func NewLibrarySyncJob(libraryService *service.LibraryService) *LibrarySyncJob {
	return &LibrarySyncJob{
		libraryService: libraryService,
	}
}

func (j *LibrarySyncJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:          LibrarySync,
		DisplayName:   "Library Sync",
		FailOnRestart: true,
	}
}

func (j *LibrarySyncJob) Run(ctx context.Context, data string) error {
	return j.libraryService.Sync(ctx)
}
