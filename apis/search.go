package apis

import (
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/pyrin"
)

type SearchArtists struct {
	Artists []Artist `json:"artists"`
}

type SearchAlbums struct {
	Albums []Album `json:"albums"`
}

type SearchTracks struct {
	Tracks []Track `json:"tracks"`
}

type SearchPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

func InstallSearchHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "SearchArtists",
			Method:       http.MethodGet,
			Path:         "/search/artists",
			ResponseType: SearchArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				artists, err := app.SearchService().SearchArtists(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Limit: 5,
					},
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

		pyrin.ApiHandler{
			Name:         "SearchAlbums",
			Path:         "/search/albums",
			Method:       http.MethodGet,
			ResponseType: SearchAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				albums, err := app.SearchService().SearchAlbums(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Limit: 5,
					},
				)
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

		pyrin.ApiHandler{
			Name:         "SearchTracks",
			Method:       http.MethodGet,
			Path:         "/search/tracks",
			ResponseType: SearchTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				tracks, err := app.SearchService().SearchTracks(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Limit: 5,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchPlaylists",
			Method:       http.MethodGet,
			Path:         "/search/playlists",
			ResponseType: SearchPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				playlists, err := app.SearchService().SearchPlaylists(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Limit: 5,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchPlaylists{
					Playlists: make([]Playlist, len(playlists)),
				}

				for i, playlist := range playlists {
					res.Playlists[i] = ConvertDBPlaylist(c, playlist)
				}

				return res, nil
			},
		},
	)
}
