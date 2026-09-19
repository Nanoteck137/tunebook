package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/tunebook/database"
)

const insertBatchSize = 400

// defaultMonthlyPlays is used when no existing history exists and the user has
// not supplied an explicit -monthly volume.
const defaultMonthlyPlays = 6000

// dayOfWeekWeights model the observed per-day listening volume for each
// weekday (indexed by time.Weekday: 0 = Sunday ... 6 = Saturday). Weekdays get
// noticeably more plays than weekends, matching the reference database.
var dayOfWeekWeights = [7]int{240, 246, 259, 290, 286, 230, 195}

// hourWeights model the observed listening volume for each hour of the day
// (0-23). Listening is spread fairly evenly throughout the day with a small
// dip in the afternoon.
var hourWeights = [24]int{
	86, 97, 105, 110, 104, 94, 89, 86,
	76, 77, 77, 80, 78, 77, 64, 72,
	61, 63, 84, 79, 78, 72, 83, 84,
}

var dialect = database.SqliteDialect()

type userRow struct {
	Id string `db:"id"`
}

type trackRow struct {
	Id       string `db:"id"`
	Duration int64  `db:"duration"`
}

type countRow struct {
	Count int `db:"count"`
}

type yearStatRow struct {
	TrackId   string `db:"track_id"`
	PlayCount int    `db:"play_count"`
	SkipCount int    `db:"skip_count"`
	PlayTime  int64  `db:"play_time"`
}

type aggregate struct {
	playCount  int
	skipCount  int
	playTime   int64
	completion int
}

func (a *aggregate) add(playTime int64, skipped bool, completion int) {
	a.playCount++
	a.playTime += playTime
	a.completion += completion
	if skipped {
		a.skipCount++
	}
}

type yearSummaryRow struct {
	Year          int     `db:"year"`
	TrackCount    int     `db:"track_count"`
	ListeningTime int64   `db:"listening_time"`
	AvgCompletion float64 `db:"avg_completion"`
	SkipCount     int     `db:"skip_count"`
	UniqueTracks  int     `db:"unique_tracks"`
	FavoritePlays int     `db:"favorite_plays"`
}

type statsAggRow struct {
	PlayCount int   `db:"play_count"`
	SkipCount int   `db:"skip_count"`
	PlayTime  int64 `db:"play_time"`
}

func main() {
	var (
		dbPath    = flag.String("db", "work/data.db", "path to the database file")
		userId    = flag.String("user", "", "target user id (default: user with the most stats)")
		from      = flag.Int("from", time.Now().Year()-5, "first year to generate")
		to        = flag.Int("to", time.Now().Year()-1, "last year to generate")
		tracks    = flag.Int("tracks", 300, "number of tracks to simulate per year")
		monthly   = flag.Int("monthly", 0, "target plays per month (default: auto-detected from the user's existing history)")
		favorites = flag.Int("favorites", 50, "number of simulated tracks to mark as favorites (0 disables)")
		gen       = flag.Bool("generate", false, "generate the user year reviews after inserting data")
		force     = flag.Bool("force", false, "overwrite existing data for the given years")
		seed      = flag.Int64("seed", 1, "random seed used for reproducible data")
	)
	flag.Parse()

	ctx := context.Background()

	db, err := database.Open(*dbPath)
	if err != nil {
		fmt.Printf("ERROR opening database: %v\n", err)
		return
	}
	defer db.Close()

	if *from > *to {
		fmt.Println("ERROR: -from must be <= -to")
		return
	}

	if *tracks <= 0 {
		fmt.Println("ERROR: -tracks must be > 0")
		return
	}

	uid, err := resolveUser(ctx, db, *userId)
	if err != nil {
		fmt.Printf("ERROR finding a user to use: %v\n", err)
		return
	}
	fmt.Printf("target user: %s\n", uid)

	selected, err := selectTracks(ctx, db, uid, *tracks)
	if err != nil {
		fmt.Printf("ERROR selecting tracks: %v\n", err)
		return
	}
	if len(selected) == 0 {
		fmt.Println("ERROR: no tracks with a duration found in the library")
		return
	}
	fmt.Printf("selected %d tracks\n", len(selected))

	monthlyPlays := *monthly
	if monthlyPlays <= 0 {
		detected, err := detectStandardMonth(ctx, db, uid)
		if err != nil {
			fmt.Printf("ERROR detecting standard month: %v\n", err)
			return
		}
		if detected <= 0 {
			detected = defaultMonthlyPlays
		}
		monthlyPlays = detected
		fmt.Printf("detected standard month volume: %d plays/month\n", monthlyPlays)
	}

	if *favorites > 0 {
		err := addFavorites(ctx, db, uid, selected, *favorites)
		if err != nil {
			fmt.Printf("ERROR adding favorites: %v\n", err)
			return
		}
	}

	rng := rand.New(rand.NewSource(*seed))
	runID := time.Now().UnixNano()

	for year := *from; year <= *to; year++ {
		err := synthesizeYear(
			ctx, db, uid, year, selected, monthlyPlays, rng, *force, *gen, runID,
		)
		if err != nil {
			fmt.Printf("ERROR generating year %d: %v\n", year, err)
			continue
		}
	}

	if err := rebuildUserStats(ctx, db, uid); err != nil {
		fmt.Printf("ERROR rebuilding user stats: %v\n", err)
		return
	}

	fmt.Println()
	printYearOverYear(ctx, db, uid)

	fmt.Println("\ndone")
}

