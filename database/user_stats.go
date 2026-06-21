package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
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
	return dialect.From("user_stats").
		Select(
			"user_stats.user_id",

			"user_stats.num_tracks_played",
			"user_stats.num_tracks_skipped",
			"user_stats.num_playlists_created",
			"user_stats.num_favorite_tracks",
			"user_stats.listening_time",
			"user_stats.last_listened_at",

			"user_stats.updated",
		)
}

func (db DB) GetUserStats(
	ctx context.Context,
	userId string,
) (UserStats, error) {
	query := UserStatsQuery().
		Where(goqu.I("user_stats.user_id").Eq(userId))

	return Single[UserStats](db, ctx, query)
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

	query := dialect.Insert("user_stats").Rows(goqu.Record{
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

func (db DB) DeleteUserStats(ctx context.Context, userId string) error {
	query := dialect.Delete("user_stats").
		Where(goqu.I("user_stats.user_id").Eq(userId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
