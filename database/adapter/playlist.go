package adapter

import (
	"go/ast"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook/tools/filter"
)

var _ filter.ResolverAdapter = (*PlaylistResolverAdapter)(nil)

type PlaylistResolverAdapter struct{}

func (a *PlaylistResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "playlists.name", filter.SortTypeAsc
}

func (a *PlaylistResolverAdapter) ResolveVariableName(
	name string,
) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlists.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlists.name",
		}, true
	case "ownerId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlists.owner_id",
		}, true
	case "ownerDisplayName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "owner.display_name",
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "playlists.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "playlists.updated",
		}, true
	case "trackCount":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "track_count.data",
		}, true
	}

	return filter.Name{}, false
}

func (a *PlaylistResolverAdapter) ResolveNameToId(
	typ, name string,
) (string, bool) {
	return "", false
}

func (a *PlaylistResolverAdapter) ResolveTable(
	typ string,
) (filter.Table, bool) {
	return filter.Table{}, false
}

func (a *PlaylistResolverAdapter) ResolveFunctionCall(
	resolver *filter.Resolver,
	name string,
	args []ast.Expr,
) (filter.FilterExpr, error) {
	return nil, filter.UnknownFunction(name)
}

var _ filter.ResolverAdapter = (*PlaylistItemResolverAdapter)(nil)

type PlaylistItemResolverAdapter struct{}

func (a *PlaylistItemResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "playlist_items.order_num", filter.SortTypeAsc
}

func (a *PlaylistItemResolverAdapter) ResolveVariableName(
	name string,
) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlist_items.id",
		}, true
	case "playlistId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlist_items.playlist_id",
		}, true
	case "trackId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlist_items.track_id",
		}, true
	case "order":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "playlist_items.order_num",
		}, true
	}

	return filter.Name{}, false
}

func (a *PlaylistItemResolverAdapter) ResolveNameToId(
	typ, name string,
) (string, bool) {
	return "", false
}

func (a *PlaylistItemResolverAdapter) ResolveTable(
	typ string,
) (filter.Table, bool) {
	return filter.Table{}, false
}

func (a *PlaylistItemResolverAdapter) ResolveFunctionCall(
	resolver *filter.Resolver,
	name string,
	args []ast.Expr,
) (filter.FilterExpr, error) {
	return nil, filter.UnknownFunction(name)
}

var _ filter.ResolverAdapter = (*PlaylistTrackResolverAdapter)(nil)

type PlaylistTrackResolverAdapter struct{}

func (a *PlaylistTrackResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "playlist_items.position", filter.SortTypeAsc
}

func (a *PlaylistTrackResolverAdapter) ResolveVariableName(
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
	case "playlistId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlist_items.playlist_id",
		}, true
	case "trackId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "playlist_items.track_id",
		}, true
	case "position":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "playlist_items.position",
		}, true
	}

	return filter.Name{}, false
}

func (a *PlaylistTrackResolverAdapter) ResolveNameToId(
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

func (a *PlaylistTrackResolverAdapter) ResolveTable(
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

func (a *PlaylistTrackResolverAdapter) ResolveFunctionCall(
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
