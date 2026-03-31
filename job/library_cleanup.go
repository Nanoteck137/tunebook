package job

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*LibraryCleanupJob)(nil)

type LibraryCleanupJob struct {
	libraryService *service.LibraryService
}

func NewLibraryCleanupJob(libraryService *service.LibraryService) *LibraryCleanupJob {
	return &LibraryCleanupJob{
		libraryService: libraryService,
	}
}

func (j *LibraryCleanupJob) Name() string {
	return LibraryCleanup
}

func (j *LibraryCleanupJob) Schedule() string {
	return ""
}

func (j *LibraryCleanupJob) Run(ctx context.Context) error {
	return j.libraryService.Cleanup(ctx)
}
