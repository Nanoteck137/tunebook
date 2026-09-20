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
	userTotalReviewsTbl = goqu.T("user_total_reviews")

	userTotalReviewTracksTbl  = goqu.T("user_total_review_tracks")
	userTotalReviewAlbumsTbl  = goqu.T("user_total_review_albums")
	userTotalReviewArtistsTbl = goqu.T("user_total_review_artists")
	userTotalReviewMonthsTbl  = goqu.T("user_total_review_months")
	userTotalReviewTagsTbl    = goqu.T("user_total_review_tags")
	userTotalReviewDecadesTbl = goqu.T("user_total_review_decades")

	userTotalReviewTrackSchema  = UserTotalReviewTrackSchema()
	userTotalReviewAlbumSchema  = UserTotalReviewAlbumSchema()
	userTotalReviewArtistSchema = UserTotalReviewArtistSchema()
	userTotalReviewTagSchema    = UserTotalReviewTagSchema()
	userTotalReviewDecadeSchema = UserTotalReviewDecadeSchema()
)

func UserTotalReviewTrackSchema() *schema.Schema {
	// TODO(patrik): Should we add the other columns from user_total_review_tracks?
	return TrackSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_total_review_tracks.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_total_review_tracks.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserTotalReviewAlbumSchema() *schema.Schema {
	// TODO(patrik): Should we add the other columns from user_total_review_albums?
	return AlbumSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_total_review_albums.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_total_review_albums.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserTotalReviewArtistSchema() *schema.Schema {
	// TODO(patrik): Should we add the other columns from user_total_review_artists?
	return ArtistSchema().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_total_review_artists.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_total_review_artists.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserTotalReviewDecadeSchema() *schema.Schema {
	// TODO(patrik): Should we add the other columns from user_total_review_decades?
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_total_review_decades.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_total_review_decades.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

func UserTotalReviewTagSchema() *schema.Schema {
	// TODO(patrik): Should we add the other columns from user_total_review_tags?
	return schema.New().
		AddField(
			"rank",
			query.TypeInt,
			schema.Column("user_total_review_tags.rank"),
		).
		AddField(
			"play_count",
			query.TypeInt,
			schema.Column("user_total_review_tags.play_count"),
		).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "rank"},
				Dir:   query.DirAsc,
			},
		)
}

type UserTotalSummary struct {
	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`
}

func (db DB) GetUserTotalSummary(
	ctx context.Context,
	userId string,
) (UserTotalSummary, error) {
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
			userTrackStatsTbl.Col("period_type").Eq("all"),
		)

	return Single[UserTotalSummary](db, ctx, query)
}

func (db DB) generateUserTotalReviewTracks(
	ctx context.Context,
	userId string,
	now int64,
) error {
	query := `
	INSERT INTO user_total_review_tracks (
		user_id, 
		track_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
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
			user_track_stats.period_type = 'all'
	`

	params := []any{userId, now, now, userId}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserTotalReviewAlbums(
	ctx context.Context,
	userId string,
	now int64,
) error {
	query := `
	INSERT INTO user_total_review_albums (
		user_id, 
		album_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
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
			user_track_stats.period_type = 'all'
		GROUP BY albums.id
	`

	params := []any{userId, now, now, userId}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserTotalReviewArtists(
	ctx context.Context,
	userId string,
	now int64,
) error {
	query := `
	INSERT INTO user_total_review_artists (
		user_id, 
		artist_id, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
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
			user_track_stats.period_type = 'all'
		GROUP BY artists.id
	`

	params := []any{userId, now, now, userId}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserTotalReviewTags(
	ctx context.Context,
	userId string,
	now int64,
) error {
	query := `
	INSERT INTO user_total_review_tags (
		user_id, 
		tag_slug, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
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
			user_track_stats.period_type = 'all'
		GROUP BY tags.slug
	`

	params := []any{userId, now, now, userId}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) generateUserTotalReviewDecades(
	ctx context.Context,
	userId string,
	now int64,
) error {
	query := `
	INSERT INTO user_total_review_decades (
		user_id, 
		decade, 
		rank, 
		play_count, 
		created, 
		updated
	)
		SELECT 
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
			user_track_stats.period_type = 'all' AND 
			tracks.year IS NOT NULL
		GROUP BY tracks.year - tracks.year % 10
	`

	params := []any{userId, now, now, userId}

	_, err := db.Exec(ctx, RawQuery{Query: query, Params: params})
	if err != nil {
		return err
	}

	return nil
}

type UserTotalHistorySummary struct {
	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
}

func (db DB) GetUserTotalHistorySummary(
	ctx context.Context,
	userId string,
) (UserTotalHistorySummary, error) {
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
		s.period_type = 'all'
  	`

	return Single[UserTotalHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId},
	})
}

func (db DB) GetUserTotalFavoritePlays(
	ctx context.Context,
	userId string,
) (int, error) {
	query := `
	SELECT 
		COALESCE(SUM(s.play_count), 0)
	FROM user_track_stats s
	JOIN user_favorites f ON f.user_id = s.user_id AND f.track_id = s.track_id
	WHERE s.user_id = ? AND 
		s.period_type = 'all'
  	`

	return Single[int](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId},
	})
}

func (db DB) GetUserTotalMonthSummary(
	ctx context.Context,
	userId string,
	month int,
) (UserTotalSummary, error) {
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
			userTrackStatsTbl.Col("period_value").Eq(month),
		)

	return Single[UserTotalSummary](db, ctx, query)
}

