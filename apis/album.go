package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

type Album struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Year      *int64 `json:"year"`
	AlbumType string `json:"albumType"`

	CoverArt types.Images `json:"coverArt"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBAlbum(c pyrin.Context, album database.Album) Album {
	allArtists := make([]ArtistInfo, len(album.FeaturingArtists.Data)+1)

	allArtists[0] = ArtistInfo{
		Id:   album.ArtistId,
		Name: album.ArtistName,
	}

	for i, v := range album.FeaturingArtists.Data {
		allArtists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}

	return Album{
		Id:        album.Id,
		Name:      album.Name,
		Year:      utils.SqlNullToInt64Ptr(album.Year),
		AlbumType: string(album.AlbumType),
		CoverArt:  ConvertAlbumCoverURL(c, album.Id),
		Artists:   allArtists,
		Tags:      utils.SplitTagString(album.Tags.String),
		Created:   formatTime(album.Created),
		Updated:   formatTime(album.Updated),
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

func handleAlbumServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrAlbumServiceAlbumNotFound):
		return AlbumNotFound()
	}

	var queryErr *database.QueryError
	if errors.As(err, &queryErr) {
		var filterErr, sortErr error
		if queryErr.Filter != nil {
			filterErr = queryErr.Filter
		}
		if queryErr.Sort != nil {
			sortErr = queryErr.Sort
		}
		return QueryError(filterErr, sortErr)
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
				queryParams := getQueryParams(q)

				albums, pageInfo, err := app.AlbumService().GetAlbums(
					ctx,
					service.GetAlbumsParams{
						Page:  pageParams,
						Query: queryParams,
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
	)
}
