package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
)

type TrackFilter struct {
	RowId int `db:"rowid"`

	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Name   string `db:"name"`
	Filter string `db:"filter"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func TrackFilterQuery() *goqu.SelectDataset {
	query := dialect.From("track_filters").
		Select(
			"track_filters.rowid",

			"track_filters.id",
			"track_filters.user_id",

			"track_filters.name",
			"track_filters.filter",

			"track_filters.created",
			"track_filters.updated",
		)

	return query
}

func (db DB) GetTrackFilterById(ctx context.Context, id, userId string) (TrackFilter, error) {
	query := TrackFilterQuery().
		Where(
			goqu.I("track_filters.id").Eq(id),
			goqu.I("track_filters.user_id").Eq(userId),
		)

	return ember.Single[TrackFilter](db.db, ctx, query)
}

func (db DB) GetTrackFiltersByUserId(ctx context.Context, userId string) ([]TrackFilter, error) {
	query := TrackFilterQuery().
		Where(goqu.I("track_filters.user_id").Eq(userId))

	return ember.Multiple[TrackFilter](db.db, ctx, query)
}

type CreateTrackFilterParams struct {
	Id     string
	UserId string

	Name   string
	Filter string

	Created int64
	Updated int64
}

func (db DB) CreateTrackFilter(ctx context.Context, params CreateTrackFilterParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateTrackFilterId()
	}

	query := dialect.Insert("track_filters").Rows(goqu.Record{
		"id":      params.Id,
		"user_id": params.UserId,

		"name":   params.Name,
		"filter": params.Filter,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type TrackFilterChanges struct {
	Name   types.Change[string]
	Filter types.Change[string]

	Created types.Change[int64]
}

func (db DB) UpdateTrackFilter(ctx context.Context, id, userId string, changes TrackFilterChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "filter", changes.Filter)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("track_filters").
		Set(record).
		Where(
			goqu.I("track_filters.id").Eq(id),
			goqu.I("track_filters.user_id").Eq(userId),
		)

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTrackFilter(ctx context.Context, id, userId string) error {
	query := dialect.Delete("track_filters").
		Where(
			goqu.I("track_filters.id").Eq(id),
			goqu.I("track_filters.user_id").Eq(userId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