func resolveUser(
	ctx context.Context,
	db *database.Database,
	requested string,
) (string, error) {
	if requested != "" {
		return requested, nil
	}

	query := dialect.From("user_track_stats").
		Select(goqu.I("user_id").As("id")).
		GroupBy("user_id").
		Order(goqu.COUNT(goqu.Star()).Desc()).
		Limit(1)

	row, err := database.Single[userRow](db, ctx, query)
	if err == nil {
		return row.Id, nil
	}

	if !errors.Is(err, database.ErrItemNotFound) {
		return "", err
	}

	fallback := dialect.From("users").
		Select("id").
		Order(goqu.I("created").Asc()).
		Limit(1)

	row, err = database.Single[userRow](db, ctx, fallback)
	if err != nil {
		return "", err
	}

	return row.Id, nil
}

func selectTracks(
	ctx context.Context,
	db *database.Database,
	userId string,
	limit int,
) ([]trackRow, error) {
	query := dialect.From(goqu.I("user_track_stats").As("uts")).
		Select(goqu.I("tracks.id"), goqu.I("tracks.duration")).
		Join(goqu.I("tracks"), goqu.On(
			goqu.I("tracks.id").Eq(goqu.I("uts.track_id")),
		)).
		Where(
			goqu.I("uts.user_id").Eq(userId),
			goqu.I("uts.period_type").Eq("all"),
			goqu.I("tracks.duration").Gt(0),
		).
		GroupBy(goqu.I("tracks.id"), goqu.I("tracks.duration")).
		Order(
			goqu.SUM(goqu.I("uts.play_count")).Desc(),
			goqu.I("tracks.id").Asc(),
		)

	played, err := database.Multiple[trackRow](db, ctx, query)
	if err != nil {
		return nil, err
	}

	tracks := played
	if len(tracks) > limit {
		tracks = tracks[:limit]
	}

	if len(tracks) < limit {
		extra := dialect.From("tracks").
			Select(goqu.I("tracks.id"), goqu.I("tracks.duration")).
			Where(
				goqu.I("tracks.duration").Gt(0),
				goqu.I("tracks.id").NotIn(
					dialect.From("user_track_stats").
						Select(goqu.I("track_id")).
						Where(goqu.I("user_id").Eq(userId)),
				),
			).
			Order(goqu.L("RANDOM()").Asc()).
			Limit(uint(limit - len(tracks)))

		padding, err := database.Multiple[trackRow](db, ctx, extra)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, padding...)
	}

	return tracks, nil
}

func addFavorites(
	ctx context.Context,
	db *database.Database,
	userId string,
	tracks []trackRow,
	count int,
) error {
	if count > len(tracks) {
		count = len(tracks)
	}

	now := time.Now().UnixMilli()
	rows := make([]goqu.Record, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, goqu.Record{
			"user_id":  userId,
			"track_id": tracks[i].Id,
			"added":    now,
		})
	}

	_, err := db.Exec(ctx, dialect.Insert("user_favorites").
		Rows(rows).
		OnConflict(goqu.DoNothing()))
	if err != nil {
		return err
	}

	fmt.Printf("marked %d tracks as favorites\n", count)
	return nil
}

