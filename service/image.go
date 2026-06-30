package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nanoteck137/tunebook/assets"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/nrednav/cuid2"
)

var imageErr = NewServiceErrCreator("image")

var (
	ErrImageServiceUnsupportedImageFormat = imageErr.New("unsupported image format")
)

const (
	MaxImageSize = 2048
)

var magickImageMapping = map[string]types.ImageFormat{
	"PNG":  types.ImageFormatPng,
	"JPEG": types.ImageFormatJpeg,
}

var newTempFileId, _ = cuid2.Init(cuid2.WithLength(16))

type ImageService struct {
	logger *slog.Logger

	db      *database.Database
	dataDir types.DataDir
}

func NewImageService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
) *ImageService {
	return &ImageService{
		logger:  logger,
		db:      db,
		dataDir: dataDir,
	}
}

func (s *ImageService) ConvertImage(
	input, outputDir, name string,
	size int,
) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := createResizedImage(input, p, size, size)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return p, nil
}

func (s *ImageService) createTempFilename(ext string) string {
	return path.Join(s.dataDir.Temp(), newTempFileId()+ext)
}

func (s *ImageService) ConvertSquareImage(
	input, outputDir, name string,
) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := createSquareImage(input, p)
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

	dest := s.createTempFilename(ext)

	src, err := assets.DefaultImagesFS.Open(filename)
	if err != nil {
		return "", err
	}
	defer src.Close()

	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}

	return dest, nil
}

type ProcessImageParams struct {
	Input       string
	Default     string
	OutputDir   string
	Size        int
	ImageFormat types.ImageFormat
}

func (s *ImageService) ProcessImage(params ProcessImageParams) (string, error) {
	ext, ok := params.ImageFormat.ToExt()
	if !ok {
		return "", ErrImageServiceUnsupportedImageFormat
	}

	size := params.Size
	if size > MaxImageSize {
		size = MaxImageSize
	}

	input := params.Input
	if input == "" {
		var err error
		input, err = s.copyDefaultToTemp(params.Default)
		if err != nil {
			return "", imageErr.Wrap("process image: copy default", err)
		}
		defer os.Remove(input)
	}

	if size == 0 {
		p, err := s.ConvertSquareImage(
			input, params.OutputDir, "original_square"+ext)
		if err != nil {
			return "", imageErr.Wrap("process image: convert square original", err)
		}

		return p, nil
	}

	name := strconv.Itoa(size)
	p, err := s.ConvertImage(input, params.OutputDir, name+ext, size)
	if err != nil {
		return "", imageErr.Wrap("process image: convert sized", err)
	}

	return p, nil
}

func (s *ImageService) getImageFormat(p string) (types.ImageFormat, error) {
	cmd := exec.Command("magick", "identify", "-ping", "-format", "%m", p)

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		details := strings.TrimSpace(errOut.String())

		s.logger.Error("get image format",
			slog.Any("err", err),
			slog.String("output", details),
		)

		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			if !execErr.Success() {
				return types.ImageFormatUnknown, nil
			}
		}

		return types.ImageFormatUnknown, err
	}

	ty := strings.TrimSpace(out.String())

	res, exists := magickImageMapping[ty]
	if !exists {
		s.logger.Warn("get image format: no mapping found", "mapping", ty)

		return types.ImageFormatUnknown, nil
	}

	return res, nil
}

func (s *ImageService) ValidateImage(p string) (types.ImageFormat, error) {
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
				return types.ImageFormatUnknown, nil
			}
		}

		return types.ImageFormatUnknown, err
	}

	ty := strings.TrimSpace(out.String())

	res, exists := magickImageMapping[ty]
	if !exists {
		return types.ImageFormatUnknown, imageErr.Newf("no mapping: %s", ty)
	}

	return res, nil
}

type DownloadCoverForPlaylistParams struct {
	PlaylistId string
	Url        string
}

