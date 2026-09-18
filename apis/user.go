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

type UpdateMeBody struct {
	DisplayName *string `json:"displayName,omitempty"`
	PictureUrl  *string `json:"pictureUrl,omitempty"`
}

func (b *UpdateMeBody) Transform() {
	b.DisplayName = anvil.StringPtr(b.DisplayName)
	b.PictureUrl = anvil.StringPtr(b.PictureUrl)
}

func (b UpdateMeBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(
			&b.DisplayName, validate.Required.When(b.DisplayName != nil)),
		validate.Field(
			&b.PictureUrl, validate.Required.When(b.PictureUrl != nil)),
	)
}

type SetQuickPlaylistBody struct {
	PlaylistId string `json:"playlistId"`
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
	case errors.Is(err, service.ErrUserServiceApiTokenNotFound):
		return ApiTokenNotFound()
	case errors.Is(err, service.ErrUserServiceUnauthorized):
		return NotAuthorized()
	case errors.Is(err, service.ErrImageServiceUnsupportedImageFormat):
		return UnsupportedImageType()
	}

	return err
}

type GetUserTopTracks struct {
	Tracks []Track `json:"tracks"`
}

type UserYearReview struct {
	Year          int   `json:"year"`
	TrackCount    int   `json:"trackCount"`
	ListeningTime int64 `json:"listeningTime"`

	AvgCompletion float64 `json:"avgCompletion"`
	SkipCount     int     `json:"skipCount"`
	UniqueTracks  int     `json:"uniqueTracks"`
	FavoritePlays int     `json:"favoritePlays"`
}

type GetAllUserYearReviews struct {
	Reviews []UserYearReview `json:"reviews"`
}

type GetUserYearReview struct {
	Review UserYearReview `json:"review"`
}

type RankedTrack struct {
	Track

	Rank      int `json:"rank"`
	PlayCount int `json:"playCount"`
}

type GetUserYearReviewTracks struct {
	Page   types.Page    `json:"page"`
	Tracks []RankedTrack `json:"tracks"`
}

type RankedAlbum struct {
	Album

	Rank      int `json:"rank"`
	PlayCount int `json:"playCount"`
}

type GetUserYearReviewAlbums struct {
	Page   types.Page    `json:"page"`
	Albums []RankedAlbum `json:"albums"`
}

type RankedArtist struct {
	Artist

	Rank      int `json:"rank"`
	PlayCount int `json:"playCount"`
}

type GetUserYearReviewArtists struct {
	Page    types.Page     `json:"page"`
	Artists []RankedArtist `json:"artists"`
}

type RankedTag struct {
	TagSlug string `json:"tagSlug"`

	Rank      int `json:"rank"`
	PlayCount int `json:"playCount"`
}

type GetUserYearReviewTags struct {
	Page types.Page  `json:"page"`
	Tags []RankedTag `json:"tags"`
}

type RankedDecade struct {
	Decade int `json:"decade"`

	Rank      int `json:"rank"`
	PlayCount int `json:"playCount"`
}

type GetUserYearReviewDecades struct {
	Page    types.Page     `json:"page"`
	Decades []RankedDecade `json:"decades"`
}

type UserYearReviewMonth struct {
	Month int `json:"month"`

	PlayCount int   `json:"playCount"`
	PlayTime  int64 `json:"playTime"`

	AvgCompletion float64 `json:"avgCompletion"`
	SkipCount     int     `json:"skipCount"`
	UniqueTracks  int     `json:"uniqueTracks"`
	FavoritePlays int     `json:"favoritePlays"`
}

type GetAllUserYearReviewMonths struct {
	Months []UserYearReviewMonth `json:"months"`
}

type GetUserYearReviewMonth struct {
	Month UserYearReviewMonth `json:"month"`
}

type GetUserYearReviewMonthTracks struct {
	Page   types.Page    `json:"page"`
	Tracks []RankedTrack `json:"tracks"`
}

type GetUserYearReviewMonthAlbums struct {
	Page   types.Page    `json:"page"`
	Albums []RankedAlbum `json:"albums"`
}

type GetUserYearReviewMonthArtists struct {
	Page    types.Page     `json:"page"`
	Artists []RankedArtist `json:"artists"`
}

type GetUserYearReviewMonthTags struct {
	Page types.Page  `json:"page"`
	Tags []RankedTag `json:"tags"`
}

