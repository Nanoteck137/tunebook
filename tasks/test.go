package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*TestTask)(nil)

const TestJob = "test-job"

type TestTask struct {
	jobService *service.JobService
}

func NewTestTask(jobService *service.JobService) *TestTask {
	return &TestTask{
		jobService: jobService,
	}
}

func (t *TestTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        TestJob,
		DisplayName: "Test Job",
		Schedule:    "",
	}
}

func (t *TestTask) Run(ctx context.Context) error {
	return t.jobService.PushJob(
		ctx,
		jobs.Test,
		service.TestJobParams{
			DurationMs: int64(30 * 1000),
		},
	)
}
