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

// type ReviewTrack struct {
// 	Rank      int   `json:"rank"`
// 	PlayCount int   `json:"playCount"`
// 	Track     Track `json:"track"`
// }
//
// type ReviewAlbum struct {
// 	Rank      int   `json:"rank"`
// 	PlayCount int   `json:"playCount"`
// 	Album     Album `json:"album"`
// }
//
// type ReviewArtist struct {
// 	Rank      int    `json:"rank"`
// 	PlayCount int    `json:"playCount"`
// 	Artist    Artist `json:"artist"`
// }
//
// type ReviewInnerTrack struct {
// 	Rank      int   `json:"rank"`
// 	PlayCount int   `json:"playCount"`
// 	PlayTime  int64 `json:"playTime"`
// 	Track     Track `json:"track"`
// }
//
// type ReviewArtistTracks struct {
// 	Artist Artist             `json:"artist"`
// 	Tracks []ReviewInnerTrack `json:"tracks"`
// }
//
// type ReviewAlbumTracks struct {
// 	Album  Album              `json:"album"`
// 	Tracks []ReviewInnerTrack `json:"tracks"`
// }
//
// type ReviewMonth struct {
// 	Month     int   `json:"month"`
// 	PlayCount int   `json:"playCount"`
// 	PlayTime  int64 `json:"playTime"`
//
// 	AvgCompletion float64 `json:"avgCompletion"`
// 	SkipCount     int     `json:"skipCount"`
// 	UniqueTracks  int     `json:"uniqueTracks"`
// 	FavoritePlays int     `json:"favoritePlays"`
//
// 	TrackCount  int `json:"trackCount"`
// 	AlbumCount  int `json:"albumCount"`
// 	ArtistCount int `json:"artistCount"`
//
// 	Tracks  []ReviewTrack  `json:"tracks"`
// 	Albums  []ReviewAlbum  `json:"albums"`
// 	Artists []ReviewArtist `json:"artists"`
//
// 	Tags    []ReviewTag    `json:"tags"`
// 	Decades []ReviewDecade `json:"decades"`
// }

// type ReviewTag struct {
// 	TagSlug   string `json:"tagSlug"`
// 	Rank      int    `json:"rank"`
// 	PlayCount int    `json:"playCount"`
// }
//
// type ReviewDecade struct {
// 	Decade    int `json:"decade"`
// 	Rank      int `json:"rank"`
// 	PlayCount int `json:"playCount"`
// }

// type GetUserYearReview struct {
// 	Review YearStat `json:"review"`
//
// 	Tracks  []ReviewTrack  `json:"tracks"`
// 	Albums  []ReviewAlbum  `json:"albums"`
// 	Artists []ReviewArtist `json:"artists"`
//
// 	TrackCount  int `json:"trackCount"`
// 	AlbumCount  int `json:"albumCount"`
// 	ArtistCount int `json:"artistCount"`
//
// 	ArtistTracks []ReviewArtistTracks `json:"artistTracks"`
// 	AlbumTracks  []ReviewAlbumTracks  `json:"albumTracks"`
//
// 	Months []ReviewMonth `json:"months"`
//
// 	Tags    []ReviewTag    `json:"tags"`
// 	Decades []ReviewDecade `json:"decades"`
// }

