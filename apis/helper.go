package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type UserCheckFunc func(user *database.User) error

func RequireAdmin(user *database.User) error {
	if user.Role != types.RoleSuperUser && user.Role != types.RoleAdmin {
		return InvalidAuth("user requires 'super_user' or 'admin' role")
	}

	return nil
}

func User(app core.App, c pyrin.Context, checks ...UserCheckFunc) (*database.User, error) {
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

func getUser(app core.App, c pyrin.Context) (*database.User, error) {
	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		ctx := context.TODO()
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(c.Request().Context(), token.UserId)
		if err != nil {
			return nil, InvalidAuth("invalid api token")
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("invalid authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Handle error better
		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		return &user, nil
	}

	return nil, InvalidAuth("invalid authorization token")
}

// TODO(patrik): Move to utils
func ConvertSqlNullString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

// TODO(patrik): Move to utils
func ConvertSqlNullInt64(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}

const (
	UNKNOWN_ARTIST_ID   = "unknown"
	UNKNOWN_ARTIST_NAME = "UNKNOWN"
)

// TODO(patrik): Cleanup
func EnsureUnknownArtistExists(ctx context.Context, db *database.Database, workDir types.WorkDir) error {
	_, err := db.GetArtistById(ctx, UNKNOWN_ARTIST_ID)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			slog.Info("Creating 'unknown' artist")
			_, err := db.CreateArtist(ctx, database.CreateArtistParams{
				Id:   UNKNOWN_ARTIST_ID,
				Name: UNKNOWN_ARTIST_NAME,
				Slug: utils.Slug(UNKNOWN_ARTIST_NAME),
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
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

const (
	IMAGE_ORIGINAL = "original.png"
	IMAGE_SMALL    = "128.png"
	IMAGE_MEDIUM   = "256.png"
	IMAGE_LARGE    = "512.png"
)

func ConvertArtistCoverURL(c pyrin.Context, artistId string, val sql.NullString) types.Images {
	first := "/files/artists/images/" + artistId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string, val sql.NullString) types.Images {
	first := "/files/albums/images/" + albumId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func ConvertPlaylistCoverURL(c pyrin.Context, playlistId string, val sql.NullString) types.Images {
	first := "/files/playlists/images/" + playlistId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}
