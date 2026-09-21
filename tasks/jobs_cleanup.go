package tasks

import (
	"context"
	"time"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*JobsCleanupTask)(nil)

const JobsCleanup = "jobs-cleanup"

var jobsRetention = 7 * 24 * time.Hour

type JobsCleanupTask struct {
	jobService *service.JobService
}

func NewJobsCleanupTask(jobService *service.JobService) *JobsCleanupTask {
	return &JobsCleanupTask{
		jobService: jobService,
	}
}

func (j *JobsCleanupTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        JobsCleanup,
		DisplayName: "Jobs Cleanup",
		Schedule:    "@daily",
	}
}

func (j *JobsCleanupTask) Run(ctx context.Context) error {
	_, err := j.jobService.CleanupOldJobs(ctx, jobsRetention)
	return err
}
