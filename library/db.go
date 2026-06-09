package library

import "path/filepath"

type ArtistEntry struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	CoverArt string   `json:"coverArt"`
	Tags     []string `json:"tags"`

	Path string `json:"path"`
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

func (e TrackEntry) GetTrackFile() string {
	if e.TrackFile == "" || e.Path == "" {
		return ""
	}

	return filepath.Join(e.Path, e.TrackFile)
}
