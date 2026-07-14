package apis

import (
	"context"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/anvil"
	"github.com/nanoteck137/validate"
)

type TrackFilter struct {
	FilterId string `json:"filterId"`
	UserId   string `json:"userId"`

	Name   string `json:"name"`
	Filter string `json:"filter"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

type GetTrackFilters struct {
	Filters []TrackFilter `json:"filters"`
}

type CreateTrackFilter struct {
	FilterId string `json:"filterId"`
}

type CreateTrackFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *CreateTrackFilterBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Filter = anvil.String(b.Filter)
}

func (b CreateTrackFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.Filter, validate.Required, validateTrackFilter),
	)
}

type UpdateTrackFilterBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b *UpdateTrackFilterBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	b.Filter = anvil.StringPtr(b.Filter)
}

func (b UpdateTrackFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		validate.Field(
			&b.Filter,
			validate.Required.When(b.Filter != nil),
			validateTrackFilter,
		),
	)
}

func InstallFilterHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTrackFilters",
			Method:       http.MethodGet,
			Path:         "/filters/tracks",
			ResponseType: GetTrackFilters{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				filters, err := app.TrackService().GetTrackFilters(
					ctx,
					service.GetTrackFiltersParams{
						UserId: user.Id,
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				res := GetTrackFilters{
					Filters: make([]TrackFilter, len(filters)),
				}

				for i, filter := range filters {
					res.Filters[i] = TrackFilter{
						FilterId: filter.Id,
						UserId:   filter.UserId,
						Name:     filter.Name,
						Filter:   filter.Filter,
						Created:  formatTime(filter.Created),
						Updated:  formatTime(filter.Updated),
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateTrackFilter",
			Method:       http.MethodPost,
			Path:         "/filters/tracks",
			ResponseType: CreateTrackFilter{},
			BodyType:     CreateTrackFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateTrackFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				filterId, err := app.TrackService().CreateTrackFilter(
					ctx,
					service.CreateTrackFilterParams{
						UserId: user.Id,
						Name:   body.Name,
						Filter: body.Filter,
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				return CreateTrackFilter{
					FilterId: filterId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateTrackFilter",
			Method:   http.MethodPatch,
			Path:     "/filters/tracks/:filterId",
			BodyType: UpdateTrackFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateTrackFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.TrackService().UpdateTrackFilter(
					ctx,
					service.UpdateTrackFilterParams{
						FilterId: c.Param("filterId"),
						UserId:   user.Id,
						Name:     body.Name,
						Filter:   body.Filter,
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteTrackFilter",
			Method: http.MethodDelete,
			Path:   "/filters/tracks/:filterId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.TrackService().DeleteTrackFilter(
					ctx,
					service.DeleteTrackFilterParams{
						FilterId: c.Param("filterId"),
						UserId:   user.Id,
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
