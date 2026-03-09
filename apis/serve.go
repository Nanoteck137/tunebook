package apis

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)

	g = router.Group("")
	g.Register(
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/static/*",
			HandlerFunc: func(c pyrin.Context) error {
				// TODO(patrik): Fix this
				f := os.DirFS("./render/static")
				fs := http.StripPrefix("/static", http.FileServerFS(f))

				fs.ServeHTTP(c.Response(), c.Request())

				return nil
			},
		},

		// TODO(patrik): Fix this
		pyrin.SpaHandler(os.DirFS("./result"), "index.html"),
	)

	g = router.Group("/media")
	g.Register(
		pyrin.NormalHandler{
			Name:   "StreamTrack",
			Method: http.MethodGet,
			Path:   "/tracks/:trackId/stream",
			HandlerFunc: func(c pyrin.Context) error {
				trackId := c.Param("trackId")

				query := c.Request().URL.Query()

				// mode=raw,smart
				// device
				// quality
				// policy
				// format
				// bitrate

				bitrate, _ := strconv.ParseInt(query.Get("bitrate"), 10, 32) 

				filename, err := app.MediaService().GetTrackStream(trackId, service.MediaStreamOptions{
					Device:  service.Device(query.Get("device")),
					Policy:  service.Policy(query.Get("policy")),
					Quality: service.Quality(query.Get("quality")),
					Format:  types.MediaFormat(query.Get("format")),
					Bitrate: int(bitrate),
				})
				if err != nil {
					// TODO(patrik): Better error handling
					return err
				}

				f := os.DirFS(path.Dir(filename))
				return pyrin.ServeFile(c, f, path.Base(filename))
			},
		},
	)

	g = router.Group("/files")
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
	)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: dwebble.AppName,
		ErrorCallback: func(err error) {
			// TODO(patrik): Handle this better
			slog.Error("API Error", "err", err)
		},
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	return s, nil
}
