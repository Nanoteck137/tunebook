package apis

import (
	"context"
	"net/http"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

type GetSystemInfo struct {
	Version string `json:"version"`
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
					Version: dwebble.Version,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RunJob",
			Method: http.MethodPost,
			Path:   "/system/job/:jobName",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				jobName := c.Param("jobName")

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				go func() {
					app.JobService().RunJob(context.Background(), jobName)
				}()

				return nil, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/sse",
			HandlerFunc: func(c pyrin.Context) error {
				// TODO(patrik): Figure out how to authenticate this, because
				// with the EventSource API you can't send custom headers
				app.Broker().ServeHTTP(c.Response(), c.Request())
				return nil
			},
		},
	)
}
