package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/types"
)

var (
	userFavoritesTbl = goqu.T("user_favorites")

	userFavoriteTrackSchema = UserFavoriteTrackSchema()
)

type UserFavorite struct {
	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	Added int64 `db:"added"`
}

type UserFavoriteTrack struct {
	Track

	Added int64 `db:"added"`
}

func UserFavoriteTrackSchema() *schema.Schema {
	return TrackSchema().
		AddField("added", query.TypeInt, schema.Column("user_favorites.added")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "added"},
				Dir:   query.DirDesc,
			},
		)
}

func UserFavoriteQuery() *goqu.SelectDataset {
	query := dialect.From(userFavoritesTbl).
		Select(
			userFavoritesTbl.Col("user_id"),
			userFavoritesTbl.Col("track_id"),

			userFavoritesTbl.Col("added"),
		)

	return query
}

type CreateUserFavoriteParams struct {
	UserId  string
	TrackId string

	Added int64
}

func (db DB) CreateUserFavorite(
	ctx context.Context, 
	params CreateUserFavoriteParams,
) error {
	if params.Added == 0 {
		params.Added = time.Now().UnixMilli()
	}

	query := dialect.Insert(userFavoritesTbl).
		Rows(goqu.Record{
			"user_id":  params.UserId,
			"track_id": params.TrackId,

			"added": params.Added,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUserFavorite(
	ctx context.Context, 
	userId, trackId string,
) error {
	query := dialect.Delete(userFavoritesTbl).
		Where(
			userFavoritesTbl.Col("user_id").Eq(userId),
			userFavoritesTbl.Col("track_id").Eq(trackId),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetUserFavoritesIds(
	ctx context.Context, 
	userId string,
) ([]string, error) {
	query := UserFavoriteQuery().
		Select(userFavoritesTbl.Col("track_id")).
		Where(userFavoritesTbl.Col("user_id").Eq(userId))

	items, err := Multiple[string](db, ctx, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type GetUserFavoriteTracksParams struct {
	UserId string
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetUserFavoriteTracks(
	ctx context.Context,
	params GetUserFavoriteTracksParams,
) ([]UserFavoriteTrack, types.Page, error) {
	var err error

	query := TrackQuery().
		SelectAppend(
			userFavoritesTbl.Col("added"),
		).
		Join(
			userFavoritesTbl,
			goqu.On(userFavoritesTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		)

	query = query.Where(userFavoritesTbl.Col("user_id").Eq(params.UserId))

	query, err = ApplyQuery(query, userFavoriteTrackSchema, QueryParams{
		Filter: params.Filter.Filter,
		Sort:   params.Filter.Sort,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, tracksTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserFavoriteTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetUserFavoriteCount(
	ctx context.Context,
	userId string,
) (int, error) {
	query := dialect.From(userFavoritesTbl).
		Select(
			goqu.COUNT(userFavoritesTbl.Col("track_id")),
		).
		Where(userFavoritesTbl.Col("user_id").Eq(userId)).
		GroupBy(userFavoritesTbl.Col("user_id"))

	return Single[int](db, ctx, query)
}
