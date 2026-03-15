package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
)

type UpdateUserSettingsBody struct {
	QuickPlaylist *string `json:"quickPlaylist,omitempty"`
}

func (b *UpdateUserSettingsBody) Transform() {
	b.QuickPlaylist = anvil.StringPtr(b.QuickPlaylist)
}

func (b UpdateUserSettingsBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.QuickPlaylist), // validate.Required.When(b.QuickPlaylist != nil),
	)
}

type TrackId struct {
	TrackId string `json:"trackId"`
}

func (b *TrackId) Transform() {
	// b.Tracks = *transform.DiscardEmptyStringEntries()
}

type GetUserQuickPlaylistItemIds struct {
	TrackIds []string `json:"trackIds"`
}

type CreateApiToken struct {
	Token string `json:"token"`
}

type CreateApiTokenBody struct {
	Name string `json:"name"`
}

func (b *CreateApiTokenBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b CreateApiTokenBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type ApiToken struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAllApiTokens struct {
	Tokens []ApiToken `json:"tokens"`
}

type GetUser struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type TrackFilter struct {
	FilterId string `json:"filterId"`
	UserId   string `json:"userId"`

	Name   string `json:"name"`
	Filter string `json:"filter"`

	// TODO(patrik): Created, Updated
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
		validate.Field(&b.Filter, validate.Required, validateFilter),
	)
}

type EditTrackFilterBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b *EditTrackFilterBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	b.Filter = anvil.StringPtr(b.Filter)
}

func (b EditTrackFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		// TODO(patrik): Test if we need When on validate filter when
		// b.Filter is nil
		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil), validateFilter),
	)
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetUser",
			Method:       http.MethodGet,
			Path:         "/users/:id",
			ResponseType: GetUser{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				user, err := app.DB().GetUserById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, UserNotFound()
					}

					return nil, err
				}

				return GetUser{
					Id:          user.Id,
					DisplayName: user.DisplayName,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateUserSettings",
			Method:   http.MethodPatch,
			Path:     "/user/settings",
			BodyType: UpdateUserSettingsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateUserSettingsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				settings := user.ToUserSettings()

				if body.QuickPlaylist != nil {
					id := *body.QuickPlaylist

					if id != "" {
						_, err := app.DB().GetPlaylistById(context.TODO(), id)
						if err != nil {
							// TODO(patrik): Handle error
							return nil, err
						}
					}

					settings.QuickPlaylist = sql.NullString{
						String: id,
						Valid:  id != "",
					}
				}

				err = app.DB().UpdateUserSettings(context.TODO(), settings)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddToUserQuickPlaylist",
			Method:   http.MethodPost,
			Path:     "/user/quickplaylist",
			BodyType: TrackId{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[TrackId](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					index, err := app.DB().GetNextPlaylistItemIndex(ctx, user.QuickPlaylist.String)
					if err != nil {
						return nil, err
					}

					// TODO(patrik): Check body.TrackId

					err = app.DB().CreatePlaylistItem(ctx, database.CreatePlaylistItemParams{
						PlaylistId: user.QuickPlaylist.String,
						TrackId:    body.TrackId,
						Order:      index,
					})
					if err != nil {
						// TODO(patrik): Handle error
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemoveItemFromUserQuickPlaylist",
			Method:   http.MethodDelete,
			Path:     "/user/quickplaylist",
			BodyType: TrackId{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[TrackId](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					err := app.DB().DeletePlaylistItem(ctx, user.QuickPlaylist.String, body.TrackId)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserQuickPlaylistItemIds",
			Method:       http.MethodGet,
			Path:         "/user/quickplaylist",
			ResponseType: GetUserQuickPlaylistItemIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					items, err := app.DB().GetPlaylistItems(ctx, user.QuickPlaylist.String)
					if err != nil {
						return nil, err
					}

					res := GetUserQuickPlaylistItemIds{
						TrackIds: make([]string, len(items)),
					}

					for i, item := range items {
						res.TrackIds[i] = item.TrackId
					}

					return res, nil
				}

				// TODO(patrik): Better error
				return nil, errors.New("No Quick Playlist set")
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateApiToken",
			Method:       http.MethodPost,
			Path:         "/user/apitoken",
			ResponseType: CreateApiToken{},
			BodyType:     CreateApiTokenBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateApiTokenBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				token, err := app.DB().CreateApiToken(ctx, database.CreateApiTokenParams{
					UserId: user.Id,
					Name:   body.Name,
				})
				if err != nil {
					return nil, err
				}

				return CreateApiToken{
					Token: token.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAllApiTokens",
			Method:       http.MethodGet,
			Path:         "/user/apitoken",
			ResponseType: GetAllApiTokens{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				tokens, err := app.DB().GetAllApiTokensForUser(ctx, user.Id)
				if err != nil {
					return nil, err
				}

				res := GetAllApiTokens{
					Tokens: make([]ApiToken, len(tokens)),
				}

				for i, token := range tokens {
					res.Tokens[i] = ApiToken{
						Id:   token.Id,
						Name: token.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteApiToken",
			Method: http.MethodDelete,
			Path:   "/user/apitoken/:id",
			Errors: []pyrin.ErrorType{ErrTypeApiTokenNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				tokenId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				token, err := app.DB().GetApiTokenById(ctx, tokenId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ApiTokenNotFound()
					}

					return nil, err
				}

				if token.UserId != user.Id {
					return nil, ApiTokenNotFound()
				}

				err = app.DB().DeleteApiToken(ctx, tokenId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackFilters",
			Method:       http.MethodGet,
			Path:         "/user/tracks/filter",
			ResponseType: GetTrackFilters{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				filters, err := app.DB().GetTrackFiltersByUserId(ctx, user.Id)
				if err != nil {
					return nil, err
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
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateTrackFilter",
			Method:       http.MethodPost,
			Path:         "/user/tracks/filter",
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

				ctx := c.Request().Context()

				filterId, err := app.DB().CreateTrackFilter(ctx, database.CreateTrackFilterParams{
					UserId: user.Id,
					Name:   body.Name,
					Filter: body.Filter,
				})
				if err != nil {
					return nil, err
				}

				return CreateTrackFilter{
					FilterId: filterId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditTrackFilter",
			Method:   http.MethodPost,
			Path:     "/user/tracks/filter/:filterId",
			BodyType: EditTrackFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				filterId := c.Param("filterId")

				body, err := pyrin.Body[EditTrackFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				dbFilter, err := app.DB().GetTrackFilterById(ctx, filterId, user.Id)
				if err != nil {
					// TODO(patrik): Handle Error
					// if errors.Is(err, database.ErrItemNotFound) {
					// 	return nil, error here pls
					// }

					return nil, err
				}

				changes := database.TrackFilterChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbFilter.Name,
					}
				}

				if body.Filter != nil {
					changes.Filter = types.Change[string]{
						Value:   *body.Filter,
						Changed: *body.Filter != dbFilter.Filter,
					}
				}

				err = app.DB().UpdateTrackFilter(ctx, dbFilter.Id, dbFilter.UserId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteTrackFilter",
			Method: http.MethodDelete,
			Path:   "/user/tracks/filter/:filterId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				filterId := c.Param("filterId")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				dbFilter, err := app.DB().GetTrackFilterById(ctx, filterId, user.Id)
				if err != nil {
					// TODO(patrik): Handle Error
					// if errors.Is(err, database.ErrItemNotFound) {
					// 	return nil, error here pls
					// }

					return nil, err
				}

				err = app.DB().DeleteTrackFilter(ctx, dbFilter.Id, dbFilter.UserId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
