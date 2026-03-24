package apis

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

func RegisterApiHandlers(app core.App, g pyrin.Group) {
	InstallSystemHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallUserHandlers(app, g)

	InstallArtistHandlers(app, g)
	InstallAlbumHandlers(app, g)
	InstallTrackHandlers(app, g)
	InstallMediaHandlers(app, g)
	InstallTagHandlers(app, g)

	InstallPlaylistHandlers(app, g)

	InstallSearchHandlers(app, g)
}

func RegisterStaticHandlers(app core.App, g pyrin.Group) {
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
}

func RegisterHandlers(app core.App, router pyrin.Router) {
	RegisterStaticHandlers(app, router.Group(""))

	RegisterApiHandlers(app, router.Group("/api/v1"))
	// TODO(patrik): Should the files be under the /api/v1 group?
	InstallFilesHandlers(app, router.Group("/files"))
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
