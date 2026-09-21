package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*AuthCleanupTask)(nil)

const AuthCleanup = "auth-cleanup"

type AuthCleanupTask struct {
	jobService *service.JobService
}

func NewAuthCleanupTask(jobService *service.JobService) *AuthCleanupTask {
	return &AuthCleanupTask{
		jobService: jobService,
	}
}

func (j *AuthCleanupTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        AuthCleanup,
		DisplayName: "Auth Cleanup",
		Schedule:    "@every 30m",
	}
}

func (j *AuthCleanupTask) Run(ctx context.Context) error {
	return j.jobService.PushJob(
		ctx,
		jobs.AuthCleanup,
		nil,
	)
}