type GetUserYearReviewMonthDecades struct {
	Page    types.Page     `json:"page"`
	Decades []RankedDecade `json:"decades"`
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAllUserYearReviewMonths",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months",
			ResponseType: GetAllUserYearReviewMonths{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				months, err := app.UserService().GetAllUserYearReviewMonths(
					ctx,
					service.GetAllUserYearReviewMonthsParams{
						UserId: c.Param("userId"),
						Year:   year,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetAllUserYearReviewMonths{
					Months: make([]UserYearReviewMonth, len(months)),
				}

				for i, month := range months {
					res.Months[i] = UserYearReviewMonth{
						Month:         month.Month,
						PlayCount:     month.PlayCount,
						PlayTime:      month.PlayTime,
						AvgCompletion: month.AvgCompletion,
						SkipCount:     month.SkipCount,
						UniqueTracks:  month.UniqueTracks,
						FavoritePlays: month.FavoritePlays,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonth",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month",
			ResponseType: GetUserYearReviewMonth{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				m, err := app.UserService().GetUserYearReviewMonth(
					ctx,
					service.GetUserYearReviewMonthParams{
						UserId: c.Param("userId"),
						Year:   year,
						Month:  month,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return GetUserYearReviewMonth{
					Month: UserYearReviewMonth{
						Month:         m.Month,
						PlayCount:     m.PlayCount,
						PlayTime:      m.PlayTime,
						AvgCompletion: m.AvgCompletion,
						SkipCount:     m.SkipCount,
						UniqueTracks:  m.UniqueTracks,
						FavoritePlays: m.FavoritePlays,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonthTracks",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month/tracks",
			ResponseType: GetUserYearReviewMonthTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				tracks, page, err := app.UserService().
					GetUserYearReviewMonthTracks(
						ctx,
						service.GetUserYearReviewMonthTracksParams{
							UserId: c.Param("userId"),
							Year:   year,
							Month:  month,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewMonthTracks{
					Page:   page,
					Tracks: make([]RankedTrack, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = RankedTrack{
						Track:     ConvertDBTrack(c, track.Track),
						Rank:      track.Rank,
						PlayCount: track.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonthAlbums",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month/albums",
			ResponseType: GetUserYearReviewMonthAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				albums, page, err := app.UserService().
					GetUserYearReviewMonthAlbums(
						ctx,
						service.GetUserYearReviewMonthAlbumsParams{
							UserId: c.Param("userId"),
							Year:   year,
							Month:  month,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewMonthAlbums{
					Page:   page,
					Albums: make([]RankedAlbum, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = RankedAlbum{
						Album:     ConvertDBAlbum(c, album.Album),
						Rank:      album.Rank,
						PlayCount: album.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonthArtists",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month/artists",
			ResponseType: GetUserYearReviewMonthArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				artists, page, err := app.UserService().
					GetUserYearReviewMonthArtists(
						ctx,
						service.GetUserYearReviewMonthArtistsParams{
							UserId: c.Param("userId"),
							Year:   year,
							Month:  month,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewMonthArtists{
					Page:    page,
					Artists: make([]RankedArtist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = RankedArtist{
						Artist:    ConvertDBArtist(c, artist.Artist),
						Rank:      artist.Rank,
						PlayCount: artist.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonthTags",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month/tags",
			ResponseType: GetUserYearReviewMonthTags{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				tags, page, err := app.UserService().
					GetUserYearReviewMonthTags(
						ctx,
						service.GetUserYearReviewMonthTagsParams{
							UserId: c.Param("userId"),
							Year:   year,
							Month:  month,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewMonthTags{
					Page: page,
					Tags: make([]RankedTag, len(tags)),
				}

				for i, tag := range tags {
					res.Tags[i] = RankedTag{
						TagSlug:   tag.TagSlug,
						Rank:      tag.Rank,
						PlayCount: tag.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewMonthDecades",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/months/:month/decades",
			ResponseType: GetUserYearReviewMonthDecades{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				month, err := parseIntParam(c, "month")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				decades, page, err := app.UserService().
					GetUserYearReviewMonthDecades(
						ctx,
						service.GetUserYearReviewMonthDecadesParams{
							UserId: c.Param("userId"),
							Year:   year,
							Month:  month,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewMonthDecades{
					Page:    page,
					Decades: make([]RankedDecade, len(decades)),
				}

				for i, decade := range decades {
					res.Decades[i] = RankedDecade{
						Decade:    decade.Decade,
						Rank:      decade.Rank,
						PlayCount: decade.PlayCount,
					}
				}

				return res, nil
			},
		},
	)

	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAllUserYearReviews",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews",
			ResponseType: GetAllUserYearReviews{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				reviews, err := app.UserService().GetAllUserYearReviews(
					ctx,
					service.GetAllUserYearReviewsParams{
						UserId: c.Param("userId"),
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetAllUserYearReviews{
					Reviews: make([]UserYearReview, 0, len(reviews)),
				}

				for _, review := range reviews {
					res.Reviews = append(res.Reviews, UserYearReview{
						Year:          review.Year,
						TrackCount:    review.TrackCount,
						ListeningTime: review.ListeningTime,
						AvgCompletion: review.AvgCompletion,
						SkipCount:     review.SkipCount,
						UniqueTracks:  review.UniqueTracks,
						FavoritePlays: review.FavoritePlays,
					})
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReview",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year",
			ResponseType: GetUserYearReview{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				review, err := app.UserService().GetUserYearReview(
					ctx,
					service.GetUserYearReviewParams{
						UserId: c.Param("userId"),
						Year:   year,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return GetUserYearReview{
					Review: UserYearReview{
						Year:          review.Year,
						TrackCount:    review.TrackCount,
						ListeningTime: review.ListeningTime,
						AvgCompletion: review.AvgCompletion,
						SkipCount:     review.SkipCount,
						UniqueTracks:  review.UniqueTracks,
						FavoritePlays: review.FavoritePlays,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewTracks",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/tracks",
			ResponseType: GetUserYearReviewTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				tracks, page, err := app.UserService().GetUserYearReviewTracks(
					ctx,
					service.GetUserYearReviewTracksParams{
						UserId: c.Param("userId"),
						Year:   year,
						Page:   pageParams,
						Query:  queryParams,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewTracks{
					Page:   page,
					Tracks: make([]RankedTrack, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = RankedTrack{
						Track:     ConvertDBTrack(c, track.Track),
						Rank:      track.Rank,
						PlayCount: track.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewAlbums",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/albums",
			ResponseType: GetUserYearReviewAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				albums, page, err := app.UserService().
					GetUserYearReviewAlbums(
						ctx,
						service.GetUserYearReviewAlbumsParams{
							UserId: c.Param("userId"),
							Year:   year,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewAlbums{
					Page:   page,
					Albums: make([]RankedAlbum, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = RankedAlbum{
						Album:     ConvertDBAlbum(c, album.Album),
						Rank:      album.Rank,
						PlayCount: album.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewArtists",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/artists",
			ResponseType: GetUserYearReviewArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				artists, page, err := app.UserService().
					GetUserYearReviewArtists(
						ctx,
						service.GetUserYearReviewArtistsParams{
							UserId: c.Param("userId"),
							Year:   year,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewArtists{
					Page:    page,
					Artists: make([]RankedArtist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = RankedArtist{
						Artist:    ConvertDBArtist(c, artist.Artist),
						Rank:      artist.Rank,
						PlayCount: artist.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewTags",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/tags",
			ResponseType: GetUserYearReviewTags{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				tags, page, err := app.UserService().GetUserYearReviewTags(
					ctx,
					service.GetUserYearReviewTagsParams{
						UserId: c.Param("userId"),
						Year:   year,
						Page:   pageParams,
						Query:  queryParams,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewTags{
					Page: page,
					Tags: make([]RankedTag, len(tags)),
				}

				for i, tag := range tags {
					res.Tags[i] = RankedTag{
						TagSlug:   tag.TagSlug,
						Rank:      tag.Rank,
						PlayCount: tag.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserYearReviewDecades",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews/:year/decades",
			ResponseType: GetUserYearReviewDecades{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				ctx := c.Request().Context()

				year, err := parseIntParam(c, "year")
				if err != nil {
					return nil, err
				}

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				decades, page, err := app.UserService().
					GetUserYearReviewDecades(
						ctx,
						service.GetUserYearReviewDecadesParams{
							UserId: c.Param("userId"),
							Year:   year,
							Page:   pageParams,
							Query:  queryParams,
						},
					)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserYearReviewDecades{
					Page:    page,
					Decades: make([]RankedDecade, len(decades)),
				}

				for i, decade := range decades {
					res.Decades[i] = RankedDecade{
						Decade:    decade.Decade,
						Rank:      decade.Rank,
						PlayCount: decade.PlayCount,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUser",
			Method:       http.MethodGet,
			Path:         "/users/:userId",
			ResponseType: GetUser{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

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
					LastListenedAt: utils.SqlNullToInt64Ptr(
						stats.LastListenedAt),
					Updated: formatTime(stats.Updated),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserTopTracks",
			Method:       http.MethodGet,
			Path:         "/users/:userId/top-tracks",
			ResponseType: GetUserTopTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()
				q := c.Request().URL.Query()

				periodType := q.Get("period")
				if periodType == "" {
					periodType = "all"
				}

				year := parseIntQuery(q, "year", 0)
				limit := parseIntQuery(q, "limit", 5)

				tracks, err := app.UserService().GetUserTopTracks(
					ctx,
					service.GetUserTopTracksParams{
						UserId:     c.Param("userId"),
						PeriodType: periodType,
						Year:       year,
						Limit:      limit,
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				res := GetUserTopTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, t := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, t.Track)
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
					PictureUrl:  body.PictureUrl,
				})
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:   "UploadUserImage",
			Method: http.MethodPost,
			Path:   "/me/image/upload",
			Spec: pyrin.FormSpec{
				Files: map[string]pyrin.FormFileSpec{
					"image": {
						NumExpected: 1,
					},
				},
			},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				files, err := pyrin.FormFiles(c, "image")
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				err = app.UserService().UploadUserImage(
					ctx,
					service.UploadUserImageParams{
						UserId: user.Id,
						File:   files[0],
					},
				)
				if err != nil {
					return nil, handleUserServiceErrors(err)
				}

				return nil, nil
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
