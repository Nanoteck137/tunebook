package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/types"
)

var (
	createTrackId = createIdGenerator(32)

	tracksTbl = goqu.T("tracks")

	trackSchema = TrackSchema()
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

	// This can be set elsewhere just to give tracks a order number,
	// so for albums it is just set to the number but for
	// playlists it set to the playlist item position
	Order *int
}

func TrackSchema() *schema.Schema {
	return schema.New().
		AddField("id", query.TypeString, schema.Column("tracks.id")).
		AddField("name", query.TypeString, schema.Column("tracks.name")).
		AddField("number", query.TypeInt, schema.Column("tracks.number"), schema.Nullable()).
		AddField("duration", query.TypeInt, schema.Column("tracks.duration"), schema.Nullable()).
		AddField("year", query.TypeInt, schema.Column("tracks.year"), schema.Nullable()).
		AddField("albumId", query.TypeString, schema.Column("tracks.album_id")).
		AddField("artistId", query.TypeString, schema.Column("tracks.artist_id")).
		AddField("albumName", query.TypeString, schema.Column("albums.name")).
		AddField("artistName", query.TypeString, schema.Column("artists.name")).
		AddField("tags", query.TypeRelation, schema.Relation("tracks_tags", "track_id", "tag_slug", query.TypeString)).
		AddField("featuringArtist", query.TypeRelation, schema.Relation("tracks_featuring_artists", "track_id", "artist_id", query.TypeString)).
		AddField("created", query.TypeInt, schema.Column("tracks.created")).
		AddField("updated", query.TypeInt, schema.Column("tracks.updated")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "name"},
				Dir:   query.DirAsc,
			},
		)
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

			albumsTbl.Col("name").As("album_name"),
			albumsTbl.Col("cover_art").As("album_cover_art"),

			artistsTbl.Col("name").As("artist_name"),
		).
		Join(
			albumsTbl,
			goqu.On(tracksTbl.Col("album_id").Eq(albumsTbl.Col("id"))),
		).
		Join(
			artistsTbl,
			goqu.On(tracksTbl.Col("artist_id").Eq(artistsTbl.Col("id"))),
		)

	query = AddTagsToQuery(query, idCol, tracksTagsTbl, "track_id")
	query = AddFeaturingArtistsToQuery(
		query, idCol, tracksFeaturingArtistsTbl, "track_id")

	return query
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

	query := dialect.Insert(tracksTbl).Rows(goqu.Record{
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

	query := dialect.Update(tracksTbl).
		Set(record).
		Where(tracksTbl.Col("id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTrack(ctx context.Context, id string) error {
	query := dialect.Delete(tracksTbl).
		Where(tracksTbl.Col("id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(tracksTbl.Col("id").Eq(id))

	return Single[Track](db, ctx, query)
}

func (db DB) GetTracksByIds(
	ctx context.Context,
	ids []string,
) ([]Track, error) {
	query := TrackQuery().Where(tracksTbl.Col("id").In(ids))

	return Multiple[Track](db, ctx, query)
}

func (db DB) GetAllTrackIds(ctx context.Context) ([]string, error) {
	query := dialect.From("tracks").
		Select(tracksTbl.Col("id"))

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByFilterParams struct {
	Filter types.FilterParams
}

func (db DB) GetTrackIdsByFilter(
	ctx context.Context,
	params GetTrackIdsByFilterParams,
) ([]string, error) {
	query := TrackQuery().Select(tracksTbl.Col("id"))

	query, err := ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter.Filter,
		Sort:   params.Filter.Sort,
	})
	if err != nil {
		return nil, err
	}

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

	// Use the new query system
	query, err = ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter.Filter,
		Sort:   params.Filter.Sort,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, tracksTbl.Col("id"))
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

type GetTrackIdsByAlbumParams struct {
	Filter string
}

func (db DB) GetTrackIdsByAlbum(
	ctx context.Context,
	albumId string,
	params GetTrackIdsByAlbumParams,
) ([]string, error) {
	var err error

	query := TrackQuery().Select(tracksTbl.Col("id"))

	// Apply the custom album filter first
	query = query.Where(tracksTbl.Col("album_id").Eq(albumId))

	// Then apply the user-provided filter and sort
	query, err = ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter,
	})
	if err != nil {
		return nil, err
	}

	// Apply default album track ordering
	query = query.Order(
		tracksTbl.Col("number").Asc().NullsLast(),
		tracksTbl.Col("name").Asc(),
	)

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByArtistParams struct {
	Filter string
}

func (db DB) GetTrackIdsByArtist(
	ctx context.Context,
	artistId string,
	params GetTrackIdsByArtistParams,
) ([]string, error) {
	var err error

	query := TrackQuery().Select(tracksTbl.Col("id"))

	// Apply the custom artist filter first
	query = query.Where(tracksTbl.Col("artist_id").Eq(artistId))

	// Then apply the user-provided filter and sort
	query, err = ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter,
	})
	if err != nil {
		return nil, err
	}

	// Apply default ordering
	query = query.Order(tracksTbl.Col("name").Asc())

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByPlaylistParams struct {
	Filter string
}

func (db DB) GetTrackIdsByPlaylist(
	ctx context.Context,
	playlistId string,
	params GetTrackIdsByPlaylistParams,
) ([]string, error) {
	var err error

	query := TrackQuery().
		Select(tracksTbl.Col("id")).
		Join(
			playlistItemsTbl,
			goqu.On(tracksTbl.Col("id").Eq(playlistItemsTbl.Col("track_id"))),
		)

	// Apply the custom playlist filter first
	query = query.Where(playlistItemsTbl.Col("playlist_id").Eq(playlistId))

	// Then apply the user-provided filter and sort
	query, err = ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter,
	})
	if err != nil {
		return nil, err
	}

	// Apply default ordering
	query = query.Order(playlistItemsTbl.Col("position").Asc())

	return Multiple[string](db, ctx, query)
}

type GetTrackIdsByUserFavoritesParams struct {
	Filter string
}

func (db DB) GetTrackIdsByUserFavorites(
	ctx context.Context,
	userId string,
	params GetTrackIdsByUserFavoritesParams,
) ([]string, error) {
	var err error

	query := TrackQuery().
		Select(tracksTbl.Col("id")).
		Join(
			userFavoritesTbl,
			goqu.On(tracksTbl.Col("id").Eq(userFavoritesTbl.Col("track_id"))),
		)

	// Apply the custom user favorites filter first
	query = query.Where(userFavoritesTbl.Col("user_id").Eq(userId))

	// Then apply the user-provided filter and sort
	query, err = ApplyQuery(query, trackSchema, QueryParams{
		Filter: params.Filter,
	})
	if err != nil {
		return nil, err
	}

	// Apply default ordering
	query = query.Order(userFavoritesTbl.Col("added").Desc())

	return Multiple[string](db, ctx, query)
}

func (db DB) GetTracksByAlbum(
	ctx context.Context,
	albumId string,
) ([]Track, error) {
	query := TrackQuery().
		Where(tracksTbl.Col("album_id").Eq(albumId)).
		Order(
			tracksTbl.Col("number").Asc().NullsLast(),
			tracksTbl.Col("name").Asc(),
		)

	return Multiple[Track](db, ctx, query)
}
