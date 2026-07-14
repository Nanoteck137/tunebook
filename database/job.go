package database

import (
	"context"
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
	Id          string `db:"id"`
	Name        string `db:"name"`
	Data        string `db:"data"`
	Status      string `db:"status"`
	Error       string `db:"error"`
	Attempts    int    `db:"attempts"`
	MaxAttempts int    `db:"max_attempts"`
	Created     int64  `db:"created"`
	Updated     int64  `db:"updated"`
}

func JobQuery() *goqu.SelectDataset {
	query := dialect.From(jobsTbl).
		Select(
			jobsTbl.Col("id"),
			jobsTbl.Col("name"),
			jobsTbl.Col("data"),
			jobsTbl.Col("status"),
			jobsTbl.Col("error"),
			jobsTbl.Col("attempts"),
			jobsTbl.Col("max_attempts"),
			jobsTbl.Col("created"),
			jobsTbl.Col("updated"),
		)

	return query
}

type CreateJobParams struct {
	Id          string
	Name        string
	Data        string
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

	query := dialect.Insert(jobsTbl).Rows(goqu.Record{
		"id":           params.Id,
		"name":         params.Name,
		"data":         params.Data,
		"status":       JobStatusPending,
		"error":        "",
		"attempts":     0,
		"max_attempts": params.MaxAttempts,
		"created":      t,
		"updated":      t,
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

func (db DB) GetPendingJobs(ctx context.Context, limit int) ([]Job, error) {
	query := JobQuery().
		Where(jobsTbl.Col("status").Eq(JobStatusPending)).
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
	Requeue bool
	Error   string
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
		record["status"] = JobStatusPending
	} else {
		record["status"] = JobStatusFailed
	}

	query := dialect.Update(jobsTbl).
		Set(record).
		Where(jobsTbl.Col("id").Eq(jobId))

	_, err := db.Exec(ctx, query)
	return err
}
