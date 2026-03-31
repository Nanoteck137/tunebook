package job

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*AuthCleanupJob)(nil)

type AuthCleanupJob struct {
	authService *service.AuthService
}

func NewAuthCleanupJob(authService *service.AuthService) *AuthCleanupJob {
	return &AuthCleanupJob{
		authService: authService,
	}
}

func (j *AuthCleanupJob) Name() string {
	return AuthCleanup
}

func (j *AuthCleanupJob) Schedule() string {
	return "@every 30m"
}

func (j *AuthCleanupJob) Run(ctx context.Context) error {
	j.authService.Cleanup()
	return nil
}
