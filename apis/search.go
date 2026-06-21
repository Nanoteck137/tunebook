package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

type SearchArtists struct {
	Page    types.Page `json:"page"`
	Artists []Artist   `json:"artists"`
}

type SearchAlbums struct {
	Page   types.Page `json:"page"`
	Albums []Album    `json:"albums"`
}

type SearchTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type SearchPlaylists struct {
	Page      types.Page `json:"page"`
	Playlists []Playlist `json:"playlists"`
}

type SearchUsers struct {
	Page  types.Page `json:"page"`
	Users []UserData `json:"users"`
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

				pageParams := getPageParams(q, 5)
				artists, page, err := app.SearchService().SearchArtists(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Page:  pageParams,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchArtists{
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
			Name:         "SearchAlbums",
			Path:         "/search/albums",
			Method:       http.MethodGet,
			ResponseType: SearchAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 5)
				albums, page, err := app.SearchService().SearchAlbums(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Page:  pageParams,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchAlbums{
					Page:   page,
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

				pageParams := getPageParams(q, 5)
				tracks, page, err := app.SearchService().SearchTracks(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Page:  pageParams,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchTracks{
					Page:   page,
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

				pageParams := getPageParams(q, 5)
				playlists, page, err := app.SearchService().SearchPlaylists(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Page:  pageParams,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchPlaylists{
					Page:      page,
					Playlists: make([]Playlist, len(playlists)),
				}

				for i, playlist := range playlists {
					res.Playlists[i] = ConvertDBPlaylist(c, playlist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchUsers",
			Method:       http.MethodGet,
			Path:         "/search/users",
			ResponseType: SearchUsers{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 5)
				users, page, err := app.SearchService().SearchUsers(
					ctx,
					service.SearchParams{
						Query: q.Get("query"),
						Page:  pageParams,
					},
				)
				if err != nil {
					return nil, err
				}

				res := SearchUsers{
					Page:  page,
					Users: make([]UserData, len(users)),
				}

				for i, user := range users {
					res.Users[i] = ConvertDBUser(c, user)
				}

				return res, nil
			},
		},
	)
}
