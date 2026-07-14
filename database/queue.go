package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	createQueueId = createIdGenerator(10)

	queuesTbl = goqu.T("queues")
)

type Queue struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`

	CurrentIndex int `db:"current_index"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func QueueQuery() *goqu.SelectDataset {
	query := dialect.From(queuesTbl).
		Select(
			queuesTbl.Col("id"),
			queuesTbl.Col("user_id"),

			queuesTbl.Col("current_index"),

			queuesTbl.Col("updated"),
			queuesTbl.Col("created"),
		)

	return query
}

type CreateQueueParams struct {
	Id     string
	UserId string

	CurrentIndex int

	Created int64
	Updated int64
}

func (db DB) CreateQueue(
	ctx context.Context,
	params CreateQueueParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createQueueId()
	}

	query := dialect.Insert(queuesTbl).Rows(goqu.Record{
		"id":     params.Id,
		"user_id": params.UserId,

		"current_index": params.CurrentIndex,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type QueueChanges struct {
	CurrentIndex Change[int]

	Created Change[int64]
}

func (db DB) UpdateQueue(
	ctx context.Context,
	queueId string,
	changes QueueChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "current_index", changes.CurrentIndex)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update(queuesTbl).
		Set(record).
		Where(queuesTbl.Col("id").Eq(queueId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteQueue(ctx context.Context, queueId string) error {
	query := dialect.Delete(queuesTbl).
		Where(queuesTbl.Col("id").Eq(queueId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetQueueById(
	ctx context.Context,
	queueId string,
) (Queue, error) {
	query := QueueQuery().
		Where(queuesTbl.Col("id").Eq(queueId))

	return Single[Queue](db, ctx, query)
}
