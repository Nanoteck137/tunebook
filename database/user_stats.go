package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	userStatsTbl = goqu.T("user_stats")
)

type UserStats struct {
	UserId string `db:"user_id"`

	NumTracksPlayed     int           `db:"num_tracks_played"`
	NumTracksSkipped    int           `db:"num_tracks_skipped"`
	NumPlaylistsCreated int           `db:"num_playlists_created"`
	NumFavoriteTracks   int           `db:"num_favorite_tracks"`
	ListeningTime       int64         `db:"listening_time"`
	LastListenedAt      sql.NullInt64 `db:"last_listened_at"`

	Updated int64 `db:"updated"`
}

func UserStatsQuery() *goqu.SelectDataset {
	return dialect.From(userStatsTbl).
		Select(
			userStatsTbl.Col("user_id"),

			userStatsTbl.Col("num_tracks_played"),
			userStatsTbl.Col("num_tracks_skipped"),
			userStatsTbl.Col("num_playlists_created"),
			userStatsTbl.Col("num_favorite_tracks"),
			userStatsTbl.Col("listening_time"),
			userStatsTbl.Col("last_listened_at"),

			userStatsTbl.Col("updated"),
		)
}

type SetUserStatsParams struct {
	UserId string

	NumTracksPlayed     int
	NumTracksSkipped    int
	NumPlaylistsCreated int
	NumFavoriteTracks   int
	ListeningTime       int64
	LastListenedAt      sql.NullInt64
}

func (db DB) SetUserStats(
	ctx context.Context,
	params SetUserStatsParams,
) error {
	updated := time.Now().UnixMilli()

	query := dialect.Insert(userStatsTbl).Rows(goqu.Record{
		"user_id": params.UserId,

		"num_tracks_played":     params.NumTracksPlayed,
		"num_tracks_skipped":    params.NumTracksSkipped,
		"num_playlists_created": params.NumPlaylistsCreated,
		"num_favorite_tracks":   params.NumFavoriteTracks,
		"listening_time":        params.ListeningTime,
		"last_listened_at":      params.LastListenedAt,

		"updated": updated,
	}).OnConflict(
		goqu.DoUpdate("user_id", goqu.Record{
			"num_tracks_played":     params.NumTracksPlayed,
			"num_tracks_skipped":    params.NumTracksSkipped,
			"num_playlists_created": params.NumPlaylistsCreated,
			"num_favorite_tracks":   params.NumFavoriteTracks,
			"listening_time":        params.ListeningTime,
			"last_listened_at":      params.LastListenedAt,

			"updated": updated,
		}),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type IncrementUserStatsParams struct {
	UserId             string
	SkipDelta          int
	ListeningTimeDelta int64
	LastListenedAt     int64
}

func (db DB) IncrementUserStats(
	ctx context.Context,
	params IncrementUserStatsParams,
) error {
	now := time.Now().UnixMilli()

	query := dialect.Insert(userStatsTbl).Rows(goqu.Record{
		"user_id":               params.UserId,
		"num_tracks_played":     1,
		"num_tracks_skipped":    params.SkipDelta,
		"num_playlists_created": 0,
		"num_favorite_tracks":   0,
		"listening_time":        params.ListeningTimeDelta,
		"last_listened_at":      params.LastListenedAt,
		"updated":               now,
	}).OnConflict(
		goqu.DoUpdate("user_id", goqu.Record{
			"num_tracks_played":  goqu.L("num_tracks_played + 1"),
			"num_tracks_skipped": goqu.L("num_tracks_skipped + ?", params.SkipDelta),
			"listening_time":     goqu.L("listening_time + ?", params.ListeningTimeDelta),
			"last_listened_at":   params.LastListenedAt,
			"updated":            now,
		}),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUserStats(ctx context.Context, userId string) error {
	query := dialect.Delete(userStatsTbl).
		Where(userStatsTbl.Col("user_id").Eq(userId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetUserStats(
	ctx context.Context,
	userId string,
) (UserStats, error) {
	query := UserStatsQuery().
		Where(userStatsTbl.Col("user_id").Eq(userId))

	return Single[UserStats](db, ctx, query)
}
