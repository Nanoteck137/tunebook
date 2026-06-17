package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type Queue struct {
	Id           string `db:"id"`
	CurrentIndex int    `db:"current_index"`
	Created      int64  `db:"created"`
	Updated      int64  `db:"updated"`
}

func (db DB) CreateQueue(ctx context.Context, userId string, currentIndex int) error {
	t := time.Now().UnixMilli()

	query := dialect.Insert("queues").
		Rows(goqu.Record{
			"id":            userId,
			"current_index": currentIndex,
			"created":       t,
			"updated":       t,
		})

	_, err := db.Exec(ctx, query)
	return err
}

type QueueChanges struct {
	CurrentIndex Change[int]
}

func (db DB) UpdateQueue(ctx context.Context, userId string, changes QueueChanges) error {
	record := goqu.Record{}

	addToRecord(record, "current_index", changes.CurrentIndex)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("queues").
		Set(record).
		Where(goqu.I("queues.id").Eq(userId))

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) GetQueueByUserId(ctx context.Context, userId string) (Queue, error) {
	query := dialect.From("queues").
		Select(
			"queues.id",
			"queues.current_index",
			"queues.created",
			"queues.updated",
		).
		Where(goqu.I("queues.id").Eq(userId))

	return Single[Queue](db, ctx, query)
}

func (db DB) DeleteQueue(ctx context.Context, userId string) error {
	query := goqu.Delete("queues").
		Where(goqu.I("queues.id").Eq(userId))

	_, err := db.Exec(ctx, query)
	return err
}
