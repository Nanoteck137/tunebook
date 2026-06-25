package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/anvil"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/nanoteck137/validate"
)

type UserData struct {
	Id string `json:"id"`

	DisplayName string `json:"displayName"`
	Role        string `json:"role"`

	Picture types.Images `json:"picture"`

	Created string `json:"created"`
}

func ConvertDBUser(c pyrin.Context, user database.User) UserData {
	return UserData{
		Id:          user.Id,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		Picture:     ConvertUserPictureURL(c, user.Id),
		Created:     formatTime(user.Created),
	}
}

type GetUser struct {
	User UserData `json:"user"`
}

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

type GetUserFavorites struct {
	Page  types.Page `json:"page"`
	Items []Track    `json:"items"`
}

type UpdateMeBody struct {
	DisplayName *string `json:"displayName,omitempty"`
}

func (b *UpdateMeBody) Transform() {
	b.DisplayName = anvil.StringPtr(b.DisplayName)
}

func (b UpdateMeBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.DisplayName, validate.Required.When(b.DisplayName != nil)),
	)
}

type GetQuickPlaylistIds struct {
	Ids []string `json:"ids"`
}

type SetQuickPlaylistBody struct {
	PlaylistId string `json:"playlistId"`
}

type GetFavoriteTrackIds struct {
	Ids []string `json:"ids"`
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
		// TODO(patrik): Test if we need '.When()' on validate filter when
		// b.Filter is nil
		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil), validateFilter),
	)
}

type GetUserStats struct {
	NumTracksPlayed     int    `json:"numTracksPlayed"`
	NumTracksSkipped    int    `json:"numTracksSkipped"`
	NumPlaylistsCreated int    `json:"numPlaylistsCreated"`
	NumFavoriteTracks   int    `json:"numFavoriteTracks"`
	ListeningTime       int64  `json:"listeningTime"`
	LastListenedAt      *int64 `json:"lastListenedAt"`

	Updated string `json:"updated"`
}

type ApiToken struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

