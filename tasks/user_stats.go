package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*UserStatsRecalculateTask)(nil)

const UserStatsRecalculate = "user-stats-recalculate"

type UserStatsRecalculateTask struct {
	userService *service.UserService
	jobService  *service.JobService
}

func NewUserStatsRecalculateTask(
	userService *service.UserService,
	jobService *service.JobService,
) *UserStatsRecalculateTask {
	return &UserStatsRecalculateTask{
		userService: userService,
		jobService:  jobService,
	}
}

func (j *UserStatsRecalculateTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        UserStatsRecalculate,
		DisplayName: "User Stats Recalculate",
		Schedule:    "",
	}
}

func (j *UserStatsRecalculateTask) Run(ctx context.Context) error {
	users, err := j.userService.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		err := j.jobService.PushJob(
			ctx,
			jobs.UserStatsUpdate,
			service.UpdateUserStatsParams{
				UserId: user.Id,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
