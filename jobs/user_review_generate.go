package jobs

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*UserReviewGenerateJob)(nil)
var _ service.Job = (*UserReviewGenerateTotalJob)(nil)

const UserReviewGenerate = "user-review-generate"
const UserReviewGenerateTotal = "user-review-generate-total"

// UserReviewGenerateYearKey returns the unique key for regenerating one
// user's review for one year, so at most one generation job per user and
// year is queued at a time.
func UserReviewGenerateYearKey(userId string, year int) string {
	return "user-review-generate-" + userId + "-" + strconv.Itoa(year)
}

// UserReviewGenerateTotalKey returns the unique key for regenerating one
// user's all-time review.
func UserReviewGenerateTotalKey(userId string) string {
	return "user-review-generate-total-" + userId
}

type UserReviewGenerateJob struct {
	userService *service.UserService
}

func NewUserReviewGenerateJob(
	userService *service.UserService,
) *UserReviewGenerateJob {
	return &UserReviewGenerateJob{
		userService: userService,
	}
}

func (j *UserReviewGenerateJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        UserReviewGenerate,
		DisplayName: "User Review Generate",
		NoTimeout:   true,
	}
}

func (j *UserReviewGenerateJob) Run(ctx context.Context, data string) error {
	var params service.GenerateUserReviewParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.userService.GenerateUserReview(ctx, params)
}

type UserReviewGenerateTotalJob struct {
	userService *service.UserService
}

func NewUserReviewGenerateTotalJob(
	userService *service.UserService,
) *UserReviewGenerateTotalJob {
	return &UserReviewGenerateTotalJob{
		userService: userService,
	}
}

func (j *UserReviewGenerateTotalJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        UserReviewGenerateTotal,
		DisplayName: "User Review Generate (All Time)",
		NoTimeout:   true,
	}
}

func (j *UserReviewGenerateTotalJob) Run(ctx context.Context, data string) error {
	var params service.GenerateUserTotalReviewParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.userService.GenerateUserTotalReview(ctx, params.UserId)
}
