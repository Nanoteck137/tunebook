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
