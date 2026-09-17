package database

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	userYearReviewsTbl = goqu.T("user_year_reviews")

	userYearReviewTracksTbl  = goqu.T("user_year_review_tracks")
	userYearReviewAlbumsTbl  = goqu.T("user_year_review_albums")
	userYearReviewArtistsTbl = goqu.T("user_year_review_artists")
	userYearReviewMonthsTbl  = goqu.T("user_year_review_months")
	userYearReviewTagsTbl    = goqu.T("user_year_review_tags")
	userYearReviewDecadesTbl = goqu.T("user_year_review_decades")

	userYearReviewMonthTracksTbl  = goqu.T("user_year_review_month_tracks")
	userYearReviewMonthAlbumsTbl  = goqu.T("user_year_review_month_albums")
	userYearReviewMonthArtistsTbl = goqu.T("user_year_review_month_artists")
	userYearReviewMonthTagsTbl    = goqu.T("user_year_review_month_tags")
	userYearReviewMonthDecadesTbl = goqu.T("user_year_review_month_decades")
)

const (
	// TODO(patrik): Rename
	TopYearReviewItems = 5
)

type UserYearReview struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`

	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
	FavoritePlays int     `db:"favorite_plays"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewTag struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	TagSlug string `db:"tag_slug"`
	Rank    int    `db:"rank"`

	PlayCount int `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewDecade struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	Decade int `db:"decade"`
	Rank   int `db:"rank"`

	PlayCount int `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewTrack struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Rank   int    `db:"rank"`

	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewAlbum struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Rank   int    `db:"rank"`

	AlbumId   string `db:"album_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewArtist struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Rank   int    `db:"rank"`

	ArtistId  string `db:"artist_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonth struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`

	PlayCount int   `db:"play_count"`
	PlayTime  int64 `db:"play_time"`

	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
	FavoritePlays int     `db:"favorite_plays"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonthTrack struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`
	Rank   int    `db:"rank"`

	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonthAlbum struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`
	Rank   int    `db:"rank"`

	AlbumId   string `db:"album_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonthArtist struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`
	Rank   int    `db:"rank"`

	ArtistId  string `db:"artist_id"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonthTag struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`

	TagSlug string `db:"tag_slug"`
	Rank    int    `db:"rank"`

	PlayCount int `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewMonthDecade struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`

	Decade int `db:"decade"`
	Rank   int `db:"rank"`

	PlayCount int `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYear struct {
	Year int `db:"year"`
}

type userYearMonth struct {
	Month     int   `db:"month"`
	PlayCount int   `db:"play_count"`
	PlayTime  int64 `db:"play_time"`
}

type UserYearArtistTrack struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
	PlayTime  int64  `db:"play_time"`
}

type UserYearAlbumTrack struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
	PlayTime  int64  `db:"play_time"`
}

type UserYearSummary struct {
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

// TODO(patrik): Move to user_track_stats
func (db DB) GetUserYearSummary(
	ctx context.Context,
	userId string,
	year int,
) (UserYearSummary, error) {
	query := dialect.From(userTrackStatsTbl).
		Select(
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_count")), 0,
			).As("track_count"),
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_time")), 0,
			).As("listening_time"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
		)

	return Single[UserYearSummary](db, ctx, query)
}

func userYearReviewTracksInsertQuery(
	userId string,
	year int,
	month int,
	now int64,
) RawQuery {
	table := "user_year_review_tracks"
	monthCol := ""
	monthSelect := ""
	periodType := "year"
	periodClause := ""
	periodParams := []any{}
	if month != 0 {
		table = "user_year_review_month_tracks"
		monthCol = ", month"
		monthSelect = ", ?"
		periodType = "month"
		periodClause = " AND user_track_stats.period_value = ?"
		periodParams = []any{month}
	}

	q := fmt.Sprintf(`INSERT INTO %s (user_id, year%s, track_id, rank, play_count, created_at, updated_at)
SELECT ?, ?%s, user_track_stats.track_id,
       ROW_NUMBER() OVER (ORDER BY user_track_stats.play_count DESC, user_track_stats.track_id ASC),
       user_track_stats.play_count, ?, ?
