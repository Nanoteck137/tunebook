package database

import (
	"context"
	"fmt"
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
	Completion  int    `db:"completion_sum"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

type UpsertUserTrackStatsParams struct {
	UserId  string
	TrackId string

	PeriodType  string
	Year        int
	PeriodValue int

	SkipDelta       int
	PlayTimeDelta   int64
	CompletionDelta int
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

		"completion_sum": params.CompletionDelta,

		"created_at": now,
		"updated_at": now,
	}).OnConflict(
		goqu.DoUpdate(
			"user_id, track_id, period_type, year, period_value",
			goqu.Record{
				"play_count": goqu.L("play_count + 1"),
				"skip_count": goqu.L("skip_count + ?", params.SkipDelta),
				"play_time":  goqu.L("play_time + ?", params.PlayTimeDelta),
				"completion_sum": goqu.L(
					"completion_sum + ?", params.CompletionDelta,
				),

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

func (db *Database) RebuildUserTrackStats(
	ctx context.Context,
	userId string,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(ctx, dialect.Delete(userTrackStatsTbl).Where(
		userTrackStatsTbl.Col("user_id").Eq(userId),
	))
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()

	// Rebuild from track_history, grouping each history entry into one of the
	// four period buckets; the year/month/quarter are derived from listened_at
	// in local time to match how the stats were originally populated.
	query := `
INSERT INTO user_track_stats (user_id, track_id, period_type, year, period_value, play_count, skip_count, play_time, completion_sum, created_at, updated_at)
SELECT h.user_id,
       h.track_id,
       '%s',
       %s,
       %s,
       COUNT(*),
       SUM(CASE WHEN h.status = 'skipped' THEN 1 ELSE 0 END),
       CAST(SUM(CAST(t.duration AS REAL) * h.percent_played / 100.0) AS INTEGER),
       SUM(h.percent_played),
       ?,
       ?
FROM track_history h
JOIN tracks t ON t.id = h.track_id
WHERE h.user_id = ?
GROUP BY h.user_id, h.track_id, %s`

	type period struct {
		periodType string
		yearExpr   string
		periodVal  string
		groupExpr  string
	}

	monthExpr := `CAST(strftime('%m', h.listened_at / 1000, 'unixepoch', 'localtime') AS INTEGER)`
	yearExpr := `CAST(strftime('%Y', h.listened_at / 1000, 'unixepoch', 'localtime') AS INTEGER)`

	periods := []period{
		{"all", "0", "0", "h.user_id, h.track_id"},
		{
			"year",
			yearExpr,
			"0",
			`h.user_id, h.track_id, ` + yearExpr,
		},
		{
			"quarter",
			yearExpr,
			`(((` + monthExpr + ` - 1) / 3) + 1)`,
			`h.user_id, h.track_id, ` + yearExpr + `, ((` + monthExpr + ` - 1) / 3) + 1`,
		},
		{
			"month",
			yearExpr,
			monthExpr,
			`h.user_id, h.track_id, ` + yearExpr + `, ` + monthExpr,
		},
	}

	for _, p := range periods {
		_, err = tx.Exec(ctx, RawQuery{
			Query: fmt.Sprintf(query, p.periodType, p.yearExpr, p.periodVal, p.groupExpr),
			Params: []any{
				now, now, userId,
			},
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
