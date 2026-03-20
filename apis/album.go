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

type Album struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Year *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBAlbum(c pyrin.Context, album database.Album) Album {
	allArtists := make([]ArtistInfo, len(album.FeaturingArtists)+1)

	allArtists[0] = ArtistInfo{
		Id:   album.ArtistId,
		Name: album.ArtistName,
	}

	for i, v := range album.FeaturingArtists {
		allArtists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}
	return Album{
		Id:       album.Id,
		Name:     album.Name,
		Year:     utils.SqlNullToInt64Ptr(album.Year),
		CoverArt: ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
		Artists:  allArtists,
		Tags:     utils.SplitString(album.Tags.String),
		Created:  album.Created,
		Updated:  album.Updated,
	}
}

type GetAlbums struct {
	Page   types.Page `json:"page"`
	Albums []Album    `json:"albums"`
}

type GetAlbumById struct {
	Album Album `json:"album"`
}

type GetAlbumTracks struct {
	Tracks []Track `json:"tracks"`
}

type SearchAlbums struct {
	Albums []Album `json:"albums"`
}

func handleAlbumServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrAlbumServiceAlbumNotFound):
		return AlbumNotFound()
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

func InstallAlbumHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAlbums",
			Path:         "/albums",
			Method:       http.MethodGet,
			ResponseType: GetAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 100)
				filterParams := getFilterParams(q)

				albums, pageInfo, err := app.AlbumService().GetAlbums(
					ctx,
					service.GetAlbumsParams{
						Page:   pageParams,
						Filter: filterParams,
					},
				)
				if err != nil {
					return nil, handleAlbumServiceErrors(err)
				}

				res := GetAlbums{
					Page:   pageInfo,
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumById",
			Method:       http.MethodGet,
			Path:         "/albums/:id",
			ResponseType: GetAlbumById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				album, err := app.AlbumService().GetAlbumById(
					ctx,
					service.GetAlbumByIdParams{
						AlbumId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, handleAlbumServiceErrors(err)
				}

				return GetAlbumById{
					Album: ConvertDBAlbum(c, album),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumTracks",
			Method:       http.MethodGet,
			Path:         "/albums/:id/tracks",
			ResponseType: GetAlbumTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				tracks, err := app.AlbumService().GetAlbumTracks(
					ctx,
					service.GetAlbumTracksParams{
						AlbumId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, err
				}

				res := GetAlbumTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchAlbums",
			Path:         "/albums/search",
			Method:       http.MethodGet,
			ResponseType: SearchAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				albums, err := app.SearchService().SearchAlbums(ctx, q.Get("query"))
				if err != nil {
					return nil, err
				}

				res := SearchAlbums{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},
	)
}
