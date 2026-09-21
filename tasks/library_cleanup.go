package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*LibraryCleanupTask)(nil)

const LibraryCleanup = "library-cleanup"

type LibraryCleanupTask struct {
	jobService *service.JobService
}

func NewLibraryCleanupTask(jobService *service.JobService) *LibraryCleanupTask {
	return &LibraryCleanupTask{
		jobService: jobService,
	}
}

func (j *LibraryCleanupTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        LibraryCleanup,
		DisplayName: "Library Cleanup",
		Schedule:    "",
	}
}

func (j *LibraryCleanupTask) Run(ctx context.Context) error {
	return j.jobService.PushJob(
		ctx,
		jobs.LibraryCleanup,
		nil,
	)
}
