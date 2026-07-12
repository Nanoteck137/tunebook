package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

type CreateSSEToken struct {
	Token string `json:"token"`
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
					Version: tunebook.Version,
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
			Name:   "CreateSseToken",
			Method: http.MethodPost,
			Path:   "/system/sse/token",
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
