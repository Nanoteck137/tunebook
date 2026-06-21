package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/types"
)

var createTrackHistoryId = createIdGenerator(32)

type TrackHistory struct {
	Track

	Id string `db:"history_id"`

	UserId  string `db:"history_user_id"`
	TrackId string `db:"history_track_id"`

	ListenedAt int64 `db:"history_listened_at"`

	PlaybackType string `db:"history_playback_type"`
	Status       string `db:"history_status"`

	PercentPlayed int `db:"history_percent_played"`

	Created int64 `db:"history_created"`
	Updated int64 `db:"history_updated"`
}

func TrackHistoryQuery() *goqu.SelectDataset {
	tbl := goqu.T("track_history")

	query := dialect.From(tbl).
		Select(
			"tracks.*",

			tbl.Col("id").As("history_id"),

			tbl.Col("user_id").As("history_user_id"),
			tbl.Col("track_id").As("history_track_id"),

			tbl.Col("listened_at").As("history_listened_at"),

			tbl.Col("playback_type").As("history_playback_type"),
			tbl.Col("status").As("history_status"),

			tbl.Col("percent_played").As("history_percent_played"),

			tbl.Col("updated").As("history_updated"),
			tbl.Col("created").As("history_created"),
		).
		Join(
			TrackQuery().As("tracks"),
			goqu.On(tbl.Col("track_id").Eq(goqu.I("tracks.id"))),
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

	PercentPlayed int

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

		"percent_played": params.PercentPlayed,

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

func (db DB) GetCompletedPlayCount(
	ctx context.Context, 
	userId string,
) (int, error) {
	tbl := goqu.T("track_history")
	query := dialect.From(tbl).
		Where(
			tbl.Col("user_id").Eq(userId),
			tbl.Col("status").Eq("completed"),
		).
		Select(goqu.COUNT("*"))

	return Single[int](db, ctx, query)
}

func (db DB) GetSkippedPlayCount(
	ctx context.Context, 
	userId string,
) (int, error) {
	tbl := goqu.T("track_history")
	query := dialect.From(tbl).
		Where(
			tbl.Col("user_id").Eq(userId),
			tbl.Col("status").Eq("skipped"),
		).
		Select(goqu.COUNT("*"))

	return Single[int](db, ctx, query)
}

func (db DB) GetCompletedListeningTime(
	ctx context.Context, 
	userId string,
) (int64, error) {
	historyTbl := goqu.T("track_history")
	tracksTbl := goqu.T("tracks")

	query := dialect.From(historyTbl).
		Select(
			goqu.COALESCE(
				goqu.SUM(tracksTbl.Col("duration")), 0,
			).As("listening_time"),
		).
		Join(
			tracksTbl,
			goqu.On(historyTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			historyTbl.Col("user_id").Eq(userId),
			historyTbl.Col("status").Eq("completed"),
		)
	
	listeningTime, err := Single[sql.NullInt64](db, ctx, query)
	if err != nil {
		return 0, err
	}

	return listeningTime.Int64, nil
}

func (db DB) GetLastListenedAt(
	ctx context.Context, 
	userId string,
) (sql.NullInt64, error) {
	tbl := goqu.T("track_history")
	query := dialect.From(tbl).
		Where(tbl.Col("user_id").Eq(userId)).
		Select(goqu.MAX(tbl.Col("listened_at")))

	return Single[sql.NullInt64](db, ctx, query)
}
