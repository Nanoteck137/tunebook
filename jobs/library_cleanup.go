package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibraryCleanupJob)(nil)

const LibraryCleanup = "library-cleanup"

type LibraryCleanupJob struct {
	libraryService *service.LibraryService
	jobService     *service.JobService
}

func NewLibraryCleanupJob(
	libraryService *service.LibraryService,
	jobService *service.JobService,
) *LibraryCleanupJob {
	return &LibraryCleanupJob{
		libraryService: libraryService,
		jobService:     jobService,
	}
}

func (j *LibraryCleanupJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        LibraryCleanup,
		DisplayName: "Library Cleanup",
	}
}

func (j *LibraryCleanupJob) Run(ctx context.Context, data string) error {
	err := j.libraryService.Cleanup(ctx)
	if err != nil {
		return err
	}

	// Cleanup deletes missing artists, albums and tracks, which would
	// otherwise linger in the search index.
	return DispatchSearchIndex(ctx, j.jobService)
}
