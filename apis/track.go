package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

type Track struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Order *int `json:"order"`

	Duration int64  `json:"duration"`
	Number   *int64 `json:"number"`
	Year     *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	AlbumId   string `json:"albumId"`
	AlbumName string `json:"albumName"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	artists := make([]ArtistInfo, len(track.FeaturingArtists.Data)+1)

	artists[0] = ArtistInfo{
		Id:   track.ArtistId,
		Name: track.ArtistName,
	}

	for i, v := range track.FeaturingArtists.Data {
		artists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}

	return Track{
		Id:        track.Id,
		Name:      track.Name,
		Order:     track.Order,
		Duration:  track.Duration,
		Number:    utils.SqlNullToInt64Ptr(track.Number),
		Year:      utils.SqlNullToInt64Ptr(track.Year),
		CoverArt:  ConvertAlbumCoverURL(c, track.AlbumId),
		AlbumId:   track.AlbumId,
		AlbumName: track.AlbumName,
		Artists:   artists,
		Tags:      utils.SplitTagString(track.Tags.String),
		Created:   formatTime(track.Created),
		Updated:   formatTime(track.Updated),
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetTrackById struct {
	Track Track `json:"track"`
}

func handleTrackServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrTrackServiceTrackNotFound):
		return TrackNotFound()
	case errors.Is(err, service.ErrTrackServiceFilterNotFound):
		return FilterNotFound()
	case errors.Is(err, service.ErrTrackServiceUnauthorized):
		return NotAuthorized()
	}

	var queryErr *database.QueryError
	if errors.As(err, &queryErr) {
		var filterErr, sortErr error
		if queryErr.Filter != nil {
			filterErr = queryErr.Filter
		}
		if queryErr.Sort != nil {
			sortErr = queryErr.Sort
		}
		return QueryError(filterErr, sortErr)
	}

	return err
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTracks",
			Method:       http.MethodGet,
			Path:         "/tracks",
			ResponseType: GetTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				ctx := c.Request().Context()

				pageParams := getPageParams(q, 100)
				queryParams := getQueryParams(q)

				tracks, p, err := app.TrackService().GetTracks(
					ctx,
					service.GetTracksParams{
						Page:     pageParams,
						Query:    queryParams,
						FilterId: q.Get("filterId"),
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				res := GetTracks{
					Page:   p,
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackById",
			Method:       http.MethodGet,
			Path:         "/tracks/:id",
			ResponseType: GetTrackById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {

				ctx := c.Request().Context()

				track, err := app.TrackService().GetTrackById(
					ctx,
					service.GetTrackByIdParams{
						TrackId: c.Param("id"),
					},
				)
				if err != nil {
					return nil, handleTrackServiceErrors(err)
				}

				return GetTrackById{
					Track: ConvertDBTrack(c, track),
				}, nil
			},
		},
	)
}
