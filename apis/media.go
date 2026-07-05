package apis

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

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

				filename, err := app.MediaService().ProcessTrackStream(
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
			Name:         "GetMediaSettings",
			Method:       http.MethodGet,
			Path:         "/media/settings",
			ResponseType: GetMediaSettings{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				mediaService := app.MediaService()

				res := GetMediaSettings{
					Formats: make(
						[]MediaFormat, 0, len(types.ValidMediaFormats)),
					DeviceSpecs: make(
						[]MediaDeviceSpec, 0, len(mediaService.DeviceSpecs)),
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
	)
}
