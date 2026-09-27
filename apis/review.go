package apis

import (
	"context"
	"log/slog"
	"time"

	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

// Review refresh dispatch lives here rather than in the services because the
// service package cannot import jobs (jobs imports service). The
// generate-playlist-image handler dispatches from the API layer for the same
// reason.
//
// Dispatch is synchronous. In the steady state the reviews are fresh, so
// this costs two primary-key reads and creates no job rows; the reads are
// cheap enough not to warrant a goroutine per playback report.

// dispatchReviewYear queues regeneration of one user's review for one year,
// unless that review was generated recently.
//
// A unique key alone is not enough: it only dedupes while a job is pending
// or running, and the job worker drains every 5s, so a user who keeps
// listening would otherwise trigger a full rebuild each time the previous
// one finished. The staleness check is what actually bounds the cost.
func dispatchReviewYear(
	ctx context.Context,
	app core.App,
	userId string,
	year int,
) {
	needs, err := app.UserService().ReviewNeedsRefresh(
		ctx,
		userId,
		year,
		service.ReviewDebounceWindowMs,
	)
	if err != nil {
		slog.Error("review refresh: check staleness",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
			slog.Int("year", year),
		)
		return
	}

	if !needs {
		return
	}

	err = app.JobService().PushJob(
		ctx,
		jobs.UserReviewGenerate,
		service.GenerateUserReviewParams{
			UserId: userId,
			Year:   year,
		},
		service.WithUniqueKey(jobs.UserReviewGenerateYearKey(userId, year)),
	)
	if err != nil {
		slog.Error("review refresh: push job",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
			slog.Int("year", year),
		)
	}
}

// dispatchReviewTotal queues regeneration of one user's all-time review,
// unless it was generated recently.
func dispatchReviewTotal(ctx context.Context, app core.App, userId string) {
	needs, err := app.UserService().TotalReviewNeedsRefresh(
		ctx,
		userId,
		service.ReviewDebounceWindowMs,
	)
	if err != nil {
		slog.Error("review refresh: check total staleness",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
		)
		return
	}

	if !needs {
		return
	}

	err = app.JobService().PushJob(
		ctx,
		jobs.UserReviewGenerateTotal,
		service.GenerateUserTotalReviewParams{
			UserId: userId,
		},
		service.WithUniqueKey(jobs.UserReviewGenerateTotalKey(userId)),
	)
	if err != nil {
		slog.Error("review refresh: push total job",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
		)
	}
}

// dispatchReviewCurrentYear refreshes the reviews that change as a result of
// playback. Plays always stamp the current year, so only that year and the
// all-time review can have moved.
func dispatchReviewCurrentYear(ctx context.Context, app core.App, userId string) {
	year := time.Now().Year()

	dispatchReviewYear(ctx, app, userId, year)
	dispatchReviewTotal(ctx, app, userId)
}

// dispatchReviewAll queues a full rebuild of one user's reviews.
//
// Favorites need this rather than a per-year dispatch:
// GetUserYearFavoritePlays joins the *current* favorites against each year's
// play history, so favoriting or unfavoriting a track changes favorite_plays
// in every year the user played it, not just the current one.
func dispatchReviewAll(ctx context.Context, app core.App, userId string) {
	needs, err := app.UserService().TotalReviewNeedsRefresh(
		ctx,
		userId,
		service.ReviewDebounceWindowMs,
	)
	if err != nil {
		slog.Error("review refresh: check staleness for rebuild",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
		)
		return
	}

	if !needs {
		return
	}

	err = app.JobService().PushJob(
		ctx,
		jobs.UserReviewRebuild,
		jobs.UserReviewRebuildParams{
			UserId: userId,
		},
		service.WithUniqueKey(jobs.UserReviewRebuildKey(userId)),
	)
	if err != nil {
		slog.Error("review refresh: push rebuild job",
			slog.String("error", err.Error()),
			slog.String("userId", userId),
		)
	}
}
