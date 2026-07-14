package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/library"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/types"
)

var (
	createAlbumId = createIdGenerator(16)

	albumsTbl = goqu.T("albums")

	albumSchema = AlbumSchema()
)

type Album struct {
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

func AlbumSchema() *schema.Schema {
	return schema.New().
		AddField("id", query.TypeString, schema.Column("albums.id")).
		AddField("name", query.TypeString, schema.Column("albums.name")).
		AddField(
			"artistId",
			query.TypeString,
			schema.Column("albums.artist_id"),
		).
		AddField(
			"coverArt",
			query.TypeString,
			schema.Column("albums.cover_art"),
			schema.Nullable(),
		).
		AddField(
			"year",
			query.TypeInt,
			schema.Column("albums.year"),
			schema.Nullable(),
		).
		AddField(
			"albumType", query.TypeString, schema.Column("albums.album_type")).
		AddField(
			"artistName", query.TypeString, schema.Column("artists.name")).
		AddField(
			"tags",
			query.TypeRelation,
			schema.Relation(
				"albums_tags", "album_id", "tag_slug", query.TypeString, "albums.id"),
		).
		AddField(
			"featuringArtists",
			query.TypeRelation,
			schema.Relation(
				"albums_featuring_artists",
				"album_id",
				"artist_id",
				query.TypeString,
				"albums.id",
			),
		).
		AddField("created", query.TypeInt, schema.Column("albums.created")).
		AddField("updated", query.TypeInt, schema.Column("albums.updated")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "name"},
				Dir:   query.DirAsc,
			},
		)
}

func AlbumQuery() *goqu.SelectDataset {
	idCol := albumsTbl.Col("id")

	query := dialect.From(albumsTbl).
		Select(
			idCol,

			albumsTbl.Col("name"),

			albumsTbl.Col("artist_id"),

			albumsTbl.Col("cover_art"),
			albumsTbl.Col("year"),
			albumsTbl.Col("album_type"),

			albumsTbl.Col("created"),
			albumsTbl.Col("updated"),

			artistsTbl.Col("name").As("artist_name"),
		).
		Join(
			artistsTbl,
			goqu.On(albumsTbl.Col("artist_id").Eq(artistsTbl.Col("id"))),
		)

	query = AddTagsToQuery(query, idCol, albumsTagsTbl, "album_id")
	query = AddFeaturingArtistsToQuery(
		query, idCol, albumsFeaturingArtistsTbl, "album_id")

	return query
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

	query := dialect.Insert(albumsTbl).
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

	query := dialect.Update(albumsTbl).
		Set(record).
		Where(albumsTbl.Col("id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteAlbum(ctx context.Context, albumId string) error {
	query := dialect.Delete(albumsTbl).
		Where(albumsTbl.Col("id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetAllAlbumIds(ctx context.Context) ([]string, error) {
	query := dialect.From(albumsTbl).
		Select(albumsTbl.Col("id"))

	return Multiple[string](db, ctx, query)
}

type GetAlbumsParams struct {
	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetAlbums(
	ctx context.Context,
	params GetAlbumsParams,
) ([]Album, types.Page, error) {
	query := AlbumQuery()

	var err error

	query, err = ApplyQuery(query, albumSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, albumsTbl.Col("id"))
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
		Where(albumsTbl.Col("artist_id").Eq(artistId))

	return Multiple[Album](db, ctx, query)
}

func (db DB) GetAlbumById(ctx context.Context, albumId string) (Album, error) {
	query := AlbumQuery().
		Where(albumsTbl.Col("id").Eq(albumId))

	return Single[Album](db, ctx, query)
}

func (db DB) GetAlbumsByIds(
	ctx context.Context,
	ids []string,
) ([]Album, error) {
	query := AlbumQuery().Where(albumsTbl.Col("id").In(ids))

	return Multiple[Album](db, ctx, query)
}
