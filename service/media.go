package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"slices"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/probe"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
)

var ErrInternalError = errors.New("internal error")

var (
	ErrMediaServiceTrackNotFound  = errors.New("media-service: track not found")
	ErrMediaServiceInvalidFormat  = errors.New("media-service: invalid format")
	ErrMediaServiceInvalidQuality = errors.New("media-service: invalid quality")
	ErrMediaServiceInvalidPolicy  = errors.New("media-service: invalid policy")
	ErrMediaServiceBitrateNotSet  = errors.New("media-service: bitrate not set")
)

func WrapInternalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInternalError, err)
}

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

func getDefaultQualityMapping() map[types.MediaFormat]QualitySpec {
	return map[types.MediaFormat]QualitySpec{
		types.MediaFormatFlac:     {},
		types.MediaFormatPcmS16LE: {},
		types.MediaFormatOpus: {
			High:   128,
			Medium: 96,
			Low:    64,
		},
		types.MediaFormatVorbis: {High: 192, Medium: 128, Low: 96},
		types.MediaFormatMp3:    {High: 320, Medium: 192, Low: 128},
		types.MediaFormatAac:    {High: 256, Medium: 192, Low: 96},
	}
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
			AllowedFormats: []types.MediaFormat{types.MediaFormatVorbis, types.MediaFormatMp3},
		},
	}
}

type MediaService struct {
	logger *slog.Logger

	db      *database.Database
	dataDir types.DataDir

	// TODO(patrik): Make this configureble through the config
	// TODO(patrik): I want to add test for this, test if every format has
	// quality specs set
	QualityMapping map[types.MediaFormat]QualitySpec
	DeviceSpecs    map[Device]DeviceSpec
}

func NewMediaService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
) *MediaService {
	return &MediaService{
		logger:         logger,
		db:             db,
		dataDir:        dataDir,
		QualityMapping: getDefaultQualityMapping(),
		DeviceSpecs:    getDefaultDeviceSpecs(),
	}
}

