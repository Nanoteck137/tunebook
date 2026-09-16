package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	userYearReviewsTbl            = goqu.T("user_year_reviews")
	userYearReviewTracksTbl       = goqu.T("user_year_review_tracks")
	userYearReviewAlbumsTbl       = goqu.T("user_year_review_albums")
	userYearReviewArtistsTbl      = goqu.T("user_year_review_artists")
	userYearReviewMonthsTbl       = goqu.T("user_year_review_months")
	userYearReviewDayTbl          = goqu.T("user_year_review_day")
	userYearReviewArtistTracksTbl = goqu.T("user_year_review_artist_tracks")
	userYearReviewAlbumTracksTbl  = goqu.T("user_year_review_album_tracks")
	userYearReviewHoursTbl        = goqu.T("user_year_review_hours")
	userYearReviewTagsTbl         = goqu.T("user_year_review_tags")
	userYearReviewDecadesTbl      = goqu.T("user_year_review_decades")
	userYearReviewMilestonesTbl   = goqu.T("user_year_review_milestones")
	userYearReviewMonthTracksTbl  = goqu.T("user_year_review_month_tracks")
	userYearReviewMonthAlbumsTbl  = goqu.T("user_year_review_month_albums")
	userYearReviewMonthArtistsTbl = goqu.T("user_year_review_month_artists")
	userYearReviewMonthHoursTbl   = goqu.T("user_year_review_month_hours")
	userYearReviewMonthTagsTbl    = goqu.T("user_year_review_month_tags")
	userYearReviewMonthDecadesTbl = goqu.T("user_year_review_month_decades")
)

const (
	TopYearReviewItems = 5
	TopYearTags        = 6
)

type UserYearReview struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`

	DaysActive        int     `db:"days_active"`
	LongestStreak     int     `db:"longest_streak"`
	AvgCompletion     float64 `db:"avg_completion"`
	SkipCount         int     `db:"skip_count"`
	UniqueTracks      int     `db:"unique_tracks"`
	FavoritePlays     int     `db:"favorite_plays"`
	PrevTrackCount    int     `db:"prev_track_count"`
	PrevListeningTime int64   `db:"prev_listening_time"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewHour struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Hour   int    `db:"hour"`

	PlayCount int `db:"play_count"`

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

type UserYearReviewMilestone struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	FirstTrackId string `db:"first_track_id"`
	LastTrackId  string `db:"last_track_id"`

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

	DaysActive    int     `db:"days_active"`
	LongestStreak int     `db:"longest_streak"`
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

