package jobs

import (
	"context"
	"encoding/json"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*UserTrackStatsRebuildJob)(nil)

const UserTrackStatsRebuild = "user-track-stats-rebuild"

type UserTrackStatsRebuildJob struct {
	userService *service.UserService
}

func NewUserTrackStatsRebuildJob(
	userService *service.UserService,
) *UserTrackStatsRebuildJob {
	return &UserTrackStatsRebuildJob{
		userService: userService,
	}
}

func (j *UserTrackStatsRebuildJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        UserTrackStatsRebuild,
		DisplayName: "User Track Stats Rebuild",
	}
}

func (j *UserTrackStatsRebuildJob) Run(ctx context.Context, data string) error {
	var params service.RebuildUserTrackStatsParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.userService.RebuildUserTrackStats(ctx, params.UserId)
}