package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibraryCleanupJob)(nil)

const LibraryCleanup = "library-cleanup"

type LibraryCleanupJob struct {
	libraryService *service.LibraryService
}

func NewLibraryCleanupJob(
	libraryService *service.LibraryService,
) *LibraryCleanupJob {
	return &LibraryCleanupJob{
		libraryService: libraryService,
	}
}

func (j *LibraryCleanupJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        LibraryCleanup,
		DisplayName: "Library Cleanup",
	}
}

func (j *LibraryCleanupJob) Run(ctx context.Context, data string) error {
	return j.libraryService.Cleanup(ctx)
}
