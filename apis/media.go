package apis

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/pyrin"
)

type MediaRef struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MediaItem struct {
	TrackId string `json:"trackId"`
	Name    string `json:"name"`

	Artists []MediaRef `json:"artists"`
	Album   MediaRef   `json:"album"`

	CoverArt types.Images `json:"coverArt"`
}

type GetMedia struct {
	Items []MediaItem `json:"items"`
}

type GetMediaCommonBody struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetMediaFromPlaylistBody struct {
	GetMediaCommonBody

	FilterId string `json:"filterId"`
}

type GetMediaFromTaglistBody struct {
	GetMediaCommonBody
}

type GetMediaFromFilterBody struct {
	GetMediaCommonBody

	Filter string `json:"filter"`
}

type GetMediaFromArtistBody struct {
	GetMediaCommonBody
}

type GetMediaFromAlbumBody struct {
	GetMediaCommonBody
}

type GetMediaFromIdsBody struct {
	GetMediaCommonBody

	TrackIds  []string `json:"trackIds"`
	KeepOrder bool     `json:"keepOrder,omitempty"`
}

func packMediaResult(c pyrin.Context, tracks []database.Track) (GetMedia, error) {
	res := GetMedia{
		Items: make([]MediaItem, len(tracks)),
	}

	for i, track := range tracks {
		artists := make([]MediaRef, len(track.FeaturingArtists.Data)+1)

		artists[0] = MediaRef{
			Id:   track.ArtistId,
			Name: track.ArtistName,
		}

		for i, v := range track.FeaturingArtists.Data {
			artists[i+1] = MediaRef{
				Id:   v.Id,
				Name: v.Name,
			}
		}

		res.Items[i] = MediaItem{
			TrackId: track.Id,
			Name:    track.Name,
			Artists: artists,
			Album: MediaRef{
				Id:   track.AlbumId,
				Name: track.AlbumName,
			},
			CoverArt: ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		}
	}

	return res, nil
}

type MediaFormat struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Ext    string `json:"ext"`

	QualityHighBitrate   int `json:"qualityHighBitrate"`
	QualityMediumBitrate int `json:"qualityMediumBitrate"`
	QualityLowBitrate    int `json:"qualityLowBitrate"`

	Order int `json:"order"`
}

type MediaDeviceSpec struct {
	Name           string   `json:"name"`
	PreferedFormat string   `json:"preferedFormat"`
	AllowedFormats []string `json:"allowedFormats"`
}

type GetMediaSettings struct {
	Formats     []MediaFormat     `json:"formats"`
	DeviceSpecs []MediaDeviceSpec `json:"deviceSpecs"`
}

type AddTrackEventBody struct {
	Position float64 `json:"position"`
	// TODO(patrik): Validate
	Source string `json:"source"`
}

