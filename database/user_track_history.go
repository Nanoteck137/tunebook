package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/types"
)

var createUserTrackHistoryId = createIdGenerator(32)

type UserTrackHistory struct {
	Id string `db:"id"`

	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	ListenedAt int64 `db:"listened_at"`

	PlaybackType string `db:"playback_type"`
	Status       string `db:"status"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func UserTrackHistoryQuery() *goqu.SelectDataset {
	query := dialect.From("user_track_history").
		Select(
			"user_track_history.id",

			"user_track_history.user_id",
			"user_track_history.track_id",

			"user_track_history.listened_at",

			"user_track_history.playback_type",
			"user_track_history.status",

			"user_track_history.updated",
			"user_track_history.created",
		).
		Join(
			TrackQuery().As("tracks"),
			goqu.On(goqu.I("user_track_history.track_id").Eq(goqu.I("tracks.id"))),
		)

	return query
}

type GetUserTrackHistoryParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetUserTrackHistory(
	ctx context.Context,
	params GetUserTrackHistoryParams,
) ([]UserTrackHistory, types.Page, error) {
	query := UserTrackHistoryQuery()

	var err error

		a := adapter.UserTrackHistoryResolverAdapter{}
	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, "user_track_history.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTrackHistory](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetUserTrackHistoryById(
	ctx context.Context,
	id string,
) (UserTrackHistory, error) {
	query := UserTrackHistoryQuery().
		Where(goqu.I("user_track_history.id").Eq(id))

	return Single[UserTrackHistory](db, ctx, query)
}

type CreateUserTrackHistoryParams struct {
	Id string

	UserId  string
	TrackId string

	ListenedAt int64

	PlaybackType string
	Status       string

	Created int64
	Updated int64
}

func (db DB) CreateUserTrackHistory(
	ctx context.Context,
	params CreateUserTrackHistoryParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createUserTrackHistoryId()
	}

	query := dialect.Insert("user_track_history").Rows(goqu.Record{
		"id": params.Id,

		"user_id":  params.UserId,
		"track_id": params.TrackId,

		"listened_at": params.ListenedAt,

		"playback_type": params.PlaybackType,
		"status":        params.Status,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

func (db DB) DeleteUserTrackHistory(ctx context.Context, user_track_historyId string) error {
	query := dialect.Delete("user_track_history").
		Where(goqu.I("user_track_history.id").Eq(user_track_historyId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
