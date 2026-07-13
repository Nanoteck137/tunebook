package jobs

import (
	"context"
	"encoding/json"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*UserStatsUpdateJob)(nil)

const UserStatsUpdate = "user-stats-update"

type UserStatsUpdateJob struct {
	userService *service.UserService
}

func NewUserStatsUpdateJob(
	userService *service.UserService,
) *UserStatsUpdateJob {
	return &UserStatsUpdateJob{
		userService: userService,
	}
}

func (j *UserStatsUpdateJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        UserStatsUpdate,
		DisplayName: "User Stats Update",
	}
}

func (j *UserStatsUpdateJob) Run(ctx context.Context, data string) error {
	var params service.UpdateUserStatsParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.userService.RecalculateUserStats(ctx, params.UserId)
}
