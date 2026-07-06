package database

import (
	"context"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/types"
)

var createQueueItemId = createIdGenerator(32)

type QueueItem struct {
	Id       string `db:"id"`
	QueueId  string `db:"queue_id"`
	UserId   string `db:"user_id"`
	TrackId  string `db:"track_id"`
	Position int    `db:"position"`
	Created  int64  `db:"created"`
}

type QueueItemTrack struct {
	Track

	Position int `db:"position"`
}

type QueueItemEntry struct {
	QueueItemId string `db:"id"`
	TrackId     string `db:"track_id"`
	Position    int    `db:"position"`
}

func QueueItemQuery() *goqu.SelectDataset {
	return dialect.From("queue_items").
		Select(
			"queue_items.id",
			"queue_items.queue_id",
			"queue_items.user_id",
			"queue_items.track_id",
			"queue_items.position",
			"queue_items.created",
		)
}

func QueueItemTrackQuery() *goqu.SelectDataset {
	tracks := TrackQuery()

	return dialect.From("queue_items").
		Select("tracks.*", "queue_items.position").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("queue_items.track_id").Eq(goqu.I("tracks.id"))),
		)
}

type CreateQueueItemsParams struct {
	QueueId string
	UserId  string
	Items   []CreateQueueItemParams
}

type CreateQueueItemParams struct {
	TrackId  string
	Position int
}

func (db DB) CreateQueueItems(
	ctx context.Context,
	params CreateQueueItemsParams,
) error {
	t := time.Now().UnixMilli()

	rows := make([]goqu.Record, len(params.Items))
	for i, item := range params.Items {
		rows[i] = goqu.Record{
			"id":       createQueueItemId(),
			"queue_id": params.QueueId,
			"user_id":  params.UserId,
			"track_id": item.TrackId,
			"position": item.Position,
			"created":  t,
		}
	}

	if len(rows) == 0 {
		return nil
	}

	query := dialect.Insert("queue_items").Rows(rows)

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) GetNextQueueItemPosition(
	ctx context.Context,
	queueId string,
	userId string,
) (int, error) {
	query := dialect.From("queue_items").
		Select("queue_items.position").
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
		).
		Order(goqu.I("queue_items.position").Desc()).
		Limit(1)

	res, err := Single[int](db, ctx, query)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return res + 1, nil
}

type GetQueueItemsParams struct {
	Page types.PageParams

	QueueId string
	UserId  string
}

func (db DB) GetQueueItems(
	ctx context.Context,
	params GetQueueItemsParams,
) ([]QueueItemTrack, types.Page, error) {
	query := QueueItemTrackQuery().
		Where(
			goqu.I("queue_items.queue_id").Eq(params.QueueId),
			goqu.I("queue_items.user_id").Eq(params.UserId),
		).
		Order(goqu.I("queue_items.position").Asc())

	page, err := buildPage(ctx, db, params.Page, query, "queue_items.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[QueueItemTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetQueueItemIds(
	ctx context.Context,
	queueId string,
	userId string,
) ([]string, error) {
	query := dialect.From("queue_items").
		Select("queue_items.track_id").
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
		).
		Order(goqu.I("queue_items.position").Asc())

	return Multiple[string](db, ctx, query)
}

func (db DB) GetQueueItemEntries(
	ctx context.Context,
	queueId string,
	userId string,
) ([]QueueItemEntry, error) {
	query := dialect.From("queue_items").
		Select(
			"queue_items.id",
			"queue_items.track_id",
			"queue_items.position",
		).
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
		).
		Order(goqu.I("queue_items.position").Asc())

	return Multiple[QueueItemEntry](db, ctx, query)
}

func (db DB) GetQueueItemAtPosition(
	ctx context.Context,
	queueId string,
	userId string,
	position int,
) (QueueItemTrack, error) {
	query := QueueItemTrackQuery().
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
			goqu.I("queue_items.position").Eq(position),
		)

	return Single[QueueItemTrack](db, ctx, query)
}

func (db DB) GetQueueItemById(
	ctx context.Context,
	itemId string,
) (QueueItem, error) {
	query := QueueItemQuery().
		Where(goqu.I("queue_items.id").Eq(itemId))

	return Single[QueueItem](db, ctx, query)
}

func (db DB) DeleteQueueItem(ctx context.Context, itemId string) error {
	query := goqu.Delete("queue_items").
		Where(goqu.I("queue_items.id").Eq(itemId))

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) ClearQueueItems(ctx context.Context, queueId string, userId string) error {
	query := goqu.Delete("queue_items").
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
		)

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) GetQueueItemCount(
	ctx context.Context,
	queueId string,
	userId string,
) (int, error) {
	query := dialect.From("queue_items").
		Select(goqu.COUNT("queue_items.id")).
		Where(
			goqu.I("queue_items.queue_id").Eq(queueId),
			goqu.I("queue_items.user_id").Eq(userId),
		)

	return Single[int](db, ctx, query)
}
