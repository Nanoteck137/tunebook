package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/nanoteck137/dwebble/database"
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

// Device  string
// Policy  string
// Quality string
//
// // Raw options
// Format types.MediaFormat
// Bitrate int

type MediaStreamOptions struct {
	Format  types.MediaFormat
	Bitrate int
}

func (s *MediaService) GetTrackStream(trackId string, opts MediaStreamOptions) (string, error) {
	if !opts.Format.IsValid() || opts.Bitrate == 0 {
		return "", errors.New("invalid stream options")
	}

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

	// mediaType := types.GetMediaTypeFromExt(fileExt)
	//
	// if mediaType == types.MediaTypeUnknown {
	// 	return pyrin.NoContentNotFound()
	// }
	//
	// // Return the original file if the filename matches the
	// // one stored inside the track
	// if track.MediaType == mediaType {
	// 	d := path.Dir(track.Filename)
	// 	filename := path.Base(track.Filename)
	//
	// 	f := os.DirFS(d)
	// 	return pyrin.ServeFile(c, f, filename)
	// }

	// Here we need to start transcoding the original track
	// media to the requested format

	cacheDir := s.workDir.Cache()
	trackCache := cacheDir.Track(track.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Tracks(),
		trackCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	// filename: format-bitrate.ext
	ext, ok := opts.Format.ToExt()
	if !ok {
		return "", errors.New("unsupported media format")
	}

	filename := fmt.Sprintf("transcode-%s-%dk", opts.Format, opts.Bitrate) + ext
	fmt.Printf("filename: %v\n", filename)

	out := path.Join(trackCache, filename)

	transcode := func() error {
		args := []string{
			"-i", track.Filename,
		}

		switch opts.Format {
		case types.MediaFormatFlac:
			args = append(args, "-codec:a", "flac", "-compression_level", "5")
		case types.MediaFormatWav:
			args = append(args, "-codec:a", "pcm_s16le")
		case types.MediaFormatOpus:
			args = append(args, "-codec:a", "libopus", "-b:a", fmt.Sprintf("%dk", opts.Bitrate), "-vbr", "on", "-compression_level", "10")
		case types.MediaFormatVorbis:
			args = append(args, "-codec:a", "libvorbis", "-b:a", fmt.Sprintf("%dk", opts.Bitrate))
		case types.MediaFormatMp3:
			args = append(args, "-codec:a", "libmp3lame", "-b:a", fmt.Sprintf("%dk", opts.Bitrate), "-q:a", "0")
		case types.MediaFormatAac:
			args = append(args, "-codec:a", "aac", "-b:a", fmt.Sprintf("%dk", opts.Bitrate), "-aac_coder", "twoloop", "-movflags", "+faststart")
		default:
			return errors.New("unsupported media format: " + string(opts.Format))
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