func (s *MediaService) getBitrateFromQuality(format types.MediaFormat, quality Quality) (int, error) {
	// TODO(patrik): Check for valid quality

	q, ok := s.QualityMapping[format]
	if !ok {
		// TODO(patrik): Better error
		// return 0, fmt.Errorf("format '%s' missing from quality mapping", format)
		return 0, WrapInternalError(errors.New("format not mapped '" + string(format) + "'"))
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

// TODO(patrik): Rename, ProcessTrackStream
func (s *MediaService) GetTrackStream(
	trackId string, 
	opts MediaStreamOptions,
) (string, error) {
	// TODO(patrik): Add a lock

	if opts.Policy == PolicyEmpty {
		opts.Policy = PolicyOriginal
	}

	if opts.Quality == QualityEmpty {
		// TODO(patrik): Set constant
		opts.Quality = QualityMedium
	}

	log := s.logger.With(
		slog.String("trackId", trackId),
		slog.Group("options",
			slog.String("policy", string(opts.Policy)),
			slog.String("device", string(opts.Device)),
			slog.String("quality", string(opts.Quality)),
			slog.String("format", string(opts.Format)),
		),
	)

	log.Info("track stream request")

	ctx := context.Background()

	track, err := s.db.GetTrackById(ctx, trackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", ErrMediaServiceTrackNotFound
		}

		return "", err
	}

	// TODO(patrik): Add checks for the track here

	format := types.MediaFormatUnknown

	switch opts.Policy {
	case PolicyOriginal:
		log.Info("track stream request using the original file")
		return track.Filename, nil
	case PolicyTranscode:
		if opts.Format == types.MediaFormatEmpty || !opts.Format.IsValid() {
			return "", ErrMediaServiceInvalidFormat
		}

		format = opts.Format
	case PolicyLossy:
		if opts.Device != DeviceEmpty {
			device, ok := s.DeviceSpecs[opts.Device]
			if !ok {
				return "", errors.New("unknown device: " + string(opts.Device))
			}

			if slices.Contains(device.AllowedFormats, track.MediaFormat) {
				format = track.MediaFormat
			}
		} else if track.MediaFormat.IsLossy() {
			format = track.MediaFormat
		} else if opts.Format != "" && opts.Format.IsValid() {
			format = opts.Format
		} else {
			format = types.MediaFormatOpus
		}
	default:
		return "", ErrMediaServiceInvalidPolicy
	}

	// NOTE(patrik): At this point we should have a valid format
	if !format.IsValid() {
		log.Error("no valid format is selected",
			slog.String("format", string(format)),
		)

		return "", errors.New("media-service: selected format is not valid: " + string(format))
	}

	// TODO(patrik): Maybe here we should still "transcode" the media but
	// just remove the metadata
	if format == track.MediaFormat {
		log.Info("track stream request selected format is matching the original track format")
		return track.Filename, nil
	}

	cacheDir := s.dataDir.Cache()
	trackCache := cacheDir.Track(track.Id)

	// Make sure that the cache directory is setup
	err = utils.CreateDirectories([]string{
		cacheDir.String(),
		cacheDir.Tracks(),
		trackCache,
	})
	if err != nil {
		return "", err
	}

	bitrate, err := s.getBitrateFromQuality(format, opts.Quality)
	if err != nil {
		return "", err
	}

	// At this point we should have a valid bitrate
	if format.IsLossy() && bitrate <= 0 {
		log.Error("bitrate not set for lossy format (something might be wrong with getBitrateFromQuality())",
			slog.String("format", string(format)),
		)

		// TODO(patrik): Better error
		return "", errors.New("media-service: bitrate not set for lossy format")
	}

	// filename: format-bitrate.ext
	ext, ok := format.ToExt()
	if !ok {
		log.Error("format has no extention",
			slog.String("format", string(format)),
		)

		return "", errors.New("media-service: format has no extention" + string(format))
	}

	filename := ""
	if !format.IsLossless() {
		filename = fmt.Sprintf("transcode-%s-%dk", format, bitrate) + ext
	} else {
		filename = fmt.Sprintf("transcode-%s-lossless", format) + ext
	}

	out := path.Join(trackCache, filename)
	tmpOut := path.Join(trackCache, "temp-"+filename)

	transcode := func() error {
		defer func() {
			if _, err := os.Stat(tmpOut); err == nil {
				os.Remove(tmpOut)
			}
		}()

		log.Info("track stream request starting transcoding process",
			slog.String("input", track.Filename),
			slog.String("output", out),
			slog.String("format", string(format)),
			slog.Int("bitrate", bitrate),
		)

		timer := utils.SimpleTimer{}
		timer.Start()

		args := []string{
			"-i", track.Filename,
			"-map", "0:a:0",
			"-vn",
			"-map_metadata", "-1",
			// TODO(patrik): Make option
			"-af", "loudnorm=I=-16:TP=-1.5:LRA=11",
		}

		switch format {
		case types.MediaFormatFlac:
			args = append(args, "-codec:a", "flac", "-compression_level", "5")
		case types.MediaFormatPcmS16LE:
			args = append(args, "-codec:a", "pcm_s16le")
		case types.MediaFormatOpus:
			args = append(args, "-codec:a", "libopus", "-b:a", fmt.Sprintf("%dk", bitrate), "-vbr", "on", "-compression_level", "10")
		case types.MediaFormatVorbis:
			args = append(args, "-codec:a", "libvorbis", "-b:a", fmt.Sprintf("%dk", bitrate))
		case types.MediaFormatMp3:
			args = append(args, "-codec:a", "libmp3lame", "-b:a", fmt.Sprintf("%dk", bitrate), "-q:a", "0")
		case types.MediaFormatAac:
			args = append(args, "-codec:a", "aac", "-b:a", fmt.Sprintf("%dk", bitrate), "-aac_coder", "twoloop", "-movflags", "+faststart")
		default:
			return errors.New("unsupported media format: " + string(format))
		}

		args = append(args, tmpOut)

		cmd := exec.Command("ffmpeg", args...)
		// TODO(patrik): Print when error
		// cmd.Stderr = os.Stderr
		// cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}

		err = os.Rename(tmpOut, out)
		if err != nil {
			return err
		}

		duration := timer.Stop()
		log.Info("track stream request transcoding done",
			slog.String("input", track.Filename),
			slog.String("output", out),
			slog.String("format", string(format)),
			slog.Int("bitrate", bitrate),
			slog.Duration("duration", duration),
		)

		return nil
	}

	_, err = os.Stat(out)
	if err != nil {
		if os.IsNotExist(err) {
			err := transcode()
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		log.Info("track stream request using the cached version")
	}

	return out, nil
}

func (s *MediaService) ProbeMedia(ctx context.Context, filepath string) (*probe.ProbeResult, error) {
	s.logger.Info("Probing media", "filepath", filepath)

	result, err := probe.ProbeMedia(ctx, filepath)
	if err != nil {
		s.logger.Info("failed to probe media", "err", err, "filepath", filepath)
		return nil, err
	}

	s.logger.Info("Probing result",
		"filepath", filepath,
		"format", result.MediaFormat,
		"duration", result.Duration,
	)

	return result, nil
}
