package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

var magickImageMapping = map[string]ImageType{
	"PNG": ImageTypePng,
	"JPEG": ImageTypeJpeg,
}

type ImageService struct {
	logger *slog.Logger

	db      *database.Database
	workDir types.WorkDir
}

func NewImageService(
	logger *slog.Logger,
	db *database.Database,
	workDir types.WorkDir,
) *ImageService {
	return &ImageService{
		logger:  logger,
		db:      db,
		workDir: workDir,
	}
}

func (s *ImageService) convertImage(input, outputDir, name string, size int) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := utils.CreateResizedImage(input, p, size, size)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return p, nil
}

func (s *ImageService) convertSquareImage(input, outputDir, name string) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := utils.CreateSquareImage(input, p)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return p, nil
}

func (s *ImageService) copyDefaultToTemp(filename string) (string, error) {
	ext := path.Ext(filename)

	dest, err := os.CreateTemp("", "default*"+ext)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	src, err := assets.DefaultImagesFS.Open(filename)
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return "", err
	}

	return dest.Name(), nil
}

// TODO(patrik): Rename to ImageFormat
// TODO(patrik): Move to types package
type ImageType string

const (
	ImageTypeEmpty   ImageType = ""
	ImageTypeUnknown ImageType = "unknown"
	ImageTypePng     ImageType = "png"
	ImageTypeJpeg    ImageType = "jpeg"
)

func (t ImageType) ToExt() (string, bool) {
	switch t {
	case ImageTypePng:
		return ".png", true
	case ImageTypeJpeg:
		return ".jpeg", true
	}

	return "", false
}

func (s *ImageService) GetImageTypeFromExt(ext string) (ImageType, bool) {
	switch ext {
	case ".png":
		return ImageTypePng, true
	case ".jpg", ".jpeg":
		return ImageTypeJpeg, true
	}

	return "", false
}

func (s *ImageService) GetAlbumImage(ctx context.Context, albumId, typ string, imageType ImageType) (string, error) {
	album, err := s.db.GetAlbumById(ctx, albumId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", errors.New("album not found")
		}

		return "", err
	}

	cacheDir := s.workDir.Cache()
	albumCache := cacheDir.Album(album.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Albums(),
		albumCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	ext, ok := imageType.ToExt()
	if !ok {
		// TODO(patrik): Better error
		return "", errors.New("unknown image type")
	}

	input := ""

	if album.CoverArt.Valid {
		input = album.CoverArt.String
	} else {
		input, err = s.copyDefaultToTemp("default_album.png")
		if err != nil {
			return "", err
		}
		defer os.Remove(input)
	}

	switch typ {
	case "original":
		return s.convertSquareImage(input, albumCache, "original_square"+ext)
	case "128":
		return s.convertImage(input, albumCache, "128"+ext, 128)
	case "256":
		return s.convertImage(input, albumCache, "256"+ext, 256)
	case "512":
		return s.convertImage(input, albumCache, "512"+ext, 512)
	}

	return "", errors.New("unknown type")
}

func (s *ImageService) GetArtistImage(ctx context.Context, artistId, typ string, imageType ImageType) (string, error) {
	artist, err := s.db.GetArtistById(ctx, artistId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", errors.New("artist not found")
		}

		return "", err
	}

	cacheDir := s.workDir.Cache()
	artistCache := cacheDir.Artist(artist.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Artists(),
		artistCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	ext, ok := imageType.ToExt()
	if !ok {
		// TODO(patrik): Better error
		return "", errors.New("unknown image type")
	}

	input := ""

	if artist.CoverArt.Valid {
		input = artist.CoverArt.String
	} else {
		input, err = s.copyDefaultToTemp("default_artist.png")
		if err != nil {
			return "", err
		}
		defer os.Remove(input)
	}

	switch typ {
	case "original":
		return s.convertSquareImage(input, artistCache, "original_square"+ext)
	case "128":
		return s.convertImage(input, artistCache, "128"+ext, 128)
	case "256":
		return s.convertImage(input, artistCache, "256"+ext, 256)
	case "512":
		return s.convertImage(input, artistCache, "512"+ext, 512)
	}

	return "", errors.New("unknown type")
}

func (s *ImageService) GetPlaylistImage(ctx context.Context, playlistId, typ string, imageType ImageType) (string, error) {
	playlist, err := s.db.GetPlaylistById(ctx, playlistId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", errors.New("playlist not found")
		}

		return "", err
	}

	cacheDir := s.workDir.Cache()
	playlistCache := cacheDir.Playlist(playlist.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Playlists(),
		playlistCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	ext, ok := imageType.ToExt()
	if !ok {
		// TODO(patrik): Better error
		return "", errors.New("unknown image type")
	}

	input := ""

	if playlist.CoverArt.Valid {
		playlistDir := s.workDir.Playlist(playlist.Id)

		input = path.Join(playlistDir, playlist.CoverArt.String)
	} else {
		input, err = s.copyDefaultToTemp("default_album.png")
		if err != nil {
			return "", err
		}
		defer os.Remove(input)
	}

	switch typ {
	case "original":
		return s.convertSquareImage(input, playlistCache, "original_square"+ext)
	case "128":
		return s.convertImage(input, playlistCache, "128"+ext, 128)
	case "256":
		return s.convertImage(input, playlistCache, "256"+ext, 256)
	case "512":
		return s.convertImage(input, playlistCache, "512"+ext, 512)
	}

	return "", errors.New("unknown type")
}

func (s *ImageService) GetUserImage(ctx context.Context, userId, typ string, imageType ImageType) (string, error) {
	user, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", errors.New("user not found")
		}

		return "", err
	}

	cacheDir := s.workDir.Cache()
	userCache := cacheDir.User(user.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Users(),
		userCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	ext, ok := imageType.ToExt()
	if !ok {
		// TODO(patrik): Better error
		return "", errors.New("unknown image type")
	}

	input := ""

	if user.Picture.Valid {
		dir := s.workDir.User(user.Id)
		input = path.Join(dir, user.Picture.String)
	} else {
		// TODO(patrik): Create a default user picture
		input, err = s.copyDefaultToTemp("default_album.png")
		if err != nil {
			return "", err
		}
		defer os.Remove(input)
	}

	switch typ {
	case "original":
		return s.convertSquareImage(input, userCache, "original_square"+ext)
	case "128":
		return s.convertImage(input, userCache, "128"+ext, 128)
	case "256":
		return s.convertImage(input, userCache, "256"+ext, 256)
	case "512":
		return s.convertImage(input, userCache, "512"+ext, 512)
	}

	return "", errors.New("unknown type")
}

func (s *ImageService) ValidateImage(p string) (ImageType, error) {

	// out, err := exec.Command("magick", "identify", "-ping", "-format", "%m", p).CombinedOutput()

	cmd := exec.Command("magick", "identify", "-ping", "-format", "%m", p)

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		details := strings.TrimSpace(errOut.String())

		s.logger.Error("failed to validate image", 
			slog.Any("err", err), 
			slog.String("output", details),
		)

		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			if !execErr.Success() {
				return ImageTypeUnknown, nil
			}
		}

		return ImageTypeUnknown, err
	}

	ty := strings.TrimSpace(out.String())

	res, exists := magickImageMapping[ty]
	if !exists {
		return ImageTypeUnknown, errors.New("no mapping: " + ty)
	}

	return res, nil
}
