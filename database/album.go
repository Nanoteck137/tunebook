package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/library"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/types"
)

var createAlbumId = createIdGenerator(16)

type Album struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Name string `db:"name"`

	ArtistId string `db:"artist_id"`

	CoverArt  sql.NullString    `db:"cover_art"`
	Year      sql.NullInt64     `db:"year"`
	AlbumType library.AlbumType `db:"album_type"`

	ArtistName string `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists JsonColumn[[]FeaturingArtist] `db:"featuring_artists"`
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
			"albums.album_type",

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
			FeaturingArtistsQuery(
				"albums_featuring_artists", 
				"album_id",
			).As("featuring_artists"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}

func (db DB) GetAllAlbumIds(ctx context.Context) ([]string, error) {
	query := dialect.From("albums").
		Select("albums.id")

	return Multiple[string](db, ctx, query)
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

	page, err := buildPage(ctx, db, params.Page, query, "albums.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[Album](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetAlbumsByArtist(
	ctx context.Context,
	artistId string,
) ([]Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.artist_id").Eq(artistId))

	return Multiple[Album](db, ctx, query)
}

func (db DB) GetAlbumById(ctx context.Context, albumId string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.id").Eq(albumId))

	return Single[Album](db, ctx, query)
}

func (db DB) GetAlbumsIn(
	ctx context.Context,
	in any,
	sort string,
) ([]Album, error) {
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

	return Multiple[Album](db, ctx, query)
}

type CreateAlbumParams struct {
	Id string

	Name string

	ArtistId string

	CoverArt  sql.NullString
	Year      sql.NullInt64
	AlbumType library.AlbumType

	Created int64
	Updated int64
}

func (db DB) CreateAlbum(
	ctx context.Context,
	params CreateAlbumParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createAlbumId()
	}

	query := dialect.Insert("albums").
		Rows(goqu.Record{
			"id": params.Id,

			"name": params.Name,

			"artist_id": params.ArtistId,

			"cover_art":  params.CoverArt,
			"year":       params.Year,
			"album_type": params.AlbumType,

			"created": params.Created,
			"updated": params.Updated,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type AlbumChanges struct {
	Name Change[string]

	ArtistId Change[string]

	CoverArt  Change[sql.NullString]
	Year      Change[sql.NullInt64]
	AlbumType Change[library.AlbumType]

	Created Change[int64]
}

func (db DB) UpdateAlbum(
	ctx context.Context,
	albumId string,
	changes AlbumChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "cover_art", changes.CoverArt)
	addToRecord(record, "year", changes.Year)
	addToRecord(record, "album_type", changes.AlbumType)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("albums").
		Set(record).
		Where(goqu.I("albums.id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteAlbum(ctx context.Context, albumId string) error {
	query := dialect.Delete("albums").
		Where(goqu.I("albums.id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
