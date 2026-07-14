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
	createTrackHistoryId = createIdGenerator(32)

	trackHistoryTbl = goqu.T("track_history")

	trackHistorySchema = TrackHistorySchema()
)

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

func TrackHistorySchema() *schema.Schema {
	return TrackSchema().
		AddField("id", query.TypeString, schema.Column("track_history.id")).
		AddField(
			"userId", 
			query.TypeString, 
			schema.Column("track_history.user_id"),
		).
		AddField(
			"trackId", 
			query.TypeString, 
			schema.Column("track_history.track_id"),
		).
		AddField(
			"listenedAt", 
			query.TypeInt, 
			schema.Column("track_history.listened_at"),
		).
		AddField(
			"playbackType", 
			query.TypeString, 
			schema.Column("track_history.playback_type"),
		).
		AddField(
			"status", 
			query.TypeString, 
			schema.Column("track_history.status"),
		).
		AddField(
			"percentPlayed", 
			query.TypeInt, 
			schema.Column("track_history.percent_played"),
		).
		AddField(
			"created", 
			query.TypeInt, 
			schema.Column("track_history.created"),
		).
		AddField(
			"updated", 
			query.TypeInt, 
			schema.Column("track_history.updated"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "listenedAt"},
				Dir:   query.DirDesc,
			},
		)
}

func TrackHistoryQuery() *goqu.SelectDataset {
	query := TrackQuery().
		SelectAppend(
			trackHistoryTbl.Col("id").As("history_id"),

			trackHistoryTbl.Col("user_id").As("history_user_id"),
			trackHistoryTbl.Col("track_id").As("history_track_id"),

			trackHistoryTbl.Col("listened_at").As("history_listened_at"),

			trackHistoryTbl.Col("playback_type").As("history_playback_type"),
			trackHistoryTbl.Col("status").As("history_status"),

			trackHistoryTbl.Col("percent_played").As("history_percent_played"),

			trackHistoryTbl.Col("updated").As("history_updated"),
			trackHistoryTbl.Col("created").As("history_created"),
		).
		Join(
			trackHistoryTbl,
			goqu.On(trackHistoryTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		)

	return query
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

	query := dialect.Insert(trackHistoryTbl).Rows(goqu.Record{
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

func (db DB) DeleteTrackHistory(
	ctx context.Context,
	trackHistoryId string,
) error {
	query := dialect.Delete(trackHistoryTbl).
		Where(trackHistoryTbl.Col("id").Eq(trackHistoryId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
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

	query = query.Where(trackHistoryTbl.Col("user_id").Eq(params.UserId))

	query, err = ApplyQuery(query, trackHistorySchema, QueryParams{
		Filter: params.Filter.Filter,
		Sort:   params.Filter.Sort,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx, db, params.Page, query, trackHistoryTbl.Col("id"))
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
		Where(trackHistoryTbl.Col("id").Eq(id))

	return Single[TrackHistory](db, ctx, query)
}

func (db DB) GetCompletedPlayCount(
	ctx context.Context,
	userId string,
) (int, error) {
	query := dialect.From(trackHistoryTbl).
		Where(
			trackHistoryTbl.Col("user_id").Eq(userId),
			trackHistoryTbl.Col("status").Eq("completed"),
		).
		Select(goqu.COUNT("*"))

	return Single[int](db, ctx, query)
}

func (db DB) GetSkippedPlayCount(
	ctx context.Context,
	userId string,
) (int, error) {
	query := dialect.From(trackHistoryTbl).
		Where(
			trackHistoryTbl.Col("user_id").Eq(userId),
			trackHistoryTbl.Col("status").Eq("skipped"),
		).
		Select(goqu.COUNT("*"))

	return Single[int](db, ctx, query)
}

func (db DB) GetCompletedListeningTime(
	ctx context.Context,
	userId string,
) (int64, error) {
	query := dialect.From(trackHistoryTbl).
		Select(
			goqu.COALESCE(
				goqu.SUM(tracksTbl.Col("duration")), 0,
			).As("listening_time"),
		).
		Join(
			tracksTbl,
			goqu.On(trackHistoryTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			trackHistoryTbl.Col("user_id").Eq(userId),
			trackHistoryTbl.Col("status").Eq("completed"),
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
	query := dialect.From(trackHistoryTbl).
		Where(trackHistoryTbl.Col("user_id").Eq(userId)).
		Select(goqu.MAX(trackHistoryTbl.Col("listened_at")))

	return Single[sql.NullInt64](db, ctx, query)
}
