package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
)

type Playlist struct {
	Id       string         `db:"id"`
	Name     string         `db:"name"`
	CoverArt sql.NullString `db:"cover_art"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type PlaylistItem struct {
	PlaylistId string `db:"playlist_id"`
	TrackId    string `db:"track_id"`
	Order      int    `db:"order_num"`
}

type OrderedTrack struct {
	Track

	Order int `db:"order_num"`
}

func PlaylistQuery() *goqu.SelectDataset {
	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
			"playlists.cover_art",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",
		)

	return query
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

func (db DB) GetPlaylistItems(ctx context.Context, playlistId string) ([]PlaylistItem, error) {
	query := dialect.From("playlist_items").
		Select(
			"playlist_items.playlist_id",
			"playlist_items.track_id",
		).
		Where(goqu.I("playlist_id").Eq(playlistId))

	return ember.Multiple[PlaylistItem](db.db, ctx, query)
}

func (db DB) GetPlaylistTracks(ctx context.Context, playlistId string) ([]OrderedTrack, error) {
	tracks := TrackQuery()

	query := dialect.From("playlist_items").
		Select("tracks.*", "playlist_items.order_num").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.order_num").Asc())

	return ember.Multiple[OrderedTrack](db.db, ctx, query)
}

func (db DB) GetNextIndex(ctx context.Context, playlistId string) (int, error) {
	query := dialect.From("playlist_items").
		Select("playlist_items.order_num").
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

func (db DB) GetPlaylistTracksPaged(ctx context.Context, playlistId string, opts FetchOptions) ([]OrderedTrack, types.Page, error) {
	tracks := TrackQuery()

	query := dialect.From("playlist_items").
		Select("tracks.*", "playlist_items.order_num").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("playlist_items.order_num").Asc())

	var err error

	countQuery := query.
		Select(goqu.COUNT("tracks.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	totalItems, err := ember.Single[int](db.db, ctx, countQuery)
	if err != nil {
		return nil, types.Page{}, err
	}

	totalPages := utils.TotalPages(opts.PerPage, totalItems)
	page := types.Page{
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	items, err := ember.Multiple[OrderedTrack](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type CreatePlaylistParams struct {
	Id       string
	Name     string
	CoverArt sql.NullString

	OwnerId string

	Created int64
	Updated int64
}

func (db DB) CreatePlaylist(ctx context.Context, params CreatePlaylistParams) (Playlist, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateId()
	}

	query := dialect.Insert("playlists").
		Rows(goqu.Record{
			"id":        id,
			"name":      params.Name,
			"cover_art": params.CoverArt,

			"owner_id": params.OwnerId,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"playlists.id",
			"playlists.name",
			"playlists.cover_art",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",
		)

	return ember.Single[Playlist](db.db, ctx, query)
}

func (db DB) AddItemToPlaylist(ctx context.Context, playlistId, trackId string, order int) error {
	query := goqu.Insert("playlist_items").
		Rows(goqu.Record{
			"playlist_id": playlistId,
			"track_id":    trackId,
			"order_num":   order,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemovePlaylistItem(ctx context.Context, playlistId, trackId string) error {
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

func (db DB) RemoveAllPlaylistItem(ctx context.Context, playlistId string) error {
	query := goqu.Delete("playlist_items").
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId))

	_, err := db.db.Exec(ctx, query)
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
