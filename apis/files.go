package apis

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
)

func InstallFilesHandlers(app core.App, g pyrin.Group) {
	g.Register(
		pyrin.NormalHandler{
			Name:   "GetAlbumImage",
			Method: http.MethodGet,
			Path:   "/albums/images/:albumId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")
				image := c.Param("image")

				ext := path.Ext(image)
				name := strings.TrimRight(image, ext)

				imageType, ok := app.ImageService().GetImageTypeFromExt(ext)
				if !ok {
					// TODO(patrik): Better error
					return errors.New("unsupported image ext")
				}

				ctx := c.Request().Context()

				p, err := app.ImageService().GetAlbumImage(ctx, albumId, name, imageType)
				if err != nil {
					return err
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetArtistImage",
			Method: http.MethodGet,
			Path:   "/artists/images/:artistId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				artistId := c.Param("artistId")
				image := c.Param("image")

				ext := path.Ext(image)
				name := strings.TrimRight(image, ext)

				imageType, ok := app.ImageService().GetImageTypeFromExt(ext)
				if !ok {
					// TODO(patrik): Better error
					return errors.New("unsupported image ext")
				}

				ctx := c.Request().Context()

				p, err := app.ImageService().GetArtistImage(ctx, artistId, name, imageType)
				if err != nil {
					return err
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetPlaylistImage",
			Method: http.MethodGet,
			Path:   "/playlists/images/:playlistId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				playlistId := c.Param("playlistId")
				image := c.Param("image")

				ext := path.Ext(image)
				name := strings.TrimRight(image, ext)

				imageType, ok := app.ImageService().GetImageTypeFromExt(ext)
				if !ok {
					// TODO(patrik): Better error
					return errors.New("unsupported image ext")
				}

				ctx := c.Request().Context()

				p, err := app.ImageService().GetPlaylistImage(ctx, playlistId, name, imageType)
				if err != nil {
					return err
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},

		pyrin.NormalHandler{
			Name:   "GetUserImage",
			Method: http.MethodGet,
			Path:   "/users/images/:userId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				userId := c.Param("userId")
				image := c.Param("image")

				ext := path.Ext(image)
				name := strings.TrimRight(image, ext)

				imageType, ok := app.ImageService().GetImageTypeFromExt(ext)
				if !ok {
					// TODO(patrik): Better error
					return errors.New("unsupported image ext")
				}

				ctx := c.Request().Context()

				p, err := app.ImageService().GetUserImage(ctx, userId, name, imageType)
				if err != nil {
					return err
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},
	)
}
