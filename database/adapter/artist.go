package adapter

import (
	"go/ast"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook/tools/filter"
)

var _ filter.ResolverAdapter = (*TrackResolverAdapter)(nil)

type ArtistResolverAdapter struct{}

func (a *ArtistResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "artists.name", filter.SortTypeAsc
}

func (a *ArtistResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.name",
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "artists.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "artists.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *ArtistResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return slug.Make(name), true
	}

	return "", false
}

func (a *ArtistResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "artists_tags",
			SelectName: "artist_id",
			WhereName:  "tag_slug",
		}, true
	}

	return filter.Table{}, false
}

func (a *ArtistResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "artists.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
