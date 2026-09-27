package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*UserReviewRebuildTask)(nil)

const UserReviewRebuild = "user-review-rebuild"

type UserReviewRebuildTask struct {
	userService *service.UserService
	jobService  *service.JobService
}

func NewUserReviewRebuildTask(
	userService *service.UserService,
	jobService *service.JobService,
) *UserReviewRebuildTask {
	return &UserReviewRebuildTask{
		userService: userService,
		jobService:  jobService,
	}
}

func (t *UserReviewRebuildTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        UserReviewRebuild,
		DisplayName: "User Review Rebuild",
		// Nightly, and deliberately unconditional. Review tables are not
		// invalidated by normal library syncs (names and artwork join live),
		// but deleting tracks cascades away the ranked rows while the
		// summary counters survive, so this is what repairs that.
		Schedule: "@daily",
	}
}

func (t *UserReviewRebuildTask) Run(ctx context.Context) error {
	users, err := t.userService.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		err := t.jobService.PushJob(
			ctx,
			jobs.UserReviewRebuild,
			jobs.UserReviewRebuildParams{
				UserId: user.Id,
			},
			service.WithUniqueKey(jobs.UserReviewRebuildKey(user.Id)),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
