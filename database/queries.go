package database

import "github.com/doug-martin/goqu/v9"

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
