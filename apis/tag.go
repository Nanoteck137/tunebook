package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
)

type Tag struct {
	Slug string `json:"slug"`
}

type GetTags struct {
	Tags []Tag `json:"tags"`
}

func InstallTagHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTags",
			Method:       http.MethodGet,
			Path:         "/tags",
			ResponseType: GetTags{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				tags, err := app.DB().GetAllTags(ctx)
				if err != nil {
					return nil, err
				}

				res := GetTags{
					Tags: make([]Tag, len(tags)),
				}
				for i, tag := range tags {
					res.Tags[i] = Tag{
						Slug: tag.Slug,
					}
				}

				return res, nil
			},
		},
	)
}
