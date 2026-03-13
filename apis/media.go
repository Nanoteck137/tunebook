package apis

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

// TODO(patrik): Change name?
type MediaResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MediaItem struct {
	Track   MediaResource   `json:"track"`
	Artists []MediaResource `json:"artists"`
	Album   MediaResource   `json:"album"`

	CoverArt types.Images `json:"coverArt"`

	MediaFormat types.MediaFormat `json:"mediaFormat"`
	MediaUrl    string            `json:"mediaUrl"`
}

type GetMedia struct {
	Items []MediaItem `json:"items"`
}

type GetMediaCommonBody struct {
	MediaFormat types.MediaFormat `json:"mediaFormat,omitempty"`

	Shuffle bool   `json:"shuffle,omitempty"`
	Sort    string `json:"sort,omitempty"`

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

// TODO(patrik): This might need some fixing, because this only returns the 
// original track stream so the user of the media api can't choose which 
// format to use
func packMediaResult(c pyrin.Context, tracks []database.Track, targetMediaFormat types.MediaFormat, shuffle bool) (GetMedia, error) {
	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
	}

	res := GetMedia{
		Items: make([]MediaItem, len(tracks)),
	}

	for i, track := range tracks {
		artists := make([]MediaResource, len(track.FeaturingArtists)+1)

		artists[0] = MediaResource{
			Id:   track.ArtistId,
			Name: track.ArtistName,
		}

		for i, v := range track.FeaturingArtists {
			artists[i+1] = MediaResource{
				Id:   v.Id,
				Name: v.Name,
			}
		}

		mediaFormat := types.MediaFormatUnknown
		if !track.MediaFormat.IsValid() {
			mediaFormat = track.MediaFormat
		}

		mediaUrl := ConvertURL(c, fmt.Sprintf("/media/tracks/%s/stream?policy=original", track.Id))

		res.Items[i] = MediaItem{
			Track: MediaResource{
				Id:   track.Id,
				Name: track.Name,
			},
			Artists: artists,
			Album: MediaResource{
				Id:   track.AlbumId,
				Name: track.AlbumName,
			},
			CoverArt:    ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
			MediaFormat: mediaFormat,
			MediaUrl:    mediaUrl,
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

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	group.Register(
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

				pretty.Println(res)

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromPlaylist",
			Method:       http.MethodPost,
			Path:         "/media/playlist/:playlistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromPlaylistBody{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				filter := ""

				if body.FilterId != "" {
					f, err := app.DB().GetPlaylistFilterById(ctx, body.FilterId, playlist.Id)
					if err != nil {
						// TODO(patrik): Handle err
						return nil, err
					}

					filter = f.Filter
				}

				tracks, err := app.DB().GetPlaylistTracks(ctx, playlist.Id, filter)
				if err != nil {
					return nil, err
				}

				// subquery := database.PlaylistTrackSubquery(playlist.Id)
				// tracks, err := app.DB().GetTracksIn(ctx, subquery, "")
				// if err != nil {
				// 	return nil, err
				// }

				// TODO(patrik): Better handling of this?
				t := make([]database.Track, len(tracks))
				for i, track := range tracks {
					t[i] = track.Track
				}

				return packMediaResult(c, t, body.MediaFormat, body.Shuffle)
			},
		},

		// TODO(patrik): Remove and remove from frontend also
		// pyrin.ApiHandler{
		// 	Name:         "GetMediaFromVirtualPlaylist",
		// 	Method:       http.MethodPost,
		// 	Path:         "/media/virtual-playlist/:virtualPlaylistId",
		// 	ResponseType: GetMedia{},
		// 	BodyType:     GetMediaFromTaglistBody{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		taglistId := c.Param("taglistId")
		//
		// 		ctx := context.TODO()
		//
		// 		body, err := pyrin.Body[GetMediaFromTaglistBody](c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		user, err := User(app, c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		_ = user
		// 		_ = body
		// 		_ = ctx
		// 		_ = taglistId
		//
		// 		// TODO(patrik): IMPLEMENT ME
		// 		// FIXME(patrik): FIXME
		// 		panic("IMPLEMENT ME")
		//
		// 		// taglist, err := app.DB().GetTaglistById(ctx, taglistId)
		// 		// if err != nil {
		// 		// 	if errors.Is(err, database.ErrItemNotFound) {
		// 		// 		return nil, TaglistNotFound()
		// 		// 	}
		// 		//
		// 		// 	return nil, err
		// 		// }
		// 		//
		// 		// if taglist.OwnerId != user.Id {
		// 		// 	return nil, TaglistNotFound()
		// 		// }
		// 		//
		// 		// tracks, err := app.DB().GetAllTracks(ctx, taglist.Filter, "")
		// 		// if err != nil {
		// 		// 	return nil, err
		// 		// }
		// 		//
		// 		// return packMediaResult(c, tracks, body.MediaFormat, body.Shuffle)
		// 		return nil, nil
		// 	},
		// },

		pyrin.ApiHandler{
			Name:         "GetMediaFromFilter",
			Method:       http.MethodPost,
			Path:         "/media/filter",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromFilterBody{},
			Errors:       []pyrin.ErrorType{},
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

				return packMediaResult(c, tracks, body.MediaFormat, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromArtist",
			Method:       http.MethodPost,
			Path:         "/media/artist/:artistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromArtistBody{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				artistId := c.Param("artistId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromArtistBody](c)
				if err != nil {
					return nil, err
				}

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

				return packMediaResult(c, tracks, body.MediaFormat, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromAlbum",
			Method:       http.MethodPost,
			Path:         "/media/album/:albumId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromAlbumBody{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				albumId := c.Param("albumId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromAlbumBody](c)
				if err != nil {
					return nil, err
				}

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				sort := body.Sort
				if sort == "" {
					sort = "sort=number,name"
				}

				subquery := database.AlbumTrackSubquery(album.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, sort)
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaFormat, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromIds",
			Method:       http.MethodPost,
			Path:         "/media/ids",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromIdsBody{},
			Errors:       []pyrin.ErrorType{},
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

				return packMediaResult(c, tracks, body.MediaFormat, body.Shuffle)
			},
		},
	)
}
