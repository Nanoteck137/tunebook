package apis

import (
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
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

	InstallHistoryHandlers(app, g)

	InstallQueueHandlers(app, g)

	InstallSearchHandlers(app, g)

	InstallFilterHandlers(app, g)
}

func RegisterHandlers(app core.App, router pyrin.Router) {
	RegisterApiHandlers(app, router.Group("/api/v1"))
	InstallFilesHandlers(app, router.Group("/files"))
	InstallStaticHandlers(app, router.Group("/"))
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		ErrorCallback: func(err error) {
			slog.Error("api error",
				slog.String("error", err.Error()),
			)
		},
		Middlewares: []pyrin.MiddlewareFunc{
			loggerMiddleware("server route"),
			corsMiddleware,
			middleware.Recoverer,
		},
	})

	RegisterHandlers(app, s)

	return s, nil
}
