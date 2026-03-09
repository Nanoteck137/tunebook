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

type Mode string

const (
	ModeEmpty    Mode = ""
	ModeRaw      Mode = "raw"
	ModeSmart    Mode = "smart"
	ModeOriginal Mode = "original"
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

type MediaStreamOptions struct {
	Mode Mode

	Device  Device
	Policy  Policy
	Quality Quality

	Format  types.MediaFormat
	Bitrate int
}

func (s *MediaService) GetTrackStream(trackId string, opts MediaStreamOptions) (string, error) {
	// if !opts.Original && (!opts.Format.IsValid() || opts.Bitrate == 0) {
	// 	return "", errors.New("invalid stream options")
	// }

	// device
	// quality
	// policy
	// format
	// bitrate

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
	bitrate := 128

	if opts.Mode == ModeEmpty {
		opts.Mode = ModeOriginal
	}

	switch opts.Policy {
	case PolicyEmpty, PolicyOriginal:
		useOriginal = true
	case PolicyTranscode:
		if !opts.Format.IsValid() {
			return "", errors.New("invalid format")
		}

		format = opts.Format

		// if opts.Bitrate <= 0 {
		// 	// TODO(patrik): Get the preferred bitrate from the format
		// 	return "", errors.New("TODO: get bitrate from format ")
		// }

		qualityMapping := map[types.MediaFormat]map[Quality]int{
			types.MediaFormatFlac:   {QualityHigh: 0, QualityMedium: 0, QualityLow: 0},
			types.MediaFormatWav:    {QualityHigh: 0, QualityMedium: 0, QualityLow: 0},
			types.MediaFormatOpus:   {QualityHigh: 128, QualityMedium: 96, QualityLow: 64},
			types.MediaFormatVorbis: {QualityHigh: 192, QualityMedium: 128, QualityLow: 96},
			types.MediaFormatMp3:    {QualityHigh: 320, QualityMedium: 192, QualityLow: 128},
			types.MediaFormatAac:    {QualityHigh: 256, QualityMedium: 192, QualityLow: 96},
		}

		if !opts.Format.IsLossless() {
			if opts.Bitrate > 0 {
				bitrate = opts.Bitrate
			} else if opts.Quality != QualityEmpty {
				// TODO(patrik): Check for valid quality

				qualities, ok := qualityMapping[format]
				if !ok {
					panic("Format missing qualities")
				}

				bitrate, ok = qualities[opts.Quality]
				if !ok {
					panic("Missing bitrate for quality: " + opts.Quality)
				}

				// TODO(patrik): Implement
				// panic("bitrate = getFromQuality(opts.Quality)")

			} else {
				qualities, ok := qualityMapping[format]
				if !ok {
					panic("Format missing qualities")
				}

				bitrate, ok = qualities[QualityMedium]
				if !ok {
					panic("Missing bitrate for medium quality for format: " + format)
				}

				// TODO
				// panic("bitrate = getDefaultBitrate(format)")
			}
		}
	case PolicyLossy:
		if track.MediaFormat.IsLossy() {
			format = track.MediaFormat
		} else {
			if opts.Format != "" && opts.Format.IsValid() {
				format = opts.Format
				if opts.Bitrate > 0 {
					bitrate = opts.Bitrate
				} else if opts.Quality != QualityEmpty {
					// TODO(patrik): Implement
					panic("bitrate = getFromQuality(opts.Quality)")
				} else {
					// TODO
					panic("bitrate = getDefaultBitrate(format)")
				}
			}

			if opts.Device != DeviceEmpty {
				// TODO
				panic("format = getBestLossyFormatForDevice(opts.Device, opts.Quality)")
			}

			// format = getFromDeviceAndQuality()
		}
	default:
		// TODO(patrik): Fix
		panic("Unknown policy: " + opts.Policy)
	}

	fmt.Printf("track.MediaFormat: %v\n", track.MediaFormat)
	fmt.Printf("format: %v\n", format)
	fmt.Printf("bitrate: %v\n", bitrate)
	fmt.Printf("useOriginal: %v\n", useOriginal)

	return "", errors.New("Testing")

	// TODO(patrik): Maybe here we should still "transcode" the media but
	// just remove the metadata
	if format == track.MediaFormat {
		useOriginal = true
	}

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
		return "", errors.New("use original")
		return track.Filename, nil
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

	if true {
		return "", errors.New("this is just a test for fun")
	}

	transcode := func() error {
		// TODO(patrik): Don't map the metadata
		args := []string{
			"-i", track.Filename,
		}

		switch opts.Format {
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
			return errors.New("unsupported media format: " + string(bitrate))
		}

		args = append(args, out)

		fmt.Printf("args: %v\n", args)

		cmd := exec.Command("ffmpeg", args...)
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
