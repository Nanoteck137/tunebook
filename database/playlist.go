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

var CreatePlaylistId = utils.CreateIdGenerator(16)

type Playlist struct {
	Id       string         `db:"id"`
	Name     string         `db:"name"`
	CoverArt sql.NullString `db:"cover_art"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	TrackCount sql.NullInt64 `db:"track_count"`
}

func PlaylistQuery() *goqu.SelectDataset {
	trackCountQuery := dialect.From("playlist_items").
		Select(
			goqu.I("playlist_items.playlist_id").As("id"),
			goqu.COUNT(goqu.I("playlist_items.track_id")).As("data"),
		).
		GroupBy(goqu.I("playlist_items.playlist_id"))

	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
			"playlists.cover_art",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",

			goqu.I("track_count.data").As("track_count"),
		).
		LeftJoin(
			trackCountQuery.As("track_count"),
			goqu.On(goqu.I("playlists.id").Eq(goqu.I("track_count.id"))),
		)

	return query
}

func (db DB) GetAllPlaylists(ctx context.Context) ([]Playlist, error) {
	query := PlaylistQuery()
	return ember.Multiple[Playlist](db.db, ctx, query)
}

func (db DB) GetPlaylistsByUser(ctx context.Context, userId string) ([]Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.owner_id").Eq(userId))

	return ember.Multiple[Playlist](db.db, ctx, query)
}

func (db DB) GetPlaylistById(ctx context.Context, id string) (Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.id").Eq(id))

	return ember.Single[Playlist](db.db, ctx, query)
}

func (db DB) GetPlaylistTrackImages(ctx context.Context, playlistId string, numImages int) ([]sql.NullString, error) {
	tracks := TrackQuery()

	query := dialect.From("playlist_items").
		Select("tracks.album_cover_art").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		GroupBy(goqu.I("tracks.album_id")).
		Order(goqu.I("playlist_items.order_num").Asc()).
		Limit(uint(numImages))

	return ember.Multiple[sql.NullString](db.db, ctx, query)
}

type CreatePlaylistParams struct {
	Id       string
	Name     string
	CoverArt sql.NullString

	OwnerId string

	Created int64
	Updated int64
}

func (db DB) CreatePlaylist(
	ctx context.Context,
	params CreatePlaylistParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()

		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = CreatePlaylistId()
	}

	query := dialect.Insert("playlists").
		Rows(goqu.Record{
			"id":        params.Id,
			"name":      params.Name,
			"cover_art": params.CoverArt,

			"owner_id": params.OwnerId,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type PlaylistChanges struct {
	Name types.Change[string]

	OwnerId types.Change[string]

	CoverArt types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db DB) UpdatePlaylist(
	ctx context.Context,
	id string,
	changes PlaylistChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "owner_id", changes.OwnerId)

	addToRecord(record, "cover_art", changes.CoverArt)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("playlists").
		Set(record).
		Where(goqu.I("playlists.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeletePlaylist(ctx context.Context, id string) error {
	query := dialect.Delete("playlists").
		Where(goqu.I("playlists.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
