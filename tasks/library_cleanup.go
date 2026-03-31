package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*LibraryCleanupTask)(nil)

type LibraryCleanupTask struct {
	libraryService *service.LibraryService
}

func NewLibraryCleanupTask(libraryService *service.LibraryService) *LibraryCleanupTask {
	return &LibraryCleanupTask{
		libraryService: libraryService,
	}
}

func (j *LibraryCleanupTask) Name() string {
	return LibraryCleanup
}

func (j *LibraryCleanupTask) Schedule() string {
	return ""
}

func (j *LibraryCleanupTask) Run(ctx context.Context) error {
	return j.libraryService.Cleanup(ctx)
}