FROM user_track_stats
WHERE user_track_stats.user_id = ? AND user_track_stats.period_type = '%s' AND user_track_stats.year = ?%s`,
		table, monthCol, monthSelect, periodType, periodClause)

	params := []any{userId, year}
	params = append(params, periodParams...)
	params = append(params, now, now, userId, year)
	params = append(params, periodParams...)

	return RawQuery{Query: q, Params: params}
}

func userYearReviewAlbumsInsertQuery(
	userId string,
	year int,
	month int,
	now int64,
) RawQuery {
	table := "user_year_review_albums"
	monthCol := ""
	monthSelect := ""
	periodType := "year"
	periodClause := ""
	periodParams := []any{}
	if month != 0 {
		table = "user_year_review_month_albums"
		monthCol = ", month"
		monthSelect = ", ?"
		periodType = "month"
		periodClause = " AND user_track_stats.period_value = ?"
		periodParams = []any{month}
	}

	q := fmt.Sprintf(`INSERT INTO %s (user_id, year%s, album_id, rank, play_count, created_at, updated_at)
SELECT ?, ?%s, albums.id,
       ROW_NUMBER() OVER (ORDER BY SUM(user_track_stats.play_count) DESC, albums.id ASC),
       SUM(user_track_stats.play_count), ?, ?
FROM user_track_stats
JOIN tracks ON tracks.id = user_track_stats.track_id
JOIN albums ON albums.id = tracks.album_id
WHERE user_track_stats.user_id = ? AND user_track_stats.period_type = '%s' AND user_track_stats.year = ?%s
GROUP BY albums.id`,
		table, monthCol, monthSelect, periodType, periodClause)

	params := []any{userId, year}
	params = append(params, periodParams...)
	params = append(params, now, now, userId, year)
	params = append(params, periodParams...)

	return RawQuery{Query: q, Params: params}
}

func userYearReviewTagsInsertQuery(
	userId string,
	year int,
	month int,
	now int64,
) RawQuery {
	table := "user_year_review_tags"
	monthCol := ""
	monthSelect := ""
	periodType := "year"
	periodClause := ""
	periodParams := []any{}
	if month != 0 {
		table = "user_year_review_month_tags"
		monthCol = ", month"
		monthSelect = ", ?"
		periodType = "month"
		periodClause = " AND user_track_stats.period_value = ?"
		periodParams = []any{month}
	}

	q := fmt.Sprintf(`INSERT INTO %s (user_id, year%s, tag_slug, rank, play_count, created_at, updated_at)
SELECT ?, ?%s, tags.slug,
       ROW_NUMBER() OVER (ORDER BY SUM(user_track_stats.play_count) DESC, tags.slug ASC),
       SUM(user_track_stats.play_count), ?, ?
FROM user_track_stats
JOIN tracks_tags ON tracks_tags.track_id = user_track_stats.track_id
JOIN tags ON tags.slug = tracks_tags.tag_slug
WHERE user_track_stats.user_id = ? AND user_track_stats.period_type = '%s' AND user_track_stats.year = ?%s
GROUP BY tags.slug`,
		table, monthCol, monthSelect, periodType, periodClause)

	params := []any{userId, year}
	params = append(params, periodParams...)
	params = append(params, now, now, userId, year)
	params = append(params, periodParams...)

	return RawQuery{Query: q, Params: params}
}

func userYearReviewDecadesInsertQuery(
	userId string,
	year int,
	month int,
	now int64,
) RawQuery {
	table := "user_year_review_decades"
	monthCol := ""
	monthSelect := ""
	periodType := "year"
	periodClause := ""
	periodParams := []any{}
	if month != 0 {
		table = "user_year_review_month_decades"
		monthCol = ", month"
		monthSelect = ", ?"
		periodType = "month"
		periodClause = " AND user_track_stats.period_value = ?"
		periodParams = []any{month}
	}

	q := fmt.Sprintf(`INSERT INTO %s (user_id, year%s, decade, rank, play_count, created_at, updated_at)
SELECT ?, ?%s, tracks.year - tracks.year %% 10,
       ROW_NUMBER() OVER (ORDER BY SUM(user_track_stats.play_count) DESC, tracks.year - tracks.year %% 10 ASC),
       SUM(user_track_stats.play_count), ?, ?
FROM user_track_stats
JOIN tracks ON tracks.id = user_track_stats.track_id
WHERE user_track_stats.user_id = ? AND user_track_stats.period_type = '%s' AND user_track_stats.year = ?%s
  AND tracks.year IS NOT NULL
