package apis

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

func InstallFilesHandlers(app core.App, g pyrin.Group) {
	g.Register(
		pyrin.NormalHandler{
			Name:   "GetAlbumImage",
			Method: http.MethodGet,
			Path:   "/albums/images/:albumId",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")

				q := c.Request().URL.Query()

				size, _ := strconv.Atoi(q.Get("size"))

				ctx := c.Request().Context()

				p, err := app.AlbumService().GetAlbumImage(
					ctx,
					service.GetAlbumImageParams{
						AlbumId:     albumId,
						Size:        size,
						ImageFormat: types.ImageFormatPng,
					},
				)
				if err != nil {
					return handleAlbumServiceErrors(err)
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetArtistImage",
			Method: http.MethodGet,
			Path:   "/artists/images/:artistId",
			HandlerFunc: func(c pyrin.Context) error {
				artistId := c.Param("artistId")

				q := c.Request().URL.Query()
				size, _ := strconv.Atoi(q.Get("size"))

				ctx := c.Request().Context()

				p, err := app.ArtistService().GetArtistImage(
					ctx,
					service.GetArtistImageParams{
						ArtistId:    artistId,
						Size:        size,
						ImageFormat: types.ImageFormatPng,
					},
				)
				if err != nil {
					return handleArtistServiceErrors(err)
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetPlaylistImage",
			Method: http.MethodGet,
			Path:   "/playlists/images/:playlistId",
			HandlerFunc: func(c pyrin.Context) error {
				playlistId := c.Param("playlistId")

				q := c.Request().URL.Query()
				size, _ := strconv.Atoi(q.Get("size"))

				ctx := c.Request().Context()

				p, err := app.PlaylistService().GetPlaylistImage(
					ctx,
					service.GetPlaylistImageParams{
						PlaylistId:  playlistId,
						Size:        size,
						ImageFormat: types.ImageFormatPng,
					},
				)
				if err != nil {
					return handlePlaylistServiceErrors(err)
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetUserImage",
			Method: http.MethodGet,
			Path:   "/users/images/:userId",
			HandlerFunc: func(c pyrin.Context) error {
				userId := c.Param("userId")

				q := c.Request().URL.Query()
				size, _ := strconv.Atoi(q.Get("size"))

				ctx := c.Request().Context()

				p, err := app.UserService().GetUserImage(
					ctx,
					service.GetUserImageParams{
						UserId:      userId,
						Size:        size,
						ImageFormat: types.ImageFormatPng,
					},
				)
				if err != nil {
					return handleUserServiceErrors(err)
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},
	)
}
