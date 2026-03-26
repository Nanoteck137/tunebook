package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/pyrin/ember"
)

type PlaylistFilter struct {
	RowId int `db:"rowid"`

	Id         string `db:"id"`
	PlaylistId string `db:"playlist_id"`

	Name   string `db:"name"`
	Filter string `db:"filter"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func PlaylistFilterQuery() *goqu.SelectDataset {
	query := dialect.From("playlist_filters").
		Select(
			"playlist_filters.rowid",

			"playlist_filters.id",
			"playlist_filters.playlist_id",

			"playlist_filters.name",
			"playlist_filters.filter",

			"playlist_filters.created",
			"playlist_filters.updated",
		)

	return query
}

func (db DB) GetPlaylistFilterById(ctx context.Context, id, playlistId string) (PlaylistFilter, error) {
	query := PlaylistFilterQuery().
		Where(
			goqu.I("playlist_filters.id").Eq(id),
			goqu.I("playlist_filters.playlist_id").Eq(playlistId),
		)

	return ember.Single[PlaylistFilter](db.db, ctx, query)
}

func (db DB) GetPlaylistFiltersByPlaylistId(ctx context.Context, playlistId string) ([]PlaylistFilter, error) {
	query := PlaylistFilterQuery().
		Where(goqu.I("playlist_filters.playlist_id").Eq(playlistId))

	return ember.Multiple[PlaylistFilter](db.db, ctx, query)
}

type CreatePlaylistFilterParams struct {
	Id         string
	PlaylistId string

	Name   string
	Filter string

	Created int64
	Updated int64
}

func (db DB) CreatePlaylistFilter(ctx context.Context, params CreatePlaylistFilterParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreatePlaylistFilterId()
	}

	query := dialect.Insert("playlist_filters").Rows(goqu.Record{
		"id":          params.Id,
		"playlist_id": params.PlaylistId,

		"name":   params.Name,
		"filter": params.Filter,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type PlaylistFilterChanges struct {
	Name   types.Change[string]
	Filter types.Change[string]

	Created types.Change[int64]
}

func (db DB) UpdatePlaylistFilter(ctx context.Context, id, playlistId string, changes PlaylistFilterChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "filter", changes.Filter)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("playlist_filters").
		Set(record).
		Where(
			goqu.I("playlist_filters.id").Eq(id),
			goqu.I("playlist_filters.playlist_id").Eq(playlistId),
		)

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeletePlaylistFilter(ctx context.Context, id, playlistId string) error {
	query := dialect.Delete("playlist_filters").
		Where(
			goqu.I("playlist_filters.id").Eq(id),
			goqu.I("playlist_filters.playlist_id").Eq(playlistId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