func synthesizeYear(
	ctx context.Context,
	db *database.Database,
	userId string,
	year int,
	tracks []trackRow,
	plays int,
	rng *rand.Rand,
	force bool,
	generate bool,
	runID int64,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UnixMilli()

	if force {
		if err := clearYear(ctx, tx, userId, year); err != nil {
			return err
		}
	} else {
		hasData, err := yearHasData(ctx, tx, userId, year)
		if err != nil {
			return err
		}
		if hasData {
			fmt.Printf("year %d already has data, skipping (use -force to overwrite)\n", year)
			return nil
		}
	}

	monthAgg := make([]map[int]*aggregate, len(tracks))
	quarterAgg := make([]map[int]*aggregate, len(tracks))
	yearAgg := make([]*aggregate, len(tracks))
	allAgg := make([]*aggregate, len(tracks))
	for i := range tracks {
		monthAgg[i] = map[int]*aggregate{}
		quarterAgg[i] = map[int]*aggregate{}
		yearAgg[i] = &aggregate{}
		allAgg[i] = &aggregate{}
	}

	historySeq := 0
	historyRows := make([]goqu.Record, 0, plays*12)

	totalPlays := 0

	for month := 1; month <= 12; month++ {
		quarter := (month-1)/3 + 1

		monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		daysInMonth := monthStart.AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

		dayWeights := make([]int, daysInMonth)
		for d := 0; d < daysInMonth; d++ {
			wd := monthStart.AddDate(0, 0, d).Weekday()
			dayWeights[d] = dayOfWeekWeights[wd]
		}

		monthPlays := jitterPlays(plays, rng)
		counts := allocatePlays(len(tracks), monthPlays, rng)

		for i, track := range tracks {
			if counts[i] == 0 {
				continue
			}

			mAgg := ensureAgg(monthAgg[i], month)
			qAgg := ensureAgg(quarterAgg[i], quarter)

			for p := 0; p < counts[i]; p++ {
				day := weightedIndex(rng, dayWeights)
				hour := weightedIndex(rng, hourWeights[:])
				minute := rng.Intn(60)
				second := rng.Intn(60)

				listenedAt := time.Date(
					year, time.Month(month), day+1, hour, minute, second, 0, time.Local,
				).UnixMilli()

				percentPlayed, status, skipped := mimicPlayback(rng)

				playTime := int64(float64(track.Duration) * float64(percentPlayed) / 100.0)

				mAgg.add(playTime, skipped, percentPlayed)
				qAgg.add(playTime, skipped, percentPlayed)
				yearAgg[i].add(playTime, skipped, percentPlayed)
				allAgg[i].add(playTime, skipped, percentPlayed)

				historySeq++
				historyRows = append(historyRows, goqu.Record{
					"id":             fmt.Sprintf("mock-%d-%d-%d", year, runID, historySeq),
					"user_id":        userId,
					"track_id":       track.Id,
					"listened_at":    listenedAt,
					"playback_type":  "normal",
					"status":         status,
					"percent_played": percentPlayed,

					"created": listenedAt,
					"updated": listenedAt,
				})
			}
		}

		totalPlays += monthPlays
	}

	statRows := make([]goqu.Record, 0, len(tracks)*4)
	allRows := make([]goqu.Record, 0, len(tracks))

	for i, track := range tracks {
		if yearAgg[i].playCount == 0 {
			continue
		}

		statRows = append(statRows, statRecord(
			userId, track.Id, "year", year, 0, yearAgg[i], now))
		allRows = append(allRows, statRecord(
			userId, track.Id, "all", 0, 0, allAgg[i], now))

		for month, agg := range monthAgg[i] {
			if agg.playCount == 0 {
				continue
			}
			statRows = append(statRows, statRecord(
				userId, track.Id, "month", year, month, agg, now))
		}

		for quarter, agg := range quarterAgg[i] {
			if agg.playCount == 0 {
				continue
			}
			statRows = append(statRows, statRecord(
				userId, track.Id, "quarter", year, quarter, agg, now))
		}
	}

	err = upsertStats(ctx, tx, allRows)
	if err != nil {
		return err
	}

	err = upsertStats(ctx, tx, statRows)
	if err != nil {
		return err
	}

	err = insertChunked(ctx, tx, "track_history", historyRows)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("year %d: %d tracks, %d plays\n",
		year, numActiveTracks(yearAgg), totalPlays)

	if generate {
		err := db.GenerateUserReview(ctx, database.GenerateUserReviewParams{
			UserId: userId,
			Year:   year,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureAgg(m map[int]*aggregate, key int) *aggregate {
	agg := m[key]
	if agg == nil {
		agg = &aggregate{}
		m[key] = agg
	}
	return agg
}

func statRecord(
	userId string,
	trackId string,
	periodType string,
	year int,
	periodValue int,
	agg *aggregate,
	now int64,
) goqu.Record {
	return goqu.Record{
		"user_id":      userId,
		"track_id":     trackId,
		"period_type":  periodType,
		"year":         year,
		"period_value": periodValue,

		"play_count": agg.playCount,
		"skip_count": agg.skipCount,
		"play_time":  agg.playTime,

		"completion_sum": agg.completion,

		"created_at": now,
		"updated_at": now,
	}
}

func upsertStats(ctx context.Context, exec database.Executor, rows []goqu.Record) error {
	if len(rows) == 0 {
		return nil
	}

	for i := 0; i < len(rows); i += insertBatchSize {
		end := min(i+insertBatchSize, len(rows))

		_, err := exec.Exec(ctx, dialect.Insert("user_track_stats").
			Rows(rows[i:end]).
			OnConflict(goqu.DoUpdate(
				"user_id, track_id, period_type, year, period_value",
				goqu.Record{
					"play_count": goqu.L("play_count + EXCLUDED.play_count"),
					"skip_count": goqu.L("skip_count + EXCLUDED.skip_count"),
					"play_time":  goqu.L("play_time + EXCLUDED.play_time"),
					"completion_sum": goqu.L(
						"completion_sum + EXCLUDED.completion_sum",
					),
					"updated_at": goqu.L("EXCLUDED.updated_at"),
				},
			)))
		if err != nil {
			return err
		}
	}

	return nil
}

func numActiveTracks(aggs []*aggregate) int {
	n := 0
	for _, a := range aggs {
		if a.playCount > 0 {
			n++
		}
	}
	return n
}

func insertChunked(
	ctx context.Context,
	exec database.Executor,
	table string,
	rows []goqu.Record,
) error {
	for i := 0; i < len(rows); i += insertBatchSize {
		end := min(i+insertBatchSize, len(rows))

		_, err := exec.Exec(ctx, dialect.Insert(table).Rows(rows[i:end]))
		if err != nil {
			return err
		}
	}

	return nil
}

func clearYear(
	ctx context.Context,
	exec database.Executor,
	userId string,
	year int,
) error {
	now := time.Now().UnixMilli()

	query := dialect.From("user_track_stats").
		Select(
			goqu.I("track_id"),
			goqu.I("play_count"),
			goqu.I("skip_count"),
			goqu.I("play_time"),
		).
		Where(
			goqu.I("user_id").Eq(userId),
			goqu.I("period_type").Eq("year"),
			goqu.I("year").Eq(year),
		)

	rows, err := database.Multiple[yearStatRow](exec, ctx, query)
	if err != nil {
		return err
	}

	if len(rows) > 0 {
		rollback := make([]goqu.Record, len(rows))
		for i, r := range rows {
			rollback[i] = goqu.Record{
				"user_id":      userId,
				"track_id":     r.TrackId,
				"period_type":  "all",
				"year":         0,
				"period_value": 0,

				"play_count": -r.PlayCount,
				"skip_count": -r.SkipCount,
				"play_time":  -r.PlayTime,

				"created_at": now,
				"updated_at": now,
			}
		}

		if err := upsertStats(ctx, exec, rollback); err != nil {
			return err
		}
	}

	_, err = exec.Exec(ctx, dialect.Delete("user_track_stats").Where(
		goqu.Ex{
			"user_id":     userId,
			"period_type": goqu.Op{"in": []string{"year", "quarter", "month"}},
			"year":        year,
		},
	))
	if err != nil {
		return err
	}

	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).UnixMilli()
	end := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).UnixMilli()

	_, err = exec.Exec(ctx, dialect.Delete("track_history").Where(
		goqu.I("user_id").Eq(userId),
		goqu.And(
			goqu.I("listened_at").Gte(start),
			goqu.I("listened_at").Lt(end),
		),
	))
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, dialect.Delete("user_year_reviews").Where(
		goqu.Ex{
			"user_id": userId,
			"year":    year,
		},
	))
	if err != nil {
		return err
	}

	return nil
}

func yearHasData(
	ctx context.Context,
	exec database.Executor,
	userId string,
	year int,
) (bool, error) {
	query := dialect.From("user_track_stats").
		Select(goqu.COUNT(goqu.Star()).As("count")).
		Where(
			goqu.Ex{
				"user_id":     userId,
				"period_type": goqu.Op{"in": []string{"year", "quarter", "month"}},
				"year":        year,
			},
		)

	row, err := database.Single[countRow](exec, ctx, query)
	if err != nil {
		return false, err
	}

	if row.Count > 0 {
		return true, nil
	}

	start, end := yearRange(year)

	query = dialect.From("track_history").
		Select(goqu.COUNT(goqu.Star()).As("count")).
		Where(
			goqu.I("user_id").Eq(userId),
			goqu.And(
				goqu.I("listened_at").Gte(start),
				goqu.I("listened_at").Lt(end),
			),
		)

	row, err = database.Single[countRow](exec, ctx, query)
	if err != nil {
		return false, err
	}

	return row.Count > 0, nil
}

// allocatePlays splits a total number of plays across tracks using a heavy
// tailed random popularity. Every year gets a different distribution so the
// top tracks/artists shift between years.
func allocatePlays(trackCount int, total int, rng *rand.Rand) []int {
	weights := make([]float64, trackCount)
	sum := 0.0
	for i := range weights {
		weights[i] = math.Pow(rng.Float64(), 3) + 0.02
		sum += weights[i]
	}

	counts := make([]int, trackCount)
	remaining := total
	for i := range counts {
		allocated := int(float64(total) * weights[i] / sum)
		if allocated > remaining {
			allocated = remaining
		}
		counts[i] = allocated
		remaining -= allocated
	}

	if remaining > 0 && trackCount > 0 {
		counts[0] += remaining
	}

	return counts
}

// detectStandardMonth computes the average number of plays per month from the
// user's existing history so that generated data matches the volume already in
// the database. It returns 0 when there is no history to infer from.
func detectStandardMonth(
	ctx context.Context,
	db *database.Database,
	userId string,
) (int, error) {
	type monthRow struct {
		Month string `db:"month"`
		Count int    `db:"count"`
	}

	query := dialect.From("track_history").
		Select(
			goqu.L("strftime('%Y-%m', listened_at / 1000, 'unixepoch')").As("month"),
			goqu.COUNT(goqu.Star()).As("count"),
		).
		Where(goqu.I("user_id").Eq(userId)).
		GroupBy(goqu.L("strftime('%Y-%m', listened_at / 1000, 'unixepoch')"))

	rows, err := database.Multiple[monthRow](db, ctx, query)
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, nil
	}

	total := 0
	for _, r := range rows {
		total += r.Count
	}

	return total / len(rows), nil
}