func (db DB) GetUserTotalMonthHistorySummary(
	ctx context.Context,
	userId string,
	month int,
) (UserTotalHistorySummary, error) {
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
		period_value = ?
  	`

	return Single[UserTotalHistorySummary](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, month},
	})
}

func (db DB) GetUserTotalMonthFavoritePlays(
	ctx context.Context,
	userId string,
	month int,
) (int, error) {
	query := `
	SELECT 
		COALESCE(SUM(s.play_count), 0)
	FROM user_track_stats s
	JOIN user_favorites f ON f.user_id = s.user_id AND f.track_id = s.track_id
	WHERE s.user_id = ? AND 
		s.period_type = 'month' AND 
		s.period_value = ?
  	`

	return Single[int](db, ctx, RawQuery{
		Query:  query,
		Params: []any{userId, month},
	})
}

func processUserTotalMonth(
	tx DB,
	ctx context.Context,
	userId string,
	month int,
	now int64,
) error {
	monthSummary, err := tx.GetUserTotalMonthSummary(ctx, userId, month)
	if err != nil {
		return err
	}

	monthHistory, err := tx.GetUserTotalMonthHistorySummary(
		ctx, userId, month)
	if err != nil {
		return err
	}

	monthFavorites, err := tx.GetUserTotalMonthFavoritePlays(
		ctx, userId, month)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, dialect.Insert(userTotalReviewMonthsTbl).Rows(
		goqu.Record{
			"user_id":        userId,
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

	return nil
}

type UserTotalReview struct {
	UserId string `db:"user_id"`

	TrackCount    int   `db:"track_count"`
	ListeningTime int64 `db:"listening_time"`

	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
	FavoritePlays int     `db:"favorite_plays"`

	CreatedAt int64 `db:"created"`
	UpdatedAt int64 `db:"updated"`
}

type GetUserTotalReviewParams struct {
	UserId string
}

func (db DB) GetUserTotalReview(
	ctx context.Context,
	params GetUserTotalReviewParams,
) (UserTotalReview, error) {
	query := dialect.From(userTotalReviewsTbl).
		Select(
			userTotalReviewsTbl.Col("user_id"),

			userTotalReviewsTbl.Col("track_count"),
			userTotalReviewsTbl.Col("listening_time"),
			userTotalReviewsTbl.Col("avg_completion"),
			userTotalReviewsTbl.Col("skip_count"),
			userTotalReviewsTbl.Col("unique_tracks"),
			userTotalReviewsTbl.Col("favorite_plays"),

			userTotalReviewsTbl.Col("created"),
			userTotalReviewsTbl.Col("updated"),
		).
		Where(userTotalReviewsTbl.Col("user_id").Eq(params.UserId))

	return Single[UserTotalReview](db, ctx, query)
}

type UserTotalReviewTrack struct {
	Track

	UserId  string `db:"user_id"`
	TrackId string `db:"track_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserTotalReviewTracksParams struct {
	UserId string

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserTotalReviewTracks(
	ctx context.Context,
	params GetUserTotalReviewTracksParams,
) ([]UserTotalReviewTrack, types.Page, error) {
	var err error

	query := TrackQuery().
		SelectAppend(
			userTotalReviewTracksTbl.Col("user_id"),
			userTotalReviewTracksTbl.Col("track_id"),

			userTotalReviewTracksTbl.Col("rank"),
			userTotalReviewTracksTbl.Col("play_count"),
		).
		Join(
			userTotalReviewTracksTbl,
			goqu.On(userTotalReviewTracksTbl.Col("track_id").Eq(tracksTbl.Col("id"))),
		).
		Where(
			userTotalReviewTracksTbl.Col("user_id").Eq(params.UserId),
		)

	query, err = ApplyQuery(query, userTotalReviewTrackSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, tracksTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTotalReviewTrack](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserTotalReviewAlbum struct {
	Album

	UserId  string `db:"user_id"`
	AlbumId string `db:"album_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserTotalReviewAlbumsParams struct {
	UserId string

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserTotalReviewAlbums(
	ctx context.Context,
	params GetUserTotalReviewAlbumsParams,
) ([]UserTotalReviewAlbum, types.Page, error) {
	var err error

	query := AlbumQuery().
		SelectAppend(
			userTotalReviewAlbumsTbl.Col("user_id"),
			userTotalReviewAlbumsTbl.Col("album_id"),

			userTotalReviewAlbumsTbl.Col("rank"),
			userTotalReviewAlbumsTbl.Col("play_count"),
		).
		Join(
			userTotalReviewAlbumsTbl,
			goqu.On(userTotalReviewAlbumsTbl.Col("album_id").Eq(albumsTbl.Col("id"))),
		).
		Where(
			userTotalReviewAlbumsTbl.Col("user_id").Eq(params.UserId),
		)

	query, err = ApplyQuery(query, userTotalReviewAlbumSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, albumsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTotalReviewAlbum](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserTotalReviewArtist struct {
	Artist

	UserId   string `db:"user_id"`
	ArtistId string `db:"artist_id"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserTotalReviewArtistsParams struct {
	UserId string

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserTotalReviewArtists(
	ctx context.Context,
	params GetUserTotalReviewArtistsParams,
) ([]UserTotalReviewArtist, types.Page, error) {
	var err error

	query := ArtistQuery().
		SelectAppend(
			userTotalReviewArtistsTbl.Col("user_id"),
			userTotalReviewArtistsTbl.Col("artist_id"),

			userTotalReviewArtistsTbl.Col("rank"),
			userTotalReviewArtistsTbl.Col("play_count"),
		).
		Join(
			userTotalReviewArtistsTbl,
			goqu.On(
				userTotalReviewArtistsTbl.Col("artist_id").
					Eq(artistsTbl.Col("id")),
			),
		).
		Where(
			userTotalReviewArtistsTbl.Col("user_id").Eq(params.UserId),
		)

	query, err = ApplyQuery(query, userTotalReviewArtistSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, artistsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTotalReviewArtist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserTotalReviewTag struct {
	UserId  string `db:"user_id"`
	TagSlug string `db:"tag_slug"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserTotalReviewTagsParams struct {
	UserId string

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserTotalReviewTags(
	ctx context.Context,
	params GetUserTotalReviewTagsParams,
) ([]UserTotalReviewTag, types.Page, error) {
	var err error

	query := dialect.From(userTotalReviewTagsTbl).
		Select(
			userTotalReviewTagsTbl.Col("user_id"),
			userTotalReviewTagsTbl.Col("tag_slug"),

			userTotalReviewTagsTbl.Col("rank"),
			userTotalReviewTagsTbl.Col("play_count"),
		).
		Where(
			userTotalReviewTagsTbl.Col("user_id").Eq(params.UserId),
		)

	query, err = ApplyQuery(query, userTotalReviewTagSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx, db, params.Page, query, userTotalReviewTagsTbl.Col("tag_slug"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTotalReviewTag](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserTotalReviewDecade struct {
	UserId string `db:"user_id"`
	Decade int    `db:"decade"`

	Rank      int `db:"rank"`
	PlayCount int `db:"play_count"`
}

type GetUserTotalReviewDecadesParams struct {
	UserId string

	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetUserTotalReviewDecades(
	ctx context.Context,
	params GetUserTotalReviewDecadesParams,
) ([]UserTotalReviewDecade, types.Page, error) {
	var err error

	query := dialect.From(userTotalReviewDecadesTbl).
		Select(
			userTotalReviewDecadesTbl.Col("user_id"),
			userTotalReviewDecadesTbl.Col("decade"),

			userTotalReviewDecadesTbl.Col("rank"),
			userTotalReviewDecadesTbl.Col("play_count"),
		).
		Where(
			userTotalReviewDecadesTbl.Col("user_id").Eq(params.UserId),
		)

	query, err = ApplyQuery(query, userTotalReviewDecadeSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(
		ctx, db, params.Page, query, userTotalReviewDecadesTbl.Col("decade"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[UserTotalReviewDecade](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

type UserTotalReviewMonth struct {
	UserId string `db:"user_id"`
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

type GetUserTotalReviewMonthsParams struct {
	UserId string
}

func (db DB) GetUserTotalReviewMonths(
	ctx context.Context,
	params GetUserTotalReviewMonthsParams,
) ([]UserTotalReviewMonth, error) {
	query := dialect.From(userTotalReviewMonthsTbl).
		Select(
			userTotalReviewMonthsTbl.Col("user_id"),
			userTotalReviewMonthsTbl.Col("month"),

			userTotalReviewMonthsTbl.Col("play_count"),
			userTotalReviewMonthsTbl.Col("play_time"),

			userTotalReviewMonthsTbl.Col("avg_completion"),
			userTotalReviewMonthsTbl.Col("skip_count"),
			userTotalReviewMonthsTbl.Col("unique_tracks"),
			userTotalReviewMonthsTbl.Col("favorite_plays"),

			userTotalReviewMonthsTbl.Col("created"),
			userTotalReviewMonthsTbl.Col("updated"),
		).
		Where(
			userTotalReviewMonthsTbl.Col("user_id").Eq(params.UserId),
		).
		Order(
			userTotalReviewMonthsTbl.Col("month").Asc(),
		)

	return Multiple[UserTotalReviewMonth](db, ctx, query)
}

type GenerateUserTotalReviewParams struct {
	UserId string
}

func (db *Database) GenerateUserTotalReview(
	ctx context.Context,
	params GenerateUserTotalReviewParams,
) error {
	now := time.Now().UnixMilli()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clean up any previously generated review for this user.
	_, err = tx.Exec(ctx, dialect.Delete(userTotalReviewsTbl).Where(goqu.Ex{
		"user_id": params.UserId,
	}))
	if err != nil {
		return err
	}

	summary, err := tx.GetUserTotalSummary(ctx, params.UserId)
	if err != nil {
		return err
	}

	historySummary, err := tx.GetUserTotalHistorySummary(ctx, params.UserId)
	if err != nil {
		return err
	}

	favorites, err := tx.GetUserTotalFavoritePlays(ctx, params.UserId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, dialect.Insert(userTotalReviewsTbl).Rows(
		goqu.Record{
			"user_id":        params.UserId,
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

	err = tx.generateUserTotalReviewTracks(
		ctx, params.UserId, now)
	if err != nil {
		return err
	}

	err = tx.generateUserTotalReviewAlbums(
		ctx, params.UserId, now)
	if err != nil {
		return err
	}

	err = tx.generateUserTotalReviewArtists(
		ctx, params.UserId, now)
	if err != nil {
		return err
	}

	err = tx.generateUserTotalReviewTags(
		ctx, params.UserId, now)
	if err != nil {
		return err
	}

	err = tx.generateUserTotalReviewDecades(
		ctx, params.UserId, now)
	if err != nil {
		return err
	}

	for m := 1; m <= 12; m++ {
		err := processUserTotalMonth(
			tx.DB, ctx, params.UserId, m, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