func (s *ImageService) DownloadCoverForPlaylist(
	ctx context.Context,
	params DownloadCoverForPlaylistParams,
) (string, error) {
	tmpPath, err := s.downloadToTempFile(params.Url)
	if err != nil {
		return "", imageErr.Wrap("download cover for playlist: download", err)
	}
	defer os.Remove(tmpPath)

	format, err := s.getImageFormat(tmpPath)
	if err != nil {
		return "", imageErr.Wrap(
			"download cover for playlist: image format", err)
	}

	playlistDir := s.dataDir.Playlist(params.PlaylistId)

	err = utils.CreateDirectories([]string{
		playlistDir,
	})
	if err != nil {
		return "", imageErr.Wrap(
			"download cover for playlist: mkdir", err)
	}

	filename, err := s.finalizeImage(tmpPath, format, playlistDir)
	if err != nil {
		return "", imageErr.Wrap(
			"download cover for playlist: finalize", err)
	}

	return filename, nil
}

type GenerateImageForPlaylistParams struct {
	PlaylistId string
}

func (s *ImageService) GenerateImageForPlaylist(
	ctx context.Context,
	params GenerateImageForPlaylistParams,
) (string, error) {
	images, err := s.db.GetPlaylistTrackImages(ctx, params.PlaylistId, 4)
	if err != nil {
		return "", imageErr.Wrap(
			"generate image for playlist: images", err)
	}

	imgs := [4]string{}

	for i, img := range images {
		if !img.Valid {
			continue
		}

		imgs[i] = img.String
	}

	playlistDir := s.dataDir.Playlist(params.PlaylistId)

	err = utils.CreateDirectories([]string{
		playlistDir,
	})
	if err != nil {
		return "", imageErr.Wrap("generate image for playlist: mkdir", err)
	}

	tmpPath := s.createTempFilename(".png")

	err = generatePlaylistCover(imgs, tmpPath, 512)
	if err != nil {
		return "", imageErr.Wrap("generate image for playlist: generate", err)
	}

	filename, err := s.finalizeImage(
		tmpPath, types.ImageFormatPng, playlistDir)
	if err != nil {
		return "", imageErr.Wrap("generate image for playlist: finalize", err)
	}

	return filename, nil
}

type UploadImageForPlaylistParams struct {
	PlaylistId string

	File *multipart.FileHeader
}

func (s *ImageService) UploadImageForPlaylist(
	ctx context.Context,
	params UploadImageForPlaylistParams,
) (string, error) {
	tmpPath, err := s.copyMultipartFileToTempFile(params.File)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for playlist: file to temp", err)
	}
	defer os.Remove(tmpPath)

	playlistDir := s.dataDir.Playlist(params.PlaylistId)

	err = utils.CreateDirectories([]string{
		playlistDir,
	})
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for playlist: mkdir", err)
	}

	format, err := s.getImageFormat(tmpPath)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for playlist: image format", err)
	}

	filename, err := s.finalizeImage(tmpPath, format, playlistDir)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for playlist: finalize", err)
	}

	return filename, nil
}

type UploadImageForUserParams struct {
	UserId string

	File *multipart.FileHeader
}

func (s *ImageService) UploadImageForUser(
	ctx context.Context,
	params UploadImageForUserParams,
) (string, error) {
	tmpPath, err := s.copyMultipartFileToTempFile(params.File)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for user: file to temp", err)
	}
	defer os.Remove(tmpPath)

	userDir := s.dataDir.User(params.UserId)

	err = utils.CreateDirectories([]string{
		userDir,
	})
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for user: mkdir", err)
	}

	format, err := s.getImageFormat(tmpPath)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for user: image format", err)
	}

	filename, err := s.finalizeImage(tmpPath, format, userDir)
	if err != nil {
		return "", imageErr.Wrap(
			"upload image for user: finalize", err)
	}

	return filename, nil
}

type DownloadPictureForUserParams struct {
	UserId string
	Url    string
}