// type GetUserYearReviewTopTracks struct {
// 	Tracks []ReviewTrack `json:"tracks"`
// }
//
// type GetUserYearReviewTopAlbums struct {
// 	Albums []ReviewAlbum `json:"albums"`
// }
//
// type GetUserYearReviewTopArtists struct {
// 	Artists []ReviewArtist `json:"artists"`
// }
//
// type GetUserYearReviewMonthTopTracks struct {
// 	Tracks []ReviewTrack `json:"tracks"`
// }
//
// type GetUserYearReviewMonthTopAlbums struct {
// 	Albums []ReviewAlbum `json:"albums"`
// }
//
// type GetUserYearReviewMonthTopArtists struct {
// 	Artists []ReviewArtist `json:"artists"`
// }

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

	// Tracks  []ReviewTrack  `json:"tracks"`
	// Albums  []ReviewAlbum  `json:"albums"`
	// Artists []ReviewArtist `json:"artists"`
	//
	// TrackCount  int `json:"trackCount"`
	// AlbumCount  int `json:"albumCount"`
	// ArtistCount int `json:"artistCount"`
	//
	// ArtistTracks []ReviewArtistTracks `json:"artistTracks"`
	// AlbumTracks  []ReviewAlbumTracks  `json:"albumTracks"`
	//
	// Months []ReviewMonth `json:"months"`
	//
	// Tags    []ReviewTag    `json:"tags"`
	// Decades []ReviewDecade `json:"decades"`
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

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAllUserYearReviews",
			Method:       http.MethodGet,
			Path:         "/users/:userId/year-reviews",
			ResponseType: GetAllUserYearReviews{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := c.Request().Context()

				reviews, err := app.UserService().GetAllUserYearReviewsParams(
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

				for i, track  := range tracks {
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

		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearStats",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/year-stats",
		// 	ResponseType: GetUserYearStats{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		stats, err := app.UserService().GetUserYearStats(
		// 			ctx,
		// 			service.GetUserYearStatsParams{
		// 				UserId: c.Param("userId"),
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearStats{
		// 			Stats: make([]YearStat, len(stats)),
		// 		}
		//
		// 		for i, s := range stats {
		// 			res.Stats[i] = YearStat{
		// 				Year:          s.Year,
		// 				TrackCount:    s.TrackCount,
		// 				ListeningTime: s.ListeningTime,
		// 			}
		//
		// 			pretty.Println(s)
		// 		}
		//
		// 		pretty.Println(res)
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReview",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year",
		// 	ResponseType: GetUserYearReview{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		review, err := app.UserService().GetUserYearReview(
		// 			ctx,
		// 			service.GetUserYearReviewParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReview{
		// 			Review: YearStat{
		// 				Year:          review.Review.Year,
		// 				TrackCount:    review.Review.TrackCount,
		// 				ListeningTime: review.Review.ListeningTime,
		//
		// 				AvgCompletion: review.Review.AvgCompletion,
		// 				SkipCount:     review.Review.SkipCount,
		// 				UniqueTracks:  review.Review.UniqueTracks,
		// 				FavoritePlays: review.Review.FavoritePlays,
		// 			},
		//
		// 			Tracks:  make([]ReviewTrack, len(review.Tracks)),
		// 			Albums:  make([]ReviewAlbum, len(review.Albums)),
		// 			Artists: make([]ReviewArtist, len(review.Artists)),
		//
		// 			TrackCount:  review.TrackCount,
		// 			AlbumCount:  review.AlbumCount,
		// 			ArtistCount: review.ArtistCount,
		//
		// 			Months: make([]ReviewMonth, len(review.MonthDetails)),
		// 		}
		//
		// 		for i, t := range review.Tracks {
		// 			res.Tracks[i] = ReviewTrack{
		// 				Rank:      t.Rank,
		// 				PlayCount: t.PlayCount,
		// 				Track:     ConvertDBTrack(c, t.Track),
		// 			}
		// 		}
		//
		// 		for i, a := range review.Albums {
		// 			res.Albums[i] = ReviewAlbum{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Album:     ConvertDBAlbum(c, a.Album),
		// 			}
		// 		}
		//
		// 		for i, a := range review.Artists {
		// 			res.Artists[i] = ReviewArtist{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Artist:    ConvertDBArtist(c, a.Artist),
		// 			}
		// 		}
		//
		// 		res.ArtistTracks = make([]ReviewArtistTracks, 0, len(review.Artists))
		// 		for _, a := range review.Artists {
		// 			inner := review.ArtistTracks[a.Artist.Id]
		//
		// 			tracks := make([]ReviewInnerTrack, len(inner))
		// 			for i, t := range inner {
		// 				tracks[i] = ReviewInnerTrack{
		// 					Rank:      t.Rank,
		// 					PlayCount: t.PlayCount,
		// 					PlayTime:  t.PlayTime,
		// 					Track:     ConvertDBTrack(c, t.Track),
		// 				}
		// 			}
		//
		// 			res.ArtistTracks = append(res.ArtistTracks, ReviewArtistTracks{
		// 				Artist: ConvertDBArtist(c, a.Artist),
		// 				Tracks: tracks,
		// 			})
		// 		}
		//
		// 		res.AlbumTracks = make([]ReviewAlbumTracks, 0, len(review.Albums))
		// 		for _, a := range review.Albums {
		// 			inner := review.AlbumTracks[a.Album.Id]
		//
		// 			tracks := make([]ReviewInnerTrack, len(inner))
		// 			for i, t := range inner {
		// 				tracks[i] = ReviewInnerTrack{
		// 					Rank:      t.Rank,
		// 					PlayCount: t.PlayCount,
		// 					PlayTime:  t.PlayTime,
		// 					Track:     ConvertDBTrack(c, t.Track),
		// 				}
		// 			}
		//
		// 			res.AlbumTracks = append(res.AlbumTracks, ReviewAlbumTracks{
		// 				Album:  ConvertDBAlbum(c, a.Album),
		// 				Tracks: tracks,
		// 			})
		// 		}
		//
		// 		for i, m := range review.MonthDetails {
		// 			tracks := make([]ReviewTrack, len(m.Tracks))
		// 			for j, t := range m.Tracks {
		// 				tracks[j] = ReviewTrack{
		// 					Rank:      t.Rank,
		// 					PlayCount: t.PlayCount,
		// 					Track:     ConvertDBTrack(c, t.Track),
		// 				}
		// 			}
		//
		// 			albums := make([]ReviewAlbum, len(m.Albums))
		// 			for j, a := range m.Albums {
		// 				albums[j] = ReviewAlbum{
		// 					Rank:      a.Rank,
		// 					PlayCount: a.PlayCount,
		// 					Album:     ConvertDBAlbum(c, a.Album),
		// 				}
		// 			}
		//
		// 			artists := make([]ReviewArtist, len(m.Artists))
		// 			for j, a := range m.Artists {
		// 				artists[j] = ReviewArtist{
		// 					Rank:      a.Rank,
		// 					PlayCount: a.PlayCount,
		// 					Artist:    ConvertDBArtist(c, a.Artist),
		// 				}
		// 			}
		//
		// 			tags := make([]ReviewTag, len(m.Tags))
		// 			for j, t := range m.Tags {
		// 				tags[j] = ReviewTag{
		// 					TagSlug:   t.TagSlug,
		// 					Rank:      t.Rank,
		// 					PlayCount: t.PlayCount,
		// 				}
		// 			}
		//
		// 			decades := make([]ReviewDecade, len(m.Decades))
		// 			for j, d := range m.Decades {
		// 				decades[j] = ReviewDecade{
		// 					Decade:    d.Decade,
		// 					Rank:      d.Rank,
		// 					PlayCount: d.PlayCount,
		// 				}
		// 			}
		//
		// 			res.Months[i] = ReviewMonth{
		// 				Month:     m.Month,
		// 				PlayCount: m.PlayCount,
		// 				PlayTime:  m.PlayTime,
		//
		// 				AvgCompletion: m.AvgCompletion,
		// 				SkipCount:     m.SkipCount,
		// 				UniqueTracks:  m.UniqueTracks,
		// 				FavoritePlays: m.FavoritePlays,
		//
		// 				TrackCount:  m.TrackCount,
		// 				AlbumCount:  m.AlbumCount,
		// 				ArtistCount: m.ArtistCount,
		//
		// 				Tracks:  tracks,
		// 				Albums:  albums,
		// 				Artists: artists,
		//
		// 				Tags:    tags,
		// 				Decades: decades,
		// 			}
		// 		}
		//
		// 		res.Tags = make([]ReviewTag, len(review.Tags))
		// 		for i, t := range review.Tags {
		// 			res.Tags[i] = ReviewTag{
		// 				TagSlug:   t.TagSlug,
		// 				Rank:      t.Rank,
		// 				PlayCount: t.PlayCount,
		// 			}
		// 		}
		//
		// 		res.Decades = make([]ReviewDecade, len(review.Decades))
		// 		for i, d := range review.Decades {
		// 			res.Decades[i] = ReviewDecade{
		// 				Decade:    d.Decade,
		// 				Rank:      d.Rank,
		// 				PlayCount: d.PlayCount,
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewTopTracks",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/top-tracks",
		// 	ResponseType: GetUserYearReviewTopTracks{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewTopTracks(
		// 			ctx,
		// 			service.GetUserYearReviewTopTracksParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewTopTracks{
		// 			Tracks: make([]ReviewTrack, len(top.Tracks)),
		// 		}
		// 		for i, t := range top.Tracks {
		// 			res.Tracks[i] = ReviewTrack{
		// 				Rank:      t.Rank,
		// 				PlayCount: t.PlayCount,
		// 				Track:     ConvertDBTrack(c, t.Track),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewTopAlbums",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/top-albums",
		// 	ResponseType: GetUserYearReviewTopAlbums{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewTopAlbums(
		// 			ctx,
		// 			service.GetUserYearReviewTopAlbumsParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewTopAlbums{
		// 			Albums: make([]ReviewAlbum, len(top.Albums)),
		// 		}
		// 		for i, a := range top.Albums {
		// 			res.Albums[i] = ReviewAlbum{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Album:     ConvertDBAlbum(c, a.Album),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewTopArtists",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/top-artists",
		// 	ResponseType: GetUserYearReviewTopArtists{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewTopArtists(
		// 			ctx,
		// 			service.GetUserYearReviewTopArtistsParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewTopArtists{
		// 			Artists: make([]ReviewArtist, len(top.Artists)),
		// 		}
		// 		for i, a := range top.Artists {
		// 			res.Artists[i] = ReviewArtist{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Artist:    ConvertDBArtist(c, a.Artist),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewMonthTopTracks",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/months/:month/top-tracks",
		// 	ResponseType: GetUserYearReviewMonthTopTracks{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		month, err := parseIntParam(c, "month")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		if month < 1 || month > 12 {
		// 			return nil, InvalidParam("month")
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewMonthTopTracks(
		// 			ctx,
		// 			service.GetUserYearReviewMonthTopTracksParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 				Month:  month,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewMonthTopTracks{
		// 			Tracks: make([]ReviewTrack, len(top.Tracks)),
		// 		}
		// 		for i, t := range top.Tracks {
		// 			res.Tracks[i] = ReviewTrack{
		// 				Rank:      t.Rank,
		// 				PlayCount: t.PlayCount,
		// 				Track:     ConvertDBTrack(c, t.Track),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewMonthTopAlbums",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/months/:month/top-albums",
		// 	ResponseType: GetUserYearReviewMonthTopAlbums{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		month, err := parseIntParam(c, "month")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		if month < 1 || month > 12 {
		// 			return nil, InvalidParam("month")
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewMonthTopAlbums(
		// 			ctx,
		// 			service.GetUserYearReviewMonthTopAlbumsParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 				Month:  month,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewMonthTopAlbums{
		// 			Albums: make([]ReviewAlbum, len(top.Albums)),
		// 		}
		// 		for i, a := range top.Albums {
		// 			res.Albums[i] = ReviewAlbum{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Album:     ConvertDBAlbum(c, a.Album),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:         "GetUserYearReviewMonthTopArtists",
		// 	Method:       http.MethodGet,
		// 	Path:         "/users/:userId/reviews/:year/months/:month/top-artists",
		// 	ResponseType: GetUserYearReviewMonthTopArtists{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		month, err := parseIntParam(c, "month")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		if month < 1 || month > 12 {
		// 			return nil, InvalidParam("month")
		// 		}
		//
		// 		top, err := app.UserService().GetUserYearReviewMonthTopArtists(
		// 			ctx,
		// 			service.GetUserYearReviewMonthTopArtistsParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 				Month:  month,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		res := GetUserYearReviewMonthTopArtists{
		// 			Artists: make([]ReviewArtist, len(top.Artists)),
		// 		}
		// 		for i, a := range top.Artists {
		// 			res.Artists[i] = ReviewArtist{
		// 				Rank:      a.Rank,
		// 				PlayCount: a.PlayCount,
		// 				Artist:    ConvertDBArtist(c, a.Artist),
		// 			}
		// 		}
		//
		// 		return res, nil
		// 	},
		// },
		//
		// pyrin.ApiHandler{
		// 	Name:   "GenerateUserReview",
		// 	Method: http.MethodPost,
		// 	Path:   "/users/:userId/reviews/:year",
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		ctx := c.Request().Context()
		//
		// 		year, err := parseIntParam(c, "year")
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		err = app.UserService().GenerateUserReview(
		// 			ctx,
		// 			service.GenerateUserReviewParams{
		// 				UserId: c.Param("userId"),
		// 				Year:   year,
		// 			},
		// 		)
		// 		if err != nil {
		// 			return nil, handleUserServiceErrors(err)
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },
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
