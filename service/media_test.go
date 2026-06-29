package service_test

import (
	"log/slog"
	"testing"

	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
)

func TestQualityMappingHasAllFormats(t *testing.T) {
	cfg := &config.Config{
		Media: config.MediaConfig{
			Opus:   config.MediaQualityConfig{High: 128, Medium: 96, Low: 64},
			Vorbis: config.MediaQualityConfig{High: 192, Medium: 128, Low: 96},
			Mp3:    config.MediaQualityConfig{High: 320, Medium: 192, Low: 128},
			Aac:    config.MediaQualityConfig{High: 256, Medium: 192, Low: 96},
		},
	}

	s := service.NewMediaService(slog.Default(), nil, types.DataDir("/tmp"), cfg)

	for _, format := range types.ValidMediaFormats {
		spec, ok := s.QualityMapping[format]
		if !ok {
			t.Errorf("Missing quality spec for format: %s", format)
			continue
		}

		if format.IsLossy() {
			if spec.High == 0 {
				t.Errorf("Lossy format %s has zero high bitrate", format)
			}
			if spec.Medium == 0 {
				t.Errorf("Lossy format %s has zero medium bitrate", format)
			}
			if spec.Low == 0 {
				t.Errorf("Lossy format %s has zero low bitrate", format)
			}
		}
	}
}
