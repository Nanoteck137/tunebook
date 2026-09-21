package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
)

func InstallJobHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetJobs",
			Method:       http.MethodGet,
			Path:         "/system/jobs",
			ResponseType: service.GetJobsResponse{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				jobs, err := app.JobService().GetJobs(
					c.Request().Context(),
					service.DefaultJobsLimit,
				)
				if err != nil {
					return nil, err
				}

				return jobs, nil
			},
		},
	)
}
