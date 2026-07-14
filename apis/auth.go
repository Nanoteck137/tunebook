package apis

import (
	"errors"
	"net/http"
	"time"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/render"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

type GetMe struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`

	Picture types.Images `json:"picture"`

	QuickPlaylist *string `json:"quickPlaylist"`
}

type AuthInitiate struct {
	RequestId string `json:"requestId"`
	AuthUrl   string `json:"authUrl"`
	Challenge string `json:"challenge"`
	ExpiresAt string `json:"expiresAt"`
}

type AuthInitiateBody struct {
	ProviderId string `json:"providerId"`
}

type AuthQuickConnectInitiate struct {
	Code      string `json:"code"`
	Challenge string `json:"challenge"`
	ExpiresAt string `json:"expiresAt"`
}

type AuthLoginWithCode struct {
	Token string `json:"token"`
}

type AuthLoginWithCodeBody struct {
	ProviderId string `json:"providerId"`
	Code       string `json:"code"`
	State      string `json:"state"`
}

type AuthFinishProvider struct {
	Token string `json:"token"`
}

type AuthFinishProviderBody struct {
	RequestId string `json:"requestId"`
	Challenge string `json:"challenge"`
}

type AuthProvider struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type GetAuthProviders struct {
	Providers []AuthProvider `json:"providers"`
}

type AuthClaimQuickConnectCodeBody struct {
	Code string `json:"code"`
}

type AuthFinishQuickConnect struct {
	Token string `json:"token"`
}

type AuthFinishQuickConnectBody struct {
	Code      string `json:"code"`
	Challenge string `json:"challenge"`
}

type AuthGetQuickConnectStatus struct {
	Status string `json:"status"`
}

type AuthGetQuickConnectStatusBody struct {
	Code      string `json:"code"`
	Challenge string `json:"challenge"`
}

type AuthGetProviderStatus struct {
	Status string `json:"status"`
}

type AuthGetProviderStatusBody struct {
	RequestId string `json:"requestId"`
	Challenge string `json:"challenge"`
}

func handleAuthServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrAuthServiceRequestNotFound):
		return RequestNotFound()
	case errors.Is(err, service.ErrAuthServiceProviderNotFound):
		return ProviderNotFound()
	case errors.Is(err, service.ErrAuthServiceChallengeMismatch):
		return ChallengeMismatch()
	}

	return err
}

func InstallAuthHandlers(app core.App, group pyrin.Group) {
	// NOTE(patrik): Provider Authentication
	group.Register(
		pyrin.ApiHandler{
			Name:         "AuthGetProviders",
			Method:       http.MethodGet,
			Path:         "/auth/providers",
			ResponseType: GetAuthProviders{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providers := app.Config().OidcProviders

				res := GetAuthProviders{
					Providers: make([]AuthProvider, 0, len(providers)),
				}

				for _, provider := range providers {
					res.Providers = append(res.Providers, AuthProvider{
						Id:          provider.Id,
						DisplayName: provider.Name,
					})
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AuthProviderInitiate",
			Method:       http.MethodPost,
			Path:         "/auth/providers/initiate",
			ResponseType: AuthInitiate{},
			BodyType:     AuthInitiateBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthInitiateBody](c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				authService := app.AuthService()

				res, err := authService.CreateProviderRequest(
					ctx, body.ProviderId)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthInitiate{
					RequestId: res.RequestId,
					AuthUrl:   res.AuthUrl,
					Challenge: res.Challenge,
					ExpiresAt: res.Expires.Format(time.RFC3339Nano),
				}, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "AuthCallback",
			Method: http.MethodGet,
			Path:   "/auth/providers/callback",
			HandlerFunc: func(c pyrin.Context) error {
				url := c.Request().URL
				state := url.Query().Get("state")
				code := url.Query().Get("code")

				authService := app.AuthService()

				err := authService.CompleteProviderRequest(state, code)
				if err != nil {
					if errors.Is(err, service.ErrAuthServiceRequestExpired) {
						render.RenderCallbackRequestExpired(c.Response())
						c.Response().WriteHeader(http.StatusOK)

						return nil
					}

					render.RenderCallbackError(c.Response())
					c.Response().WriteHeader(http.StatusOK)

					return nil
				}

				render.RenderCallbackSuccess(c.Response())
				c.Response().WriteHeader(http.StatusOK)

				return nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AuthFinishProvider",
			Path:         "/auth/providers/finish",
			Method:       http.MethodPost,
			ResponseType: AuthFinishProvider{},
			BodyType:     AuthFinishProviderBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthFinishProviderBody](c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				authService := app.AuthService()

				token, err := authService.CreateAuthTokenForProvider(
					ctx, body.RequestId, body.Challenge)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthFinishProvider{
					Token: token,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AuthGetProviderStatus",
			Path:         "/auth/provider/status",
			Method:       http.MethodPost,
			ResponseType: AuthGetProviderStatus{},
			BodyType:     AuthGetProviderStatusBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthGetProviderStatusBody](c)
				if err != nil {
					return nil, err
				}

				authService := app.AuthService()

				status, err := authService.CheckProviderRequestStatus(
					body.RequestId, body.Challenge)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthGetProviderStatus{
					Status: string(status),
				}, nil
			},
		},
	)

	// NOTE(patrik): Quick Connect Authentication
	group.Register(
		pyrin.ApiHandler{
			Name:         "AuthQuickConnectInitiate",
			Method:       http.MethodPost,
			Path:         "/auth/quick-connect/initiate",
			ResponseType: AuthQuickConnectInitiate{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				authService := app.AuthService()

				res, err := authService.CreateQuickConnectRequest()
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthQuickConnectInitiate{
					Code:      res.Code,
					Challenge: res.Challenge,
					ExpiresAt: res.Expires.Format(time.RFC3339Nano),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AuthClaimQuickConnectCode",
			Method:   http.MethodPost,
			Path:     "/auth/quick-connect/claim",
			BodyType: AuthClaimQuickConnectCodeBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthClaimQuickConnectCodeBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				authService := app.AuthService()

				err = authService.CompleteQuickConnectRequest(
					body.Code, user.Id)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AuthGetQuickConnectStatus",
			Path:         "/auth/quick-connect/status",
			Method:       http.MethodPost,
			ResponseType: AuthGetQuickConnectStatus{},
			BodyType:     AuthGetQuickConnectStatusBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthGetQuickConnectStatusBody](c)
				if err != nil {
					return nil, err
				}

				authService := app.AuthService()

				status, err := authService.CheckQuickConnectRequestStatus(
					body.Code, body.Challenge)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthGetQuickConnectStatus{
					Status: string(status),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AuthFinishQuickConnect",
			Path:         "/auth/quick-connect/finish",
			Method:       http.MethodPost,
			ResponseType: AuthFinishQuickConnect{},
			BodyType:     AuthFinishQuickConnectBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AuthFinishQuickConnectBody](c)
				if err != nil {
					return nil, err
				}

				authService := app.AuthService()

				token, err := authService.CreateAuthTokenForQuickConnect(
					body.Code, body.Challenge)
				if err != nil {
					return nil, handleAuthServiceErrors(err)
				}

				return AuthFinishQuickConnect{
					Token: token,
				}, nil
			},
		},
	)

	// NOTE(patrik): Other Authentication related stuff
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetMe",
			Path:         "/auth/me",
			Method:       http.MethodGet,
			ResponseType: GetMe{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				return GetMe{
					Id:          user.Id,
					Email:       user.Email,
					DisplayName: user.DisplayName,
					Role:        user.Role,
					Picture:     ConvertUserPictureURL(c, user.Id),
					QuickPlaylist: utils.SqlNullToStringPtr(
						user.QuickPlaylist),
				}, nil
			},
		},
	)
}
