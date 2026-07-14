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
	createPlaylistId = createIdGenerator(16)

	playlistsTbl = goqu.T("playlists")

	playlistSchema = PlaylistSchema()
)

type Playlist struct {
	Id       string         `db:"id"`
	Name     string         `db:"name"`
	CoverArt sql.NullString `db:"cover_art"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	OwnerDisplayName string         `db:"owner_display_name"`
	OwnerPicture     sql.NullString `db:"owner_picture"`

	TrackCount int64 `db:"track_count"`
}

func PlaylistSchema() *schema.Schema {
	return schema.New().
		AddField("id", query.TypeString, schema.Column("playlists.id")).
		AddField("name", query.TypeString, schema.Column("playlists.name")).
		AddField(
			"ownerId",
			query.TypeString,
			schema.Column("playlists.owner_id"),
		).
		AddField(
			"coverArt",
			query.TypeString,
			schema.Column("playlists.cover_art"),
			schema.Nullable(),
		).
		AddField(
			"ownerDisplayName",
			query.TypeString,
			schema.Column("owner.display_name"),
		).
		// TODO(patrik): Is this correct? I need to test this later
		AddField("trackCount", query.TypeInt, schema.Column("track_count.data")).
		AddField("created", query.TypeInt, schema.Column("playlists.created")).
		AddField("updated", query.TypeInt, schema.Column("playlists.updated")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "name"},
				Dir:   query.DirAsc,
			},
		)
}

func PlaylistQuery() *goqu.SelectDataset {
	idCol := playlistsTbl.Col("id")

	trackCountQuery := dialect.From(playlistItemsTbl).
		Select(
			playlistItemsTbl.Col("playlist_id").As("id"),
			goqu.COUNT(playlistItemsTbl.Col("track_id")).As("data"),
		).
		GroupBy(playlistItemsTbl.Col("playlist_id"))

	query := dialect.From(playlistsTbl).
		Select(
			idCol,

			playlistsTbl.Col("name"),
			playlistsTbl.Col("cover_art"),

			playlistsTbl.Col("owner_id"),

			playlistsTbl.Col("created"),
			playlistsTbl.Col("updated"),

			goqu.I("owner.display_name").As("owner_display_name"),
			goqu.I("owner.picture").As("owner_picture"),

			goqu.COALESCE(goqu.I("track_count.data"), 0).As("track_count"),
		).
		Join(
			UserQuery().As("owner"),
			goqu.On(playlistsTbl.Col("owner_id").Eq(goqu.I("owner.id"))),
		).
		LeftJoin(
			trackCountQuery.As("track_count"),
			goqu.On(idCol.Eq(goqu.I("track_count.id"))),
		)

	return query
}

type CreatePlaylistParams struct {
	Id       string
	Name     string
	CoverArt sql.NullString

	OwnerId string

	Created int64
	Updated int64
}

func (db DB) CreatePlaylist(
	ctx context.Context,
	params CreatePlaylistParams,
) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = createPlaylistId()
	}

	query := dialect.Insert(playlistsTbl).Rows(goqu.Record{
		"id":        params.Id,
		"name":      params.Name,
		"cover_art": params.CoverArt,

		"owner_id": params.OwnerId,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return "", err
	}

	return params.Id, nil
}

type PlaylistChanges struct {
	Name Change[string]

	OwnerId Change[string]

	CoverArt Change[sql.NullString]

	Created Change[int64]
}

func (db DB) UpdatePlaylist(
	ctx context.Context,
	playlistId string,
	changes PlaylistChanges,
) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "owner_id", changes.OwnerId)

	addToRecord(record, "cover_art", changes.CoverArt)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update(playlistsTbl).
		Set(record).
		Where(playlistsTbl.Col("id").Eq(playlistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeletePlaylist(ctx context.Context, playlistId string) error {
	query := dialect.Delete(playlistsTbl).
		Where(playlistsTbl.Col("id").Eq(playlistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetAllPlaylistIds(ctx context.Context) ([]string, error) {
	query := dialect.From(playlistsTbl).
		Select(playlistsTbl.Col("id"))

	return Multiple[string](db, ctx, query)
}

type GetPlaylistsParams struct {
	Page  types.PageParams
	Query types.QueryParams
}

func (db DB) GetPlaylists(
	ctx context.Context,
	params GetPlaylistsParams,
) ([]Playlist, types.Page, error) {
	query := PlaylistQuery()

	var err error

	query, err = ApplyQuery(query, playlistSchema, params.Query)
	if err != nil {
		return nil, types.Page{}, err
	}

	page, err := buildPage(ctx, db, params.Page, query, playlistsTbl.Col("id"))
	if err != nil {
		return nil, types.Page{}, err
	}

	query = applyPageParams(params.Page, query)

	items, err := Multiple[Playlist](db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetPlaylistsByIds(
	ctx context.Context,
	ids []string,
) ([]Playlist, error) {
	query := PlaylistQuery().Where(playlistsTbl.Col("id").In(ids))

	return Multiple[Playlist](db, ctx, query)
}

func (db DB) GetPlaylistById(
	ctx context.Context,
	playlistId string,
) (Playlist, error) {
	query := PlaylistQuery().
		Where(playlistsTbl.Col("id").Eq(playlistId))

	return Single[Playlist](db, ctx, query)
}

func (db DB) GetUserPlaylistCount(
	ctx context.Context,
	userId string,
) (int, error) {
	query := dialect.From(playlistsTbl).
		Select(
			goqu.COUNT(playlistsTbl.Col("id")),
		).
		Where(playlistsTbl.Col("owner_id").Eq(userId)).
		GroupBy(playlistsTbl.Col("owner_id"))

	return Single[int](db, ctx, query)
}
