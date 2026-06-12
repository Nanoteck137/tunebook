package adapter

import (
	"go/ast"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook/tools/filter"
)

var _ filter.ResolverAdapter = (*UserTrackHistoryResolverAdapter)(nil)

type UserTrackHistoryResolverAdapter struct{}

func (a *UserTrackHistoryResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "user_track_history.listened_at", filter.SortTypeDesc
}

func (a *UserTrackHistoryResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_track_history.id",
		}, true
	case "userId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_track_history.user_id",
		}, true
	case "trackId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_track_history.track_id",
		}, true
	case "listenedAt":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "user_track_history.listened_at",
		}, true
	case "playbackType":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_track_history.playback_type",
		}, true
	case "status":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_track_history.status",
		}, true
	case "trackName":
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
			Name: "user_track_history.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "user_track_history.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *UserTrackHistoryResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return slug.Make(name), true
	case "featuringArtists":
		return name, true
	}

	return "", false
}

func (a *UserTrackHistoryResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
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

func (a *UserTrackHistoryResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "tracks.id", args)
	case "hasFeaturingArtist":
		return resolver.InTable(name, "featuringArtists", "tracks.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
