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
			Name:         "RunJob",
			Method:       http.MethodPost,
			Path:         "/system/job/:jobName",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				jobName := c.Param("jobName")

				go func() {
					app.JobService().RunJob(context.Background(), jobName)
				}()

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SyncLibrary",
			Method:       http.MethodPost,
			Path:         "/system/library",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				go func() {
					app.LibraryService().Sync()
				}()

				return nil, nil
			},
		},

		// TODO(patrik): Better name?
		// pyrin.ApiHandler{
		// 	Name:   "CleanupLibrary",
		// 	Method: http.MethodPost,
		// 	Path:   "/system/library/cleanup",
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		if syncHandler.isSyncing.Load() {
		// 			return nil, errors.New("library is syncing")
		// 		}
		//
		// 		err := syncHandler.Cleanup(app)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/sse",
			HandlerFunc: func(c pyrin.Context) error {
				app.Broker().ServeHTTP(c.Response(), c.Request())
				return nil
			},
		},
	)
}
