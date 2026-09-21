package database

import (
	"context"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
)

const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
)

var (
	createJobId = createIdGenerator(16)

	jobsTbl = goqu.T("jobs")
)

type Job struct {
	Id            string `db:"id"`
	Name          string `db:"name"`
	Data          string `db:"data"`
	UniqueKey     string `db:"unique_key"`
	Status        string `db:"status"`
	Error         string `db:"error"`
	Attempts      int    `db:"attempts"`
	MaxAttempts   int    `db:"max_attempts"`
	NextAttemptAt int64  `db:"next_attempt_at"`
	Created       int64  `db:"created"`
	Updated       int64  `db:"updated"`
}

func JobQuery() *goqu.SelectDataset {
	query := dialect.From(jobsTbl).
		Select(
			jobsTbl.Col("id"),
			jobsTbl.Col("name"),
			jobsTbl.Col("data"),
			goqu.COALESCE(jobsTbl.Col("unique_key"), "").As("unique_key"),
			jobsTbl.Col("status"),
			jobsTbl.Col("error"),
			jobsTbl.Col("attempts"),
			jobsTbl.Col("max_attempts"),
			jobsTbl.Col("next_attempt_at"),
			jobsTbl.Col("created"),
			jobsTbl.Col("updated"),
		)

	return query
}

type CreateJobParams struct {
	Id          string
	Name        string
	Data        string
	UniqueKey   string
	MaxAttempts int
}

func (db DB) CreateJob(
	ctx context.Context,
	params CreateJobParams,
) (string, error) {
	t := time.Now().UnixMilli()

	if params.Id == "" {
		params.Id = createJobId()
	}

	if params.MaxAttempts <= 0 {
		params.MaxAttempts = 1
	}

	var uniqueKey any
	if params.UniqueKey != "" {
		uniqueKey = params.UniqueKey
	}

	query := dialect.Insert(jobsTbl).Rows(goqu.Record{
		"id":              params.Id,
		"name":            params.Name,
		"data":            params.Data,
		"unique_key":      uniqueKey,
		"status":          JobStatusPending,
		"error":           "",
		"attempts":        0,
		"max_attempts":    params.MaxAttempts,
		"next_attempt_at": t,
		"created":         t,
		"updated":         t,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) GetJobById(ctx context.Context, jobId string) (Job, error) {
	query := JobQuery().
		Where(jobsTbl.Col("id").Eq(jobId))

	return Single[Job](db, ctx, query)
}

// HasActiveJobWithUniqueKey reports whether a job with the given unique key is
// currently pending or running. It is used to dedupe jobs that should only
// ever be queued once at a time (e.g. library sync, one playlist image
// generation per playlist).
func (db DB) HasActiveJobWithUniqueKey(
	ctx context.Context,
	uniqueKey string,
) (bool, error) {
	if uniqueKey == "" {
		return false, nil
	}

	query := JobQuery().
		Where(
			jobsTbl.Col("unique_key").Eq(uniqueKey),
			jobsTbl.Col("status").In(
				JobStatusPending,
				JobStatusRunning,
			),
		).
		Limit(1)

	_, err := Single[Job](db, ctx, query)
	if errors.Is(err, ErrItemNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db DB) GetJobs(ctx context.Context, limit int) ([]Job, error) {
	if limit <= 0 {
		limit = 50
	}

	query := JobQuery().
		Order(jobsTbl.Col("created").Desc()).
		Limit(uint(limit))

	return Multiple[Job](db, ctx, query)
}

func (db DB) GetPendingJobs(ctx context.Context, limit int) ([]Job, error) {
	now := time.Now().UnixMilli()

	query := JobQuery().
		Where(
			jobsTbl.Col("status").Eq(JobStatusPending),
			jobsTbl.Col("next_attempt_at").Lte(now),
		).
		Order(jobsTbl.Col("created").Asc()).
		Limit(uint(limit))

	return Multiple[Job](db, ctx, query)
}

func (db DB) ClaimJob(ctx context.Context, jobId string) error {
	query := dialect.Update(jobsTbl).
		Set(goqu.Record{
			"status":   JobStatusRunning,
			"attempts": goqu.L("attempts + 1"),
			"updated":  time.Now().UnixMilli(),
		}).
		Where(jobsTbl.Col("id").Eq(jobId))

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) CompleteJob(ctx context.Context, jobId string) error {
	query := dialect.Update(jobsTbl).
		Set(goqu.Record{
			"status":  JobStatusCompleted,
			"error":   "",
			"updated": time.Now().UnixMilli(),
		}).
		Where(jobsTbl.Col("id").Eq(jobId))

	_, err := db.Exec(ctx, query)
	return err
}

type FailJobParams struct {
	Requeue       bool
	Error         string
	NextAttemptAt int64
}

func (db DB) FailJob(
	ctx context.Context,
	jobId string,
	params FailJobParams,
) error {
	record := goqu.Record{
		"error":   params.Error,
		"updated": time.Now().UnixMilli(),
	}

	if params.Requeue {
		nextAttemptAt := params.NextAttemptAt
		if nextAttemptAt <= 0 {
			nextAttemptAt = time.Now().UnixMilli()
		}

		record["status"] = JobStatusPending
		record["next_attempt_at"] = nextAttemptAt
	} else {
		record["status"] = JobStatusFailed
		record["next_attempt_at"] = 0
	}

	query := dialect.Update(jobsTbl).
		Set(record).
		Where(jobsTbl.Col("id").Eq(jobId))

	_, err := db.Exec(ctx, query)
	return err
}

// GetStuckJobs returns jobs left in the "running" state. A job found in this
// state on startup is assumed to have been interrupted by a crash or restart.
func (db DB) GetStuckJobs(ctx context.Context) ([]Job, error) {
	query := JobQuery().
		Where(jobsTbl.Col("status").Eq(JobStatusRunning))

	return Multiple[Job](db, ctx, query)
}

// RequeueJob resets a job from "running" back to "pending" so it is retried,
// scheduling it for its next attempt at nextAttemptAt (unix millis).
func (db DB) RequeueJob(
	ctx context.Context,
	jobId string,
	nextAttemptAt int64,
) error {
	query := dialect.Update(jobsTbl).
		Set(goqu.Record{
			"status":          JobStatusPending,
			"next_attempt_at": nextAttemptAt,
			"updated":         time.Now().UnixMilli(),
		}).
		Where(jobsTbl.Col("id").Eq(jobId))

	_, err := db.Exec(ctx, query)
	return err
}

// DeleteOldJobs removes finished jobs (completed or failed) last updated before
// olderThan (unix millis).
func (db DB) DeleteOldJobs(ctx context.Context, olderThan int64) (int64, error) {
	query := dialect.Delete(jobsTbl).
		Where(
			jobsTbl.Col("status").In(
				JobStatusCompleted,
				JobStatusFailed,
			),
			jobsTbl.Col("updated").Lt(olderThan),
		)

	res, err := db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
