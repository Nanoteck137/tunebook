package apis

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
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

		pyrin.SpaHandler(os.DirFS("./result"), "index.html"),
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
		pyrin.NormalHandler{
			Name:   "GetTrackFile",
			Method: http.MethodGet,
			Path:   "/tracks/:trackId/:file",
			HandlerFunc: func(c pyrin.Context) error {
				trackId := c.Param("trackId")
				file := c.Param("file")

				fileExt := path.Ext(file)

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, trackId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return pyrin.NoContentNotFound()
					}

					return err
				}

				mediaType := types.GetMediaTypeFromExt(fileExt)

				if mediaType == types.MediaTypeUnknown {
					return pyrin.NoContentNotFound()
				}

				// Return the original file if the filename matches the
				// one stored inside the track
				if track.MediaType == mediaType {
					d := path.Dir(track.Filename)
					filename := path.Base(track.Filename)

					f := os.DirFS(d)
					return pyrin.ServeFile(c, f, filename)
				}

				// Here we need to start transcoding the original track
				// media to the requested format

				cacheDir := app.WorkDir().Cache()
				trackCache := cacheDir.Track(track.Id)

				// Make sure that the cache directory is setup
				dirs := []string{
					cacheDir.String(),
					cacheDir.Tracks(),
					trackCache,
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return err
					}
				}

				switch mediaType {
				case types.MediaTypeMp3:
					name := "track.mp3"
					p := path.Join(trackCache, name)

					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							cmd := exec.Command("ffmpeg", "-i", track.Filename, "-b:a", "320k", p)
							err := cmd.Run()
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}

					f := os.DirFS(trackCache)
					return pyrin.ServeFile(c, f, name)
				case types.MediaTypeOggOpus:
					name := "track.opus"
					p := path.Join(trackCache, name)

					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							cmd := exec.Command("ffmpeg", "-i", track.Filename, "-b:a", "96k", p)
							err := cmd.Run()
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}

					f := os.DirFS(trackCache)
					return pyrin.ServeFile(c, f, name)
				case types.MediaTypeOggVorbis:
					name := "track.ogg"
					p := path.Join(trackCache, name)

					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							cmd := exec.Command("ffmpeg", "-i", track.Filename, "-b:a", "96k", p)
							err := cmd.Run()
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}

					f := os.DirFS(trackCache)
					return pyrin.ServeFile(c, f, name)
				case types.MediaTypeAac:
					name := "track.aac"
					p := path.Join(trackCache, name)

					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							cmd := exec.Command("ffmpeg", "-i", track.Filename, "-codec:a", "aac", "-vn", "-b:a", "128k", p)
							cmd.Stderr = os.Stderr
							err := cmd.Run()
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}

					f := os.DirFS(trackCache)
					return pyrin.ServeFile(c, f, name)
				}

				return pyrin.NoContentNotFound()
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
