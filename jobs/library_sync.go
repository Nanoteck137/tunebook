package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibrarySyncJob)(nil)

const LibrarySync = "library-sync"

type LibrarySyncJob struct {
	libraryService *service.LibraryService
	jobService     *service.JobService
}

func NewLibrarySyncJob(
	libraryService *service.LibraryService,
	jobService *service.JobService,
) *LibrarySyncJob {
	return &LibrarySyncJob{
		libraryService: libraryService,
		jobService:     jobService,
	}
}

func (j *LibrarySyncJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:          LibrarySync,
		DisplayName:   "Library Sync",
		FailOnRestart: true,
		NoTimeout:     true,
	}
}

func (j *LibrarySyncJob) Run(ctx context.Context, data string) error {
	err := j.libraryService.Sync(ctx)
	if err != nil {
		return err
	}

	// The sync changed artists, albums and tracks, so refresh the search
	// index. Only on success: a failed sync may have left the library in a
	// partially updated state, and reindexing that is not obviously better
	// than leaving the previous index to be revisited on the next run.
	return DispatchSearchIndex(ctx, j.jobService)
}
