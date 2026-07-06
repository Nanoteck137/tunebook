package adapter

import (
	"go/ast"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook/tools/filter"
)

var _ filter.ResolverAdapter = (*UserFavoriteTrackResolverAdapter)(nil)

type UserFavoriteTrackResolverAdapter struct{}

func (a *UserFavoriteTrackResolverAdapter) DefaultSort(
) (string, filter.SortType) {
	return "user_favorites.added", filter.SortTypeDesc
}

func (a *UserFavoriteTrackResolverAdapter) ResolveVariableName(
	name string,
) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.name",
		}, true
	case "number":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.number",
			Nullable: true,
		}, true
	case "duration":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.duration",
			Nullable: true,
		}, true
	case "year":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.year",
			Nullable: true,
		}, true
	case "albumId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.album_id",
		}, true
	case "artistId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.artist_id",
		}, true
	case "albumName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.name",
		}, true
	case "artistName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.name",
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "tracks.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "tracks.updated",
		}, true
	case "userId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_favorites.user_id",
		}, true
	case "trackId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_favorites.track_id",
		}, true
	case "added":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "user_favorites.added",
		}, true
	}

	return filter.Name{}, false
}

func (a *UserFavoriteTrackResolverAdapter) ResolveNameToId(
	typ, name string,
) (string, bool) {
	switch typ {
	case "tags":
		return slug.Make(name), true
	case "featuringArtists":
		return name, true
	}

	return "", false
}

func (a *UserFavoriteTrackResolverAdapter) ResolveTable(
	typ string,
) (filter.Table, bool) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "tracks_tags",
			SelectName: "track_id",
			WhereName:  "tag_slug",
		}, true
	case "featuringArtists":
		return filter.Table{
			Name:       "tracks_featuring_artists",
			SelectName: "track_id",
			WhereName:  "artist_id",
		}, true
	}

	return filter.Table{}, false
}

func (a *UserFavoriteTrackResolverAdapter) ResolveFunctionCall(
	resolver *filter.Resolver, 
	name string, 
	args []ast.Expr,
) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "tracks.id", args)
	case "hasFeaturingArtist":
		return resolver.InTable(name, "featuringArtists", "tracks.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