type UserYearReviewMonthHour struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`
	Hour   int    `db:"hour"`

	PlayCount int `db:"play_count"`

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

type UserYearReviewDay struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	Day       string `db:"day"`
	PlayCount int    `db:"play_count"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewArtistTrack struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	ArtistId string `db:"artist_id"`
	TrackId  string `db:"track_id"`
	Rank     int    `db:"rank"`

	PlayCount int   `db:"play_count"`
	PlayTime  int64 `db:"play_time"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYearReviewAlbumTrack struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	AlbumId string `db:"album_id"`
	TrackId string `db:"track_id"`
	Rank    int    `db:"rank"`

	PlayCount int   `db:"play_count"`
	PlayTime  int64 `db:"play_time"`

	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

type UserYear struct {
	Year int `db:"year"`
}

type userYearSummary struct {
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

type userYearTopTrack struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
}

type userYearTopAlbum struct {
	AlbumId   string `db:"album_id"`
	PlayCount int    `db:"play_count"`
}

type userYearTopArtist struct {
	ArtistId  string `db:"artist_id"`
	PlayCount int    `db:"play_count"`
}

type userYearMonth struct {
	Month     int   `db:"month"`
	PlayCount int   `db:"play_count"`
	PlayTime  int64 `db:"play_time"`
}

type userYearMostPlayedDay struct {
	DayOfWeek int `db:"day_of_week"`
	PlayCount int `db:"play_count"`
}

type userYearArtistTrack struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
	PlayTime  int64  `db:"play_time"`
}

type userYearAlbumTrack struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
	PlayTime  int64  `db:"play_time"`
}

type userYearHistorySummary struct {
	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
}

type userYearActiveDay struct {
	Day string `db:"day"`
}

type userYearHour struct {
	Hour      int `db:"hour"`
	PlayCount int `db:"play_count"`
}

type userYearDecade struct {
	Decade    int `db:"decade"`
	PlayCount int `db:"play_count"`
}

type userYearTag struct {
	TagSlug   string `db:"tag_slug"`
	PlayCount int    `db:"play_count"`
}

type userYearFavoritePlays struct {
	PlayCount int `db:"play_count"`
}

type userYearPrevSummary struct {
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

var yearDays = []string{
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}

func longestRun(days []userYearActiveDay) int {
	best := 0
	run := 0
	var prev time.Time

	for _, d := range days {
		t, err := time.Parse("2006-01-02", d.Day)
		if err != nil {
			run = 0
			continue
		}

		if !prev.IsZero() && t.Sub(prev).Hours() == 24 {
			run++
		} else {
			run = 1
		}

		if run > best {
			best = run
		}

		prev = t
	}

	return best
}

func GetUserYearSummaryQuery(userId string, year int) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
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
}

func GetUserYearTopTracksQuery(
	userId string,
	year int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("track_id"),
			userTrackStatsTbl.Col("play_count"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
		).
		Order(
			userTrackStatsTbl.Col("play_count").Desc(),
			userTrackStatsTbl.Col("track_id").Asc(),
		)
}

func GetUserYearTopAlbumsQuery(
	userId string,
	year int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			albumsTbl.Col("id").As("album_id"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Join(
			albumsTbl,
			goqu.On(tracksTbl.Col("album_id").Eq(albumsTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
		).
		GroupBy(albumsTbl.Col("id")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			albumsTbl.Col("id").Asc(),
		)
}

func GetUserYearTopArtistsQuery(
	userId string,
	year int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			artistsTbl.Col("id").As("artist_id"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Join(
			artistsTbl,
			goqu.On(tracksTbl.Col("artist_id").Eq(artistsTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
		).
		GroupBy(artistsTbl.Col("id")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			artistsTbl.Col("id").Asc(),
		)
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

func GetUserYearMostPlayedDayQuery(
	userId string,
	year int,
) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT ((strftime('%w', datetime(listened_at / 1000, 'unixepoch', 'localtime')) + 6) % 7) AS day_of_week,
       COUNT(*) AS play_count
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
GROUP BY day_of_week
ORDER BY play_count DESC
LIMIT 1`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
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

func GetUserYearHistorySummaryQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT ROUND(COALESCE(AVG(percent_played), 0), 1) AS avg_completion,
       COALESCE(SUM(CASE WHEN status = 'skipped' THEN 1 ELSE 0 END), 0) AS skip_count,
       COALESCE(COUNT(DISTINCT track_id), 0) AS unique_tracks
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearActiveDaysQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT DISTINCT date(datetime(listened_at / 1000, 'unixepoch', 'localtime')) AS day
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
ORDER BY day ASC`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearFirstTrackQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT track_id
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
ORDER BY listened_at ASC
LIMIT 1`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearLastTrackQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT track_id
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
ORDER BY listened_at DESC
LIMIT 1`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearHoursQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT CAST(strftime('%H', datetime(listened_at / 1000, 'unixepoch', 'localtime')) AS INTEGER) AS hour,
       COUNT(*) AS play_count
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
GROUP BY hour`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearFavoritePlaysQuery(userId string, year int) RawQuery {
	start, end := yearRange(year)

	query := `
SELECT COUNT(*) AS play_count
FROM track_history h
JOIN user_favorites f ON f.user_id = h.user_id AND f.track_id = h.track_id
WHERE h.user_id = ?
  AND h.listened_at >= ?
  AND h.listened_at < ?`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearPrevSummaryQuery(userId string, year int) RawQuery {
	query := `
SELECT COALESCE(SUM(play_count), 0) AS track_count,
       COALESCE(SUM(play_time), 0) AS listening_time
FROM user_track_stats
WHERE user_id = ?
  AND period_type = 'year'
  AND year = ?`

	return RawQuery{
		Query:  query,
		Params: []any{userId, year - 1},
	}
}

func GetUserYearDecadesQuery(userId string, year int) RawQuery {
	query := `
SELECT (tracks.year - tracks.year % 10) AS decade,
       SUM(uts.play_count) AS play_count
FROM user_track_stats uts
JOIN tracks ON tracks.id = uts.track_id
WHERE uts.user_id = ?
  AND uts.period_type = 'year'
  AND uts.year = ?
  AND tracks.year IS NOT NULL
GROUP BY decade
ORDER BY play_count DESC, decade ASC`

	return RawQuery{
		Query:  query,
		Params: []any{userId, year},
	}
}

func GetUserYearTagsQuery(userId string, year int, limit int) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			tagsTbl.Col("slug").As("tag_slug"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTagsTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTagsTbl.Col("track_id"))),
		).
		Join(
			tagsTbl,
			goqu.On(tracksTagsTbl.Col("tag_slug").Eq(tagsTbl.Col("slug"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
			userTrackStatsTbl.Col("year").Eq(year),
		).
		GroupBy(tagsTbl.Col("slug")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			tagsTbl.Col("slug").Asc(),
		).
		Limit(uint(limit))
}

func GetUserYearMonthSummaryQuery(
	userId string,
	year int,
	month int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
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
}

func GetUserYearMonthTopTracksQuery(
	userId string,
	year int,
	month int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			userTrackStatsTbl.Col("track_id"),
			userTrackStatsTbl.Col("play_count"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
			userTrackStatsTbl.Col("period_value").Eq(month),
		).
		Order(
			userTrackStatsTbl.Col("play_count").Desc(),
			userTrackStatsTbl.Col("track_id").Asc(),
		)
}

func GetUserYearMonthTopAlbumsQuery(
	userId string,
	year int,
	month int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			albumsTbl.Col("id").As("album_id"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Join(
			albumsTbl,
			goqu.On(tracksTbl.Col("album_id").Eq(albumsTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
			userTrackStatsTbl.Col("period_value").Eq(month),
		).
		GroupBy(albumsTbl.Col("id")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			albumsTbl.Col("id").Asc(),
		)
}

func GetUserYearMonthTopArtistsQuery(
	userId string,
	year int,
	month int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			artistsTbl.Col("id").As("artist_id"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Join(
			artistsTbl,
			goqu.On(tracksTbl.Col("artist_id").Eq(artistsTbl.Col("id"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
			userTrackStatsTbl.Col("period_value").Eq(month),
		).
		GroupBy(artistsTbl.Col("id")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			artistsTbl.Col("id").Asc(),
		)
}

func GetUserYearMonthHistorySummaryQuery(
	userId string,
	year int,
	month int,
) RawQuery {
	start, end := monthRange(year, month)

	query := `
SELECT ROUND(COALESCE(AVG(percent_played), 0), 1) AS avg_completion,
       COALESCE(SUM(CASE WHEN status = 'skipped' THEN 1 ELSE 0 END), 0) AS skip_count,
       COALESCE(COUNT(DISTINCT track_id), 0) AS unique_tracks
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearMonthActiveDaysQuery(
	userId string,
	year int,
	month int,
) RawQuery {
	start, end := monthRange(year, month)

	query := `
SELECT DISTINCT date(datetime(listened_at / 1000, 'unixepoch', 'localtime')) AS day
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
ORDER BY day ASC`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearMonthFavoritePlaysQuery(
	userId string,
	year int,
	month int,
) RawQuery {
	start, end := monthRange(year, month)

	query := `
SELECT COUNT(*) AS play_count
FROM track_history h
JOIN user_favorites f ON f.user_id = h.user_id AND f.track_id = h.track_id
WHERE h.user_id = ?
  AND h.listened_at >= ?
  AND h.listened_at < ?`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearMonthHoursQuery(
	userId string,
	year int,
	month int,
) RawQuery {
	start, end := monthRange(year, month)

	query := `
SELECT CAST(strftime('%H', datetime(listened_at / 1000, 'unixepoch', 'localtime')) AS INTEGER) AS hour,
       COUNT(*) AS play_count
FROM track_history
WHERE user_id = ?
  AND listened_at >= ?
  AND listened_at < ?
GROUP BY hour`

	return RawQuery{
		Query:  query,
		Params: []any{userId, start, end},
	}
}

func GetUserYearMonthTagsQuery(
	userId string,
	year int,
	month int,
	limit int,
) *goqu.SelectDataset {
	return dialect.From(userTrackStatsTbl).
		Select(
			tagsTbl.Col("slug").As("tag_slug"),
			goqu.SUM(userTrackStatsTbl.Col("play_count")).As("play_count"),
		).
		Join(
			tracksTagsTbl,
			goqu.On(userTrackStatsTbl.Col("track_id").Eq(tracksTagsTbl.Col("track_id"))),
		).
		Join(
			tagsTbl,
			goqu.On(tracksTagsTbl.Col("tag_slug").Eq(tagsTbl.Col("slug"))),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(userId),
			userTrackStatsTbl.Col("period_type").Eq("month"),
			userTrackStatsTbl.Col("year").Eq(year),
			userTrackStatsTbl.Col("period_value").Eq(month),
		).
		GroupBy(tagsTbl.Col("slug")).
		Order(
			goqu.SUM(userTrackStatsTbl.Col("play_count")).Desc(),
			tagsTbl.Col("slug").Asc(),
		).
		Limit(uint(limit))
}

func GetUserYearMonthDecadesQuery(
	userId string,
	year int,
	month int,
) RawQuery {
	query := `
SELECT (tracks.year - tracks.year % 10) AS decade,
       SUM(uts.play_count) AS play_count
FROM user_track_stats uts
JOIN tracks ON tracks.id = uts.track_id
WHERE uts.user_id = ?
  AND uts.period_type = 'month'
  AND uts.year = ?
  AND uts.period_value = ?
  AND tracks.year IS NOT NULL
GROUP BY decade
ORDER BY play_count DESC, decade ASC`

	return RawQuery{
		Query:  query,
		Params: []any{userId, year, month},
	}
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

	summary, err := Single[userYearSummary](tx, ctx, GetUserYearSummaryQuery(userId, year))
	if err != nil {
		return err
	}

	topTracks, err := Multiple[userYearTopTrack](
		tx, ctx, GetUserYearTopTracksQuery(userId, year))
	if err != nil {
		return err
	}

	topAlbums, err := Multiple[userYearTopAlbum](
		tx, ctx, GetUserYearTopAlbumsQuery(userId, year))
	if err != nil {
		return err
	}

	topArtists, err := Multiple[userYearTopArtist](
		tx, ctx, GetUserYearTopArtistsQuery(userId, year))
	if err != nil {
		return err
	}

	day, err := Single[userYearMostPlayedDay](
		tx, ctx, GetUserYearMostPlayedDayQuery(userId, year))
	if err != nil {
		if !errors.Is(err, ErrItemNotFound) {
			return err
		}

		day = userYearMostPlayedDay{}
	}

	historySummary, err := Single[userYearHistorySummary](
		tx, ctx, GetUserYearHistorySummaryQuery(userId, year))
	if err != nil {
		return err
	}

	activeDays, err := Multiple[userYearActiveDay](
		tx, ctx, GetUserYearActiveDaysQuery(userId, year))
	if err != nil {
		return err
	}

	hours, err := Multiple[userYearHour](
		tx, ctx, GetUserYearHoursQuery(userId, year))
	if err != nil {
		return err
	}

	favorites, err := Single[userYearFavoritePlays](
		tx, ctx, GetUserYearFavoritePlaysQuery(userId, year))
	if err != nil {
		return err
	}

	firstTrackId, err := Single[string](
		tx, ctx, GetUserYearFirstTrackQuery(userId, year))
	if err != nil {
		if !errors.Is(err, ErrItemNotFound) {
			return err
		}

		firstTrackId = ""
	}

	lastTrackId, err := Single[string](
		tx, ctx, GetUserYearLastTrackQuery(userId, year))
	if err != nil {
		if !errors.Is(err, ErrItemNotFound) {
			return err
		}

		lastTrackId = ""
	}

	tags, err := Multiple[userYearTag](
		tx, ctx, GetUserYearTagsQuery(userId, year, TopYearTags))
	if err != nil {
		return err
	}

	decades, err := Multiple[userYearDecade](
		tx, ctx, GetUserYearDecadesQuery(userId, year))
	if err != nil {
		return err
	}

	prevSummary, err := Single[userYearPrevSummary](
		tx, ctx, GetUserYearPrevSummaryQuery(userId, year))
	if err != nil {
		return err
	}

	daysActive := len(activeDays)
	longestStreak := longestRun(activeDays)

	// Clean up any previously generated review for this year.
	for _, tbl := range []any{
		userYearReviewArtistTracksTbl,
		userYearReviewAlbumTracksTbl,
		userYearReviewTracksTbl,
		userYearReviewAlbumsTbl,
		userYearReviewArtistsTbl,
		userYearReviewMonthsTbl,
		userYearReviewMonthTracksTbl,
		userYearReviewMonthAlbumsTbl,
		userYearReviewMonthArtistsTbl,
		userYearReviewMonthHoursTbl,
		userYearReviewMonthTagsTbl,
		userYearReviewMonthDecadesTbl,
		userYearReviewDayTbl,
		userYearReviewHoursTbl,
		userYearReviewTagsTbl,
		userYearReviewDecadesTbl,
		userYearReviewMilestonesTbl,
		userYearReviewsTbl,
	} {
		_, err := tx.Exec(ctx, dialect.Delete(tbl).Where(goqu.Ex{
			"user_id": userId,
			"year":    year,
		}))
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, dialect.Insert(userYearReviewsTbl).Rows(goqu.Record{
		"user_id":             userId,
		"year":                year,
		"track_count":         summary.TrackCount,
		"listening_time":      summary.ListeningTime,
		"days_active":         daysActive,
		"longest_streak":      longestStreak,
		"avg_completion":      historySummary.AvgCompletion,
		"skip_count":          historySummary.SkipCount,
		"unique_tracks":       historySummary.UniqueTracks,
		"favorite_plays":      favorites.PlayCount,
		"prev_track_count":    prevSummary.TrackCount,
		"prev_listening_time": prevSummary.ListeningTime,

		"created_at": now,
		"updated_at": now,
	}))
	if err != nil {
		return err
	}

	trackRows := make([]goqu.Record, len(topTracks))
	for i, t := range topTracks {
		trackRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"track_id":   t.TrackId,
			"rank":       i + 1,
			"play_count": t.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(trackRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewTracksTbl).Rows(trackRows))
		if err != nil {
			return err
		}
	}

	albumRows := make([]goqu.Record, len(topAlbums))
	for i, a := range topAlbums {
		albumRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"album_id":   a.AlbumId,
			"rank":       i + 1,
			"play_count": a.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(albumRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewAlbumsTbl).Rows(albumRows))
		if err != nil {
			return err
		}
	}

	artistRows := make([]goqu.Record, len(topArtists))
	for i, a := range topArtists {
		artistRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"artist_id":  a.ArtistId,
			"rank":       i + 1,
			"play_count": a.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(artistRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewArtistsTbl).Rows(artistRows))
		if err != nil {
			return err
		}
	}

	monthRows := make([]goqu.Record, 12)
	monthPlayCount := make([]int, 12)
	monthPlayTime := make([]int64, 12)
	monthDaysActive := make([]int, 12)
	monthLongestStreak := make([]int, 12)
	monthAvgCompletion := make([]float64, 12)
	monthSkipCount := make([]int, 12)
	monthUniqueTracks := make([]int, 12)
	monthFavoritePlays := make([]int, 12)

	monthTrackRows := make([][]goqu.Record, 12)
	monthAlbumRows := make([][]goqu.Record, 12)
	monthArtistRows := make([][]goqu.Record, 12)
	monthHourRows := make([][]goqu.Record, 12)
	monthTagRows := make([][]goqu.Record, 12)
	monthDecadeRows := make([][]goqu.Record, 12)

	for m := 1; m <= 12; m++ {
		monthSummary, err := Single[userYearSummary](
			tx, ctx, GetUserYearMonthSummaryQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthTopTracks, err := Multiple[userYearTopTrack](
			tx, ctx, GetUserYearMonthTopTracksQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthTopAlbums, err := Multiple[userYearTopAlbum](
			tx, ctx, GetUserYearMonthTopAlbumsQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthTopArtists, err := Multiple[userYearTopArtist](
			tx, ctx, GetUserYearMonthTopArtistsQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthHistory, err := Single[userYearHistorySummary](
			tx, ctx, GetUserYearMonthHistorySummaryQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthActiveDays, err := Multiple[userYearActiveDay](
			tx, ctx, GetUserYearMonthActiveDaysQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthFavorites, err := Single[userYearFavoritePlays](
			tx, ctx, GetUserYearMonthFavoritePlaysQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthHours, err := Multiple[userYearHour](
			tx, ctx, GetUserYearMonthHoursQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthTags, err := Multiple[userYearTag](
			tx, ctx, GetUserYearMonthTagsQuery(userId, year, m, TopYearTags))
		if err != nil {
			return err
		}

		monthDecades, err := Multiple[userYearDecade](
			tx, ctx, GetUserYearMonthDecadesQuery(userId, year, m))
		if err != nil {
			return err
		}

		monthPlayCount[m-1] = monthSummary.TrackCount
		monthPlayTime[m-1] = monthSummary.ListeningTime
		monthDaysActive[m-1] = len(monthActiveDays)
		monthLongestStreak[m-1] = longestRun(monthActiveDays)
		monthAvgCompletion[m-1] = monthHistory.AvgCompletion
		monthSkipCount[m-1] = monthHistory.SkipCount
		monthUniqueTracks[m-1] = monthHistory.UniqueTracks
		monthFavoritePlays[m-1] = monthFavorites.PlayCount

		trackRows := make([]goqu.Record, len(monthTopTracks))
		for i, t := range monthTopTracks {
			trackRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"track_id":   t.TrackId,
				"rank":       i + 1,
				"play_count": t.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthTrackRows[m-1] = trackRows

		albumRows := make([]goqu.Record, len(monthTopAlbums))
		for i, a := range monthTopAlbums {
			albumRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"album_id":   a.AlbumId,
				"rank":       i + 1,
				"play_count": a.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthAlbumRows[m-1] = albumRows

		artistRows := make([]goqu.Record, len(monthTopArtists))
		for i, a := range monthTopArtists {
			artistRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"artist_id":  a.ArtistId,
				"rank":       i + 1,
				"play_count": a.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthArtistRows[m-1] = artistRows

		hourRows := make([]goqu.Record, len(monthHours))
		for i, h := range monthHours {
			if h.Hour < 0 || h.Hour > 23 {
				return fmt.Errorf("invalid hour: %d", h.Hour)
			}

			hourRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"hour":       h.Hour,
				"play_count": h.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthHourRows[m-1] = hourRows

		tagRows := make([]goqu.Record, len(monthTags))
		for i, t := range monthTags {
			tagRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"tag_slug":   t.TagSlug,
				"rank":       i + 1,
				"play_count": t.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthTagRows[m-1] = tagRows

		decadeRows := make([]goqu.Record, len(monthDecades))
		for i, d := range monthDecades {
			decadeRows[i] = goqu.Record{
				"user_id":    userId,
				"year":       year,
				"month":      m,
				"decade":     d.Decade,
				"rank":       i + 1,
				"play_count": d.PlayCount,

				"created_at": now,
				"updated_at": now,
			}
		}
		monthDecadeRows[m-1] = decadeRows
	}

	for i, playTime := range monthPlayTime {
		monthRows[i] = goqu.Record{
			"user_id":        userId,
			"year":           year,
			"month":          i + 1,
			"play_count":     monthPlayCount[i],
			"play_time":      playTime,
			"days_active":    monthDaysActive[i],
			"longest_streak": monthLongestStreak[i],
			"avg_completion": monthAvgCompletion[i],
			"skip_count":     monthSkipCount[i],
			"unique_tracks":  monthUniqueTracks[i],
			"favorite_plays": monthFavoritePlays[i],

			"created_at": now,
			"updated_at": now,
		}
	}
	_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthsTbl).Rows(monthRows))
	if err != nil {
		return err
	}

	for _, rows := range monthTrackRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthTracksTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	for _, rows := range monthAlbumRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthAlbumsTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	for _, rows := range monthArtistRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthArtistsTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	for _, rows := range monthHourRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthHoursTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	for _, rows := range monthTagRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthTagsTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	for _, rows := range monthDecadeRows {
		if len(rows) == 0 {
			continue
		}

		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMonthDecadesTbl).Rows(rows))
		if err != nil {
			return err
		}
	}

	if day.PlayCount > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewDayTbl).Rows(goqu.Record{
			"user_id":    userId,
			"year":       year,
			"day":        yearDays[day.DayOfWeek],
			"play_count": day.PlayCount,

			"created_at": now,
			"updated_at": now,
		}))
		if err != nil {
			return err
		}
	}

	artistTrackRows := make([]goqu.Record, 0, TopYearReviewItems)
	for _, a := range topArtists[:min(len(topArtists), TopYearReviewItems)] {
		tracks, err := Multiple[userYearArtistTrack](
			tx, ctx, GetUserYearArtistTracksQuery(userId, year, a.ArtistId))
		if err != nil {
			return err
		}

		for i, t := range tracks {
			artistTrackRows = append(artistTrackRows, goqu.Record{
				"user_id":    userId,
				"year":       year,
				"artist_id":  a.ArtistId,
				"track_id":   t.TrackId,
				"rank":       i + 1,
				"play_count": t.PlayCount,
				"play_time":  t.PlayTime,

				"created_at": now,
				"updated_at": now,
			})
		}
	}
	if len(artistTrackRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewArtistTracksTbl).Rows(artistTrackRows))
		if err != nil {
			return err
		}
	}

	albumTrackRows := make([]goqu.Record, 0, TopYearReviewItems)
	for _, a := range topAlbums[:min(len(topAlbums), TopYearReviewItems)] {
		tracks, err := Multiple[userYearAlbumTrack](
			tx, ctx, GetUserYearAlbumTracksQuery(userId, year, a.AlbumId))
		if err != nil {
			return err
		}

		for i, t := range tracks {
			albumTrackRows = append(albumTrackRows, goqu.Record{
				"user_id":    userId,
				"year":       year,
				"album_id":   a.AlbumId,
				"track_id":   t.TrackId,
				"rank":       i + 1,
				"play_count": t.PlayCount,
				"play_time":  t.PlayTime,

				"created_at": now,
				"updated_at": now,
			})
		}
	}
	if len(albumTrackRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewAlbumTracksTbl).Rows(albumTrackRows))
		if err != nil {
			return err
		}
	}

	hourRows := make([]goqu.Record, len(hours))
	for i, h := range hours {
		if h.Hour < 0 || h.Hour > 23 {
			return fmt.Errorf("invalid hour: %d", h.Hour)
		}

		hourRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"hour":       h.Hour,
			"play_count": h.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(hourRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewHoursTbl).Rows(hourRows))
		if err != nil {
			return err
		}
	}

	tagRows := make([]goqu.Record, len(tags))
	for i, t := range tags {
		tagRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"tag_slug":   t.TagSlug,
			"rank":       i + 1,
			"play_count": t.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(tagRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewTagsTbl).Rows(tagRows))
		if err != nil {
			return err
		}
	}

	decadeRows := make([]goqu.Record, len(decades))
	for i, d := range decades {
		decadeRows[i] = goqu.Record{
			"user_id":    userId,
			"year":       year,
			"decade":     d.Decade,
			"rank":       i + 1,
			"play_count": d.PlayCount,

			"created_at": now,
			"updated_at": now,
		}
	}
	if len(decadeRows) > 0 {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewDecadesTbl).Rows(decadeRows))
		if err != nil {
			return err
		}
	}

	if firstTrackId != "" || lastTrackId != "" {
		_, err = tx.Exec(ctx, dialect.Insert(userYearReviewMilestonesTbl).Rows(goqu.Record{
			"user_id":        userId,
			"year":           year,
			"first_track_id": firstTrackId,
			"last_track_id":  lastTrackId,

			"created_at": now,
			"updated_at": now,
		}))
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
	userYearReviewsTbl.Col("days_active"),
	userYearReviewsTbl.Col("longest_streak"),
	userYearReviewsTbl.Col("avg_completion"),
	userYearReviewsTbl.Col("skip_count"),
	userYearReviewsTbl.Col("unique_tracks"),
	userYearReviewsTbl.Col("favorite_plays"),
	userYearReviewsTbl.Col("prev_track_count"),
	userYearReviewsTbl.Col("prev_listening_time"),
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
			userYearReviewMonthsTbl.Col("days_active"),
			userYearReviewMonthsTbl.Col("longest_streak"),
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

func (db DB) GetUserYearReviewMonthHours(
	ctx context.Context,
	userId string,
	year int,
	month int,
) ([]UserYearReviewMonthHour, error) {
	query := dialect.From(userYearReviewMonthHoursTbl).
		Select(
			userYearReviewMonthHoursTbl.Col("user_id"),
			userYearReviewMonthHoursTbl.Col("year"),
			userYearReviewMonthHoursTbl.Col("month"),
			userYearReviewMonthHoursTbl.Col("hour"),
			userYearReviewMonthHoursTbl.Col("play_count"),
			userYearReviewMonthHoursTbl.Col("created_at"),
			userYearReviewMonthHoursTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMonthHoursTbl.Col("user_id").Eq(userId),
			userYearReviewMonthHoursTbl.Col("year").Eq(year),
			userYearReviewMonthHoursTbl.Col("month").Eq(month),
		).
		Order(userYearReviewMonthHoursTbl.Col("hour").Asc())

	return Multiple[UserYearReviewMonthHour](db, ctx, query)
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

func (db DB) GetUserYearReviewArtistTracks(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewArtistTrack, error) {
	query := dialect.From(userYearReviewArtistTracksTbl).
		Select(
			userYearReviewArtistTracksTbl.Col("user_id"),
			userYearReviewArtistTracksTbl.Col("year"),
			userYearReviewArtistTracksTbl.Col("artist_id"),
			userYearReviewArtistTracksTbl.Col("track_id"),
			userYearReviewArtistTracksTbl.Col("rank"),
			userYearReviewArtistTracksTbl.Col("play_count"),
			userYearReviewArtistTracksTbl.Col("play_time"),
			userYearReviewArtistTracksTbl.Col("created_at"),
			userYearReviewArtistTracksTbl.Col("updated_at"),
		).
		Where(
			userYearReviewArtistTracksTbl.Col("user_id").Eq(userId),
			userYearReviewArtistTracksTbl.Col("year").Eq(year),
		).
		Order(
			userYearReviewArtistTracksTbl.Col("artist_id").Asc(),
			userYearReviewArtistTracksTbl.Col("rank").Asc(),
		)

	return Multiple[UserYearReviewArtistTrack](db, ctx, query)
}

func (db DB) GetUserYearReviewAlbumTracks(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewAlbumTrack, error) {
	query := dialect.From(userYearReviewAlbumTracksTbl).
		Select(
			userYearReviewAlbumTracksTbl.Col("user_id"),
			userYearReviewAlbumTracksTbl.Col("year"),
			userYearReviewAlbumTracksTbl.Col("album_id"),
			userYearReviewAlbumTracksTbl.Col("track_id"),
			userYearReviewAlbumTracksTbl.Col("rank"),
			userYearReviewAlbumTracksTbl.Col("play_count"),
			userYearReviewAlbumTracksTbl.Col("play_time"),
			userYearReviewAlbumTracksTbl.Col("created_at"),
			userYearReviewAlbumTracksTbl.Col("updated_at"),
		).
		Where(
			userYearReviewAlbumTracksTbl.Col("user_id").Eq(userId),
			userYearReviewAlbumTracksTbl.Col("year").Eq(year),
		).
		Order(
			userYearReviewAlbumTracksTbl.Col("album_id").Asc(),
			userYearReviewAlbumTracksTbl.Col("rank").Asc(),
		)

	return Multiple[UserYearReviewAlbumTrack](db, ctx, query)
}

func (db DB) GetUserYearReviewDay(
	ctx context.Context,
	userId string,
	year int,
) (UserYearReviewDay, error) {
	query := dialect.From(userYearReviewDayTbl).
		Select(
			userYearReviewDayTbl.Col("user_id"),
			userYearReviewDayTbl.Col("year"),
			userYearReviewDayTbl.Col("day"),
			userYearReviewDayTbl.Col("play_count"),
			userYearReviewDayTbl.Col("created_at"),
			userYearReviewDayTbl.Col("updated_at"),
		).
		Where(
			userYearReviewDayTbl.Col("user_id").Eq(userId),
			userYearReviewDayTbl.Col("year").Eq(year),
		)

	return Single[UserYearReviewDay](db, ctx, query)
}

func (db DB) GetUserYearReviewHours(
	ctx context.Context,
	userId string,
	year int,
) ([]UserYearReviewHour, error) {
	query := dialect.From(userYearReviewHoursTbl).
		Select(
			userYearReviewHoursTbl.Col("user_id"),
			userYearReviewHoursTbl.Col("year"),
			userYearReviewHoursTbl.Col("hour"),
			userYearReviewHoursTbl.Col("play_count"),
			userYearReviewHoursTbl.Col("created_at"),
			userYearReviewHoursTbl.Col("updated_at"),
		).
		Where(
			userYearReviewHoursTbl.Col("user_id").Eq(userId),
			userYearReviewHoursTbl.Col("year").Eq(year),
		).
		Order(userYearReviewHoursTbl.Col("hour").Asc())

	return Multiple[UserYearReviewHour](db, ctx, query)
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

func (db DB) GetUserYearReviewMilestones(
	ctx context.Context,
	userId string,
	year int,
) (UserYearReviewMilestone, error) {
	query := dialect.From(userYearReviewMilestonesTbl).
		Select(
			userYearReviewMilestonesTbl.Col("user_id"),
			userYearReviewMilestonesTbl.Col("year"),
			userYearReviewMilestonesTbl.Col("first_track_id"),
			userYearReviewMilestonesTbl.Col("last_track_id"),
			userYearReviewMilestonesTbl.Col("created_at"),
			userYearReviewMilestonesTbl.Col("updated_at"),
		).
		Where(
			userYearReviewMilestonesTbl.Col("user_id").Eq(userId),
			userYearReviewMilestonesTbl.Col("year").Eq(year),
		)

	return Single[UserYearReviewMilestone](db, ctx, query)
}