func handleMediaServiceErrors(err error) error {
	switch {
	case errors.Is(err, service.ErrMediaServiceTrackNotFound):
		return TrackNotFound()
	case errors.Is(err, service.ErrMediaServiceInvalidFormat):
		// TODO(patrik): Better error
		return &pyrin.Error{
			Code:    400,
			Type:    "MEDIA_INVALID_FORMAT",
			Message: "Invalid media format",
		}
	case errors.Is(err, service.ErrMediaServiceInvalidQuality):
		// TODO(patrik): Better error
		return &pyrin.Error{
			Code:    400,
			Type:    "MEDIA_INVALID_QUALITY",
			Message: "Invalid media quality",
		}
	case errors.Is(err, service.ErrMediaServiceInvalidPolicy):
		// TODO(patrik): Better error
		return &pyrin.Error{
			Code:    400,
			Type:    "MEDIA_INVALID_POLICY",
			Message: "Invalid media policy",
		}
	case errors.Is(err, service.ErrMediaServiceBitrateNotSet):
		// TODO(patrik): Better error
		return &pyrin.Error{
			Code:    400,
			Type:    "MEDIA_BITRATE_NOT_SET",
			Message: "Bitrate not set",
		}
	}

	return err
}

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.NormalHandler{
			Name:   "StreamTrack",
			Method: http.MethodGet,
			Path:   "/media/stream/tracks/:trackId",
			HandlerFunc: func(c pyrin.Context) error {
				trackId := c.Param("trackId")

				query := c.Request().URL.Query()

				filename, err := app.MediaService().GetTrackStream(
					trackId,
					service.MediaStreamOptions{
						Device:  service.Device(query.Get("device")),
						Policy:  service.Policy(query.Get("policy")),
						Quality: service.Quality(query.Get("quality")),
						Format:  types.MediaFormat(query.Get("format")),
					},
				)
				if err != nil {
					return handleMediaServiceErrors(err)
				}

				f := os.DirFS(filepath.Dir(filename))
				return pyrin.ServeFile(c, f, filepath.Base(filename))
			},
		},

		pyrin.ApiHandler{
			Name:         "AddTrackEvent",
			Method:       http.MethodPost,
			Path:         "/media/event/track/:trackId",
			ResponseType: nil,
			BodyType:     AddTrackEventBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				trackId := c.Param("trackId")

				body, err := pyrin.Body[AddTrackEventBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					// NOTE(patrik): This is expected behavior
					return nil, nil
				}

				track, err := app.DB().GetTrackById(ctx, trackId)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				percent := math.Min(body.Position/float64(track.Duration), 1.0)
				percent = math.Round(percent*100) / 100

				_, err = app.DB().CreateUserListeningEvent(
					ctx,
					database.CreateUserListeningEventParams{
						UserId:     user.Id,
						TrackId:    track.Id,
						ListenedAt: time.Now().UnixMilli(),
						Percent:    percent,
						PositionMs: int64(body.Position * 1000),
						Source:     body.Source,
					},
				)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaSettings",
			Method:       http.MethodGet,
			Path:         "/media/settings",
			ResponseType: GetMediaSettings{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				mediaService := app.MediaService()

				res := GetMediaSettings{
					Formats:     make([]MediaFormat, 0, len(types.ValidMediaFormats)),
					DeviceSpecs: make([]MediaDeviceSpec, 0, len(mediaService.DeviceSpecs)),
				}

				mappings := mediaService.QualityMapping

				// TODO(patrik): Should we handle checking for format?

				for _, format := range types.ValidMediaFormats {
					// TODO(patrik): Handle exists?
					quality, _ := mappings[format]
					info, _ := types.MediaFormatInfos[format]

					res.Formats = append(res.Formats, MediaFormat{
						Name:                 info.Name,
						Format:               string(format),
						Ext:                  info.Ext,
						QualityHighBitrate:   quality.High,
						QualityMediumBitrate: quality.Medium,
						QualityLowBitrate:    quality.Low,
						Order:                info.Order,
					})
				}

				sort.SliceStable(res.Formats, func(i, j int) bool {
					return res.Formats[i].Order < res.Formats[j].Order
				})

				for _, spec := range mediaService.DeviceSpecs {
					r := MediaDeviceSpec{
						Name:           spec.Name,
						PreferedFormat: string(spec.PreferedFormat),
						AllowedFormats: make([]string, len(spec.AllowedFormats)),
					}

					for i, f := range spec.AllowedFormats {
						r.AllowedFormats[i] = string(f)
					}

					res.DeviceSpecs = append(res.DeviceSpecs, r)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromFilter",
			Method:       http.MethodPost,
			Path:         "/media/filter",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromFilterBody](c)
				if err != nil {
					return nil, err
				}

				tracks, err := app.DB().GetAllTracks(ctx, body.Filter, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromArtist",
			Method:       http.MethodPost,
			Path:         "/media/artist/:artistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromArtistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				artistId := c.Param("artistId")

				ctx := context.TODO()

				// body, err := pyrin.Body[GetMediaFromArtistBody](c)
				// if err != nil {
				// 	return nil, err
				// }

				artist, err := app.DB().GetArtistById(ctx, artistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
				}

				subquery := database.ArtistTrackSubquery(artist.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromAlbum",
			Method:       http.MethodPost,
			Path:         "/media/album/:albumId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromAlbumBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				albumId := c.Param("albumId")

				ctx := context.TODO()

				// body, err := pyrin.Body[GetMediaFromAlbumBody](c)
				// if err != nil {
				// 	return nil, err
				// }

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				// sort := body.Sort
				// if sort == "" {
				// 	sort = "sort=number,name"
				// }

				subquery := database.AlbumTrackSubquery(album.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromIds",
			Method:       http.MethodPost,
			Path:         "/media/ids",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromIdsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromIdsBody](c)
				if err != nil {
					return nil, err
				}

				tracks, err := app.DB().GetTracksIn(ctx, body.TrackIds, "")
				if err != nil {
					return nil, err
				}

				if body.KeepOrder {
					trackMap := make(map[string]database.Track)
					for _, t := range tracks {
						trackMap[t.Id] = t
					}

					tracks = make([]database.Track, 0, len(body.TrackIds))
					for _, v := range body.TrackIds {
						track, exists := trackMap[v]
						if !exists {
							continue
						}

						tracks = append(tracks, track)
					}
				}

				return packMediaResult(c, tracks)
			},
		},
	)
}
