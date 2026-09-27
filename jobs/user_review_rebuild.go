package jobs

import (
	"context"
	"encoding/json"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*UserReviewRebuildJob)(nil)

const UserReviewRebuild = "user-review-rebuild"

// UserReviewRebuildKey returns the unique key for rebuilding all of one
// user's reviews, so at most one rebuild per user is queued at a time.
func UserReviewRebuildKey(userId string) string {
	return "user-review-rebuild-" + userId
}

type UserReviewRebuildParams struct {
	UserId string
}

type UserReviewRebuildJob struct {
	userService *service.UserService
}

func NewUserReviewRebuildJob(
	userService *service.UserService,
) *UserReviewRebuildJob {
	return &UserReviewRebuildJob{
		userService: userService,
	}
}

func (j *UserReviewRebuildJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        UserReviewRebuild,
		DisplayName: "User Review Rebuild",
		NoTimeout:   true,
	}
}

func (j *UserReviewRebuildJob) Run(ctx context.Context, data string) error {
	var params UserReviewRebuildParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.userService.RebuildUserReviews(ctx, params.UserId)
}
