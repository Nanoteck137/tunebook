package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
)

const (
	ErrTypeInvalidAuth             pyrin.ErrorType = "INVALID_AUTH"
	ErrTypeArtistNotFound          pyrin.ErrorType = "ARTIST_NOT_FOUND"
	ErrTypeAlbumNotFound           pyrin.ErrorType = "ALBUM_NOT_FOUND"
	ErrTypeTrackNotFound           pyrin.ErrorType = "TRACK_NOT_FOUND"
	ErrTypeVirtualPlaylistNotFound pyrin.ErrorType = "VIRTUAL_PLAYLIST_NOT_FOUND"
	ErrTypeApiTokenNotFound        pyrin.ErrorType = "API_TOKEN_NOT_FOUND"
	ErrTypeQueueNotFound           pyrin.ErrorType = "QUEUE_NOT_FOUND"

	ErrTypeInvalidFilter      pyrin.ErrorType = "INVALID_FILTER"
	ErrTypeInvalidSort        pyrin.ErrorType = "INVALID_SORT"
	ErrTypeUserAlreadyExists  pyrin.ErrorType = "USER_ALREADY_EXISTS"
	ErrTypeUserNotFound       pyrin.ErrorType = "USER_NOT_FOUND"
	ErrTypeInvalidCredentials pyrin.ErrorType = "INVALID_CREDENTIALS"

	ErrTypePlaylistNotFound        pyrin.ErrorType = "PLAYLIST_NOT_FOUND"
	ErrTypePlaylistAlreadyHasTrack pyrin.ErrorType = "PLAYLIST_ALREADY_HAS_TRACK"

	ErrTypeFilterNotFound pyrin.ErrorType = "FILTER_NOT_FOUND"

	ErrTypeHistoryNotFound      pyrin.ErrorType = "HISTORY_NOT_FOUND"
	ErrTypeUnsupportedImageType pyrin.ErrorType = "UNSUPPORTED_IMAGE_TYPE"

	ErrTypeChallengeMismatch pyrin.ErrorType = "CHALLENGE_MISMATCH"
	ErrTypeProviderNotFound  pyrin.ErrorType = "PROVIDER_NOT_FOUND"
	ErrTypeRequestNotFound   pyrin.ErrorType = "REQUEST_NOT_FOUND"

	ErrTypeMediaInvalidFormat         pyrin.ErrorType = "MEDIA_INVALID_FORMAT"
	ErrTypeMediaInvalidQuality        pyrin.ErrorType = "MEDIA_INVALID_QUALITY"
	ErrTypeMediaInvalidPolicy         pyrin.ErrorType = "MEDIA_INVALID_POLICY"
	ErrTypeMediaBitrateNotSet         pyrin.ErrorType = "MEDIA_BITRATE_NOT_SET"

	ErrTypePlaylistItemNotFound       pyrin.ErrorType = "PLAYLIST_ITEM_NOT_FOUND"
	ErrTypePlaylistAnchorTrackNotFound pyrin.ErrorType = "PLAYLIST_ANCHOR_TRACK_NOT_FOUND"
	ErrTypeNotAuthorized              pyrin.ErrorType = "NOT_AUTHORIZED"
)

func InvalidAuth(message string) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidAuth,
		Message: "Invalid auth: " + message,
	}
}

func ArtistNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeArtistNotFound,
		Message: "Artist not found",
	}
}

func AlbumNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeAlbumNotFound,
		Message: "Album not found",
	}
}

func TrackNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeTrackNotFound,
		Message: "Track not found",
	}
}

func VirtualPlaylistNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeVirtualPlaylistNotFound,
		Message: "Virtual Playlist not found",
	}
}

func ApiTokenNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeApiTokenNotFound,
		Message: "Api Token not found",
	}
}

func QueueNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeQueueNotFound,
		Message: "Queue not found",
	}
}

func InvalidFilter(err error) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidFilter,
		Message: err.Error(),
	}
}

func InvalidSort(err error) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidSort,
		Message: err.Error(),
	}
}

func UserAlreadyExists() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeUserAlreadyExists,
		Message: "User already exists",
	}
}

func UserNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeUserNotFound,
		Message: "User not found",
	}
}

func InvalidCredentials() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusUnauthorized,
		Type:    ErrTypeInvalidCredentials,
		Message: "Invalid Credentials",
	}
}

func PlaylistNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypePlaylistNotFound,
		Message: "Playlist not found",
	}
}

func FilterNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeFilterNotFound,
		Message: "Filter not found",
	}
}

func HistoryNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeHistoryNotFound,
		Message: "History not found",
	}
}

func PlaylistAlreadyHasTrack() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypePlaylistAlreadyHasTrack,
		Message: "Playlist already has track",
	}
}

func UnsupportedImageType() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeUnsupportedImageType,
		Message: "Unsupported image type",
	}
}

func ChallengeMismatch() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeChallengeMismatch,
		Message: "Challenge mismatch",
	}
}

func ProviderNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeProviderNotFound,
		Message: "Provider not found",
	}
}

func RequestNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeRequestNotFound,
		Message: "Request not found",
	}
}

func MediaInvalidFormat() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeMediaInvalidFormat,
		Message: "Invalid media format",
	}
}

func MediaInvalidQuality() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeMediaInvalidQuality,
		Message: "Invalid media quality",
	}
}

func MediaInvalidPolicy() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeMediaInvalidPolicy,
		Message: "Invalid media policy",
	}
}

func MediaBitrateNotSet() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeMediaBitrateNotSet,
		Message: "Bitrate not set",
	}
}

func PlaylistItemNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypePlaylistItemNotFound,
		Message: "Playlist item not found",
	}
}

func PlaylistAnchorTrackNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypePlaylistAnchorTrackNotFound,
		Message: "Anchor track not found",
	}
}

func NotAuthorized() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusForbidden,
		Type:    ErrTypeNotAuthorized,
		Message: "Not authorized",
	}
}