type GetApiTokens struct {
	Tokens []ApiToken `json:"tokens"`
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

func handleUserServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrUserServiceUserNotFound):
		return UserNotFound()
	case errors.Is(err, service.ErrUserServicePlaylistNotFound):
		return PlaylistNotFound()
	case errors.Is(err, service.ErrUserServiceTrackNotFound):
		return TrackNotFound()
	case errors.Is(err, service.ErrUserServiceUnauthorized):
		// TODO(patrik): Custom error
		return UserNotFound()
	}

	return err
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetUser",
			Method:       http.MethodGet,
			Path:         "/users/:userId",
			ResponseType: GetUser{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := app.UserService().GetUserById(
					ctx,
					service.GetUserByIdParams{
						UserId: c.Param("userId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return GetUser{
					User: ConvertDBUser(c, user),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserStats",
			Method:       http.MethodGet,
			Path:         "/users/:userId/stats",
			ResponseType: GetUserStats{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				stats, err := app.UserService().GetUserStats(
					ctx,
					service.GetUserStatsParams{
						UserId: c.Param("userId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return GetUserStats{
					NumTracksPlayed:     stats.NumTracksPlayed,
					NumTracksSkipped:    stats.NumTracksSkipped,
					NumPlaylistsCreated: stats.NumPlaylistsCreated,
					NumFavoriteTracks:   stats.NumFavoriteTracks,
					ListeningTime:       stats.ListeningTime,
					LastListenedAt:      utils.SqlNullToInt64Ptr(
						stats.LastListenedAt),
					Updated:             formatTime(stats.Updated),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserTrackFavorites",
			Method:       http.MethodGet,
			Path:         "/users/:userId/favorites/tracks",
			ResponseType: GetUserFavorites{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				items, page, err := app.UserService().GetFavoriteTracks(
					ctx,
					service.GetFavoriteTracksParams{
						UserId:   c.Param("userId"),
						Page:     getPageParams(q, 100),
						Filter:   getFilterParams(q),
						FilterId: q.Get("filterId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserFavorites{
					Page:  page,
					Items: make([]Track, len(items)),
				}

				for i, item := range items {
					res.Items[i] = ConvertDBTrack(c, item.Track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserTrackFilters",
			Method:       http.MethodGet,
			Path:         "/users/:userId/filters/tracks",
			ResponseType: GetTrackFilters{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				userId := c.Param("userId")

				ctx := c.Request().Context()

				filters, err := app.UserService().GetTrackFilters(
					ctx,
					service.GetTrackFiltersParams{
						UserId: userId,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
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
	)

	group.Register(
		pyrin.ApiHandler{
			Name:     "UpdateMe",
			Method:   http.MethodPatch,
			Path:     "/me",
			BodyType: UpdateMeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateMeBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().UpdateMe(ctx, service.UpdateMeParams{
					UserId:      user.Id,
					DisplayName: body.DisplayName,
				})
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetQuickPlaylistIds",
			Method:       http.MethodGet,
			Path:         "/me/quickplaylist",
			ResponseType: GetQuickPlaylistIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// user, err := User(app, c)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// ctx := context.TODO()
				//
				// if user.QuickPlaylist.Valid {
				// 	items, err := app.DB().GetPlaylistItems(ctx, user.QuickPlaylist.String)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				//
				// 	res := GetUserQuickPlaylistItemIds{
				// 		TrackIds: make([]string, len(items)),
				// 	}
				//
				// 	for i, item := range items {
				// 		res.TrackIds[i] = item.TrackId
				// 	}
				//
				// 	return res, nil
				// }
				//
				// // TODO(patrik): Better error
				// return nil, errors.New("No Quick Playlist set")
				return GetQuickPlaylistIds{
					Ids: []string{},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SetQuickPlaylist",
			Method:   http.MethodPost,
			Path:     "/me/quickplaylist",
			BodyType: SetQuickPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[SetQuickPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().SetQuickPlaylist(
					ctx,
					service.SetQuickPlaylistParams{
						UserId:     user.Id,
						PlaylistId: body.PlaylistId,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetFavoriteTrackIds",
			Method:       http.MethodGet,
			Path:         "/me/favorites/tracks/ids",
			ResponseType: GetFavoriteTrackIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				items, err := app.UserService().GetFavoriteTrackIds(
					ctx,
					service.GetFavoriteTrackIdsParams{
						UserId: user.Id,
					},
				)
				if err != nil {
					return nil, err
				}

				return GetFavoriteTrackIds{
					Ids: items,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "FavoriteTrack",
			Method: http.MethodPost,
			Path:   "/me/favorites/tracks/:trackId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().FavoriteTrack(
					ctx,
					service.FavoriteTrackParams{
						UserId:  user.Id,
						TrackId: c.Param("trackId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "UnfavoriteTrack",
			Method: http.MethodDelete,
			Path:   "/me/favorites/tracks/:trackId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().UnfavoriteTrack(
					ctx,
					service.UnfavoriteTrackParams{
						UserId:  user.Id,
						TrackId: c.Param("trackId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "CreateTrackFilter",
			Method:   http.MethodPost,
			Path:     "/me/filters/tracks",
			BodyType: CreateTrackFilterBody{},
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

				filterId, err := app.UserService().CreateTrackFilter(
					ctx,
					service.CreateTrackFilterParams{
						UserId: user.Id,
						Name:   body.Name,
						Filter: body.Filter,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return CreateTrackFilter{
					FilterId: filterId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateTrackFilter",
			Method:   http.MethodPatch,
			Path:     "/me/filters/tracks/:filterId",
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

				ctx := c.Request().Context()

				err = app.UserService().UpdateTrackFilter(
					ctx,
					service.UpdateTrackFilterParams{
						FilterId: c.Param("filterId"),
						UserId:   user.Id,
						Name:     body.Name,
						Filter:   body.Filter,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteTrackFilter",
			Method: http.MethodDelete,
			Path:   "/me/filters/tracks/:filterId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().DeleteTrackFilter(
					ctx,
					service.DeleteTrackFilterParams{
						FilterId: c.Param("filterId"),
						UserId:   user.Id,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetApiTokens",
			Method:       http.MethodGet,
			Path:         "/me/apitokens",
			ResponseType: GetApiTokens{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				tokens, err := app.UserService().GetApiTokens(
					ctx,
					service.GetApiTokensParams{
						UserId: user.Id,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetApiTokens{
					Tokens: make([]ApiToken, len(tokens)),
				}

				for i, token := range tokens {
					res.Tokens[i] = ApiToken{
						Id:      token.Id,
						Name:    token.Name,
						Created: formatTime(token.Created),
						Updated: formatTime(token.Updated),
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateApiToken",
			Method:       http.MethodPost,
			Path:         "/me/apitokens",
			ResponseType: CreateApiToken{},
			BodyType:     CreateApiTokenBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateApiTokenBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				tokenId, err := app.UserService().CreateApiToken(
					ctx,
					service.CreateApiTokenParams{
						UserId: user.Id,
						Name:   body.Name,
					},
				)
				if err != nil {
					return nil, err
				}

				return CreateApiToken{
					Token: tokenId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteApiToken",
			Method: http.MethodDelete,
			Path:   "/me/apitokens/:tokenId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = app.UserService().DeleteApiToken(
					ctx,
					service.DeleteApiTokenParams{
						TokenId: c.Param("tokenId"),
						UserId:  user.Id,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},
	)
}
