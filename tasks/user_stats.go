package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*UserStatsRecalculateTask)(nil)

type UserStatsRecalculateTask struct {
	userService *service.UserService
}

func NewUserStatsRecalculateTask(userService *service.UserService) *UserStatsRecalculateTask {
	return &UserStatsRecalculateTask{
		userService: userService,
	}
}

func (j *UserStatsRecalculateTask) Name() string {
	return UserStatsRecalculate
}

func (j *UserStatsRecalculateTask) Schedule() string {
	return ""
}

func (j *UserStatsRecalculateTask) Run(ctx context.Context) error {
	return j.userService.RecalculateAllUserStats(ctx)
}
