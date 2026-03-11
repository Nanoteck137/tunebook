package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
)

type VirtualPlaylist struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Name string `db:"name"`

	OwnerId string `db:"owner_id"`

	PlaylistId sql.NullString `db:"playlist_id"`

	Filter string `db:"filter"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func VirtualPlaylistQuery() *goqu.SelectDataset {
	query := dialect.From("virtual_playlists").
		Select(
			"virtual_playlists.rowid",

			"virtual_playlists.id",

			"virtual_playlists.name",

			// TODO(patrik): Fetch more from owner
			"virtual_playlists.owner_id",

			"virtual_playlists.playlist_id",

			"virtual_playlists.filter",

			"virtual_playlists.created",
			"virtual_playlists.updated",
		)

	return query
}

type CreateVirtualPlaylistParams struct {
	Id string

	Name string

	OwnerId string

	PlaylistId sql.NullString

	Filter string

	Created int64
	Updated int64
}

func (db DB) CreateVirtualPlaylist(ctx context.Context, params CreateVirtualPlaylistParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateVirtualPlaylistId()
	}

	query := dialect.Insert("virtual_playlists").
		Rows(goqu.Record{
			"id": params.Id,

			"name": params.Name,

			"owner_id": params.OwnerId,

			"playlist_id": params.PlaylistId,

			"filter": params.Filter,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) GetVirtualPlaylistByUser(ctx context.Context, userId string) ([]VirtualPlaylist, error) {
	query := VirtualPlaylistQuery().
		Where(goqu.I("virtual_playlists.owner_id").Eq(userId))

	return ember.Multiple[VirtualPlaylist](db.db, ctx, query)
}

func (db DB) GetVirtualPlaylistForPlaylist(ctx context.Context, playlistId string) ([]VirtualPlaylist, error) {
	query := VirtualPlaylistQuery().
		Where(goqu.I("virtual_playlists.playlist_id").Eq(playlistId))

	return ember.Multiple[VirtualPlaylist](db.db, ctx, query)
}

func (db DB) GetVirtualPlaylistById(ctx context.Context, id string) (VirtualPlaylist, error) {
	query := VirtualPlaylistQuery().
		Where(goqu.I("virtual_playlists.id").Eq(id))

	return ember.Single[VirtualPlaylist](db.db, ctx, query)
}

type VirtualPlaylistChanges struct {
	Name types.Change[string]

	OwnerId types.Change[string]

	PlaylistId types.Change[sql.NullString]

	Filter types.Change[string]

	Created types.Change[int64]
}

func (db DB) UpdateVirtualPlaylist(ctx context.Context, id string, changes VirtualPlaylistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "owner_id", changes.OwnerId)

	addToRecord(record, "playlist_id", changes.PlaylistId)

	addToRecord(record, "filter", changes.Filter)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("virtual_playlists").
		Set(record).
		Where(goqu.I("virtual_playlists.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteVirtualPlaylist(ctx context.Context, id string) error {
	query := dialect.Delete("virtual_playlists").
		Where(goqu.I("virtual_playlists.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
