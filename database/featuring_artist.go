package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type FeaturingArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var (
	tracksFeaturingArtistsTbl = goqu.T("tracks_featuring_artists")
	albumsFeaturingArtistsTbl = goqu.T("albums_featuring_artists")
)

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
					artistsTbl.Col("id"),
					"name",
					artistsTbl.Col("name"),
				),
			).As("artists"),
		).
		Join(
			artistsTbl,
			goqu.On(tbl.Col("artist_id").Eq(artistsTbl.Col("id"))),
		).
		GroupBy(tbl.Col(idColName))
}

func AddFeaturingArtistsToQuery(
	query *goqu.SelectDataset,
	id exp.IdentifierExpression,
	objectTable exp.IdentifierExpression,
	idCol any,
) *goqu.SelectDataset {
	tableIdCol := objectTable.Col(idCol)

	subQuery := dialect.From(objectTable).
		Select(
			tableIdCol.As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",
					"id",
					artistsTbl.Col("id"),
					"name",
					artistsTbl.Col("name"),
				),
			).As("data"),
		).
		Join(
			artistsTbl,
			goqu.On(objectTable.Col("artist_id").Eq(artistsTbl.Col("id"))),
		).
		GroupBy(tableIdCol).
		As("featuring_artists")

	query = query.
		SelectAppend(
			goqu.I("featuring_artists.data").As("featuring_artists"),
		).
		LeftJoin(
			subQuery,
			goqu.On(goqu.I("featuring_artists.id").Eq(id)),
		)

	return query
}

func (db DB) addFeaturingArtist(
	ctx context.Context, 
	junctionTable exp.IdentifierExpression, 
	idCol, id, artistId string,
) error {
	query := dialect.Insert(junctionTable).
		Rows(goqu.Record{
			idCol:      id,
			"artist_id": artistId,
		})

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) removeAllFeaturingArtists(
	ctx context.Context, 
	junctionTable exp.IdentifierExpression, 
	idCol, id string,
) error {
	query := dialect.Delete(junctionTable).
		Where(goqu.I(idCol).Eq(id))

	_, err := db.Exec(ctx, query)
	return err
}

func (db DB) AddFeaturingArtistToTrack(
	ctx context.Context,
	trackId, artistId string,
) error {
	return db.addFeaturingArtist(
		ctx, tracksFeaturingArtistsTbl, "track_id", trackId, artistId)
}

func (db DB) RemoveAllTrackFeaturingArtists(
	ctx context.Context,
	trackId string,
) error {
	return db.removeAllFeaturingArtists(
		ctx, tracksFeaturingArtistsTbl, "track_id", trackId)
}

func (db DB) AddFeaturingArtistToAlbum(
	ctx context.Context,
	albumId, artistId string,
) error {
	return db.addFeaturingArtist(
		ctx, albumsFeaturingArtistsTbl, "album_id", albumId, artistId)
}

func (db DB) RemoveAllAlbumFeaturingArtists(
	ctx context.Context,
	albumId string,
) error {
	return db.removeAllFeaturingArtists(
		ctx, albumsFeaturingArtistsTbl, "album_id", albumId)
}
