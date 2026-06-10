package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/types"
)

type PlaylistItem struct {
	RowId int `db:"rowid"`

	PlaylistId string `db:"playlist_id"`
	TrackId    string `db:"track_id"`

	Position int `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type PlaylistItemTrack struct {
	Track

	Position int `db:"position"`
}

func PlaylistItemQuery() *goqu.SelectDataset {
	query := dialect.From("playlist_items").
		Select(
			"playlist_items.rowid",

			"playlist_items.playlist_id",
			"playlist_items.track_id",

			"playlist_items.position",

			"playlist_items.created",
			"playlist_items.updated",
		)

	return query
}

func (db DB) GetAllPlaylistItems(ctx context.Context) ([]PlaylistItem, error) {
	query := PlaylistItemQuery()
	return Multiple[PlaylistItem](db, ctx, query)
}

func (db DB) GetPlaylistItems(ctx context.Context, playlistId string) ([]PlaylistItem, error) {
	query := PlaylistItemQuery().
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.position").Asc())

	return Multiple[PlaylistItem](db, ctx, query)
}

func (db DB) GetPlaylistTrackImages(
	ctx context.Context,
	playlistId string,
	numImages int,
) ([]sql.NullString, error) {
	tracks := TrackQuery()

	query := dialect.From("playlist_items").
		Select("tracks.album_cover_art").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		GroupBy(goqu.I("tracks.album_id")).
		Order(goqu.I("playlist_items.position").Asc()).
		Limit(uint(numImages))

	return Multiple[sql.NullString](db, ctx, query)
}

func (db DB) GetNextPlaylistItemIndex(ctx context.Context, playlistId string) (int, error) {
	query := dialect.From("playlist_items").
		Select("playlist_items.position").
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.position").Desc()).
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

type GetPlaylistTracksParams struct {
	PlaylistId string
	Page       types.PageParams
	Filter     types.FilterParams
}

func (db DB) GetPlaylistTracks(
	ctx context.Context,
	params GetPlaylistTracksParams,
) ([]PlaylistItemTrack, types.Page, error) {
	tracks := TrackQuery()

	var err error

	query := dialect.From("playlist_items").
		Select("tracks.*", "playlist_items.position").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		)

	a := adapter.PlaylistTrackResolverAdapter{}
	query, err = applyFilterParamsCustom(
		params.Filter,
		&a,
		query,
		goqu.I("playlist_items.playlist_id").Eq(params.PlaylistId),
	)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, "tracks.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[PlaylistItemTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type GetPlaylistItemIdsParams struct {
	PlaylistId string
}

func (db DB) GetPlaylistItemIds(
	ctx context.Context,
	params GetPlaylistItemIdsParams,
) ([]string, error) {
	query := dialect.From("playlist_items").
		Select("playlist_items.track_id").
		Where(goqu.I("playlist_items.playlist_id").Eq(params.PlaylistId))

	items, err := Multiple[string](db, ctx, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreatePlaylistItemParams struct {
	PlaylistId string
	TrackId    string

	Position int

	Created int64
	Updated int64
}

func (db DB) CreatePlaylistItem(ctx context.Context, params CreatePlaylistItemParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("playlist_items").
		Rows(goqu.Record{
			"playlist_id": params.PlaylistId,
			"track_id":    params.TrackId,

			"position": params.Position,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type PlaylistItemChanges struct {
	Position types.Change[int]

	Created types.Change[int64]
}

func (db DB) UpdatePlaylistItem(ctx context.Context, playlistId, trackId string, changes PlaylistItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "position", changes.Position)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("playlist_items").
		Set(record).
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("playlist_items.track_id").Eq(trackId),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeletePlaylistItem(ctx context.Context, playlistId, trackId string) error {
	query := goqu.Delete("playlist_items").
		Where(goqu.And(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("playlist_items.track_id").Eq(trackId),
		))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetPlaylistItemByTrackId(ctx context.Context, playlistId, trackId string) (PlaylistItem, error) {
	query := PlaylistItemQuery().
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("playlist_items.track_id").Eq(trackId),
		)

	return Single[PlaylistItem](db, ctx, query)
}

func (db DB) ReorderPlaylistItemsAfterDelete(ctx context.Context, playlistId string, deletedPosition int) error {
	query := goqu.Update("playlist_items").
		Set(goqu.Record{
			"position": goqu.L("position - 1"),
		}).
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("position").Gt(deletedPosition),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