GROUP BY tracks.year - tracks.year %% 10`,
		table, monthCol, monthSelect, periodType, periodClause)

	params := []any{userId, year}
	params = append(params, periodParams...)
	params = append(params, now, now, userId, year)
	params = append(params, periodParams...)

	return RawQuery{Query: q, Params: params}
}

func userYearReviewArtistsInsertQuery(
	userId string,
	year int,
	month int,
	now int64,
) RawQuery {
	table := "user_year_review_artists"
	monthCol := ""
	monthSelect := ""
	periodType := "year"
	periodClause := ""
	periodParams := []any{}
	if month != 0 {
		table = "user_year_review_month_artists"
		monthCol = ", month"
		monthSelect = ", ?"
		periodType = "month"
		periodClause = " AND user_track_stats.period_value = ?"
		periodParams = []any{month}
	}

	q := fmt.Sprintf(`INSERT INTO %s (user_id, year%s, artist_id, rank, play_count, created_at, updated_at)
SELECT ?, ?%s, artists.id,
       ROW_NUMBER() OVER (ORDER BY SUM(user_track_stats.play_count) DESC, artists.id ASC),
       SUM(user_track_stats.play_count), ?, ?
FROM user_track_stats
JOIN tracks ON tracks.id = user_track_stats.track_id
JOIN artists ON artists.id = tracks.artist_id
WHERE user_track_stats.user_id = ? AND user_track_stats.period_type = '%s' AND user_track_stats.year = ?%s
GROUP BY artists.id`,
		table, monthCol, monthSelect, periodType, periodClause)

	params := []any{userId, year}
	params = append(params, periodParams...)
	params = append(params, now, now, userId, year)
	params = append(params, periodParams...)

	return RawQuery{Query: q, Params: params}
}

func GetUserYearMonthsQuery(userId string, year int) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("period_value").As("month"),
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_count")), 0,
			).As("play_count"),
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_time")), 0,
			).As("play_time"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
		).
		GroupBy(userTrackStatsTbl.Col("period_value")).
		Order(userTrackStatsTbl.Col("period_value").Asc())
}

func GetUserYearArtistTracksQuery(
	userId string,
	year int,
	artistId string,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("track_id"),
			userTrackStatsTbl.Col("play_count"),
			userTrackStatsTbl.Col("play_time"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
			tracksTbl.Col("artist_id").Eq(artistId),
		).
		Order(
			userTrackStatsTbl.Col("play_count").Desc(),
			userTrackStatsTbl.Col("track_id").Asc(),
		)
}

func GetUserYearAlbumTracksQuery(
	userId string,
	year int,
	albumId string,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("track_id"),
			userTrackStatsTbl.Col("play_count"),
			userTrackStatsTbl.Col("play_time"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
			tracksTbl.Col("album_id").Eq(albumId),
		).
		Order(
			userTrackStatsTbl.Col("play_count").Desc(),
			userTrackStatsTbl.Col("track_id").Asc(),
		)
}

func yearRange(year int) (int64, int64) {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).
		UnixMilli()
	end := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).
		UnixMilli()

	return start, end
}

func monthRange(year int, month int) (int64, int64) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).
		UnixMilli()

	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	end := time.Date(nextYear, time.Month(nextMonth), 1, 0, 0, 0, 0, time.Local).
		UnixMilli()

	return start, end
}

type UserYearHistorySummary struct {
	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
}

func (db DB) GetUserYearHistorySummary(
	ctx context.Context,
	userId string,
	year int,
) (UserYearHistorySummary, error) {
	start, end := yearRange(year)

	query := `
SELECT ROUND(COALESCE(AVG(percent_played), 0), 1) AS avg_completion,
       COALESCE(SUM(CASE WHEN status = 'skipped' THEN 1 ELSE 0 END), 0) AS skip_count,
       COALESCE(COUNT(DISTINCT track_id), 0) AS unique_tracks
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?`

	return Single[UserYearHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	})
}

type UserYearFavoritePlays struct {
	PlayCount int `db:"play_count"`
}

func (db DB) GetUserYearFavoritePlays(
	ctx context.Context,
	userId string,
	year int,
) (UserYearFavoritePlays, error) {
	start, end := yearRange(year)

	query := `
