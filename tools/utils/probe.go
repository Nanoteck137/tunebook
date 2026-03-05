package utils

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/nanoteck137/dwebble/types"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func convertMapKeysToLowercase(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type TrackInfo struct {
	Path string

	Duration int
	Tags     map[string]string
}

type ProbeResult struct {
	Tags        ffprobe.Tags
	MediaFormat types.MediaFormat
	Duration    float64
}

func ProbeTrack(filepath string) (ProbeResult, error) {
	slog.Info("Probing track", "filepath", filepath)

	ctx := context.TODO()

	probe, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return ProbeResult{}, err
	}

	var tags ffprobe.Tags
	hasGlobalTags := probe.Format.FormatName != "ogg"

	audioStream := probe.FirstAudioStream()
	if audioStream == nil {
		return ProbeResult{}, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = probe.Format.TagList
	} else {
		tags = audioStream.TagList
	}

	tags = convertMapKeysToLowercase(tags)

	duration, err := strconv.ParseFloat(audioStream.Duration, 32)
	if err != nil {
		return ProbeResult{}, err
	}

	mediaType := types.MediaFormatUnknown

	// TODO(patrik): Move to helper
	// TODO(patrik): Add more types
	switch probe.Format.FormatName {
	case "flac":
		mediaType = types.MediaFormatFlac
	case "ogg":
		switch audioStream.CodecName {
		case "opus":
			mediaType = types.MediaFormatOpus
		case "vorbis":
			mediaType = types.MediaFormatVorbis
		}
	}

	return ProbeResult{
		Tags:        tags,
		MediaFormat: mediaType,
		Duration:    duration,
	}, nil
}
