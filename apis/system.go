package apis

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
)

type GetSystemInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	StartedAt int64  `json:"startedAt"`
}

type CreateSSEToken struct {
	Token string `json:"token"`
}

type GetSystemStats struct {
	Users        int `json:"users"`
	Artists      int `json:"artists"`
	Albums       int `json:"albums"`
	Tracks       int `json:"tracks"`
	Playlists    int `json:"playlists"`
	Favorites    int `json:"favorites"`
	TrackFilters int `json:"trackFilters"`
	Queues       int `json:"queues"`

	TotalPlays         int   `json:"totalPlays"`
	TotalListeningTime int64 `json:"totalListeningTime"`

	DataDir      string `json:"dataDir"`
	DatabaseFile string `json:"databaseFile"`
	DatabaseSize int64  `json:"databaseSize"`
}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version:   tunebook.Version,
					Commit:    tunebook.Commit,
					StartedAt: tunebook.BootTime.UnixMilli(),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetSystemStats",
			Method:       http.MethodGet,
			Path:         "/system/stats",
			ResponseType: GetSystemStats{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				stats, err := app.DB().GetStats(c.Request().Context())
				if err != nil {
					return nil, err
				}

				databaseFile := app.FilesystemService().DatabaseFile()
				var databaseSize int64
				if info, err := os.Stat(databaseFile); err == nil {
					databaseSize = info.Size()
				}

				return GetSystemStats{
					Users:        stats.Users,
					Artists:      stats.Artists,
					Albums:       stats.Albums,
					Tracks:       stats.Tracks,
					Playlists:    stats.Playlists,
					Favorites:    stats.Favorites,
					TrackFilters: stats.TrackFilters,
					Queues:       stats.Queues,

					TotalPlays:         stats.TotalPlays,
					TotalListeningTime: stats.TotalListeningTime,

					DataDir:      app.FilesystemService().DataDir().String(),
					DatabaseFile: databaseFile,
					DatabaseSize: databaseSize,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RunTask",
			Method: http.MethodPost,
			Path:   "/system/task/:taskName",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taskName := c.Param("taskName")

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				go func() {
					app.TaskService().RunTask(context.Background(), taskName)
				}()

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateSseToken",
			Method:       http.MethodPost,
			Path:         "/system/sse/token",
			ResponseType: CreateSSEToken{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				token, err := app.AuthService().CreateSSEToken(user.Id)
				if err != nil {
					return nil, err
				}

				return CreateSSEToken{Token: token}, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/sse",
			HandlerFunc: func(c pyrin.Context) error {
				tokenString := c.Request().URL.Query().Get("token")
				if tokenString == "" {
					return InvalidAuth("missing sse token")
				}

				userId, err := app.AuthService().ValidateSSEToken(tokenString)
				if err != nil {
					if errors.Is(err, service.ErrAuthServiceRequestExpired) {
						return InvalidAuth("sse token expired")
					}

					return InvalidAuth("invalid sse token")
				}

				app.Broker().ServeHTTP(c.Response(), c.Request(), userId)
				return nil
			},
		},
	)
}
