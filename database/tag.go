package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
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
