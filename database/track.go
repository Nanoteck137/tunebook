package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/types"
)

var (
	createTrackId = createIdGenerator(32)

	tracksTbl                 = goqu.T("tracks")
	tracksTagsTbl             = goqu.T("tracks_tags")
	tracksFeaturingArtistsTbl = goqu.T("tracks_featuring_artists")
)

type Track struct {
	Id string `db:"id"`

	Filename     string            `db:"filename"`
	ModifiedTime int64             `db:"modified_time"`
	MediaFormat  types.MediaFormat `db:"media_format"`

	Name string `db:"name"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	Duration int64         `db:"duration"`
	Number   sql.NullInt64 `db:"number"`
	Year     sql.NullInt64 `db:"year"`

	AlbumName     string         `db:"album_name"`
	AlbumCoverArt sql.NullString `db:"album_cover_art"`

	ArtistName string `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists JsonColumn[[]FeaturingArtist] `db:"featuring_artists"`

	Order *int
}

func SqlGroupConcat(col any, seperator string) exp.SQLFunctionExpression {
	return goqu.Func("group_concat", col, seperator)
}

func ObjectTagsQuery(objectTagTable, idCol string) *goqu.SelectDataset {
	// TODO(patrik): Replace with tagsTable
	tagsTable := goqu.T("tags")
	tagsTableSlugCol := tagsTable.Col("slug")

	// table := goqu.T("tracks_tags")
	table := goqu.T(objectTagTable)
	tableSlugCol := table.Col("tag_slug")
	tableIdCol := table.Col(idCol)

	query := dialect.From(table).
		Select(
			tableIdCol.As("id"),
			SqlGroupConcat(tagsTableSlugCol, ",").As("data"),
		).
		Join(
			tagsTable,
			goqu.On(tableSlugCol.Eq(tagsTableSlugCol)),
		).
		GroupBy(tableIdCol)

	return query
}

func (db DB) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.id").Eq(id))

	return Single[Track](db, ctx, query)
}

func (db DB) GetTracksByIds(
	ctx context.Context, 
	ids []string,
) ([]Track, error) {
	query := TrackQuery().Where(goqu.I("tracks.id").In(ids))

	return Multiple[Track](db, ctx, query)
}

func (db DB) GetTrackIdsByFilter(
	ctx context.Context, 
	filterStr string,
) ([]string, error) {
	query := TrackQuery().Select(goqu.I("tracks.id"))

	a := adapter.TrackResolverAdapter{}
	query, err := applyFilterParams(types.FilterParams{
		Filter: filterStr,
	}, &a, query)
	if err != nil {
		return nil, err
	}

	return Multiple[string](db, ctx, query)
}

func (db DB) GetAllTrackIds(ctx context.Context) ([]string, error) {
	query := dialect.From("tracks").
		Select("tracks.id")

	return Multiple[string](db, ctx, query)
}

type GetTracksParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (db DB) GetTracks(
	ctx context.Context,
	params GetTracksParams,
) ([]Track, types.Page, error) {
	query := TrackQuery()

	var err error

	a := adapter.TrackResolverAdapter{}
	query, err = applyFilterParams(params.Filter, &a, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, "tracks.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[Track](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetTrackIdsByAlbum(
	ctx context.Context, 
	albumId, filterStr string,
) ([]string, error) {
	var err error

	query := TrackQuery().Select("tracks.id")

	r := filter.New(&adapter.TrackResolverAdapter{})
	query, err = applyFilterCustom(
		query, r, filterStr, goqu.I("tracks.album_id").Eq(albumId))
	if err != nil {
		return nil, err
	}

	query = query.Order(
		goqu.I("tracks.number").Asc().NullsLast(),
		goqu.I("tracks.name").Asc(),
	)

	s, _, _ := query.ToSQL()
	fmt.Printf("s: %v\n", s)

	// dialect.From("tracks").
	// 	Select("tracks.id").
	// 	Where(goqu.I("tracks.album_id").Eq(albumId)).
	// 	Order(
	// 		goqu.I("tracks.number").Asc().NullsLast(),
	// 		goqu.I("tracks.name").Asc(),
	// 	)

	return Multiple[string](db, ctx, query)
}

func (db DB) GetTrackIdsByArtist(
	ctx context.Context, 
	artistId, filterStr string,
) ([]string, error) {
	var err error

	query := TrackQuery().Select("tracks.id")

	r := filter.New(&adapter.TrackResolverAdapter{})
	query, err = applyFilterCustom(
		query, r, filterStr, goqu.I("tracks.artist_id").Eq(artistId))
	if err != nil {
		return nil, err
	}

	query = query.Order(goqu.I("tracks.name").Asc())

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByPlaylistParams struct {
	PlaylistId string
	FilterStr  string
}

func (db DB) GetTrackIdsByPlaylist(
	ctx context.Context, 
	params GetTrackIdsByPlaylistParams,
) ([]string, error) {
	var err error

	query := TrackQuery().Select("tracks.id").
		Join(
			goqu.I("playlist_items"), 
			goqu.On(tracksTbl.Col("id").Eq(goqu.I("playlist_items.track_id"))),
		)

	r := filter.New(&adapter.TrackResolverAdapter{})
	query, err = applyFilterCustom(
		query, r, params.FilterStr, 
		goqu.I("playlist_items.playlist_id").Eq(params.PlaylistId))
	if err != nil {
		return nil, err
	}

	query = query.Order(goqu.I("playlist_items.position").Asc())

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByUserFavoritesParams struct {
	UserId    string
	FilterStr string
}

func (db DB) GetTrackIdsByUserFavorites(
	ctx context.Context, 
	params GetTrackIdsByUserFavoritesParams,
) ([]string, error) {
	var err error

	query := TrackQuery().Select("tracks.id").
		Join(
			goqu.I("user_favorites"), 
			goqu.On(tracksTbl.Col("id").Eq(goqu.I("user_favorites.track_id"))),
		)

	r := filter.New(&adapter.TrackResolverAdapter{})
	query, err = applyFilterCustom(
		query, r, params.FilterStr, 
		goqu.I("user_favorites.user_id").Eq(params.UserId))
	if err != nil {
		return nil, err
	}

	query = query.Order(goqu.I("user_favorites.added").Desc().NullsLast())

	return Multiple[string](db, ctx, query)
}

func (db DB) GetTracksByAlbum(
	ctx context.Context, 
	albumId string,
) ([]Track, error) {
	query := TrackQuery().
		Where(
			goqu.I("tracks.id").In(
				goqu.From("tracks").
					Select("tracks.id").
					Where(goqu.I("tracks.album_id").Eq(albumId)),
			),
		).
		Order(
			goqu.I("tracks.number").Asc().NullsLast(),
			goqu.I("tracks.name").Asc(),
		)

	return Multiple[Track](db, ctx, query)
}

func (db DB) GetTracksIn(
	ctx context.Context, 
	in any, 
	sort string,
) ([]Track, error) {
	query := TrackQuery().
		Where(
			goqu.I("tracks.id").In(in),
		)

	a := adapter.TrackResolverAdapter{}
	resolver := filter.New(&a)

	query, err := applySort(query, resolver, sort)
	if err != nil {
		return nil, err
	}

	return Multiple[Track](db, ctx, query)
}

type CreateTrackParams struct {
	Id string

	Filename     string
	ModifiedTime int64
	MediaFormat  types.MediaFormat

	Name string

	AlbumId  string
	ArtistId string

	Duration int64
	Number   sql.NullInt64
	Year     sql.NullInt64

	Created int64
	Updated int64
}

func (db DB) CreateTrack(
	ctx context.Context, 
	params CreateTrackParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createTrackId()
	}

	query := dialect.Insert("tracks").Rows(goqu.Record{
		"id": params.Id,

		"filename":      params.Filename,
		"modified_time": params.ModifiedTime,
		"media_format":  params.MediaFormat,

		"name": params.Name,

		"album_id":  params.AlbumId,
		"artist_id": params.ArtistId,

		"duration": params.Duration,
		"number":   params.Number,
		"year":     params.Year,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type TrackChanges struct {
	Filename     Change[string]
	ModifiedTime Change[int64]
	MediaFormat  Change[types.MediaFormat]

	Name Change[string]

	AlbumId  Change[string]
	ArtistId Change[string]

	Duration Change[int64]
	Number   Change[sql.NullInt64]
	Year     Change[sql.NullInt64]

	Created Change[int64]
}

func (db DB) UpdateTrack(
	ctx context.Context, 
	id string, 
	changes TrackChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "filename", changes.Filename)
	addToRecord(record, "modified_time", changes.ModifiedTime)
	addToRecord(record, "media_format", changes.MediaFormat)

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "album_id", changes.AlbumId)
	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "duration", changes.Duration)
	addToRecord(record, "number", changes.Number)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("tracks").
		Set(record).
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTrack(ctx context.Context, id string) error {
	query := dialect.Delete("tracks").
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) AddTrackTag(ctx context.Context, tagSlug, trackId string) error {
	query := dialect.Insert(tracksTagsTbl).
		Rows(goqu.Record{
			"track_id": trackId,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) RemoveAllTrackTags(ctx context.Context, trackId string) error {
	query := dialect.Delete(tracksTagsTbl).
		Where(goqu.I("track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) RemoveAllTrackFeaturingArtists(
	ctx context.Context, 
	trackId string,
) error {
	query := dialect.Delete(tracksFeaturingArtistsTbl).
		Where(tracksFeaturingArtistsTbl.Col("track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) AddFeaturingArtistToTrack(
	ctx context.Context, 
	trackId, artistId string,
) error {
	query := dialect.Insert(tracksFeaturingArtistsTbl).
		Rows(goqu.Record{
			"track_id":  trackId,
			"artist_id": artistId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func TrackQuery() *goqu.SelectDataset {
	idCol := tracksTbl.Col("id")

	query := dialect.From(tracksTbl).
		Select(
			idCol,

			tracksTbl.Col("filename"),
			tracksTbl.Col("modified_time"),
			tracksTbl.Col("media_format"),

			tracksTbl.Col("name"),

			tracksTbl.Col("album_id"),
			tracksTbl.Col("artist_id"),

			tracksTbl.Col("number"),
			tracksTbl.Col("duration"),
			tracksTbl.Col("year"),

			tracksTbl.Col("created"),
			tracksTbl.Col("updated"),

			// TODO(patrik): Replace with albumsTable.Col("name")
			goqu.I("albums.name").As("album_name"),
			// TODO(patrik): Replace with albumsTable.Col("cover_art")
			goqu.I("albums.cover_art").As("album_cover_art"),

			// TODO(patrik): Replace with artistsTable.Col("name")
			goqu.I("artists.name").As("artist_name"),

			goqu.I("tags.data").As("tags"),
			goqu.I("featuring_artists.artists").As("featuring_artists"),
		).
		Join(
			// TODO(patrik): Replace with albumsTable
			goqu.I("albums"),
			// TODO(patrik): Replace with albumsTable.Col("id")
			goqu.On(tracksTbl.Col("album_id").Eq(goqu.I("albums.id"))),
		).
		Join(
			// TODO(patrik): Replace with artistsTable
			goqu.I("artists"),
			// TODO(patrik): Replace with artistsTable.Col("id")
			goqu.On(tracksTbl.Col("artist_id").Eq(goqu.I("artists.id"))),
		).
		LeftJoin(
			// TODO(patrik): Fix
			ObjectTagsQuery("tracks_tags", "track_id").As("tags"),
			goqu.On(idCol.Eq(goqu.I("tags.id"))),
		).
		LeftJoin(
			// TODO(patrik): Fix
			FeaturingArtistsQuery(
				"tracks_featuring_artists", 
				"track_id",
			).As("featuring_artists"),
			goqu.On(idCol.Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}
