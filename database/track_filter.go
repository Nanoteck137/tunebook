package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	createTrackFilterId = createIdGenerator(8)

	trackFiltersTbl = goqu.T("track_filters")
)

type TrackFilter struct {
	Id string `db:"id"`

	UserId string `db:"user_id"`

	Name   string `db:"name"`
	Filter string `db:"filter"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func TrackFilterQuery() *goqu.SelectDataset {
	query := dialect.From(trackFiltersTbl).
		Select(
			trackFiltersTbl.Col("id"),

			trackFiltersTbl.Col("user_id"),

			trackFiltersTbl.Col("name"),
			trackFiltersTbl.Col("filter"),

			trackFiltersTbl.Col("created"),
			trackFiltersTbl.Col("updated"),
		)

	return query
}

type CreateTrackFilterParams struct {
	Id     string
	UserId string

	Name   string
	Filter string

	Created int64
	Updated int64
}

func (db DB) CreateTrackFilter(
	ctx context.Context,
	params CreateTrackFilterParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createTrackFilterId()
	}

	query := dialect.Insert(trackFiltersTbl).Rows(goqu.Record{
		"id": params.Id,

		"user_id": params.UserId,

		"name":   params.Name,
		"filter": params.Filter,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type TrackFilterChanges struct {
	UserId Change[string]

	Name   Change[string]
	Filter Change[string]

	Created Change[int64]
}

func (db DB) UpdateTrackFilter(
	ctx context.Context,
	id string,
	changes TrackFilterChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "user_id", changes.UserId)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "filter", changes.Filter)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update(trackFiltersTbl).
		Set(record).
		Where(trackFiltersTbl.Col("id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTrackFilter(ctx context.Context, id string) error {
	query := dialect.Delete(trackFiltersTbl).
		Where(trackFiltersTbl.Col("id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetTrackFilterById(
	ctx context.Context,
	id string,
) (TrackFilter, error) {
	query := TrackFilterQuery().
		Where(trackFiltersTbl.Col("id").Eq(id))

	return Single[TrackFilter](db, ctx, query)
}

func (db DB) GetTrackFiltersByUserId(
	ctx context.Context,
	userId string,
) ([]TrackFilter, error) {
	query := TrackFilterQuery().
		Where(trackFiltersTbl.Col("user_id").Eq(userId))

	return Multiple[TrackFilter](db, ctx, query)
}
