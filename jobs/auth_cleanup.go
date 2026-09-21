package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*AuthCleanupJob)(nil)

const AuthCleanup = "auth-cleanup"

type AuthCleanupJob struct {
	authService *service.AuthService
}

func NewAuthCleanupJob(authService *service.AuthService) *AuthCleanupJob {
	return &AuthCleanupJob{
		authService: authService,
	}
}

func (j *AuthCleanupJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        AuthCleanup,
		DisplayName: "Auth Cleanup",
	}
}

func (j *AuthCleanupJob) Run(ctx context.Context, data string) error {
	j.authService.Cleanup()
	return nil
}