SELECT COUNT(*) AS play_count
FROM track_history h
JOIN user_favorites f ON f.user_id = h.user_id AND f.track_id = h.track_id
WHERE h.user_id = ?
  AND h.listened_at >= ?
  AND h.listened_at < ?`

	return Single[UserYearFavoritePlays](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	})
}

func (db DB) GetUserYearMonthSummary(
	ctx context.Context,
	userId string,
	year int,
	month int,
) (UserYearSummary, error) {
	query := dialect.From(userTrackStatsTbl).
		Select(
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_count")), 0,
			).As("track_count"),
			goqu.COALESCE(
				goqu.SUM(userTrackStatsTbl.Col("play_time")), 0,
			).As("listening_time"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
			userTrackStatsTbl.Col("period_value").Eq(month),
		)

	return Single[UserYearSummary](db, ctx, query)
}

func (db DB) GetUserYearMonthHistorySummary(
	ctx context.Context,
	userId string,
	year int,
	month int,
) (UserYearHistorySummary, error) {
	start, end := monthRange(year, month)

	query := `
SELECT ROUND(COALESCE(AVG(percent_played), 0), 1) AS avg_completion,
       COALESCE(SUM(CASE WHEN status = 'skipped' THEN 1 ELSE 0 END), 0) AS skip_count,
       COALESCE(COUNT(DISTINCT track_id), 0) AS unique_tracks
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?`

	return Single[UserYearHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	})
}

func (db DB) GetUserYearMonthFavoritePlays(
	ctx context.Context,
	userId string,
	year int,
	month int,
) (UserYearFavoritePlays, error) {
	start, end := monthRange(year, month)

	query := `
