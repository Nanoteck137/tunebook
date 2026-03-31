package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database/adapter"
	"github.com/nanoteck137/tunebook/tools/filter"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/pyrin/ember"
)

type Track struct {
	RowId int `db:"rowid"`

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

	FeaturingArtists ember.JsonColumn[[]FeaturingArtist] `db:"featuring_artists"`

	Order *int
}

// TODO(patrik): Move
func FeaturingArtistsQuery(table, idColName string) *goqu.SelectDataset {
	tbl := goqu.T(table)

	return dialect.From(tbl).
		Select(
			tbl.Col(idColName).As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",
					"id",
					goqu.I("artists.id"),
					"name",
					goqu.I("artists.name"),
				),
			).As("artists"),
		).
		Join(
			goqu.I("artists"),
			goqu.On(tbl.Col("artist_id").Eq(goqu.I("artists.id"))),
		).
		GroupBy(tbl.Col(idColName))
}

// TODO(patrik): Use goqu.T more
func TrackQuery() *goqu.SelectDataset {
	tags := dialect.From("tracks_tags").
		Select(
			goqu.I("tracks_tags.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("tracks_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("tracks_tags.track_id"))

	query := dialect.From("tracks").
		Select(
			"tracks.rowid",

			"tracks.id",

			"tracks.filename",
			"tracks.modified_time",
			"tracks.media_format",

			"tracks.name",

			"tracks.album_id",
			"tracks.artist_id",

			"tracks.number",
			"tracks.duration",
			"tracks.year",

			"tracks.created",
			"tracks.updated",

			goqu.I("albums.name").As("album_name"),
			goqu.I("albums.cover_art").As("album_cover_art"),

			goqu.I("artists.name").As("artist_name"),

			goqu.I("tags.tags").As("tags"),

			goqu.I("featuring_artists.artists").As("featuring_artists"),
		).
		Join(
			goqu.I("albums"),
			goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id"))),
		).
		Join(
			goqu.I("artists"),
			goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id"))),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("tags.track_id"))),
		).
		LeftJoin(
			FeaturingArtistsQuery("tracks_featuring_artists", "track_id").As("featuring_artists"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}

func (db DB) GetAllTracksByArtistId(ctx context.Context, artistId string) ([]Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.artist_id").Eq(artistId))

	return ember.Multiple[Track](db.db, ctx, query)
}

func (db DB) GetTracksByIds(ctx context.Context, ids []string) ([]Track, error) {
	query := TrackQuery().Where(goqu.I("tracks.id").In(ids))

	return ember.Multiple[Track](db.db, ctx, query)
}

func (db DB) GetTracksByAlbumForPlay(ctx context.Context, albumId string) ([]NewTrackQueryItem, error) {
	query := NewTrackQuery().
		Where(goqu.I("tracks.album_id").Eq(albumId)).
		Order(goqu.I("tracks.number").Asc().NullsLast(), goqu.I("tracks.name").Asc()).
		As("tracks")

	return ember.Multiple[NewTrackQueryItem](db.db, ctx, query)
}

func (db DB) GetAllTrackIds(ctx context.Context) ([]string, error) {
	query := dialect.From("tracks").
		Select("tracks.id")

	return ember.Multiple[string](db.db, ctx, query)
}

// TODO(patrik): Move
type FetchOptions struct {
	Filter  string
	Sort    string
	PerPage int
	Page    int
}

func (db DB) GetAllTracks(ctx context.Context, filterStr, sortStr string) ([]Track, error) {
	query := TrackQuery()

	var err error

	a := adapter.TrackResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Track](db.db, ctx, query)
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

	page, err := buildPage(ctx, db.db, params.Page, query, "tracks.id")
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := ember.Multiple[Track](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetTracksByAlbum(ctx context.Context, albumId string) ([]Track, error) {
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

	return ember.Multiple[Track](db.db, ctx, query)
}

func AlbumTrackSubquery(albumId string) *goqu.SelectDataset {
	return goqu.From("tracks").
		Select("tracks.id").
		Where(goqu.I("tracks.album_id").Eq(albumId))
}

func ArtistTrackSubquery(artistId string) *goqu.SelectDataset {
	tbl := goqu.T("tracks_featuring_artists")
	return goqu.From("tracks").
		Select("tracks.id").
		FullOuterJoin(
			tbl,
			goqu.On(
				goqu.I("tracks.id").Eq(tbl.Col("track_id")),
			),
		).
		Where(
			goqu.Or(
				goqu.I("tracks.artist_id").Eq(artistId),
				tbl.Col("artist_id").Eq(artistId),
			),
		)
}

func PlaylistTrackSubquery(playlistId string) *goqu.SelectDataset {
	return goqu.From("playlist_items").
		Select("playlist_items.track_id").
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
		)
}

func (db DB) GetTracksIn(ctx context.Context, in any, sort string) ([]Track, error) {
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

	return ember.Multiple[Track](db.db, ctx, query)
}

func (db DB) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.id").Eq(id))

	return ember.Single[Track](db.db, ctx, query)
}

func (db DB) GetTrackByNameAndAlbum(ctx context.Context, name string, albumId string) (Track, error) {
	query := TrackQuery().
		Where(
			goqu.And(
				goqu.I("tracks.name").Eq(name),
				goqu.I("tracks.album_id").Eq(albumId),
			),
		)

	return ember.Single[Track](db.db, ctx, query)
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

func (db DB) CreateTrack(ctx context.Context, params CreateTrackParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateId()
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

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type TrackChanges struct {
	Filename     types.Change[string]
	ModifiedTime types.Change[int64]
	MediaFormat  types.Change[types.MediaFormat]

	Name      types.Change[string]

	AlbumId  types.Change[string]
	ArtistId types.Change[string]

	Duration types.Change[int64]
	Number   types.Change[sql.NullInt64]
	Year     types.Change[sql.NullInt64]

	Created types.Change[int64]
}

func (db DB) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
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

	ds := dialect.Update("tracks").
		Set(record).
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTrack(ctx context.Context, id string) error {
	query := dialect.Delete("tracks").
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) AddTagToTrack(ctx context.Context, tagSlug, trackId string) error {
	ds := dialect.Insert("tracks_tags").
		Rows(goqu.Record{
			"track_id": trackId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveAllTagsFromTrack(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_tags").
		Where(goqu.I("track_id").Eq(trackId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) RemoveAllTrackFeaturingArtists(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_featuring_artists").
		Where(goqu.I("tracks_featuring_artists.track_id").Eq(trackId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db DB) AddFeaturingArtistToTrack(ctx context.Context, trackId, artistId string) error {
	query := dialect.Insert("tracks_featuring_artists").
		Rows(goqu.Record{
			"track_id":  trackId,
			"artist_id": artistId,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveFeaturingArtistFromTrack(ctx context.Context, trackId, artistId string) error {
	query := goqu.Delete("tracks_featuring_artists").
		Where(
			goqu.And(
				goqu.I("tracks_featuring_artists.track_id").Eq(trackId),
				goqu.I("tracks_featuring_artists.artist_id").Eq(artistId),
			),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
