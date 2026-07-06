package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
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
	tbl := goqu.T("user_track_stats")
	query := dialect.From(tbl).
		Select(
			goqu.COALESCE(goqu.SUM(tbl.Col("play_count")), 0).As("num_tracks_played"),
			goqu.COALESCE(goqu.SUM(tbl.Col("skip_count")), 0).As("num_tracks_skipped"),
			goqu.COALESCE(goqu.SUM(tbl.Col("play_time")), 0).As("play_time"),
		).
		Where(
			tbl.Col("user_id").Eq(userId),
			tbl.Col("period_type").Eq("all"),
		)

	return Single[UserTrackStatsAgg](db, ctx, query)
}

func (db DB) UpsertUserTrackStats(
	ctx context.Context,
	params UpsertUserTrackStatsParams,
) error {
	now := time.Now().UnixMilli()

	query := dialect.Insert("user_track_stats").Rows(goqu.Record{
		"user_id":      params.UserId,
		"track_id":     params.TrackId,

		"period_type":  params.PeriodType,
		"year":         params.Year,
		"period_value": params.PeriodValue,

		"play_count":   1,
		"skip_count":   params.SkipDelta,
		"play_time":    params.PlayTimeDelta,

		"created_at":   now,
		"updated_at":   now,
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
