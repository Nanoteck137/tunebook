package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var createQueueId = createIdGenerator(10)

type Queue struct {
	Id string `db:"id"`

	CurrentIndex int `db:"current_index"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func QueueQuery() *goqu.SelectDataset {
	query := dialect.From("queues").
		Select(
			"queues.id",

			"queues.current_index",

			"queues.updated",
			"queues.created",
		)

	return query
}

func (db DB) GetQueueById(
	ctx context.Context,
	queueId string,
) (Queue, error) {
	query := QueueQuery().
		Where(goqu.I("queues.id").Eq(queueId))

	return Single[Queue](db, ctx, query)
}

type CreateQueueParams struct {
	Id string

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

	query := dialect.Insert("queues").Rows(goqu.Record{
		"id": params.Id,

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
	Name Change[string]

	CurrentIndex Change[int]

	Created Change[int64]
}

func (db DB) UpdateQueue(
	ctx context.Context,
	queueId string,
	changes QueueChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "current_index", changes.CurrentIndex)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("queues").
		Set(record).
		Where(goqu.I("queues.id").Eq(queueId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteQueue(ctx context.Context, queueId string) error {
	query := dialect.Delete("queues").
		Where(goqu.I("queues.id").Eq(queueId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
