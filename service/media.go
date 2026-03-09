package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
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
	PreferedFormat types.MediaFormat
	AllowedFormats []types.MediaFormat
}

type MediaService struct {
	db      *database.Database
	workDir types.WorkDir
}

func NewMediaService(db *database.Database, workDir types.WorkDir) *MediaService {
	return &MediaService{
		db:      db,
		workDir: workDir,
	}
}

type MediaStreamOptions struct {
	Device  Device
	Policy  Policy
	Quality Quality

	Format  types.MediaFormat
	Bitrate int
}

func (s *MediaService) GetTrackStream(trackId string, opts MediaStreamOptions) (string, error) {
	// TODO(patrik): Add a lock

	ctx := context.Background()

	track, err := s.db.GetTrackById(ctx, trackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			// TODO(patrik): Better error
			return "", errors.New("track not found")
		}

		return "", err
	}

	// TODO(patrik): Add checks for the track here

	pretty.Println(opts)

	useOriginal := false
	format := types.MediaFormatUnknown
	bitrate := 0

	// TODO(patrik): I want to add test for this, test if every format has
	// quality specs set
	qualityMapping := map[types.MediaFormat]QualitySpec{
		types.MediaFormatFlac: {},
		types.MediaFormatWav:  {},
		types.MediaFormatOpus: {
			High:   128,
			Medium: 96,
			Low:    64,
		},
		types.MediaFormatVorbis: {High: 192, Medium: 128, Low: 96},
		types.MediaFormatMp3:    {High: 320, Medium: 192, Low: 128},
		types.MediaFormatAac:    {High: 256, Medium: 192, Low: 96},
	}

	devices := map[Device]DeviceSpec{
		DeviceIOS: DeviceSpec{
			PreferedFormat: types.MediaFormatAac,
			AllowedFormats: []types.MediaFormat{types.MediaFormatMp3},
		},
		DeviceAndroid: DeviceSpec{
			PreferedFormat: types.MediaFormatOpus,
			AllowedFormats: []types.MediaFormat{types.MediaFormatVorbis, types.MediaFormatMp3},
		},
	}

	switch opts.Policy {
	case PolicyEmpty, PolicyOriginal:
		useOriginal = true
	case PolicyTranscode:
		if !opts.Format.IsValid() {
			return "", errors.New("invalid format")
		}

		format = opts.Format

		if !opts.Format.IsLossless() {
			if opts.Bitrate > 0 {
				bitrate = opts.Bitrate
			} else if opts.Quality != QualityEmpty {
				// TODO(patrik): Check for valid quality

				quality, ok := qualityMapping[format]
				if !ok {
					return "", fmt.Errorf("format '%s' missing from quality mapping", format)
				}

				bitrate, ok = quality.MapFromQuality(opts.Quality)
				if !ok {
					return "", errors.New("invalid quality: " + string(opts.Quality))
				}
			} else {
				quality, ok := qualityMapping[format]
				if !ok {
					panic("Format missing qualities")
				}

				// TODO(patrik): Make constant for default quality
				bitrate = quality.Medium
			}
		}
	case PolicyLossy:
		if opts.Device != DeviceEmpty {
			device, ok := devices[opts.Device]
			if !ok {
				return "", errors.New("unknown device: " + string(opts.Device))
			}

			format = device.PreferedFormat

			for _, allowed := range device.AllowedFormats {
				if track.MediaFormat == allowed {
					format = track.MediaFormat
					break
				}
			}

			if opts.Quality != QualityEmpty {
				// TODO(patrik): Check for valid quality

				quality, ok := qualityMapping[format]
				if !ok {
					return "", fmt.Errorf("format '%s' missing from quality mapping", format)
				}

				bitrate, ok = quality.MapFromQuality(opts.Quality)
				if !ok {
					return "", errors.New("invalid quality: " + string(opts.Quality))
				}
			} else {
				quality, ok := qualityMapping[format]
				if !ok {
					return "", fmt.Errorf("format '%s' missing from quality mapping", format)
				}

				// TODO(patrik): Make constant for default quality
				bitrate = quality.Medium
			}
		} else if track.MediaFormat.IsLossy() {
			format = track.MediaFormat
		} else if opts.Format != "" && opts.Format.IsValid() {
			format = opts.Format

			if opts.Bitrate > 0 {
				bitrate = opts.Bitrate
			} else if opts.Quality != QualityEmpty {
				// TODO(patrik): Check for valid quality

				quality, ok := qualityMapping[format]
				if !ok {
					return "", fmt.Errorf("format '%s' missing from quality mapping", format)
				}

				bitrate, ok = quality.MapFromQuality(opts.Quality)
				if !ok {
					return "", errors.New("invalid quality: " + string(opts.Quality))
				}
			} else {
				quality, ok := qualityMapping[format]
				if !ok {
					return "", fmt.Errorf("format '%s' missing from quality mapping", format)
				}

				// TODO(patrik): Make constant for default quality
				bitrate = quality.Medium
			}
		} else {
			format = types.MediaFormatOpus

			quality, ok := qualityMapping[format]
			// TODO(patrik): If this ever happens, there is a bug in the
			// code, and not user error
			if !ok {
				panic("Interal: Format missing qualities")
			}

			// TODO(patrik): Make constant for default quality
			bitrate = quality.Medium
		}
	default:
		// TODO(patrik): Fix
		panic("Unknown policy: " + opts.Policy)
	}

	fmt.Printf("track.MediaFormat: %v\n", track.MediaFormat)
	fmt.Printf("format: %v\n", format)
	fmt.Printf("bitrate: %v\n", bitrate)
	fmt.Printf("useOriginal: %v\n", useOriginal)

	// TODO(patrik): Maybe here we should still "transcode" the media but
	// just remove the metadata
	if format == track.MediaFormat {
		useOriginal = true
	}

	fmt.Printf("useOriginal (after track check): %v\n", useOriginal)

	cacheDir := s.workDir.Cache()
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

	if useOriginal {
		return track.Filename, nil
	}

	if format.IsLossy() && bitrate <= 0 {
		// TODO(patrik): Better error
		return "", errors.New("bitrate not set")
	}

	// filename: format-bitrate.ext
	ext, ok := format.ToExt()
	if !ok {
		return "", errors.New("unsupported media format: " + string(format))
	}

	filename := ""
	if !format.IsLossless() {
		filename = fmt.Sprintf("transcode-%s-%dk", format, bitrate) + ext
	} else {
		filename = fmt.Sprintf("transcode-%s-lossless", format) + ext
	}
	fmt.Printf("filename: %v\n", filename)

	out := path.Join(trackCache, filename)

	transcode := func() error {
		// TODO(patrik): Don't map the metadata
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
		case types.MediaFormatWav:
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

		args = append(args, out)

		fmt.Printf("args: %v\n", args)

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}

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
	}

	return out, nil
}
