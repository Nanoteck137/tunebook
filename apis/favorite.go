package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

type GetUserFavorites struct {
	Page  types.Page `json:"page"`
	Items []Track    `json:"items"`
}

type GetFavoriteTrackIds struct {
	Ids []string `json:"ids"`
}

func InstallFavoriteHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetUserTrackFavoritesById",
			Method:       http.MethodGet,
			Path:         "/favorites/users/:userId/tracks",
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
			Name:         "GetUserTrackFavorites",
			Method:       http.MethodGet,
			Path:         "/favorites/tracks",
			ResponseType: GetUserFavorites{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				items, page, err := app.UserService().GetFavoriteTracks(
					ctx,
					service.GetFavoriteTracksParams{
						UserId:   user.Id,
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
			Name:         "GetFavoriteTrackIds",
			Method:       http.MethodGet,
			Path:         "/favorites/tracks/ids",
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
			Path:   "/favorites/tracks/:trackId",
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
			Path:   "/favorites/tracks/:trackId",
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
	)
}
