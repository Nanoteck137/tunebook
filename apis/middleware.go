package apis

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nanoteck137/tunebook/utils"
)

func loggerMiddleware(logName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sr := &utils.StatusRecorder{ResponseWriter: w, Status: http.StatusOK}
			next.ServeHTTP(sr, r)

			slog.LogAttrs(r.Context(), slog.LevelInfo, logName,
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", sr.Status),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

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
