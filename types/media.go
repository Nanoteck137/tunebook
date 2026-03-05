package types

type MediaFormat string

const (
	MediaFormatUnknown MediaFormat = "unknown"
	MediaFormatFlac    MediaFormat = "flac"
	MediaFormatWav     MediaFormat = "wav"
	MediaFormatOpus    MediaFormat = "opus"
	MediaFormatVorbis  MediaFormat = "vorbis"
	MediaFormatMp3     MediaFormat = "mp3"
	MediaFormatAac     MediaFormat = "aac"
)

func (m MediaFormat) ToExt() (string, bool) {
	switch m {
	case MediaFormatFlac:
		return ".flac", true
	case MediaFormatWav:
		return ".wav", true
	case MediaFormatOpus:
		return ".opus", true
	case MediaFormatVorbis:
		return ".ogg", true
	case MediaFormatMp3:
		return ".mp3", true
	case MediaFormatAac:
		return ".aac", true
	}

	return "", false
}

func (m MediaFormat) IsValid() bool {
	switch m {
	case MediaFormatFlac:
		return true
	case MediaFormatWav:
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

func GetMediaFormatFromExt(ext string) MediaFormat {
	switch ext {
	case ".flac":
		return MediaFormatFlac
	case ".wav":
		return MediaFormatWav
	case ".opus":
		return MediaFormatOpus
	case ".ogg":
		return MediaFormatVorbis
	case ".mp3":
		return MediaFormatMp3
	case ".aac":
		return MediaFormatAac
	}

	return MediaFormatUnknown
}