func (s *ImageService) DownloadPictureForUser(
	ctx context.Context,
	params DownloadPictureForUserParams,
) (string, error) {
	tmpPath, err := s.downloadToTempFile(params.Url)
	if err != nil {
		return "", imageErr.Wrap(
			"download picture for user: download", err)
	}
	defer os.Remove(tmpPath)

	userDir := s.dataDir.User(params.UserId)

	err = utils.CreateDirectories([]string{
		userDir,
	})
	if err != nil {
		return "", imageErr.Wrap(
			"download picture for user: mkdir", err)
	}

	format, err := s.getImageFormat(tmpPath)
	if err != nil {
		return "", imageErr.Wrap(
			"download picture for user: image format", err)
	}

	picture, err := s.finalizeImage(tmpPath, format, userDir)
	if err != nil {
		return "", imageErr.Wrap(
			"download picture for user: finalize", err)
	}

	return picture, nil
}

func (s *ImageService) copyReaderToTempFile(r io.Reader) (string, error) {
	tmpPath := s.createTempFilename("")
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	_, err = io.Copy(tmp, r)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

func (s *ImageService) copyMultipartFileToTempFile(
	file *multipart.FileHeader,
) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	tmpPath, err := s.copyReaderToTempFile(f)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

func (s *ImageService) downloadToTempFile(
	url string,
) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// TODO(patrik): Check for success

	tmpPath, err := s.copyReaderToTempFile(resp.Body)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

func (s *ImageService) finalizeImage(
	tmpPath string,
	imageFormat types.ImageFormat,
	dir string,
) (string, error) {
	imageExt, ok := imageFormat.ToExt()
	if !ok {
		return "", ErrImageServiceUnsupportedImageFormat
	}

	hash, err := hashFile(tmpPath)
	if err != nil {
		return "", err
	}

	name := hash + imageExt
	output := filepath.Join(dir, name)
	err = os.Rename(tmpPath, output)
	if err != nil {
		return "", imageErr.Wrap("promote temp file", err)
	}

	return name, nil
}

func createSquareImage(src, dest string) error {
	cmd := exec.Command(
		"magick", src,
		"-gravity", "Center",
		"-extent", "%[fx:min(w,h)]x%[fx:min(w,h)]",
		dest,
	)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createResizedImage(src string, dest string, width, height int) error {
	args := []string{
		src,
		"-resize", fmt.Sprintf("%dx%d^", width, height),
		"-gravity", "Center",
		"-extent", fmt.Sprintf("%dx%d", width, height),
		dest,
	}

	cmd := exec.Command("magick", args...)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func generatePlaylistCover(images [4]string, output string, tileSize int) error {
	if len(images) == 0 {
		return fmt.Errorf("at least one image is required")
	}

	size := fmt.Sprintf("%dx%d", tileSize, tileSize)

	buildTile := func(img string) []string {
		if img == "" {
			return []string{"(", "xc:black", "-resize", size, ")"}
		}

		return []string{
			"(",
			img,
			"-resize", size + "^",
			"-gravity", "center",
			"-extent", size,
			")",
		}
	}

	args := []string{}
	for _, img := range images {
		args = append(args, buildTile(img)...)
	}

	args = append(args,
		"(", "-clone", "0-1", "+append", ")",
		"(", "-clone", "2-3", "+append", ")",
		"-delete", "0-3",
		"-append",
		output,
	)

	cmd := exec.Command("magick", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func convertImage(src string, dest string) error {
	args := []string{
		"convert",
		src,
		dest,
	}

	cmd := exec.Command("magick", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func hashFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func getImageExtFromContentType(contentType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", imageErr.Wrap("failed to parse content type", err)
	}

	switch mediaType {
	case "image/png":
		return ".png", nil
	case "image/jpeg":
		return ".jpeg", nil
	default:
		return "", imageErr.Newf("unsupported media type: %s", mediaType)
	}
}