SELECT COUNT(*) AS play_count
FROM track_history h
JOIN user_favorites f ON f.user_id = h.user_id AND f.track_id = h.track_id
WHERE h.user_id = ?
  AND h.listened_at >= ?
  AND h.listened_at < ?`

	return Single[UserYearFavoritePlays](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	})
}

// TODO(patrik): This should be moved to the service package, UserService or ReviewService
func processUserYearMonth(
	tx DB,
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	monthSummary, err := tx.GetUserYearMonthSummary(ctx, userId, year, month)
	if err != nil {
		return err
	}

	monthHistory, err := tx.GetUserYearMonthHistorySummary(ctx, userId, year, month)
	if err != nil {
		return err
	}

	monthFavorites, err := tx.GetUserYearMonthFavoritePlays(ctx, userId, year, month)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthsTbl).Rows(goqu.Record{
		"user_id":        userId,
		"year":           year,
		"month":          month,
		"play_count":     monthSummary.TrackCount,
		"play_time":      monthSummary.ListeningTime,
		"avg_completion": monthHistory.AvgCompletion,
		"skip_count":     monthHistory.SkipCount,
		"unique_tracks":  monthHistory.UniqueTracks,
		"favorite_plays": monthFavorites.PlayCount,

		"created_at": now,
		"updated_at": now,
	}))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewTracksInsertQuery(userId, year, month, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewAlbumsInsertQuery(userId, year, month, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewArtistsInsertQuery(userId, year, month, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewTagsInsertQuery(userId, year, month, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewDecadesInsertQuery(userId, year, month, now))
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GenerateUserReview(
	ctx context.Context,
	userId string,
	year int,
) error {
	now := time.Now().UnixMilli()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clean up any previously generated review for this year.
	_, err = tx.Exec(ctx, dialect.Delete(userYearReviewsTbl).Where(goqu.Ex{
		"user_id": userId,
		"year":    year,
	}))
	if err != nil {
		return err
	}

	// TODO(patrik): Use this if we don't have the foreign keys.
	// for _, tbl := range []any{
	// 	userYearReviewTracksTbl,
	// 	userYearReviewAlbumsTbl,
	// 	userYearReviewArtistsTbl,
	// 	userYearReviewMonthsTbl,
	// 	userYearReviewMonthTracksTbl,
	// 	userYearReviewMonthAlbumsTbl,
	// 	userYearReviewMonthArtistsTbl,
	// 	userYearReviewMonthTagsTbl,
	// 	userYearReviewMonthDecadesTbl,
	// 	userYearReviewTagsTbl,
	// 	userYearReviewDecadesTbl,
	// 	userYearReviewsTbl,
	// } {
	// 	_, err := tx.Exec(ctx, dialect.Delete(tbl).Where(goqu.Ex{
	// 		"user_id": userId,
	// 		"year":    year,
	// 	}))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	summary, err := tx.GetUserYearSummary(ctx, userId, year)
	if err != nil {
		return err
	}

	historySummary, err := tx.GetUserYearHistorySummary(ctx, userId, year)
	if err != nil {
		return err
	}

	favorites, err := tx.GetUserYearFavoritePlays(ctx, userId, year)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, dialect.Insert(userYearReviewsTbl).Rows(goqu.Record{
		"user_id":        userId,
		"year":           year,
		"track_count":    summary.TrackCount,
		"listening_time": summary.ListeningTime,
		"avg_completion": historySummary.AvgCompletion,
		"skip_count":     historySummary.SkipCount,
		"unique_tracks":  historySummary.UniqueTracks,
		"favorite_plays": favorites.PlayCount,

		"created_at": now,
		"updated_at": now,
	}))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewTracksInsertQuery(userId, year, 0, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewAlbumsInsertQuery(userId, year, 0, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewArtistsInsertQuery(userId, year, 0, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewTagsInsertQuery(userId, year, 0, now))
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userYearReviewDecadesInsertQuery(userId, year, 0, now))
	if err != nil {
		return err
	}

	for m := 1; m <= 12; m++ {
		err := processUserYearMonth(tx.DB, ctx, userId, year, m, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db DB) GetUserStatsYears(ctx context.Context, userId string) ([]UserYear, error) {
	query := dialect.From(userTrackStatsTbl).
		SelectDistinct(userTrackStatsTbl.Col("year")).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
		).
		Order(userTrackStatsTbl.Col("year").Desc())

	return Multiple[UserYear](db, ctx, query)
}

var userYearReviewColumns = []any{
	userYearReviewsTbl.Col("user_id"),
	userYearReviewsTbl.Col("year"),
	userYearReviewsTbl.Col("track_count"),
	userYearReviewsTbl.Col("listening_time"),
	userYearReviewsTbl.Col("avg_completion"),
	userYearReviewsTbl.Col("skip_count"),
	userYearReviewsTbl.Col("unique_tracks"),
	userYearReviewsTbl.Col("favorite_plays"),
	userYearReviewsTbl.Col("created_at"),
	userYearReviewsTbl.Col("updated_at"),
}

func (db DB) GetUserYearReviews(
	ctx context.Context,
	userId string,
) ([]UserYearReview, error) {
	query := dialect.From(userYearReviewsTbl).
		Select(userYearReviewColumns...).
		Where(userYearReviewsTbl.Col("user_id").Eq(userId)).
		Order(userYearReviewsTbl.Col("year").Desc())

	return Multiple[UserYearReview](db, ctx, query)
}

func (db DB) GetUserYearReview(
	ctx context.Context,
	userId string,
	year int,
) (UserYearReview, error) {
	query := dialect.From(userYearReviewsTbl).
		Select(userYearReviewColumns...).
		Where(
			userYearReviewsTbl.Col("user_id").Eq(userId),
			userYearReviewsTbl.Col("year").Eq(year),
		)

	return Single[UserYearReview](db, ctx, query)
}

func (db DB) GetUserYearReviewTracks(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewTrack, error) {
	query := dialect.From(userYearReviewTracksTbl).
		Select(
			userYearReviewTracksTbl.Col("user_id"),
			userYearReviewTracksTbl.Col("year"),
			userYearReviewTracksTbl.Col("rank"),
			userYearReviewTracksTbl.Col("track_id"),
			userYearReviewTracksTbl.Col("play_count"),
			userYearReviewTracksTbl.Col("created_at"),
			userYearReviewTracksTbl.Col("updated_at"),
		).
		Where(
			userYearReviewTracksTbl.Col("user_id").Eq(userId),
			userYearReviewTracksTbl.Col("year").Eq(year),
		).
		Order(userYearReviewTracksTbl.Col("rank").Asc())

	return Multiple[UserYearReviewTrack](db, ctx, query)
}

func (db DB) GetUserYearReviewAlbums(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewAlbum, error) {
	query := dialect.From(userYearReviewAlbumsTbl).
		Select(
			userYearReviewAlbumsTbl.Col("user_id"),
			userYearReviewAlbumsTbl.Col("year"),
			userYearReviewAlbumsTbl.Col("rank"),
			userYearReviewAlbumsTbl.Col("album_id"),
			userYearReviewAlbumsTbl.Col("play_count"),
			userYearReviewAlbumsTbl.Col("created_at"),
			userYearReviewAlbumsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewAlbumsTbl.Col("user_id").Eq(userId),
			userYearReviewAlbumsTbl.Col("year").Eq(year),
		).
		Order(userYearReviewAlbumsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewAlbum](db, ctx, query)
}

func (db DB) GetUserYearReviewArtists(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewArtist, error) {
	query := dialect.From(userYearReviewArtistsTbl).
		Select(
			userYearReviewArtistsTbl.Col("user_id"),
			userYearReviewArtistsTbl.Col("year"),
			userYearReviewArtistsTbl.Col("rank"),
			userYearReviewArtistsTbl.Col("artist_id"),
			userYearReviewArtistsTbl.Col("play_count"),
			userYearReviewArtistsTbl.Col("created_at"),
			userYearReviewArtistsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewArtistsTbl.Col("user_id").Eq(userId),
			userYearReviewArtistsTbl.Col("year").Eq(year),
		).
		Order(userYearReviewArtistsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewArtist](db, ctx, query)
}

func (db DB) GetUserYearReviewMonths(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewMonth, error) {
	query := dialect.From(userYearReviewMonthsTbl).
		Select(
			userYearReviewMonthsTbl.Col("user_id"),
			userYearReviewMonthsTbl.Col("year"),
			userYearReviewMonthsTbl.Col("month"),
			userYearReviewMonthsTbl.Col("play_count"),
			userYearReviewMonthsTbl.Col("play_time"),
			userYearReviewMonthsTbl.Col("avg_completion"),
			userYearReviewMonthsTbl.Col("skip_count"),
			userYearReviewMonthsTbl.Col("unique_tracks"),
			userYearReviewMonthsTbl.Col("favorite_plays"),
			userYearReviewMonthsTbl.Col("created_at"),
			userYearReviewMonthsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthsTbl.Col("user_id").Eq(userId),
			userYearReviewMonthsTbl.Col("year").Eq(year),
		).
		Order(userYearReviewMonthsTbl.Col("month").Asc())

	return Multiple[UserYearReviewMonth](db, ctx, query)
}

func (db DB) GetUserYearReviewMonthTracks(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthTrack, error) {
	query := dialect.From(userYearReviewMonthTracksTbl).
		Select(
			userYearReviewMonthTracksTbl.Col("user_id"),
			userYearReviewMonthTracksTbl.Col("year"),
			userYearReviewMonthTracksTbl.Col("month"),
			userYearReviewMonthTracksTbl.Col("rank"),
			userYearReviewMonthTracksTbl.Col("track_id"),
			userYearReviewMonthTracksTbl.Col("play_count"),
			userYearReviewMonthTracksTbl.Col("created_at"),
			userYearReviewMonthTracksTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthTracksTbl.Col("user_id").Eq(userId),
			userYearReviewMonthTracksTbl.Col("year").Eq(year),
			userYearReviewMonthTracksTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthTracksTbl.Col("rank").Asc())

	return Multiple[UserYearReviewMonthTrack](db, ctx, query)
}

func (db DB) GetUserYearReviewMonthAlbums(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthAlbum, error) {
	query := dialect.From(userYearReviewMonthAlbumsTbl).
		Select(
			userYearReviewMonthAlbumsTbl.Col("user_id"),
			userYearReviewMonthAlbumsTbl.Col("year"),
			userYearReviewMonthAlbumsTbl.Col("month"),
			userYearReviewMonthAlbumsTbl.Col("rank"),
			userYearReviewMonthAlbumsTbl.Col("album_id"),
			userYearReviewMonthAlbumsTbl.Col("play_count"),
			userYearReviewMonthAlbumsTbl.Col("created_at"),
			userYearReviewMonthAlbumsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthAlbumsTbl.Col("user_id").Eq(userId),
			userYearReviewMonthAlbumsTbl.Col("year").Eq(year),
			userYearReviewMonthAlbumsTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthAlbumsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewMonthAlbum](db, ctx, query)
}

func (db DB) GetUserYearReviewMonthArtists(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthArtist, error) {
	query := dialect.From(userYearReviewMonthArtistsTbl).
		Select(
			userYearReviewMonthArtistsTbl.Col("user_id"),
			userYearReviewMonthArtistsTbl.Col("year"),
			userYearReviewMonthArtistsTbl.Col("month"),
			userYearReviewMonthArtistsTbl.Col("rank"),
			userYearReviewMonthArtistsTbl.Col("artist_id"),
			userYearReviewMonthArtistsTbl.Col("play_count"),
			userYearReviewMonthArtistsTbl.Col("created_at"),
			userYearReviewMonthArtistsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthArtistsTbl.Col("user_id").Eq(userId),
			userYearReviewMonthArtistsTbl.Col("year").Eq(year),
			userYearReviewMonthArtistsTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthArtistsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewMonthArtist](db, ctx, query)
}

func (db DB) GetUserYearReviewMonthTags(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthTag, error) {
	query := dialect.From(userYearReviewMonthTagsTbl).
		Select(
			userYearReviewMonthTagsTbl.Col("user_id"),
			userYearReviewMonthTagsTbl.Col("year"),
			userYearReviewMonthTagsTbl.Col("month"),
			userYearReviewMonthTagsTbl.Col("tag_slug"),
			userYearReviewMonthTagsTbl.Col("rank"),
			userYearReviewMonthTagsTbl.Col("play_count"),
			userYearReviewMonthTagsTbl.Col("created_at"),
			userYearReviewMonthTagsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthTagsTbl.Col("user_id").Eq(userId),
			userYearReviewMonthTagsTbl.Col("year").Eq(year),
			userYearReviewMonthTagsTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthTagsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewMonthTag](db, ctx, query)
}

func (db DB) GetUserYearReviewMonthDecades(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthDecade, error) {
	query := dialect.From(userYearReviewMonthDecadesTbl).
		Select(
			userYearReviewMonthDecadesTbl.Col("user_id"),
			userYearReviewMonthDecadesTbl.Col("year"),
			userYearReviewMonthDecadesTbl.Col("month"),
			userYearReviewMonthDecadesTbl.Col("decade"),
			userYearReviewMonthDecadesTbl.Col("rank"),
			userYearReviewMonthDecadesTbl.Col("play_count"),
			userYearReviewMonthDecadesTbl.Col("created_at"),
			userYearReviewMonthDecadesTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthDecadesTbl.Col("user_id").Eq(userId),
			userYearReviewMonthDecadesTbl.Col("year").Eq(year),
			userYearReviewMonthDecadesTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthDecadesTbl.Col("rank").Asc())

	return Multiple[UserYearReviewMonthDecade](db, ctx, query)
}

func (db DB) GetUserYearArtistTracks(
	ctx context.Context,
	userId string,
	year int,
	artistId string,
) ([]UserYearArtistTrack, error) {
	return Multiple[UserYearArtistTrack](
		db, ctx, GetUserYearArtistTracksQuery(userId, year, artistId))
}

func (db DB) GetUserYearAlbumTracks(
	ctx context.Context,
	userId string,
	year int,
	albumId string,
) ([]UserYearAlbumTrack, error) {
	return Multiple[UserYearAlbumTrack](
		db, ctx, GetUserYearAlbumTracksQuery(userId, year, albumId))
}

func (db DB) GetUserYearReviewTags(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewTag, error) {
	query := dialect.From(userYearReviewTagsTbl).
		Select(
			userYearReviewTagsTbl.Col("user_id"),
			userYearReviewTagsTbl.Col("year"),
			userYearReviewTagsTbl.Col("tag_slug"),
			userYearReviewTagsTbl.Col("rank"),
			userYearReviewTagsTbl.Col("play_count"),
			userYearReviewTagsTbl.Col("created_at"),
			userYearReviewTagsTbl.Col("updated_at"),
		).
		Where(
			userYearReviewTagsTbl.Col("user_id").Eq(userId),
			userYearReviewTagsTbl.Col("year").Eq(year),
		).
		Order(userYearReviewTagsTbl.Col("rank").Asc())

	return Multiple[UserYearReviewTag](db, ctx, query)
}

func (db DB) GetUserYearReviewDecades(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewDecade, error) {
	query := dialect.From(userYearReviewDecadesTbl).
		Select(
			userYearReviewDecadesTbl.Col("user_id"),
			userYearReviewDecadesTbl.Col("year"),
			userYearReviewDecadesTbl.Col("decade"),
			userYearReviewDecadesTbl.Col("rank"),
			userYearReviewDecadesTbl.Col("play_count"),
			userYearReviewDecadesTbl.Col("created_at"),
			userYearReviewDecadesTbl.Col("updated_at"),
		).
		Where(
			userYearReviewDecadesTbl.Col("user_id").Eq(userId),
			userYearReviewDecadesTbl.Col("year").Eq(year),
		).
		Order(userYearReviewDecadesTbl.Col("rank").Asc())

	return Multiple[UserYearReviewDecade](db, ctx, query)
}
