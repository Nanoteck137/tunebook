package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/types"
)

var createPlaylistId = createIdGenerator(16)

type Playlist struct {
	Id       string         `db:"id"`
	Name     string         `db:"name"`
	CoverArt sql.NullString `db:"cover_art"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	OwnerDisplayName string         `db:"owner_display_name"`
	OwnerPicture     sql.NullString `db:"owner_picture"`

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

			goqu.I("owner.display_name").As("owner_display_name"),
			goqu.I("owner.picture").As("owner_picture"),

			goqu.I("track_count.data").As("track_count"),
		).
		Join(
			UserQuery().As("owner"),
			goqu.On(goqu.I("playlists.owner_id").Eq(goqu.I("owner.id"))),
		).
		LeftJoin(
			trackCountQuery.As("track_count"),
			goqu.On(goqu.I("playlists.id").Eq(goqu.I("track_count.id"))),
		)

	return query
}

type GetPlaylistsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetPlaylists(
	ctx context.Context,
	params GetPlaylistsParams,
) ([]Playlist, types.Page, error) {
	query := PlaylistQuery()

	var err error

	a := adapter.PlaylistResolverAdapter{}
	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, "playlists.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[Playlist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetPlaylistsIn(
	ctx context.Context,
	in any,
	sort string,
) ([]Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.id").In(in))

	a := adapter.PlaylistResolverAdapter{}
	resolver := filter.New(&a)

	query, err := applySort(query, resolver, sort)
	if err != nil {
		return nil, err
	}

	return Multiple[Playlist](db, ctx, query)
}

func (db DB) GetPlaylistById(
	ctx context.Context,
	playlistId string,
) (Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.id").Eq(playlistId))

	return Single[Playlist](db, ctx, query)
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
		params.Id = createPlaylistId()
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

	_, err := db.Exec(ctx, query)
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
	playlistId string,
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

	query := dialect.Update("playlists").
		Set(record).
		Where(goqu.I("playlists.id").Eq(playlistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeletePlaylist(ctx context.Context, playlistId string) error {
	query := dialect.Delete("playlists").
		Where(goqu.I("playlists.id").Eq(playlistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
