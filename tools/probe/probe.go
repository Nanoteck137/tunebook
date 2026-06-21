package probe

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nanoteck137/tunebook/types"
)

type Tags map[string]any

type ffprobeStream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	Duration  string `json:"duration"`
	Tags      Tags   `json:"tags"`
}

type ffprobeFormat struct {
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Tags       Tags   `json:"tags"`
}

type ffprobeOutput struct {
	Streams []ffprobeStream `json:"streams"`
	Format  ffprobeFormat   `json:"format"`
}

func convertMapKeysToLowercase(m Tags) Tags {
	res := make(Tags)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type ProbeResult struct {
	Tags        Tags
	MediaFormat types.MediaFormat
	Duration    time.Duration
}

func ProbeMedia(ctx context.Context, filepath string) (*ProbeResult, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filepath,
	}

	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result ffprobeOutput
	err = json.Unmarshal(out, &result)
	if err != nil {
		return nil, err
	}

	var tags Tags
	hasGlobalTags := result.Format.FormatName != "ogg"

	var audioStream *ffprobeStream
	for i, s := range result.Streams {
		if s.CodecType == "audio" {
			audioStream = &result.Streams[i]
			break
		}
	}

	if audioStream == nil {
		return nil, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = result.Format.Tags
	} else {
		tags = audioStream.Tags
	}

	if tags == nil {
		tags = make(Tags)
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
