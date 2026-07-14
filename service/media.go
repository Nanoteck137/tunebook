package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"slices"
	"sync"

	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/probe"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var mediaErr = NewServiceErrCreator("media")

var (
	ErrMediaServiceTrackNotFound  = mediaErr.New("track not found")
	ErrMediaServiceInvalidFormat  = mediaErr.New("invalid format")
	ErrMediaServiceInvalidQuality = mediaErr.New("invalid quality")
	ErrMediaServiceInvalidPolicy  = mediaErr.New("invalid policy")
	ErrMediaServiceBitrateNotSet  = mediaErr.New("bitrate not set")
)

type Device string

const (
	DeviceEmpty   Device = ""
	DeviceAndroid Device = "android"
	DeviceIOS     Device = "ios"
)

type Policy string

const (
	PolicyEmpty     Policy = ""
	PolicyOriginal  Policy = "original"
	PolicyTranscode Policy = "transcode"
	PolicyLossy     Policy = "lossy"
)

type Quality string

const (
	QualityEmpty  Quality = ""
	QualityLow    Quality = "low"
	QualityMedium Quality = "medium"
	QualityHigh   Quality = "high"

	DefaultQuality = QualityMedium
)

type QualitySpec struct {
	High   int
	Medium int
	Low    int
}

func (s QualitySpec) MapFromQuality(q Quality) (int, bool) {
	switch q {
	case QualityHigh:
		return s.High, true
	case QualityMedium:
		return s.Medium, true
	case QualityLow:
		return s.Low, true
	}

	return 0, false
}

type DeviceSpec struct {
	Name           string
	PreferedFormat types.MediaFormat
	AllowedFormats []types.MediaFormat
}

func getDefaultDeviceSpecs() map[Device]DeviceSpec {
	return map[Device]DeviceSpec{
		DeviceIOS: {
			Name:           "IOS",
			PreferedFormat: types.MediaFormatAac,
			AllowedFormats: []types.MediaFormat{types.MediaFormatMp3},
		},
		DeviceAndroid: {
			Name:           "Android",
			PreferedFormat: types.MediaFormatOpus,
			AllowedFormats: []types.MediaFormat{
				types.MediaFormatVorbis, types.MediaFormatMp3},
		},
	}
}

type MediaService struct {
	logger *slog.Logger

	db         *database.Database
	filesystem *FilesystemService

	QualityMapping map[types.MediaFormat]QualitySpec
	DeviceSpecs    map[Device]DeviceSpec

	audioNormalization bool

	transcodeLocks sync.Map
}

func NewMediaService(
	logger *slog.Logger,
	db *database.Database,
	filesystem *FilesystemService,
	cfg *config.Config,
) *MediaService {
	s := &MediaService{
		logger:             logger,
		db:                 db,
		filesystem:         filesystem,
		DeviceSpecs:        getDefaultDeviceSpecs(),
		audioNormalization: cfg.Media.AudioNormalization,
	}

	s.QualityMapping = map[types.MediaFormat]QualitySpec{
		types.MediaFormatFlac:     {},
		types.MediaFormatPcmS16LE: {},
		types.MediaFormatOpus: {
			High:   cfg.Media.Opus.High,
			Medium: cfg.Media.Opus.Medium,
			Low:    cfg.Media.Opus.Low,
		},
		types.MediaFormatVorbis: {
			High:   cfg.Media.Vorbis.High,
			Medium: cfg.Media.Vorbis.Medium,
			Low:    cfg.Media.Vorbis.Low,
		},
		types.MediaFormatMp3: {
			High:   cfg.Media.Mp3.High,
			Medium: cfg.Media.Mp3.Medium,
			Low:    cfg.Media.Mp3.Low,
		},
		types.MediaFormatAac: {
			High:   cfg.Media.Aac.High,
			Medium: cfg.Media.Aac.Medium,
			Low:    cfg.Media.Aac.Low,
		},
	}

	return s
}

func (s *MediaService) getBitrateFromQuality(
	format types.MediaFormat,
	quality Quality,
) (int, error) {
	switch quality {
	case QualityLow, QualityMedium, QualityHigh:
	default:
		return 0, ErrMediaServiceInvalidQuality
	}

	q, ok := s.QualityMapping[format]
	if !ok {
		return 0, mediaErr.Newf("format not mapped '%s'", format)
	}

	bitrate, ok := q.MapFromQuality(quality)
	if !ok {
		return 0, ErrMediaServiceInvalidQuality
	}

	return bitrate, nil
}

type MediaStreamOptions struct {
	Policy  Policy
	Device  Device
	Quality Quality

	Format types.MediaFormat
}

