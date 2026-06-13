package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/types"
)

var createTrackHistoryId = createIdGenerator(32)

type TrackHistory struct {
	Id string `db:"id"`

	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	ListenedAt int64 `db:"listened_at"`

	PlaybackType string `db:"playback_type"`
	Status       string `db:"status"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func TrackHistoryQuery() *goqu.SelectDataset {
	query := dialect.From("track_history").
		Select(
			"track_history.id",

			"track_history.user_id",
			"track_history.track_id",

			"track_history.listened_at",

			"track_history.playback_type",
			"track_history.status",

			"track_history.updated",
			"track_history.created",
		).
		Join(
			TrackQuery().As("tracks"),
			goqu.On(goqu.I("track_history.track_id").Eq(goqu.I("tracks.id"))),
		)

	return query
}

type GetTrackHistoryParams struct {
	UserId string
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetTrackHistory(
	ctx context.Context,
	params GetTrackHistoryParams,
) ([]TrackHistory, types.Page, error) {
	query := TrackHistoryQuery()

	var err error

	a := adapter.TrackHistoryResolverAdapter{}
	query, err = applyFilterParamsCustom(
		params.Filter, 
		&a, 
		query, 
		goqu.I("track_history.user_id").Eq(params.UserId),
	)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, "track_history.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[TrackHistory](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetTrackHistoryById(
	ctx context.Context,
	id string,
) (TrackHistory, error) {
	query := TrackHistoryQuery().
		Where(goqu.I("track_history.id").Eq(id))

	return Single[TrackHistory](db, ctx, query)
}

type CreateTrackHistoryParams struct {
	Id string

	UserId  string
	TrackId string

	ListenedAt int64

	PlaybackType string
	Status       string

	Created int64
	Updated int64
}

func (db DB) CreateTrackHistory(
	ctx context.Context,
	params CreateTrackHistoryParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createTrackHistoryId()
	}

	query := dialect.Insert("track_history").Rows(goqu.Record{
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

func (db DB) DeleteTrackHistory(ctx context.Context, track_historyId string) error {
	query := dialect.Delete("track_history").
		Where(goqu.I("track_history.id").Eq(track_historyId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