// weightedIndex picks an index from weights with probability proportional to
// each entry.
func weightedIndex(rng *rand.Rand, weights []int) int {
	sum := 0
	for _, w := range weights {
		sum += w
	}

	r := rng.Intn(sum)
	for i, w := range weights {
		if r < w {
			return i
		}
		r -= w
	}

	return len(weights) - 1
}

// jitterPlays applies mild random variation to the standard monthly volume so
// consecutive months aren't identical, while staying close to the reference
// quantity.
func jitterPlays(base int, rng *rand.Rand) int {
	const low = 0.9
	const high = 1.1
	factor := low + rng.Float64()*(high-low)
	return int(float64(base) * factor)
}

// mimicPlayback reproduces the completion/skip pattern observed in the
// reference database: almost every play is a completed listen at 100%, with a
// tiny number of skips (-0.1%) and the rare partially completed listen.
//
// It returns the percent played, the status, and whether the play was skipped.
func mimicPlayback(rng *rand.Rand) (int, string, bool) {
	if rng.Float64() < 0.001 {
		return 10 + rng.Intn(70), "skipped", true
	}

	if rng.Float64() < 0.001 {
		return 94 + rng.Intn(6), "completed", false
	}

	return 100, "completed", false
}

func rebuildUserStats(
	ctx context.Context,
	db *database.Database,
	userId string,
) error {
	agg, err := database.Single[statsAggRow](db, ctx,
		dialect.From("user_track_stats").
			Select(
				goqu.COALESCE(goqu.SUM("play_count"), 0).As("play_count"),
				goqu.COALESCE(goqu.SUM("skip_count"), 0).As("skip_count"),
				goqu.COALESCE(goqu.SUM("play_time"), 0).As("play_time"),
			).
			Where(
				goqu.I("user_id").Eq(userId),
				goqu.I("period_type").Eq("all"),
			))
	if err != nil {
		return err
	}

	favorites, err := database.Single[countRow](db, ctx,
		dialect.From("user_favorites").
			Select(goqu.COUNT(goqu.Star()).As("count")).
			Where(goqu.I("user_id").Eq(userId)))
	if err != nil {
		return err
	}

	type lastRow struct {
		Last sql.NullInt64 `db:"last"`
	}
	last, err := database.Single[lastRow](db, ctx,
		dialect.From("track_history").
			Select(goqu.MAX("listened_at").As("last")).
			Where(goqu.I("user_id").Eq(userId)))
	if err != nil {
		return err
	}

	type playlistsRow struct {
		NumPlaylistsCreated int `db:"num_playlists_created"`
	}
	playlists := 0
	existing, err := database.Single[playlistsRow](db, ctx,
		dialect.From("user_stats").
			Select("num_playlists_created").
			Where(goqu.I("user_id").Eq(userId)))
	if err == nil {
		playlists = existing.NumPlaylistsCreated
	} else if !errors.Is(err, database.ErrItemNotFound) {
		return err
	}

	return db.SetUserStats(ctx, database.SetUserStatsParams{
		UserId:              userId,
		NumTracksPlayed:     agg.PlayCount,
		NumTracksSkipped:    agg.SkipCount,
		NumPlaylistsCreated: playlists,
		NumFavoriteTracks:   favorites.Count,
		ListeningTime:       agg.PlayTime,
		LastListenedAt:      last.Last,
	})
}

