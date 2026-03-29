package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/types"
)

type UserFavorite struct {
	RowId int `db:"rowid"`

	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	Added int64 `db:"added"`
}

type UserFavoriteTrack struct {
	Track

	Added int64 `db:"added"`
}

func UserFavoriteQuery() *goqu.SelectDataset {
	query := dialect.From("user_favorites").
		Select(
			"user_favorites.rowid",

			"user_favorites.user_id",
			"user_favorites.track_id",

			"user_favorites.added",
		)

	return query
}

func (db DB) GetUserFavoritesIds(ctx context.Context, userId string) ([]string, error) {
	query := UserFavoriteQuery().
		Select("user_favorites.track_id").
		Where(goqu.I("user_favorites.user_id").Eq(userId))

	items, err := ember.Multiple[string](db.db, ctx, query)
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
	tracks := TrackQuery()

	var err error

	query := dialect.From("user_favorites").
		Select("tracks.*", "user_favorites.added").
		Join(
			tracks.As("tracks"),
			goqu.On(goqu.I("user_favorites.track_id").Eq(goqu.I("tracks.id"))),
		)

	// TODO(patrik): Replace with custom
	a := adapter.PlaylistTrackResolverAdapter{}
	query, err = applyFilterParamsCustom(
		params.Filter,
		&a,
		query,
		goqu.I("user_favorites.user_id").Eq(params.UserId),
	)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db.db, params.Page, query, "tracks.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[UserFavoriteTrack](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type CreateUserFavoriteParams struct {
	UserId  string
	TrackId string

	Added int64
}

func (db DB) CreateUserFavorite(ctx context.Context, params CreateUserFavoriteParams) error {
	if params.Added == 0 {
		params.Added = time.Now().UnixMilli()
	}

	query := dialect.Insert("user_favorites").
		Rows(goqu.Record{
			"user_id":  params.UserId,
			"track_id": params.TrackId,

			"added": params.Added,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUserFavorite(ctx context.Context, userId, trackId string) error {
	query := goqu.Delete("user_favorites").
		Where(
			goqu.I("user_favorites.user_id").Eq(userId),
			goqu.I("user_favorites.track_id").Eq(trackId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
