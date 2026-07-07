package apis

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

type UserCheckFunc func(user *database.User) error

func RequireAdmin(user *database.User) error {
	if user.Role != types.RoleSuperUser && user.Role != types.RoleAdmin {
		return InvalidAuth("user requires 'super_user' or 'admin' role")
	}

	return nil
}

func User(
	app core.App,
	c pyrin.Context,
	checks ...UserCheckFunc,
) (*database.User, error) {
	user, err := getUser(app, c)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		err := check(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func parseAuthHeader(authHeader string) string {
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 {
		return ""
	}

	if splits[0] != "Bearer" {
		return ""
	}

	return splits[1]
}

func getUser(app core.App, c pyrin.Context) (*database.User, error) {
	ctx := c.Request().Context()

	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(ctx, token.UserId)
		if err != nil {
			return nil, InvalidAuth("invalid api token")
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := parseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("missing or malformed authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, InvalidAuth("token expired")
		}

		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	err = jwtValidator.Validate(token.Claims)
	if err != nil {
		return nil, InvalidAuth("invalid authorization token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, InvalidAuth("invalid authorization token")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return nil, InvalidAuth("invalid authorization token")
	}

	user, err := app.DB().GetUserById(ctx, userId)
	if err != nil {
		return nil, InvalidAuth("invalid authorization token")
	}

	return &user, nil
}

func ConvertURL(c pyrin.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func handleUrl(c pyrin.Context, base string) types.Images {
	return types.Images{
		Original: ConvertURL(c, base),
		Small:    ConvertURL(c, base+"?size=128"),
		Medium:   ConvertURL(c, base+"?size=256"),
		Large:    ConvertURL(c, base+"?size=512"),
	}
}

func ConvertArtistCoverURL(c pyrin.Context, artistId string) types.Images {
	return handleUrl(c, "/files/artists/images/"+artistId)
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string) types.Images {
	return handleUrl(c, "/files/albums/images/"+albumId)
}

func ConvertPlaylistCoverURL(c pyrin.Context, playlistId string) types.Images {
	return handleUrl(c, "/files/playlists/images/"+playlistId)
}

func ConvertUserPictureURL(c pyrin.Context, userId string) types.Images {
	return handleUrl(c, "/files/users/images/"+userId)
}

func getPageParams(q url.Values, defaultPerPage int) types.PageParams {
	perPage := defaultPerPage
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return types.PageParams{
		PerPage: perPage,
		Page:    page,
	}
}

func getFilterParams(q url.Values) types.FilterParams {
	return types.FilterParams{
		Filter: q.Get("filter"),
		Sort:   q.Get("sort"),
	}
}

func formatTime(ms int64) string {
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}
