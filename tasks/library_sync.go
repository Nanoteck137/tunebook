package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*LibrarySyncTask)(nil)

const LibrarySync = "library-sync"

type LibrarySyncTask struct {
	jobService *service.JobService
}

func NewLibrarySyncTask(jobService *service.JobService) *LibrarySyncTask {
	return &LibrarySyncTask{
		jobService: jobService,
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
	return j.jobService.PushJob(
		ctx,
		jobs.LibrarySync,
		nil,
		service.WithUniqueKey(jobs.LibrarySync),
	)
}
