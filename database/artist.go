package database

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
)

type Artist struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Name string `db:"name"`

	CoverArt sql.NullString `db:"cover_art"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`
}

func ArtistQuery() *goqu.SelectDataset {
	tags := dialect.From("artists_tags").
		Select(
			goqu.I("artists_tags.artist_id").As("artist_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("artists_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("artists_tags.artist_id"))

	query := dialect.From("artists").
		Select(
			"artists.rowid",

			"artists.id",

			"artists.name",

			"artists.cover_art",

			"artists.updated",
			"artists.created",

			goqu.I("tags.tags").As("tags"),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("artists.id").Eq(goqu.I("tags.artist_id"))),
		)

	return query
}

type GetAllArtistsParams struct {
	Filter types.FilterParams
}

func (db DB) GetAllArtists(
	ctx context.Context,
	params GetAllArtistsParams,
) ([]Artist, error) {
	query := ArtistQuery()

	var err error

	a := adapter.ArtistResolverAdapter{}
	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Artist](db.db, ctx, query)
}

// TODO(patrik): Move
func applyFilterParams(
	params types.FilterParams,
	adapter filter.ResolverAdapter,
	query *goqu.SelectDataset,
) (*goqu.SelectDataset, error) {
	resolver := filter.New(adapter)

	query, err := applyFilter(query, resolver, params.Filter)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, params.Sort)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func applyFilterParamsCustom(
	params types.FilterParams,
	adapter filter.ResolverAdapter,
	query *goqu.SelectDataset,
	customWhere exp.Expression,
) (*goqu.SelectDataset, error) {
	resolver := filter.New(adapter)

	query, err := applyFilterCustom(query, resolver, params.Filter, customWhere)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, params.Sort)
	if err != nil {
		return nil, err
	}

	return query, nil
}

// TODO(patrik): Move
func buildPage(
	ctx context.Context,
	db ember.DB,
	params types.PageParams,
	query *goqu.SelectDataset,
	countCol any,
) (types.Page, error) {
	countQuery := query.Select(goqu.COUNT(countCol))

	totalItems, err := ember.Single[int](db, ctx, countQuery)
	if err != nil {
		return types.Page{}, err
	}

	return types.Page{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		TotalPages: utils.TotalPages(params.PerPage, totalItems),
	}, nil
}

// TODO(patrik): Move
func applyPageParams(
	params types.PageParams,
	query *goqu.SelectDataset,
) *goqu.SelectDataset {
	return query.
		Limit(uint(params.PerPage)).
		Offset(uint(params.Page * params.PerPage))
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

	a := adapter.ArtistResolverAdapter{}
	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db.db, params.Page, query, "artists.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[Artist](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetArtistById(ctx context.Context, id string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.I("artists.id").Eq(id))

	return ember.Single[Artist](db.db, ctx, query)
}

func (db DB) GetArtistByName(ctx context.Context, name string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.Func("LOWER", goqu.I("artists.name")).Eq(strings.ToLower(name)))

	return ember.Single[Artist](db.db, ctx, query)
}

func (db DB) GetArtistsIn(ctx context.Context, in any, sort string) ([]Artist, error) {
	query := ArtistQuery().
		Where(
			goqu.I("artists.id").In(in),
		)

	a := adapter.ArtistResolverAdapter{}
	resolver := filter.New(&a)

	query, err := applySort(query, resolver, sort)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Artist](db.db, ctx, query)
}

func (db DB) GetAllArtistIds(ctx context.Context) ([]string, error) {
	query := dialect.From("artists").Select("artists.id")
	return ember.Multiple[string](db.db, ctx, query)
}

type CreateArtistParams struct {
	Id string

	Name string

	CoverArt sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateArtist(ctx context.Context, params CreateArtistParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateArtistId()
	}

	query := dialect.Insert("artists").Rows(goqu.Record{
		"id": params.Id,

		"name": params.Name,

		"cover_art": params.CoverArt,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type ArtistChanges struct {
	Name types.Change[string]

	CoverArt types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db DB) UpdateArtist(ctx context.Context, id string, changes ArtistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "cover_art", changes.CoverArt)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("artists").
		Set(record).
		Where(goqu.I("artists.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteArtist(ctx context.Context, id string) error {
	query := dialect.Delete("artists").
		Where(goqu.I("artists.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to AddArtistTag, same with track
func (db DB) AddTagToArtist(ctx context.Context, tagSlug, artistId string) error {
	ds := dialect.Insert("artists_tags").
		Rows(goqu.Record{
			"artist_id": artistId,
			"tag_slug":  tagSlug,
		})

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to RemoveAllArtistTags, same with track
func (db DB) RemoveAllTagsFromArtist(ctx context.Context, artistId string) error {
	query := dialect.Delete("artists_tags").
		Where(goqu.I("artist_id").Eq(artistId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
