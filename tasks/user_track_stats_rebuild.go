package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*UserTrackStatsRebuildTask)(nil)

type UserTrackStatsRebuildTask struct {
	userService *service.UserService
	jobService  *service.JobService
}

func NewUserTrackStatsRebuildTask(
	userService *service.UserService,
	jobService *service.JobService,
) *UserTrackStatsRebuildTask {
	return &UserTrackStatsRebuildTask{
		userService: userService,
		jobService:  jobService,
	}
}

func (j *UserTrackStatsRebuildTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        jobs.UserTrackStatsRebuild,
		DisplayName: "User Track Stats Rebuild",
		Schedule:    "",
	}
}

func (j *UserTrackStatsRebuildTask) Run(ctx context.Context) error {
	users, err := j.userService.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		err := j.jobService.PushJob(
			ctx,
			jobs.UserTrackStatsRebuild,
			service.RebuildUserTrackStatsParams{
				UserId: user.Id,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
