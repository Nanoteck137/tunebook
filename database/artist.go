package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/types"
)

var (
	createArtistId = createIdGenerator(10)

	artistsTbl = goqu.T("artists")

	artistSchema = ArtistSchema()
)

type Artist struct {
	Id string `db:"id"`

	Name string `db:"name"`

	CoverArt sql.NullString `db:"cover_art"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`
}

func ArtistSchema() *schema.Schema {
	return schema.New().
		AddField("id", query.TypeString, schema.Column("artists.id")).
		AddField("name", query.TypeString, schema.Column("artists.name")).
		AddField(
			"coverArt",
			query.TypeString,
			schema.Column("artists.cover_art"),
			schema.Nullable(),
		).
		AddField(
			"tags",
			query.TypeRelation,
			schema.Relation(
				"artists_tags", "artist_id", "tag_slug", query.TypeString, "artists.id"),
		).
		AddField("created", query.TypeInt, schema.Column("artists.created")).
		AddField("updated", query.TypeInt, schema.Column("artists.updated")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "name"},
				Dir:   query.DirAsc,
			},
		)
}

func ArtistQuery() *goqu.SelectDataset {
	idCol := artistsTbl.Col("id")

	query := dialect.From(artistsTbl).
		Select(
			idCol,

			artistsTbl.Col("name"),

			artistsTbl.Col("cover_art"),

			artistsTbl.Col("created"),
			artistsTbl.Col("updated"),
		)

	query = AddTagsToQuery(query, idCol, artistsTagsTbl, "artist_id")

	return query
}

type CreateArtistParams struct {
	Id string

	Name string

	CoverArt sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateArtist(
	ctx context.Context,
	params CreateArtistParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createArtistId()
	}

	query := dialect.Insert(artistsTbl).Rows(goqu.Record{
		"id": params.Id,

		"name": params.Name,

		"cover_art": params.CoverArt,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type ArtistChanges struct {
	Name Change[string]

	CoverArt Change[sql.NullString]

	Created Change[int64]
}

func (db DB) UpdateArtist(
	ctx context.Context,
	artistId string,
	changes ArtistChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "cover_art", changes.CoverArt)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update(artistsTbl).
		Set(record).
		Where(artistsTbl.Col("id").Eq(artistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteArtist(ctx context.Context, artistId string) error {
	query := dialect.Delete(artistsTbl).
		Where(artistsTbl.Col("id").Eq(artistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type GetArtistsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetArtists(
	ctx context.Context,
	params GetArtistsParams,
) ([]Artist, types.Page, error) {
	query := ArtistQuery()

	var err error

	query, err = ApplyQuery(query, artistSchema, QueryParams{
		Filter: params.Filter.Filter,
		Sort:   params.Filter.Sort,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, artistsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[Artist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetArtistById(
	ctx context.Context,
	artistId string,
) (Artist, error) {
	query := ArtistQuery().
		Where(artistsTbl.Col("id").Eq(artistId))

	return Single[Artist](db, ctx, query)
}

func (db DB) GetArtistsByIds(
	ctx context.Context,
	ids []string,
) ([]Artist, error) {
	query := ArtistQuery().Where(artistsTbl.Col("id").In(ids))

	return Multiple[Artist](db, ctx, query)
}

func (db DB) GetAllArtistIds(ctx context.Context) ([]string, error) {
	query := dialect.From(artistsTbl).Select(artistsTbl.Col("id"))
	return Multiple[string](db, ctx, query)
}
