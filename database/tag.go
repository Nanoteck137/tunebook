package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	artistsTagsTbl = goqu.T("artists_tags")
	albumsTagsTbl  = goqu.T("albums_tags")
	tracksTagsTbl  = goqu.T("tracks_tags")
)

type Tag struct {
	Slug string `db:"slug"`
}

func TagQuery() *goqu.SelectDataset {
	query := dialect.From("tags").
		Select(
			"tags.slug",
		).
		Prepared(true)

	return query
}

func (db DB) GetAllTags(ctx context.Context) ([]Tag, error) {
	query := TagQuery()

	return Multiple[Tag](db, ctx, query)
}

func (db DB) GetTagBySlug(ctx context.Context, slug string) (Tag, error) {
	query := TagQuery().
		Where(goqu.I("tags.slug").Eq(slug))

	return Single[Tag](db, ctx, query)
}

func (db DB) CreateTag(ctx context.Context, slug string) error {
	query := dialect.Insert("tags").
		Rows(goqu.Record{
			"slug": slug,
		}).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) addTag(ctx context.Context, junctionTable exp.IdentifierExpression, idCol, id, tagSlug string) error {
	query := dialect.Insert(junctionTable).
		Rows(goqu.Record{
			idCol:     id,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) removeAllTags(ctx context.Context, junctionTable exp.IdentifierExpression, idCol, id string) error {
	query := dialect.Delete(junctionTable).
		Where(goqu.I(idCol).Eq(id))

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) AddTrackTag(ctx context.Context, tagSlug, trackId string) error {
	return db.addTag(ctx, tracksTagsTbl, "track_id", trackId, tagSlug)
}

func (db DB) RemoveAllTrackTags(ctx context.Context, trackId string) error {
	return db.removeAllTags(ctx, tracksTagsTbl, "track_id", trackId)
}

func (db DB) AddAlbumTag(ctx context.Context, tagSlug, albumId string) error {
	return db.addTag(ctx, albumsTagsTbl, "album_id", albumId, tagSlug)
}

func (db DB) RemoveAllAlbumTags(ctx context.Context, albumId string) error {
	return db.removeAllTags(ctx, albumsTagsTbl, "album_id", albumId)
}

func (db DB) AddArtistTag(ctx context.Context, tagSlug, artistId string) error {
	return db.addTag(ctx, artistsTagsTbl, "artist_id", artistId, tagSlug)
}

func (db DB) RemoveAllArtistTags(ctx context.Context, artistId string) error {
	return db.removeAllTags(ctx, artistsTagsTbl, "artist_id", artistId)
}

func AddTagsToQuery(
	query *goqu.SelectDataset,
	id exp.IdentifierExpression,
	objectTagTable exp.IdentifierExpression,
	idCol any,
) *goqu.SelectDataset {
	tagsTable := goqu.T("tags")
	tagsTableSlugCol := tagsTable.Col("slug")

	tableSlugCol := objectTagTable.Col("tag_slug")
	tableIdCol := objectTagTable.Col(idCol)

	subQuery := dialect.From(objectTagTable).
		Select(
			tableIdCol.As("id"),
			SqlGroupConcat(tagsTableSlugCol, ",").As("data"),
		).
		Join(
			tagsTable,
			goqu.On(tableSlugCol.Eq(tagsTableSlugCol)),
		).
		GroupBy(tableIdCol).
		As("tags")

	query = query.
		SelectAppend(
			goqu.I("tags.data").As("tags"),
		).
		LeftJoin(
			subQuery,
			goqu.On(goqu.I("tags.id").Eq(id)),
		)

	return query
}
