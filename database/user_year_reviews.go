package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/types"
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

	userYearReviewTrackSchema  = UserYearReviewTrackSchema()
	userYearReviewAlbumSchema  = UserYearReviewAlbumSchema()
	userYearReviewArtistSchema = UserYearReviewArtistSchema()
	userYearReviewTagSchema    = UserYearReviewTagSchema()
	userYearReviewDecadeSchema = UserYearReviewDecadeSchema()

	userYearReviewMonthTrackSchema  = UserYearReviewMonthTrackSchema()
	userYearReviewMonthAlbumSchema  = UserYearReviewMonthAlbumSchema()
	userYearReviewMonthArtistSchema = UserYearReviewMonthArtistSchema()
	userYearReviewMonthTagSchema    = UserYearReviewMonthTagSchema()
	userYearReviewMonthDecadeSchema = UserYearReviewMonthDecadeSchema()
)

func UserYearReviewTrackSchema() *schema.Schema {
	return TrackSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_tracks.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_tracks.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewAlbumSchema() *schema.Schema {
	return AlbumSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_albums.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_albums.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewArtistSchema() *schema.Schema {
	return ArtistSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_artists.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_artists.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewDecadeSchema() *schema.Schema {
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_decades.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_decades.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewTagSchema() *schema.Schema {
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_tags.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_tags.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewMonthTrackSchema() *schema.Schema {
	return TrackSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_month_tracks.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_month_tracks.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewMonthAlbumSchema() *schema.Schema {
	return AlbumSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_month_albums.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_month_albums.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewMonthArtistSchema() *schema.Schema {
	return ArtistSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_month_artists.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_month_artists.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewMonthTagSchema() *schema.Schema {
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_month_tags.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_month_tags.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserYearReviewMonthDecadeSchema() *schema.Schema {
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_year_review_month_decades.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_year_review_month_decades.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

type UserYearSummary struct {
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

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

func (db DB) generateUserYearReviewTracks(
	ctx context.Context,
	userId string,
	year int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_tracks (
		user_id, 
		year, 
		track_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?,
			?,
			user_track_stats.track_id,
			ROW_NUMBER() OVER (
				ORDER BY user_track_stats.play_count DESC, 
					user_track_stats.track_id ASC
			),
			user_track_stats.play_count,
			?,
			?
		FROM user_track_stats
		WHERE 
			user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'year' AND 
			user_track_stats.year = ?
	`

	params := []any{userId, year, now, now, userId, year}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewMonthTracks(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_month_tracks (
		user_id, 
		year, 
		month, 
		track_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?,
			?,
			?,
			user_track_stats.track_id,
			ROW_NUMBER() OVER (
				ORDER BY user_track_stats.play_count DESC,
					user_track_stats.track_id ASC
			),
			user_track_stats.play_count,
			?,
			?
		FROM user_track_stats
		WHERE 
			user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'month' AND 
			user_track_stats.year = ? AND 
			user_track_stats.period_value = ?
	`

	params := []any{userId, year, month, now, now, userId, year, month}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewAlbums(
	ctx context.Context,
	userId string,
	year int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_albums (
		user_id, 
		year, 
		album_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			albums.id,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
				albums.id ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		JOIN albums ON albums.id = tracks.album_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'year' AND 
			user_track_stats.year = ?
		GROUP BY albums.id
	`

	params := []any{userId, year, now, now, userId, year}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewMonthAlbums(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_month_albums (
		user_id, 
		year, 
		month, 
		album_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			?, 
			albums.id,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
				albums.id ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		JOIN albums ON albums.id = tracks.album_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'month' AND 
			user_track_stats.year = ? AND 
			user_track_stats.period_value = ?
		GROUP BY albums.id
	`

	params := []any{userId, year, month, now, now, userId, year, month}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewArtists(
	ctx context.Context,
	userId string,
	year int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_artists (
		user_id, 
		year, 
		artist_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			artists.id,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
					artists.id ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		JOIN artists ON artists.id = tracks.artist_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'year' AND 
			user_track_stats.year = ?
		GROUP BY artists.id
	`

	params := []any{userId, year, now, now, userId, year}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewMonthArtists(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_month_artists (
		user_id, 
		year, 
		month, 
		artist_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			?, 
			artists.id,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC,
					artists.id ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		JOIN artists ON artists.id = tracks.artist_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'month' AND 
			user_track_stats.year = ? AND 
			user_track_stats.period_value = ?
		GROUP BY artists.id
	`

	params := []any{userId, year, month, now, now, userId, year, month}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewTags(
	ctx context.Context,
	userId string,
	year int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_tags (
		user_id, 
		year, 
		tag_slug, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			tags.slug,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
					tags.slug ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks_tags ON tracks_tags.track_id = user_track_stats.track_id
		JOIN tags ON tags.slug = tracks_tags.tag_slug
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'year' AND 
			user_track_stats.year = ?
		GROUP BY tags.slug
	`

	params := []any{userId, year, now, now, userId, year}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return err
}

func (db DB) generateUserYearReviewMonthTags(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_month_tags (
		user_id, 
		year, 
		month, 
		tag_slug, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			?, 
			tags.slug,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
					tags.slug ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks_tags ON tracks_tags.track_id = user_track_stats.track_id
		JOIN tags ON tags.slug = tracks_tags.tag_slug
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'month' AND 
			user_track_stats.year = ? AND 
			user_track_stats.period_value = ?
		GROUP BY tags.slug
	`

	params := []any{userId, year, month, now, now, userId, year, month}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserYearReviewDecades(
	ctx context.Context,
	userId string,
	year int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_decades (
		user_id, 
		year, 
		decade, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			tracks.year - tracks.year % 10,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
					tracks.year - tracks.year % 10 ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'year' AND 
			user_track_stats.year = ? AND 
			tracks.year IS NOT NULL
		GROUP BY tracks.year - tracks.year % 10
	`

	params := []any{userId, year, now, now, userId, year}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return err
}

func (db DB) generateUserYearReviewMonthDecades(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	query := `
	INSERT INTO user_year_review_month_decades (
		user_id, 
		year, 
		month, 
		decade, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
			?, 
			?, 
			?, 
			tracks.year - tracks.year % 10,
			ROW_NUMBER() OVER (
				ORDER BY SUM(user_track_stats.play_count) DESC, 
				tracks.year - tracks.year % 10 ASC
			),
			SUM(user_track_stats.play_count), 
			?, 
			?
		FROM user_track_stats
		JOIN tracks ON tracks.id = user_track_stats.track_id
		WHERE user_track_stats.user_id = ? AND 
			user_track_stats.period_type = 'month' AND 
			user_track_stats.year = ? AND 
			user_track_stats.period_value = ? AND 
			tracks.year IS NOT NULL
		GROUP BY tracks.year - tracks.year % 10
	`

	params := []any{userId, year, month, now, now, userId, year, month}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
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
	query := `
	SELECT 
		ROUND(
			COALESCE(
				CAST(SUM(s.completion_sum) AS REAL) / SUM(s.play_count), 0
			), 1
		) AS avg_completion,
		COALESCE(
			SUM(s.skip_count), 0
		) AS skip_count,
		COALESCE(
			COUNT(DISTINCT s.track_id), 0
		) AS unique_tracks
	FROM user_track_stats s
	WHERE s.user_id = ? AND 
		s.period_type = 'year' AND 
		s.year = ?
  	`

	return Single[UserYearHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, year},
	})
}

func (db DB) GetUserYearFavoritePlays(
	ctx context.Context,
	userId string,
	year int,
) (int, error) {
	query := `
	SELECT 
		COALESCE(SUM(s.play_count), 0)
	FROM user_track_stats s
	JOIN user_favorites f ON f.user_id = s.user_id AND f.track_id = s.track_id
	WHERE s.user_id = ? AND 
		s.period_type = 'year' AND 
		s.year = ?
  	`

	return Single[int](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, year},
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
	query := `
	SELECT 
		ROUND(
			COALESCE(
				CAST(SUM(completion_sum) AS REAL) / SUM(play_count), 0
			), 1
		) AS avg_completion,
		COALESCE(
			SUM(skip_count), 0
		) AS skip_count,
		COALESCE(
			COUNT(DISTINCT track_id), 0
		) AS unique_tracks
	FROM user_track_stats
	WHERE user_id = ? AND 
		period_type = 'month' AND 
		year = ? AND 
		period_value = ?
  	`

	return Single[UserYearHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, year, month},
	})
}

func (db DB) GetUserYearMonthFavoritePlays(
	ctx context.Context,
	userId string,
	year int,
	month int,
) (int, error) {
	query := `
	SELECT 
		COALESCE(SUM(s.play_count), 0)
	FROM user_track_stats s
	JOIN user_favorites f ON f.user_id = s.user_id AND f.track_id = s.track_id
	WHERE s.user_id = ? AND 
		s.period_type = 'month' AND 
		s.year = ? AND 
		s.period_value = ?
  	`

	return Single[int](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, year, month},
	})
}

func (db DB) processUserYearMonth(
	ctx context.Context,
	userId string,
	year int,
	month int,
	now int64,
) error {
	monthSummary, err := db.GetUserYearMonthSummary(ctx, userId, year, month)
	if err != nil {
		return err
	}

	monthHistory, err := db.GetUserYearMonthHistorySummary(
		ctx, userId, year, month)
	if err != nil {
		return err
	}

	monthFavorites, err := db.GetUserYearMonthFavoritePlays(
		ctx, userId, year, month)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, dialect.Insert(userYearReviewMonthsTbl).Rows(
		goqu.Record{
			"user_id":        userId,
			"year":           year,
			"month":          month,
			"play_count":     monthSummary.TrackCount,
			"play_time":      monthSummary.ListeningTime,
			"avg_completion": monthHistory.AvgCompletion,
			"skip_count":     monthHistory.SkipCount,
			"unique_tracks":  monthHistory.UniqueTracks,
			"favorite_plays": monthFavorites,

			"created": now,
			"updated": now,
		},
	))
	if err != nil {
		return err
	}

	err = db.generateUserYearReviewMonthTracks(
		ctx, userId, year, month, now)
	if err != nil {
		return err
	}

	err = db.generateUserYearReviewMonthAlbums(
		ctx, userId, year, month, now)
	if err != nil {
		return err
	}

	err = db.generateUserYearReviewMonthArtists(
		ctx, userId, year, month, now)
	if err != nil {
		return err
	}

	err = db.generateUserYearReviewMonthTags(
		ctx, userId, year, month, now)
	if err != nil {
		return err
	}

	err = db.generateUserYearReviewMonthDecades(
		ctx, userId, year, month, now)
	if err != nil {
		return err
	}

	return nil
}

type GenerateUserReviewParams struct {
	UserId string
	Year   int
}

func (db *Database) GenerateUserReview(
	ctx context.Context,
	params GenerateUserReviewParams,
) error {
	now := time.Now().UnixMilli()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clean up any previously generated review for this year.
	_, err = tx.Exec(ctx, dialect.Delete(userYearReviewsTbl).Where(goqu.Ex{
		"user_id": params.UserId,
		"year":    params.Year,
	}))
	if err != nil {
		return err
	}

	summary, err := tx.GetUserYearSummary(ctx, params.UserId, params.Year)
	if err != nil {
		return err
	}

	historySummary, err := tx.GetUserYearHistorySummary(
		ctx, params.UserId, params.Year)
	if err != nil {
		return err
	}

	favorites, err := tx.GetUserYearFavoritePlays(
		ctx, params.UserId, params.Year)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, dialect.Insert(userYearReviewsTbl).Rows(
		goqu.Record{
			"user_id":        params.UserId,
			"year":           params.Year,
			"track_count":    summary.TrackCount,
			"listening_time": summary.ListeningTime,
			"avg_completion": historySummary.AvgCompletion,
			"skip_count":     historySummary.SkipCount,
			"unique_tracks":  historySummary.UniqueTracks,
			"favorite_plays": favorites,

			"created": now,
			"updated": now,
		},
	))
	if err != nil {
		return err
	}

	err = tx.generateUserYearReviewTracks(
		ctx, params.UserId, params.Year, now)
	if err != nil {
		return err
	}

	err = tx.generateUserYearReviewAlbums(
		ctx, params.UserId, params.Year, now)
	if err != nil {
		return err
	}

	err = tx.generateUserYearReviewArtists(
		ctx, params.UserId, params.Year, now)
	if err != nil {
		return err
	}

	err = tx.generateUserYearReviewTags(
		ctx, params.UserId, params.Year, now)
	if err != nil {
		return err
	}

	err = tx.generateUserYearReviewDecades(
		ctx, params.UserId, params.Year, now)
	if err != nil {
		return err
	}

	for m := 1; m <= 12; m++ {
		err := tx.processUserYearMonth(ctx, params.UserId, params.Year, m, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type GetUserStatsYearsParams struct {
	UserId string
}

func (db DB) GetUserStatsYears(
	ctx context.Context,
	params GetUserStatsYearsParams,
) ([]int, error) {
	query := dialect.From(userTrackStatsTbl).
		SelectDistinct(
			userTrackStatsTbl.Col("year"),
		).
		Where(
			userTrackStatsTbl.Col("user_id").Eq(params.UserId),
			userTrackStatsTbl.Col("period_type").Eq("year"),
		).
		Order(userTrackStatsTbl.Col("year").Desc())

	return Multiple[int](db, ctx, query)
}

type UserYearReview struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`

	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`

	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
	FavoritePlays int     `db:"favorite_plays"`

	CreatedAt int64 `db:"created"`
	UpdatedAt int64 `db:"updated"`
}

type GetUserYearReviewsParams struct {
	UserId string
}

func (db DB) GetUserYearReviews(
	ctx context.Context,
	params GetUserYearReviewsParams,
) ([]UserYearReview, error) {
	query := dialect.From(userYearReviewsTbl).
		Select(
			userYearReviewsTbl.Col("user_id"),
			userYearReviewsTbl.Col("year"),

			userYearReviewsTbl.Col("track_count"),
			userYearReviewsTbl.Col("listening_time"),
			userYearReviewsTbl.Col("avg_completion"),
			userYearReviewsTbl.Col("skip_count"),
			userYearReviewsTbl.Col("unique_tracks"),
			userYearReviewsTbl.Col("favorite_plays"),

			userYearReviewsTbl.Col("created"),
			userYearReviewsTbl.Col("updated"),
		).
		Where(userYearReviewsTbl.Col("user_id").Eq(params.UserId)).
		Order(userYearReviewsTbl.Col("year").Desc())

	return Multiple[UserYearReview](db, ctx, query)
}

type GetUserYearReviewParams struct {
	UserId string
	Year   int
}

func (db DB) GetUserYearReview(
	ctx context.Context,
	params GetUserYearReviewParams,
) (UserYearReview, error) {
	query := dialect.From(userYearReviewsTbl).
		Select(
			userYearReviewsTbl.Col("user_id"),
			userYearReviewsTbl.Col("year"),

			userYearReviewsTbl.Col("track_count"),
			userYearReviewsTbl.Col("listening_time"),
			userYearReviewsTbl.Col("avg_completion"),
			userYearReviewsTbl.Col("skip_count"),
			userYearReviewsTbl.Col("unique_tracks"),
			userYearReviewsTbl.Col("favorite_plays"),

			userYearReviewsTbl.Col("created"),
			userYearReviewsTbl.Col("updated"),
		).
		Where(
			userYearReviewsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewsTbl.Col("year").Eq(params.Year),
		)

	return Single[UserYearReview](db, ctx, query)
}

type UserYearReviewTrack struct {
	Track

	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	TrackId string `db:"track_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewTracksParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewTracks(
	ctx context.Context,
	params GetUserYearReviewTracksParams,
) ([]UserYearReviewTrack, types.Page, error) {
	var err error

	query := TrackQuery().
		SelectAppend(
			userYearReviewTracksTbl.Col("user_id"),
			userYearReviewTracksTbl.Col("year"),
			userYearReviewTracksTbl.Col("track_id"),

			userYearReviewTracksTbl.Col("rank"),
			userYearReviewTracksTbl.Col("play_count"),
		).
		Join(
			userYearReviewTracksTbl,
			goqu.On(userYearReviewTracksTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			userYearReviewTracksTbl.Col("user_id").Eq(params.UserId),
			userYearReviewTracksTbl.Col("year").Eq(params.Year),
		)

	query, err = ApplyQuery(query, userYearReviewTrackSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, tracksTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewAlbum struct {
	Album

	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	AlbumId string `db:"album_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewAlbumsParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewAlbums(
	ctx context.Context,
	params GetUserYearReviewAlbumsParams,
) ([]UserYearReviewAlbum, types.Page, error) {
	var err error

	query := AlbumQuery().
		SelectAppend(
			userYearReviewAlbumsTbl.Col("user_id"),
			userYearReviewAlbumsTbl.Col("year"),
			userYearReviewAlbumsTbl.Col("album_id"),

			userYearReviewAlbumsTbl.Col("rank"),
			userYearReviewAlbumsTbl.Col("play_count"),
		).
		Join(
			userYearReviewAlbumsTbl,
			goqu.On(userYearReviewAlbumsTbl.Col("album_id").Eq(albumsTbl.Col("id"))),
		).
		Where(
			userYearReviewAlbumsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewAlbumsTbl.Col("year").Eq(params.Year),
		)

	query, err = ApplyQuery(query, userYearReviewAlbumSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, albumsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewAlbum](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewArtist struct {
	Artist

	UserId   string `db:"user_id"`
	Year     int    `db:"year"`
	ArtistId string `db:"artist_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewArtistsParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewArtists(
	ctx context.Context,
	params GetUserYearReviewArtistsParams,
) ([]UserYearReviewArtist, types.Page, error) {
	var err error

	query := ArtistQuery().
		SelectAppend(
			userYearReviewArtistsTbl.Col("user_id"),
			userYearReviewArtistsTbl.Col("year"),
			userYearReviewArtistsTbl.Col("artist_id"),

			userYearReviewArtistsTbl.Col("rank"),
			userYearReviewArtistsTbl.Col("play_count"),
		).
		Join(
			userYearReviewArtistsTbl,
			goqu.On(
				userYearReviewArtistsTbl.Col("artist_id").
					Eq(artistsTbl.Col("id")),
			),
		).
		Where(
			userYearReviewArtistsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewArtistsTbl.Col("year").Eq(params.Year),
		)

	query, err = ApplyQuery(query, userYearReviewArtistSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, artistsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewArtist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewTag struct {
	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	TagSlug string `db:"tag_slug"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewTagsParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewTags(
	ctx context.Context,
	params GetUserYearReviewTagsParams,
) ([]UserYearReviewTag, types.Page, error) {
	var err error

	query := dialect.From(userYearReviewTagsTbl).
		Select(
			userYearReviewTagsTbl.Col("user_id"),
			userYearReviewTagsTbl.Col("year"),
			userYearReviewTagsTbl.Col("tag_slug"),

			userYearReviewTagsTbl.Col("rank"),
			userYearReviewTagsTbl.Col("play_count"),
		).
		Where(
			userYearReviewTagsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewTagsTbl.Col("year").Eq(params.Year),
		)

	query, err = ApplyQuery(query, userYearReviewTagSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx, db, params.Page, query, userYearReviewTagsTbl.Col("tag_slug"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewTag](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewDecade struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Decade int    `db:"decade"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewDecadesParams struct {
	UserId string
	Year   int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewDecades(
	ctx context.Context,
	params GetUserYearReviewDecadesParams,
) ([]UserYearReviewDecade, types.Page, error) {
	var err error

	query := dialect.From(userYearReviewDecadesTbl).
		Select(
			userYearReviewDecadesTbl.Col("user_id"),
			userYearReviewDecadesTbl.Col("year"),
			userYearReviewDecadesTbl.Col("decade"),

			userYearReviewDecadesTbl.Col("rank"),
			userYearReviewDecadesTbl.Col("play_count"),
		).
		Where(
			userYearReviewDecadesTbl.Col("user_id").Eq(params.UserId),
			userYearReviewDecadesTbl.Col("year").Eq(params.Year),
		)

	query, err = ApplyQuery(query, userYearReviewDecadeSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx, db, params.Page, query, userYearReviewDecadesTbl.Col("decade"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewDecade](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
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

	CreatedAt int64 `db:"created"`
	UpdatedAt int64 `db:"updated"`
}

type GetUserYearReviewMonthsParams struct {
	UserId string
	Year   int
}

func (db DB) GetUserYearReviewMonths(
	ctx context.Context,
	params GetUserYearReviewMonthsParams,
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

			userYearReviewMonthsTbl.Col("created"),
			userYearReviewMonthsTbl.Col("updated"),
		).
		Where(
			userYearReviewMonthsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthsTbl.Col("year").Eq(params.Year),
		).
		Order(
			userYearReviewMonthsTbl.Col("month").Asc(),
		)

	return Multiple[UserYearReviewMonth](db, ctx, query)
}

type GetUserYearReviewMonthParams struct {
	UserId string
	Year   int
	Month  int
}

func (db DB) GetUserYearReviewMonth(
	ctx context.Context,
	params GetUserYearReviewMonthParams,
) (UserYearReviewMonth, error) {
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

			userYearReviewMonthsTbl.Col("created"),
			userYearReviewMonthsTbl.Col("updated"),
		).
		Where(
			userYearReviewMonthsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthsTbl.Col("year").Eq(params.Year),
			userYearReviewMonthsTbl.Col("month").Eq(params.Month),
		).
		Order(
			userYearReviewMonthsTbl.Col("month").Asc(),
		)

	return Single[UserYearReviewMonth](db, ctx, query)
}

type UserYearReviewMonthTrack struct {
	Track

	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	Month   int    `db:"month"`
	TrackId string `db:"track_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewMonthTracksParams struct {
	UserId string
	Year   int
	Month  int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewMonthTracks(
	ctx context.Context,
	params GetUserYearReviewMonthTracksParams,
) ([]UserYearReviewMonthTrack, types.Page, error) {
	var err error

	query := TrackQuery().
		SelectAppend(
			userYearReviewMonthTracksTbl.Col("user_id"),
			userYearReviewMonthTracksTbl.Col("year"),
			userYearReviewMonthTracksTbl.Col("month"),
			userYearReviewMonthTracksTbl.Col("track_id"),

			userYearReviewMonthTracksTbl.Col("rank"),
			userYearReviewMonthTracksTbl.Col("play_count"),
		).
		Join(
			userYearReviewMonthTracksTbl,
			goqu.On(
				userYearReviewMonthTracksTbl.Col("track_id").
					Eq(tracksTbl.Col("id")),
			),
		).
		Where(
			userYearReviewMonthTracksTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthTracksTbl.Col("year").Eq(params.Year),
			userYearReviewMonthTracksTbl.Col("month").Eq(params.Month),
		)

	query, err = ApplyQuery(
		query, userYearReviewMonthTrackSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, tracksTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewMonthTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewMonthAlbum struct {
	Album

	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	Month   int    `db:"month"`
	AlbumId string `db:"album_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewMonthAlbumsParams struct {
	UserId string
	Year   int
	Month  int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewMonthAlbums(
	ctx context.Context,
	params GetUserYearReviewMonthAlbumsParams,
) ([]UserYearReviewMonthAlbum, types.Page, error) {
	var err error

	query := AlbumQuery().
		SelectAppend(
			userYearReviewMonthAlbumsTbl.Col("user_id"),
			userYearReviewMonthAlbumsTbl.Col("year"),
			userYearReviewMonthAlbumsTbl.Col("month"),
			userYearReviewMonthAlbumsTbl.Col("album_id"),

			userYearReviewMonthAlbumsTbl.Col("rank"),
			userYearReviewMonthAlbumsTbl.Col("play_count"),
		).
		Join(
			userYearReviewMonthAlbumsTbl,
			goqu.On(
				userYearReviewMonthAlbumsTbl.Col("album_id").
					Eq(albumsTbl.Col("id")),
			),
		).
		Where(
			userYearReviewMonthAlbumsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthAlbumsTbl.Col("year").Eq(params.Year),
			userYearReviewMonthAlbumsTbl.Col("month").Eq(params.Month),
		)

	query, err = ApplyQuery(
		query, userYearReviewMonthAlbumSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, albumsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewMonthAlbum](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewMonthArtist struct {
	Artist

	UserId   string `db:"user_id"`
	Year     int    `db:"year"`
	Month    int    `db:"month"`
	ArtistId string `db:"artist_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewMonthArtistsParams struct {
	UserId string
	Year   int
	Month  int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewMonthArtists(
	ctx context.Context,
	params GetUserYearReviewMonthArtistsParams,
) ([]UserYearReviewMonthArtist, types.Page, error) {
	var err error

	query := ArtistQuery().
		SelectAppend(
			userYearReviewMonthArtistsTbl.Col("user_id"),
			userYearReviewMonthArtistsTbl.Col("year"),
			userYearReviewMonthArtistsTbl.Col("month"),
			userYearReviewMonthArtistsTbl.Col("artist_id"),

			userYearReviewMonthArtistsTbl.Col("rank"),
			userYearReviewMonthArtistsTbl.Col("play_count"),
		).
		Join(
			userYearReviewMonthArtistsTbl,
			goqu.On(
				userYearReviewMonthArtistsTbl.Col("artist_id").
					Eq(artistsTbl.Col("id")),
			),
		).
		Where(
			userYearReviewMonthArtistsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthArtistsTbl.Col("year").Eq(params.Year),
			userYearReviewMonthArtistsTbl.Col("month").Eq(params.Month),
		)

	query, err = ApplyQuery(
		query, userYearReviewMonthArtistSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, artistsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewMonthArtist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewMonthTag struct {
	UserId  string `db:"user_id"`
	Year    int    `db:"year"`
	Month   int    `db:"month"`
	TagSlug string `db:"tag_slug"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserYearReviewMonthTagsParams struct {
	UserId string
	Year   int
	Month  int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewMonthTags(
	ctx context.Context,
	params GetUserYearReviewMonthTagsParams,
) ([]UserYearReviewMonthTag, types.Page, error) {
	var err error

	query := dialect.From(userYearReviewMonthTagsTbl).
		Select(
			userYearReviewMonthTagsTbl.Col("user_id"),
			userYearReviewMonthTagsTbl.Col("year"),
			userYearReviewMonthTagsTbl.Col("month"),
			userYearReviewMonthTagsTbl.Col("tag_slug"),

			userYearReviewMonthTagsTbl.Col("rank"),
			userYearReviewMonthTagsTbl.Col("play_count"),
		).
		Where(
			userYearReviewMonthTagsTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthTagsTbl.Col("year").Eq(params.Year),
			userYearReviewMonthTagsTbl.Col("month").Eq(params.Month),
		)

	query, err = ApplyQuery(query, userYearReviewMonthTagSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx,
		db,
		params.Page,
		query,
		userYearReviewMonthTagsTbl.Col("tag_slug"),
	)
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewMonthTag](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserYearReviewMonthDecade struct {
	UserId string `db:"user_id"`
	Year   int    `db:"year"`
	Month  int    `db:"month"`
	Decade int    `db:"decade"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}


type GetUserYearReviewMonthDecadesParams struct {
	UserId string
	Year   int
	Month  int

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserYearReviewMonthDecades(
	ctx context.Context,
	params GetUserYearReviewMonthDecadesParams,
) ([]UserYearReviewMonthDecade, types.Page, error) {
	var err error

	query := dialect.From(userYearReviewMonthDecadesTbl).
		Select(
			userYearReviewMonthDecadesTbl.Col("user_id"),
			userYearReviewMonthDecadesTbl.Col("year"),
			userYearReviewMonthDecadesTbl.Col("month"),
			userYearReviewMonthDecadesTbl.Col("decade"),

			userYearReviewMonthDecadesTbl.Col("rank"),
			userYearReviewMonthDecadesTbl.Col("play_count"),
		).
		Where(
			userYearReviewMonthDecadesTbl.Col("user_id").Eq(params.UserId),
			userYearReviewMonthDecadesTbl.Col("year").Eq(params.Year),
			userYearReviewMonthDecadesTbl.Col("month").Eq(params.Month),
		)

	query, err = ApplyQuery(
		query, userYearReviewMonthDecadeSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx,
		db,
		params.Page,
		query,
		userYearReviewMonthDecadesTbl.Col("decade"),
	)
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserYearReviewMonthDecade](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}
