package library

import (
	"path/filepath"

	"github.com/nanoteck137/validate"
)

type ArtistEntry struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	CoverArt string   `json:"coverArt"`
	Tags     []string `json:"tags"`

	Path string `json:"path"`
}

func (e ArtistEntry) Validate() error {
	return validate.ValidateStruct(&e,
		validate.Field(&e.Id, validate.Required),
		validate.Field(&e.Name, validate.Required),
		validate.Field(&e.Path, validate.Required),
	)
}

func (e ArtistEntry) GetCoverArt() string {
	if e.CoverArt == "" || e.Path == "" {
		return ""
	}

	return filepath.Join(e.Path, e.CoverArt)
}

type AlbumEntry struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	CoverArt           string   `json:"coverArt"`
	Year               int64    `json:"year"`
	ArtistId           string   `json:"artistId"`
	FeaturingArtistIds []string `json:"featuringArtistIds"`
	Tags               []string `json:"tags"`

	Path string `json:"path"`
}

func (e AlbumEntry) Validate() error {
	return validate.ValidateStruct(&e,
		validate.Field(&e.Id, validate.Required),
		validate.Field(&e.Name, validate.Required),
		validate.Field(&e.ArtistId, validate.Required),
		validate.Field(&e.Path, validate.Required),
	)
}

func (e AlbumEntry) GetCoverArt() string {
	if e.CoverArt == "" || e.Path == "" {
		return ""
	}

	return filepath.Join(e.Path, e.CoverArt)
}

type TrackEntry struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	TrackFile          string   `json:"trackFile"`
	Number             int64    `json:"number"`
	Year               int64    `json:"year"`
	Tags               []string `json:"tags"`
	AlbumId            string   `json:"albumId"`
	ArtistId           string   `json:"artistId"`
	FeaturingArtistIds []string `json:"featuringArtistIds"`

	Path string `json:"path"`
}

func (e TrackEntry) Validate() error {
	return validate.ValidateStruct(&e,
		validate.Field(&e.Id, validate.Required),
		validate.Field(&e.Name, validate.Required),
		validate.Field(&e.TrackFile, validate.Required),
		validate.Field(&e.AlbumId, validate.Required),
		validate.Field(&e.ArtistId, validate.Required),
		validate.Field(&e.Path, validate.Required),
	)
}

func (e TrackEntry) GetTrackFile() string {
	if e.TrackFile == "" || e.Path == "" {
		return ""
	}

	return filepath.Join(e.Path, e.TrackFile)
}
