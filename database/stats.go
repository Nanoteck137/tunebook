package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type Stats struct {
	Users        int `db:"users"`
	Artists      int `db:"artists"`
	Albums       int `db:"albums"`
	Tracks       int `db:"tracks"`
	Playlists    int `db:"playlists"`
	Favorites    int `db:"favorites"`
	TrackFilters int `db:"track_filters"`
	Queues       int `db:"queues"`

	TotalPlays         int   `db:"total_plays"`
	TotalListeningTime int64 `db:"total_listening_time"`
}

func (db DB) GetStats(ctx context.Context) (Stats, error) {
	query := dialect.From().
		Select(
			goqu.L("(select count(*) from users)").As("users"),
			goqu.L("(select count(*) from artists)").As("artists"),
			goqu.L("(select count(*) from albums)").As("albums"),
			goqu.L("(select count(*) from tracks)").As("tracks"),
			goqu.L("(select count(*) from playlists)").As("playlists"),
			goqu.L("(select count(*) from user_favorites)").As("favorites"),
			goqu.L("(select count(*) from track_filters)").As("track_filters"),
			goqu.L("(select count(*) from queues)").As("queues"),
			goqu.L("(select coalesce(sum(play_count), 0) from user_track_stats where period_type = 'all')").As("total_plays"),
			goqu.L("(select coalesce(sum(play_time), 0) from user_track_stats where period_type = 'all')").As("total_listening_time"),
		)

	return Single[Stats](db, ctx, query)
}