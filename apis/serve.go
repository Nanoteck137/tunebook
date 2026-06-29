package apis

import (
	"log/slog"
	"net/http"
	"os"
	"time"

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
	)

	// TODO(patrik): I don't like this
	if app != nil {
		webDir := app.Config().WebDir
		if webDir != "" {
			g.Register(
				// TODO(patrik): Fix this
				pyrin.SpaHandler(os.DirFS(webDir), "index.html"),
			)
		}
	}
}

func RegisterHandlers(app core.App, router pyrin.Router) {
	RegisterStaticHandlers(app, router.Group("/"))

	RegisterApiHandlers(app, router.Group("/api/v1"))
	// TODO(patrik): Should the files be under the /api/v1 group?
	InstallFilesHandlers(app, router.Group("/files"))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// TODO(patrik): Move
func loggerMiddleware(logName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sr, r)

			slog.LogAttrs(r.Context(), slog.LevelInfo, logName,
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", sr.status),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

// TODO(patrik): Move
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		ErrorCallback: func(err error) {
			// TODO(patrik): Handle this better
			slog.Error("API Error", "err", err)
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