func (s *MediaService) ProcessTrackStream(
	trackId string,
	opts MediaStreamOptions,
) (string, error) {
	lock, _ := s.transcodeLocks.LoadOrStore(trackId, &sync.Mutex{})
	mu := lock.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	if opts.Policy == PolicyEmpty {
		opts.Policy = PolicyOriginal
	}

	if opts.Quality == QualityEmpty {
		opts.Quality = DefaultQuality
	}

	ctx := context.Background()

	track, err := s.db.GetTrackById(ctx, trackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", ErrMediaServiceTrackNotFound
		}

		return "", err
	}

	if track.Filename == "" {
		return "", mediaErr.New("track has no file")
	}

	if !track.MediaFormat.IsValid() {
		return "", mediaErr.Newf(
			"invalid media format: %s", track.MediaFormat)
	}

	if opts.Policy == PolicyOriginal {
		s.logger.Info("using original file", "trackId", trackId)
		return track.Filename, nil
	}

	format, err := s.resolveStreamFormat(track, opts)
	if err != nil {
		return "", err
	}

	if format == track.MediaFormat {
		s.logger.Info(
			"format matches original", "trackId", trackId, "format", format)
		return track.Filename, nil
	}

	err = s.filesystem.EnsureTranscodingTrackDir(trackId)
	if err != nil {
		return "", err
	}

	bitrate, err := s.getBitrateFromQuality(format, opts.Quality)
	if err != nil {
		return "", err
	}

	if format.IsLossy() && bitrate <= 0 {
		return "", ErrMediaServiceBitrateNotSet
	}

	ext, ok := format.ToExt()
	if !ok {
		return "", mediaErr.Newf("format has no extension: %s", format)
	}

	filename := fmt.Sprintf("transcode-%s-%dk%s", format, bitrate, ext)
	if format.IsLossless() {
		filename = fmt.Sprintf("transcode-%s-lossless%s", format, ext)
	}

	out := s.filesystem.TranscodingPath(trackId, filename)

	exists, err := s.filesystem.FileExists(out)
	if err != nil {
		return "", err
	}

	if !exists {
		err := s.transcodeTrack(transcodeTrackParams{
			track:   track,
			format:  format,
			bitrate: bitrate,
			out:     out,
		})
		if err != nil {
			return "", err
		}
	} else {
		s.logger.Info("using cached version", "trackId", trackId)
	}

	return out, nil
}

func (s *MediaService) resolveStreamFormat(
	track database.Track,
	opts MediaStreamOptions,
) (types.MediaFormat, error) {
	switch opts.Policy {
	case PolicyTranscode:
		if opts.Format == types.MediaFormatEmpty || !opts.Format.IsValid() {
			return "", ErrMediaServiceInvalidFormat
		}

		return opts.Format, nil
	case PolicyLossy:
		format := types.MediaFormatOpus

		if opts.Device != DeviceEmpty {
			device, ok := s.DeviceSpecs[opts.Device]
			if !ok {
				return "", mediaErr.Newf("device not found: %s", opts.Device)
			}

			if slices.Contains(device.AllowedFormats, track.MediaFormat) {
				format = track.MediaFormat
			}
		} else if track.MediaFormat.IsLossy() {
			format = track.MediaFormat
		} else if opts.Format != "" && opts.Format.IsValid() {
			format = opts.Format
		}

		if !format.IsValid() {
			return "", mediaErr.Newf("invalid format: %s", format)
		}

		return format, nil
	default:
		return "", ErrMediaServiceInvalidPolicy
	}
}

type transcodeTrackParams struct {
	track   database.Track
	format  types.MediaFormat
	bitrate int
	out     string
}

func (s *MediaService) transcodeTrack(
	params transcodeTrackParams,
) error {
	tmpOut := path.Join(path.Dir(params.out), "temp-"+path.Base(params.out))
	defer os.RemoveAll(tmpOut)

	s.logger.Info("starting transcode",
		"input", params.track.Filename,
		"output", params.out,
		"format", params.format,
		"bitrate", params.bitrate,
	)

	timer := utils.SimpleTimer{}
	timer.Start()

	args := []string{
		"-i", params.track.Filename,
		"-map", "0:a:0",
		"-vn",
		"-map_metadata", "-1",
	}

	if s.audioNormalization {
		args = append(args, "-af", "loudnorm=I=-16:TP=-1.5:LRA=11")
	}

	switch params.format {
	case types.MediaFormatFlac:
		args = append(args, "-codec:a", "flac", "-compression_level", "5")
	case types.MediaFormatPcmS16LE:
		args = append(args, "-codec:a", "pcm_s16le")
	case types.MediaFormatOpus:
		args = append(
			args,
			"-codec:a", "libopus",
			"-b:a", fmt.Sprintf("%dk", params.bitrate),
			"-vbr", "on",
			"-compression_level", "10",
		)
	case types.MediaFormatVorbis:
		args = append(
			args,
			"-codec:a", "libvorbis",
			"-b:a", fmt.Sprintf("%dk", params.bitrate),
		)
	case types.MediaFormatMp3:
		args = append(
			args,
			"-codec:a", "libmp3lame",
			"-b:a", fmt.Sprintf("%dk", params.bitrate),
			"-q:a", "0",
		)
	case types.MediaFormatAac:
		args = append(
			args,
			"-codec:a", "aac",
			"-b:a", fmt.Sprintf("%dk", params.bitrate),
			"-aac_coder", "twoloop",
			"-movflags", "+faststart",
		)
	default:
		return mediaErr.Newf("unsupported media format: %s", params.format)
	}

	args = append(args, tmpOut)

	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		s.logger.Error("ffmpeg failed",
			"stderr", stderr.String(),
			"err", err,
		)
		return err
	}

	err = os.Rename(tmpOut, params.out)
	if err != nil {
		return err
	}

	duration := timer.Stop()
	s.logger.Info("transcode done",
		"input", params.track.Filename,
		"output", params.out,
		"format", params.format,
		"bitrate", params.bitrate,
		"duration", duration,
	)

	return nil
}

func (s *MediaService) ProbeMedia(
	ctx context.Context,
	filepath string,
) (*probe.ProbeResult, error) {
	result, err := probe.ProbeMedia(ctx, filepath)
	if err != nil {
		s.logger.Error(
			"probe media", "err", err, "filepath", filepath)
		return nil, err
	}

	s.logger.Debug("probe media",
		"filepath", filepath,
		"format", result.MediaFormat,
		"duration", result.Duration,
	)

	return result, nil
}
