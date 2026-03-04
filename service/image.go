package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path"

	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type ImageService struct {
	db      *database.Database
	workDir types.WorkDir
}

func NewImageService(db *database.Database, workDir types.WorkDir) *ImageService {
	return &ImageService{
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

type ImageType string

const (
	ImageTypePng  ImageType = "png"
	ImageTypeJpeg ImageType = "jpeg"
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

	// TODO(patrik): Move error
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
