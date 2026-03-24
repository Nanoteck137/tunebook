package apis

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
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

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	artists := make([]ArtistInfo, len(track.FeaturingArtists)+1)

	artists[0] = ArtistInfo{
		Id:   track.ArtistId,
		Name: track.ArtistName,
	}

	for i, v := range track.FeaturingArtists {
		artists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}

	return Track{
		Id:        track.Id,
		Name:      track.Name,
		Duration:  track.Duration,
		Number:    utils.SqlNullToInt64Ptr(track.Number),
		Year:      utils.SqlNullToInt64Ptr(track.Year),
		CoverArt:  ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumId:   track.AlbumId,
		AlbumName: track.AlbumName,
		Artists:   artists,
		Tags:      utils.SplitString(track.Tags.String),
		Created:   track.Created,
		Updated:   track.Updated,
		Order:     track.Order,
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetTrackById struct {
	Track Track `json:"track"`
}

// TODO(patrik): Remove after PlaylistService
func getPageOptions(q url.Values) database.FetchOptions {
	perPage := 100
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

	return database.FetchOptions{
		Filter:  q.Get("filter"),
		Sort:    q.Get("sort"),
		PerPage: perPage,
		Page:    page,
	}
}

func handleTrackServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrTrackServiceTrackNotFound):
		return TrackNotFound()
	}

	var invalidFilter *service.InvalidFilterError
	if errors.As(err, &invalidFilter) {
		return InvalidFilter(errors.New(invalidFilter.Message))
	}

	var invalidSort *service.InvalidSortError
	if errors.As(err, &invalidSort) {
		return InvalidSort(errors.New(invalidSort.Message))
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
				filterParams := getFilterParams(q)

				tracks, p, err := app.TrackService().GetTracks(
					ctx,
					service.GetTracksParams{
						Page:     pageParams,
						Filter:   filterParams,
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
