package types

type MediaFormat string

const (
	MediaFormatUnknown  MediaFormat = "unknown"
	MediaFormatEmpty    MediaFormat = ""
	MediaFormatFlac     MediaFormat = "flac"
	MediaFormatPcmS16LE MediaFormat = "pcm_s16le"
	MediaFormatOpus     MediaFormat = "opus"
	MediaFormatVorbis   MediaFormat = "vorbis"
	MediaFormatMp3      MediaFormat = "mp3"
	MediaFormatAac      MediaFormat = "aac"
)

func (m MediaFormat) ToExt() (string, bool) {
	switch m {
	case MediaFormatFlac:
		return ".flac", true
	case MediaFormatPcmS16LE:
		return ".wav", true
	case MediaFormatOpus:
		return ".opus", true
	case MediaFormatVorbis:
		return ".ogg", true
	case MediaFormatMp3:
		return ".mp3", true
	case MediaFormatAac:
		return ".m4a", true
	}

	return "", false
}

func (m MediaFormat) IsValid() bool {
	switch m {
	case MediaFormatFlac:
		return true
	case MediaFormatPcmS16LE:
		return true
	case MediaFormatOpus:
		return true
	case MediaFormatVorbis:
		return true
	case MediaFormatMp3:
		return true
	case MediaFormatAac:
		return true
	}

	return false
}

func (m MediaFormat) IsLossless() bool {
	switch m {
	case MediaFormatFlac:
		return true
	case MediaFormatPcmS16LE:
		return true
	}

	return false
}

func (m MediaFormat) IsLossy() bool {
	switch m {
	case MediaFormatOpus:
		return true
	case MediaFormatVorbis:
		return true
	case MediaFormatMp3:
		return true
	case MediaFormatAac:
		return true
	}

	return false
}
