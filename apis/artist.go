package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type ArtistInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	CoverArt types.Images `json:"coverArt"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBArtist(c pyrin.Context, artist database.Artist) Artist {
	return Artist{
		Id:       artist.Id,
		Name:     artist.Name,
		CoverArt: ConvertArtistCoverURL(c, artist.Id, artist.CoverArt),
		Tags:     utils.SplitString(artist.Tags.String),
		Created:  artist.Created,
		Updated:  artist.Updated,
	}
}

type GetArtists struct {
	Page    types.Page `json:"page"`
	Artists []Artist   `json:"artists"`
}

type GetArtistById struct {
	Artist Artist `json:"artist"`
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type SearchArtists struct {
	Artists []Artist `json:"artists"`
}

func handleArtistServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrArtistServiceArtistNotFound):
		return ArtistNotFound()
	}

	var invalidFilter *service.InvalidFilterError
	if errors.As(err, &invalidFilter) {
		return InvalidFilter(errors.New(invalidFilter.Message))
	}

	var invalidSort *service.InvalidSortError
	if errors.As(err, &invalidSort) {
		return InvalidSort(errors.New(invalidSort.Message))
	}

	return err
}

func InstallArtistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetArtists",
			Method:       http.MethodGet,
			Path:         "/artists",
			ResponseType: GetArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				artists, page, err := app.ArtistService().GetArtists(
					ctx,
					service.GetArtistsParams{
						Page:   pageParams,
						Filter: filterParams,
					},
				)
				if err != nil {
					return nil, handleArtistServiceErrors(err)
				}

				res := GetArtists{
					Page:    page,
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = ConvertDBArtist(c, artist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetArtistById",
			Method:       http.MethodGet,
			Path:         "/artists/:id",
			ResponseType: GetArtistById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				artist, err := app.ArtistService().GetArtistById(
					ctx,
					service.GetArtistByIdParams{
						ArtistId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, handleArtistServiceErrors(err)
				}

				return GetArtistById{
					Artist: ConvertDBArtist(c, artist),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetArtistAlbums",
			Method:       http.MethodGet,
			Path:         "/artists/:id/albums",
			ResponseType: GetArtistAlbumsById{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				albums, err := app.ArtistService().GetArtistAlbums(
					ctx,
					service.GetArtistAlbumsParams{
						ArtistId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, handleArtistServiceErrors(err)
				}

				res := GetArtistAlbumsById{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		// TODO(patrik): Move from /artists/search to /search/artists
		pyrin.ApiHandler{
			Name:         "SearchArtists",
			Method:       http.MethodGet,
			Path:         "/artists/search",
			ResponseType: SearchArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				artists, err := app.SearchService().SearchArtists(
					ctx,
					q.Get("query"),
				)
				if err != nil {
					return nil, err
				}

				res := SearchArtists{
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = ConvertDBArtist(c, artist)
				}

				return res, nil
			},
		},
	)
}
