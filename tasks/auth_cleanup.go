package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*AuthCleanupTask)(nil)

const AuthCleanup = "auth-cleanup"

type AuthCleanupTask struct {
	authService *service.AuthService
}

func NewAuthCleanupTask(authService *service.AuthService) *AuthCleanupTask {
	return &AuthCleanupTask{
		authService: authService,
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
	j.authService.Cleanup()
	return nil
}