func yearRange(year int) (int64, int64) {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).UnixMilli()
	end := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).UnixMilli()
	return start, end
}

func printYearOverYear(
	ctx context.Context,
	db *database.Database,
	userId string,
) {
	rows, err := database.Multiple[yearSummaryRow](db, ctx,
		dialect.From("user_year_reviews").
			Select(
				"year",
				"track_count",
				"listening_time",
				"avg_completion",
				"skip_count",
				"unique_tracks",
				"favorite_plays",
			).
			Where(goqu.I("user_id").Eq(userId)).
			Order(goqu.I("year").Asc()))
	if err != nil {
		fmt.Printf("could not load year reviews: %v\n", err)
		return
	}

	if len(rows) == 0 {
		fmt.Println("no year reviews generated yet")
		return
	}

	fmt.Println("Year over year review")
	fmt.Printf("%-6s %10s %10s %10s %8s %10s %10s %10s\n",
		"Year", "Plays", "Hours", "Unique", "Skipped", "Skip%", "AvgComp", "FavPlays")

	for _, r := range rows {
		skipRate := 0.0
		if r.TrackCount > 0 {
			skipRate = float64(r.SkipCount) / float64(r.TrackCount) * 100
		}

		fmt.Printf("%-6d %10d %10.1f %10d %8d %9.1f%% %9.1f%% %10d\n",
			r.Year,
			r.TrackCount,
			float64(r.ListeningTime)/3600.0,
			r.UniqueTracks,
			r.SkipCount,
			skipRate,
			r.AvgCompletion,
			r.FavoritePlays,
		)
	}
}
