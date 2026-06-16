package probe

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/nanoteck137/tunebook/types"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func convertMapKeysToLowercase(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type ProbeResult struct {
	Tags        ffprobe.Tags
	MediaFormat types.MediaFormat
	Duration    time.Duration
}

func ProbeMedia(ctx context.Context, filepath string) (*ProbeResult, error) {
	probe, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return nil, err
	}

	var tags ffprobe.Tags
	hasGlobalTags := probe.Format.FormatName != "ogg"

	audioStream := probe.FirstAudioStream()
	if audioStream == nil {
		// TODO(patrik): Better error?
		return nil, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = probe.Format.TagList
	} else {
		tags = audioStream.TagList
	}

	tags = convertMapKeysToLowercase(tags)

	dur, err := strconv.ParseFloat(audioStream.Duration, 32)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(dur * float64(time.Second))

	mediaFormat := types.MediaFormatUnknown
	switch audioStream.CodecName {
	case "flac":
		mediaFormat = types.MediaFormatFlac
	case "pcm_s16le":
		mediaFormat = types.MediaFormatPcmS16LE
	case "opus":
		mediaFormat = types.MediaFormatOpus
	case "vorbis":
		mediaFormat = types.MediaFormatVorbis
	case "mp3":
		mediaFormat = types.MediaFormatMp3
	case "aac":
		mediaFormat = types.MediaFormatAac
	}

	return &ProbeResult{
		Tags:        tags,
		MediaFormat: mediaFormat,
		Duration:    duration,
	}, nil
}
