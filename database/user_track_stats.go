package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	userTrackStatsTbl = goqu.T("user_track_stats")
)

type UserTrackStats struct {
	UserId      string `db:"user_id"`
	TrackId     string `db:"track_id"`
	PeriodType  string `db:"period_type"`
	Year        int    `db:"year"`
	PeriodValue int    `db:"period_value"`
	PlayCount   int    `db:"play_count"`
	SkipCount   int    `db:"skip_count"`
	PlayTime    int64  `db:"play_time"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

type UpsertUserTrackStatsParams struct {
	UserId  string
	TrackId string

	PeriodType  string
	Year        int
	PeriodValue int

	SkipDelta     int
	PlayTimeDelta int64
}

type UserTrackStatsAgg struct {
	NumTracksPlayed  int   `db:"num_tracks_played"`
	NumTracksSkipped int   `db:"num_tracks_skipped"`
	PlayTime         int64 `db:"play_time"`
}

func (db DB) GetUserTrackStatsAgg(
	ctx context.Context,
	userId string,
) (UserTrackStatsAgg, error) {
	query := dialect.From(userTrackStatsTbl).
		Select(
			goqu.COALESCE(goqu.SUM(userTrackStatsTbl.Col("play_count")), 0).
				As("num_tracks_played"),
			goqu.COALESCE(goqu.SUM(userTrackStatsTbl.Col("skip_count")), 0).
				As("num_tracks_skipped"),
			goqu.COALESCE(goqu.SUM(userTrackStatsTbl.Col("play_time")), 0).
				As("play_time"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("all"),
		)

	return Single[UserTrackStatsAgg](db, ctx, query)
}

type UserTopTrack struct {
	Track
	PlayCount int `db:"play_count"`
}

type GetUserTopTracksParams struct {
	UserId     string
	PeriodType string
	Year       int
	Limit      int
}

func (db DB) GetUserTopTracks(
	ctx context.Context,
	params GetUserTopTracksParams,
) ([]UserTopTrack, error) {
	query := TrackQuery().
		SelectAppend(
			userTrackStatsTbl.Col("play_count"),
		).
		Join(
			userTrackStatsTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(params.UserId),
			userTrackStatsTbl.Col("period_type").Eq(params.PeriodType),
		).
		Order(userTrackStatsTbl.Col("play_count").Desc()).
		Limit(uint(params.Limit))

	if params.Year != 0 {
		query = query.Where(userTrackStatsTbl.Col("year").Eq(params.Year))
	}

	return Multiple[UserTopTrack](db, ctx, query)
}

type UserYearStats struct {
	Year          int   `db:"year"`
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

func (db DB) GetUserYearStats(
	ctx context.Context,
	userId string,
) ([]UserYearStats, error) {
	query := dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("year"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("track_count"),
			goqu.SUM(userTrackStatsTbl.Col("play_time")).As("listening_time"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
		).
		GroupBy(userTrackStatsTbl.Col("year")).
		Order(userTrackStatsTbl.Col("year").Desc())

	return Multiple[UserYearStats](db, ctx, query)
}

func (db DB) UpsertUserTrackStats(
	ctx context.Context,
	params UpsertUserTrackStatsParams,
) error {
	now := time.Now().UnixMilli()

	query := dialect.Insert(userTrackStatsTbl).Rows(goqu.Record{
		"user_id":  params.UserId,
		"track_id": params.TrackId,

		"period_type":  params.PeriodType,
		"year":         params.Year,
		"period_value": params.PeriodValue,

		"play_count": 1,
		"skip_count": params.SkipDelta,
		"play_time":  params.PlayTimeDelta,

		"created_at": now,
		"updated_at": now,
	}).OnConflict(
		goqu.DoUpdate(
			"user_id, track_id, period_type, year, period_value",
			goqu.Record{
				"play_count": goqu.L("play_count + 1"),
				"skip_count": goqu.L("skip_count + ?", params.SkipDelta),
				"play_time":  goqu.L("play_time + ?", params.PlayTimeDelta),

				"updated_at": now,
			},
		),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
