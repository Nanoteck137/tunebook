package service

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"testing"

	"github.com/nanoteck137/tunebook/types"
)

func createTestImage(t *testing.T, path string, format string) {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	switch format {
	case "png":
		err = png.Encode(f, img)
	case "jpeg":
		err = jpeg.Encode(f, img, nil)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetImageFormat(t *testing.T) {
	dir, err := os.MkdirTemp("", "tunebook-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	logger := slog.Default()
	fs := NewFilesystemService(logger, dir)
	s := NewImageService(logger, nil, fs)

	t.Run("png", func(t *testing.T) {
		p := dir + "/test.png"
		createTestImage(t, p, "png")

		got, err := s.getImageFormat(p)
		if err != nil {
			t.Fatal(err)
		}
		if got != types.ImageFormatPng {
			t.Errorf("got %q, want %q", got, types.ImageFormatPng)
		}
	})

	t.Run("jpeg", func(t *testing.T) {
		p := dir + "/test.jpg"
		createTestImage(t, p, "jpeg")

		got, err := s.getImageFormat(p)
		if err != nil {
			t.Fatal(err)
		}
		if got != types.ImageFormatJpeg {
			t.Errorf("got %q, want %q", got, types.ImageFormatJpeg)
		}
	})

	t.Run("not_an_image", func(t *testing.T) {
		p := dir + "/test.txt"
		if err := os.WriteFile(p, []byte("not an image"), 0644); err != nil {
			t.Fatal(err)
		}

		got, err := s.getImageFormat(p)
		if err != nil {
			t.Fatal(err)
		}
		if got != types.ImageFormatUnknown {
			t.Errorf("got %q, want %q", got, types.ImageFormatUnknown)
		}
	})
}
