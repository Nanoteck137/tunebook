package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*LibrarySyncTask)(nil)

type LibrarySyncTask struct {
	libraryService *service.LibraryService
}

func NewLibrarySyncTask(libraryService *service.LibraryService) *LibrarySyncTask {
	return &LibrarySyncTask{
		libraryService: libraryService,
	}
}

func (j *LibrarySyncTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        LibrarySync,
		DisplayName: "Library Sync",
		Schedule:    "",
	}
}

func (j *LibrarySyncTask) Run(ctx context.Context) error {
	return j.libraryService.Sync()
}
