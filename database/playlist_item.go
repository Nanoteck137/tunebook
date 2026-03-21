package database

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
)

type PlaylistItem struct {
	RowId int `db:"rowid"`

	PlaylistId string `db:"playlist_id"`
	TrackId    string `db:"track_id"`

	Order int `db:"order_num"`
}

type OrderedTrack struct {
	Track

	Order int `db:"order_num"`
}

func PlaylistItemQuery() *goqu.SelectDataset {
	query := dialect.From("playlist_items").
		Select(
			"playlist_items.rowid",

			"playlist_items.playlist_id",
			"playlist_items.track_id",

			"playlist_items.order_num",
		)

	return query
}

func (db DB) GetAllPlaylistItems(ctx context.Context) ([]PlaylistItem, error) {
	query := PlaylistItemQuery()
	return ember.Multiple[PlaylistItem](db.db, ctx, query)
}

func (db DB) GetPlaylistItems(ctx context.Context, playlistId string) ([]PlaylistItem, error) {
	query := PlaylistItemQuery().
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.order_num").Asc())

	return ember.Multiple[PlaylistItem](db.db, ctx, query)
}

func (db DB) GetNextPlaylistItemIndex(ctx context.Context, playlistId string) (int, error) {
	query := dialect.From("playlist_items").
		Select("playlist_items.order_num").
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.order_num").Desc()).
		Limit(1)

	res, err := ember.Single[int](db.db, ctx, query)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return res + 1, nil
}

// func (db DB) GetPlaylistTracks(ctx context.Context, playlistId, filterStr string) ([]OrderedTrack, error) {
// 	tracks := TrackQuery()
//
// 	var err error
//
// 	a := adapter.TrackResolverAdapter{}
// 	resolver := filter.New(&a)
//
// 	tracks, err = applyFilter(tracks, resolver, filterStr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// tracks, err = applySort(tracks, resolver, sortStr)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
//
// 	query := dialect.From("playlist_items").
// 		Select("tracks.*", "playlist_items.order_num").
// 		Join(
// 			tracks.As("tracks"),
// 			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
// 		).
// 		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
// 		Order(goqu.I("playlist_items.order_num").Asc())
//
// 	return ember.Multiple[OrderedTrack](db.db, ctx, query)
// }

type GetPlaylistTracksParams struct {
	PlaylistId string
	Page       types.PageParams
	Filter     types.FilterParams
}

func (db DB) GetPlaylistTracks(
	ctx context.Context, 
	params GetPlaylistTracksParams,
) ([]OrderedTrack, types.Page, error) {
	tracks := TrackQuery()

	var err error

	// TODO(patrik): Fix this, using order for playlist_items

	a := adapter.TrackResolverAdapter{}
	tracks, err = applyFilterParams(params.Filter, &a, tracks)
	if err != nil {
		return nil, types.Page{}, err
	}

	query := dialect.From("playlist_items").
		Select("tracks.*", "playlist_items.order_num").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).
		Where(goqu.I("playlist_items.playlist_id").Eq(params.PlaylistId)).
		Order(goqu.I("playlist_items.order_num").Asc())

	page, err := buildPage(ctx, db.db, params.Page, query, "tracks.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[OrderedTrack](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type CreatePlaylistItemParams struct {
	PlaylistId string
	TrackId    string

	Order int
}

func (db DB) CreatePlaylistItem(ctx context.Context, params CreatePlaylistItemParams) error {
	query := dialect.Insert("playlist_items").
		Rows(goqu.Record{
			"playlist_id": params.PlaylistId,
			"track_id":    params.TrackId,

			"order_num": params.Order,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type PlaylistItemChanges struct {
	Order types.Change[int]
}

func (db DB) UpdatePlaylistItem(ctx context.Context, playlistId, trackId string, changes PlaylistItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "order_num", changes.Order)

	if len(record) == 0 {
		return nil
	}

	query := dialect.Update("playlist_items").
		Set(record).
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("playlist_items.track_id").Eq(trackId),
		)

	_, err := db.db.Exec(ctx, query)
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

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
