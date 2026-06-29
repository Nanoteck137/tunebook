package types

type ImageFormat string

const (
	ImageFormatEmpty   ImageFormat = ""
	ImageFormatUnknown ImageFormat = "unknown"
	ImageFormatPng     ImageFormat = "png"
	ImageFormatJpeg    ImageFormat = "jpeg"
)

func (f ImageFormat) IsValid() bool {
	switch f {
	case ImageFormatPng:
		return true
	case ImageFormatJpeg:
		return true
	}

	return false
}

func (f ImageFormat) ToExt() (string, bool) {
	switch f {
	case ImageFormatPng:
		return ".png", true
	case ImageFormatJpeg:
		return ".jpeg", true
	}

	return "", false
}
