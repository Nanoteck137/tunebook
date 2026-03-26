package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/pyrin/ember"
)

type Album struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Name string `db:"name"`

	ArtistId string `db:"artist_id"`

	CoverArt sql.NullString `db:"cover_art"`
	Year     sql.NullInt64  `db:"year"`

	ArtistName string `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	// TODO(patrik): Change to JsonColumn
	FeaturingArtists FeaturingArtists `db:"featuring_artists"`
}

func AlbumQuery() *goqu.SelectDataset {
	tags := dialect.From("albums_tags").
		Select(
			goqu.I("albums_tags.album_id").As("album_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("albums_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("albums_tags.album_id"))

	query := dialect.From("albums").
		Select(
			"albums.rowid",

			"albums.id",

			"albums.name",

			"albums.artist_id",

			"albums.cover_art",
			"albums.year",

			"albums.created",
			"albums.updated",

			goqu.I("artists.name").As("artist_name"),

			goqu.I("tags.tags").As("tags"),

			goqu.I("featuring_artists.artists").As("featuring_artists"),
		).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("tags.album_id"))),
		).
		LeftJoin(
			FeaturingArtistsQuery("albums_featuring_artists", "album_id").As("featuring_artists"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}

func (db DB) GetAllAlbumIds(ctx context.Context) ([]string, error) {
	query := dialect.From("albums").
		Select("albums.id")

	return ember.Multiple[string](db.db, ctx, query)
}

func (db DB) GetAllAlbums(ctx context.Context, filterStr string, sortStr string) ([]Album, error) {
	query := AlbumQuery()

	var err error

	a := adapter.AlbumResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Album](db.db, ctx, query)
}

type GetAlbumsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetAlbums(
	ctx context.Context,
	params GetAlbumsParams,
) ([]Album, types.Page, error) {
	query := AlbumQuery()

	var err error

	a := adapter.AlbumResolverAdapter{}

	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db.db, params.Page, query, "albums.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[Album](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetAlbumsByArtist(ctx context.Context, artistId string) ([]Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.artist_id").Eq(artistId))

	return ember.Multiple[Album](db.db, ctx, query)
}

func (db DB) GetAlbumById(ctx context.Context, id string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.id").Eq(id))

	return ember.Single[Album](db.db, ctx, query)
}

func (db DB) GetAlbumByName(ctx context.Context, name string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.name").Eq(name))

	return ember.Single[Album](db.db, ctx, query)
}

func (db DB) GetAlbumsIn(ctx context.Context, in any, sort string) ([]Album, error) {
	query := AlbumQuery().
		Where(
			goqu.I("albums.id").In(in),
		)

	a := adapter.AlbumResolverAdapter{}
	resolver := filter.New(&a)

	query, err := applySort(query, resolver, sort)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Album](db.db, ctx, query)
}

type CreateAlbumParams struct {
	Id string

	Name string

	ArtistId string

	CoverArt sql.NullString
	Year     sql.NullInt64

	Created int64
	Updated int64
}

func (db DB) CreateAlbum(ctx context.Context, params CreateAlbumParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateAlbumId()
	}

	query := dialect.Insert("albums").
		Rows(goqu.Record{
			"id": params.Id,

			"name": params.Name,

			"artist_id": params.ArtistId,

			"cover_art": params.CoverArt,
			"year":      params.Year,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) DeleteAlbum(ctx context.Context, id string) error {
	query := dialect.Delete("albums").
		Where(goqu.I("albums.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type AlbumChanges struct {
	Name      types.Change[string]

	ArtistId types.Change[string]

	CoverArt types.Change[sql.NullString]
	Year     types.Change[sql.NullInt64]

	Created types.Change[int64]
}

func (db DB) UpdateAlbum(ctx context.Context, id string, changes AlbumChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "cover_art", changes.CoverArt)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("albums").
		Set(record).
		Where(goqu.I("albums.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to AddAlbumTag, same with track
func (db DB) AddTagToAlbum(ctx context.Context, tagSlug, albumId string) error {
	ds := dialect.Insert("albums_tags").
		Rows(goqu.Record{
			"album_id": albumId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to RemoveAllAlbumTags, same with track
func (db DB) RemoveAllTagsFromAlbum(ctx context.Context, albumId string) error {
	query := dialect.Delete("albums_tags").
		Where(goqu.I("album_id").Eq(albumId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) RemoveAllAlbumFeaturingArtists(ctx context.Context, albumId string) error {
	query := dialect.Delete("albums_featuring_artists").
		Where(
			goqu.I("albums_featuring_artists.album_id").Eq(albumId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) AddFeaturingArtistToAlbum(ctx context.Context, albumId, artistId string) error {
	query := dialect.Insert("albums_featuring_artists").
		Rows(goqu.Record{
			"album_id":  albumId,
			"artist_id": artistId,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveFeaturingArtistFromAlbum(ctx context.Context, albumId, artistId string) error {
	query := goqu.Delete("albums_featuring_artists").
		Where(
			goqu.And(
				goqu.I("albums_featuring_artists.album_id").Eq(albumId),
				goqu.I("albums_featuring_artists.artist_id").Eq(artistId),
			),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
