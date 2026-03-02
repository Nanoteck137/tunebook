package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"sync/atomic"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/pelletier/go-toml/v2"
)

type ArtistEntry struct {
	Id string `json:"id"`

	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	CoverArt string   `json:"coverArt"`
	Tags     []string `json:"tags"`

	Path string `json:"path"`
}

func (e ArtistEntry) GetCoverArt() string {
	if e.CoverArt == "" || e.Path == "" {
		return ""
	}

	return path.Join(e.Path, e.CoverArt)
}

type AlbumEntry struct {
	Id string `json:"id"`

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

	return path.Join(e.Path, e.CoverArt)
}

type TrackEntry struct {
	Id string `json:"id"`

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

	return path.Join(e.Path, e.TrackFile)
}

type MetadataGeneral struct {
	Cover     string   `json:"cover" toml:"cover"`
	Tags      []string `json:"tags" toml:"tags"`
	TrackTags []string `json:"trackTags" toml:"trackTags"`
	Year      int64    `json:"year" toml:"year"`
}

type MetadataAlbum struct {
	Id      string   `json:"id" toml:"id"`
	Name    string   `json:"name" toml:"name"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type MetadataTrack struct {
	Id      string   `json:"id" toml:"id"`
	File    string   `json:"file" toml:"file"`
	Name    string   `json:"name" toml:"name"`
	Number  int64    `json:"number" toml:"number"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

// TODO(patrik): Rename
type Metadata struct {
	General MetadataGeneral `json:"general" toml:"general"`
	Album   MetadataAlbum   `json:"album" toml:"album"`
	Tracks  []MetadataTrack `json:"tracks" toml:"tracks"`

	Path string `json:"-" toml:"-"`
}

type Album struct {
	Path     string
	Metadata Metadata
}

type ArtistMetadata struct {
	Id string `json:"id" toml:"id"`

	Slug  string   `json:"slug" toml:"slug"`
	Name  string   `json:"name" toml:"name"`
	Cover string   `json:"cover" toml:"cover"`
	Tags  []string `json:"tags" toml:"tags"`

	Path string `json:"-" toml:"-"`
}

func (a ArtistMetadata) CoverPath() string {
	if a.Cover == "" {
		return ""
	}

	return path.Join(a.Path, a.Cover)
}

type LibraryService struct {
	db     *database.Database
	config *config.Config

	searchService *SearchService

	syncRunning atomic.Bool
}

func NewLibraryService(db *database.Database, config *config.Config, searchService *SearchService) *LibraryService {
	return &LibraryService{
		db:            db,
		config:        config,
		searchService: searchService,
		syncRunning:   atomic.Bool{},
	}
}

func readAlbum(p string) (Album, error) {
	metadataPath := path.Join(p, "album.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return Album{}, err
	}

	var metadata Metadata
	err = toml.Unmarshal(data, &metadata)
	if err != nil {
		return Album{}, err
	}

	if metadata.General.Cover != "" {
		metadata.General.Cover = path.Join(p, metadata.General.Cover)
	}

	for i, t := range metadata.Tracks {
		metadata.Tracks[i].File = path.Join(p, t.File)
	}

	return Album{
		Path:     p,
		Metadata: metadata,
	}, nil
}

func readArtistMetadata(p string) (ArtistMetadata, error) {
	metadataPath := path.Join(p, "artist.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return ArtistMetadata{}, err
	}

	var metadata ArtistMetadata
	err = toml.Unmarshal(data, &metadata)
	if err != nil {
		return ArtistMetadata{}, err
	}

	metadata.Path = p

	return metadata, nil
}

func setArtistTags(ctx context.Context, db database.DB, artistId string, tags []string) error {
	err := db.RemoveAllTagsFromArtist(ctx, artistId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToArtist(ctx, slug, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncArtists(ctx context.Context, libraryDir string) error {
	p := path.Join(libraryDir, ".library", "artists")
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var entry ArtistEntry
		err := decoder.Decode(&entry)
		if err != nil {
			return err
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		coverArt := entry.GetCoverArt()

		dbArtist, err := s.db.GetArtistById(ctx, entry.Id)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				_, err = s.db.CreateArtist(ctx, database.CreateArtistParams{
					Id:   entry.Id,
					Slug: entry.Slug,
					Name: entry.Name,
					Picture: sql.NullString{
						String: coverArt,
						Valid:  coverArt != "",
					},
				})
				if err != nil {
					return err
				}
			}
		} else {
			changes := database.ArtistChanges{}

			changes.Name = types.Change[string]{
				Value:   entry.Name,
				Changed: entry.Name != dbArtist.Name,
			}

			changes.Slug = types.Change[string]{
				Value:   entry.Slug,
				Changed: entry.Slug != dbArtist.Slug,
			}

			changes.Picture = types.Change[sql.NullString]{
				Value: sql.NullString{
					String: coverArt,
					Valid:  coverArt != "",
				},
				Changed: coverArt != dbArtist.Picture.String,
			}

			err := s.db.UpdateArtist(ctx, dbArtist.Id, changes)
			if err != nil {
				return err
			}
		}

		err = setArtistTags(
			ctx,
			s.db.DB,
			entry.Id,
			entry.Tags,
		)
		if err != nil {
			return err
		}

		err = s.searchService.UpdateArtist(ctx, entry.Id)
		if err != nil {
			// TODO(patrik): Better error
			return err
		}
	}

	return nil
}

func setAlbumFeaturingArtists(ctx context.Context, db database.DB, albumId string, artistIds []string) error {
	err := db.RemoveAllAlbumFeaturingArtists(ctx, albumId)
	if err != nil {
		return err
	}

	for _, artistId := range artistIds {
		err = db.AddFeaturingArtistToAlbum(ctx, albumId, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func setAlbumTags(ctx context.Context, db database.DB, albumId string, tags []string) error {
	err := db.RemoveAllTagsFromAlbum(ctx, albumId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToAlbum(ctx, slug, albumId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncAlbums(ctx context.Context, libraryDir string) error {
	p := path.Join(libraryDir, ".library", "albums")
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var entry AlbumEntry
		err := decoder.Decode(&entry)
		if err != nil {
			return err
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		coverArt := entry.GetCoverArt()

		dbAlbum, err := s.db.GetAlbumById(ctx, entry.Id)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				_, err = s.db.CreateAlbum(ctx, database.CreateAlbumParams{
					Id:       entry.Id,
					Name:     entry.Name,
					ArtistId: entry.ArtistId,
					CoverArt: sql.NullString{
						String: coverArt,
						Valid:  coverArt != "",
					},
					Year: sql.NullInt64{
						Int64: entry.Year,
						Valid: entry.Year != 0,
					},
				})
				if err != nil {
					return err
				}
			}
		} else {
			changes := database.AlbumChanges{}

			changes.Name = types.Change[string]{
				Value:   entry.Name,
				Changed: entry.Name != dbAlbum.Name,
			}

			changes.ArtistId = types.Change[string]{
				Value:   entry.ArtistId,
				Changed: entry.ArtistId != dbAlbum.ArtistId,
			}

			changes.CoverArt = types.Change[sql.NullString]{
				Value: sql.NullString{
					String: coverArt,
					Valid:  coverArt != "",
				},
				Changed: coverArt != dbAlbum.CoverArt.String,
			}

			changes.Year = types.Change[sql.NullInt64]{
				Value: sql.NullInt64{
					Int64: entry.Year,
					Valid: entry.Year != 0,
				},
				Changed: entry.Year != dbAlbum.Year.Int64,
			}

			err = s.db.UpdateAlbum(ctx, dbAlbum.Id, changes)
			if err != nil {
				return fmt.Errorf("failed to update album: %w", err)
			}
		}

		err = setAlbumFeaturingArtists(
			ctx,
			s.db.DB,
			entry.Id,
			entry.FeaturingArtistIds,
		)
		if err != nil {
			return err
		}

		err = setAlbumTags(
			ctx,
			s.db.DB,
			entry.Id,
			entry.Tags,
		)
		if err != nil {
			return err
		}

		err = s.searchService.UpdateAlbum(ctx, entry.Id)
		if err != nil {
			// TODO(patrik): Better error
			return err
		}
	}

	return nil
}

func setTrackFeaturingArtists(ctx context.Context, db database.DB, trackId string, artistIds []string) error {
	err := db.RemoveAllTrackFeaturingArtists(ctx, trackId)
	if err != nil {
		return err
	}

	for _, artistId := range artistIds {
		err = db.AddFeaturingArtistToTrack(ctx, trackId, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func setTrackTags(ctx context.Context, db database.DB, trackId string, tags []string) error {
	err := db.RemoveAllTagsFromTrack(ctx, trackId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToTrack(ctx, slug, trackId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncTracks(ctx context.Context, libraryDir string) error {
	p := path.Join(libraryDir, ".library", "tracks")
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var entry TrackEntry
		err := decoder.Decode(&entry)
		if err != nil {
			return err
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		trackFile := entry.GetTrackFile()

		stat, err := os.Stat(trackFile)
		if err != nil {
			// TODO(patrik): Better error
			return fmt.Errorf("failed to stat track file (%s): %w", trackFile, err)
		}

		modifiedTime := stat.ModTime().UnixMilli()

		dbTrack, err := s.db.GetTrackById(ctx, entry.Id)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				// Id:           track.Id,
				// Filename:     track.File,
				// ModifiedTime: modifiedTime,
				// MediaType:    probeResult.MediaType,
				// Name:         track.Name,
				// OtherName:    sql.NullString{},
				// AlbumId:      dbAlbum.Id,
				// ArtistId:     artist,
				// Duration:     int64(probeResult.Duration),
				// Number: sql.NullInt64{
				// 	Int64: track.Number,
				// 	Valid: track.Number != 0,
				// },
				// Year: sql.NullInt64{
				// 	Int64: track.Year,
				// 	Valid: track.Year != 0,
				// },

				probeResult, err := utils.ProbeTrack(trackFile)
				if err != nil {
					// TODO(patrik): Better error
					return fmt.Errorf("failed to probe track file (%s): %w", trackFile, err)
				}

				_, err = s.db.CreateTrack(ctx, database.CreateTrackParams{
					Id:           entry.Id,
					Filename:     trackFile,
					ModifiedTime: modifiedTime,
					MediaType:    probeResult.MediaType,
					Name:         entry.Name,
					AlbumId:      entry.AlbumId,
					ArtistId:     entry.ArtistId,
					Duration:     int64(probeResult.Duration),
					Number: sql.NullInt64{
						Int64: entry.Number,
						Valid: entry.Number != 0,
					},
					Year: sql.NullInt64{
						Int64: entry.Year,
						Valid: entry.Year != 0,
					},
				})
				if err != nil {
					return err
				}
			}
		} else {
			changes := database.TrackChanges{}

			if modifiedTime > dbTrack.ModifiedTime || dbTrack.Filename != trackFile {
				probeResult, err := utils.ProbeTrack(trackFile)
				if err != nil {
					// TODO(patrik): Better error
					return fmt.Errorf("failed to probe track file (%s): %w", trackFile, err)
				}

				dur := int64(probeResult.Duration)
				changes.Duration = types.Change[int64]{
					Value:   dur,
					Changed: dur != dbTrack.Duration,
				}

				changes.MediaType = types.Change[types.MediaType]{
					Value:   probeResult.MediaType,
					Changed: probeResult.MediaType != dbTrack.MediaType,
				}

				changes.ModifiedTime = types.Change[int64]{
					Value:   modifiedTime,
					Changed: modifiedTime != dbTrack.ModifiedTime,
				}
			}

			changes.Filename = types.Change[string]{
				Value:   trackFile,
				Changed: trackFile != dbTrack.Filename,
			}

			changes.Name = types.Change[string]{
				Value:   entry.Name,
				Changed: entry.Name != dbTrack.Name,
			}

			changes.AlbumId = types.Change[string]{
				Value:   entry.AlbumId,
				Changed: entry.AlbumId != dbTrack.AlbumId,
			}

			changes.ArtistId = types.Change[string]{
				Value:   entry.ArtistId,
				Changed: entry.ArtistId != dbTrack.ArtistId,
			}

			changes.Number = types.Change[sql.NullInt64]{
				Value: sql.NullInt64{
					Int64: entry.Number,
					Valid: entry.Number != 0,
				},
				Changed: entry.Number != dbTrack.Number.Int64,
			}

			changes.Year = types.Change[sql.NullInt64]{
				Value: sql.NullInt64{
					Int64: entry.Year,
					Valid: entry.Year != 0,
				},
				Changed: entry.Year != dbTrack.Year.Int64,
			}

			err = s.db.UpdateTrack(ctx, dbTrack.Id, changes)
			if err != nil {
				return err
			}
		}

		err = setTrackFeaturingArtists(
			ctx,
			s.db.DB,
			entry.Id,
			entry.FeaturingArtistIds,
		)
		if err != nil {
			return err
		}

		err = setTrackTags(
			ctx,
			s.db.DB,
			entry.Id,
			entry.Tags,
		)
		if err != nil {
			return err
		}

		err = s.searchService.UpdateTrack(ctx, entry.Id)
		if err != nil {
			// TODO(patrik): Better error
			return err
		}
	}

	return nil
}

func (s *LibraryService) runSync() error {
	p := s.config.LibraryDir

	// Library:
	//  Root Library Dir:
	//   .library
	//     init
	//     artists
	//     albums

	syncHelper := SyncHelper{
		artists: map[string]string{},
		albums:  map[string]struct{}{},
		tracks:  map[string]struct{}{},
	}

	_ = syncHelper

	slog.Info("Starting library sync...")
	defer slog.Info("Stopped library sync")

	ctx := context.TODO()

	// TODO(patrik): Check for deleted artists
	err := s.syncArtists(ctx, p)
	if err != nil {
		return err
	}

	// TODO(patrik): Check for deleted albums
	err = s.syncAlbums(ctx, p)
	if err != nil {
		return err
	}

	// TODO(patrik): Check for deleted tracks
	err = s.syncTracks(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (s *LibraryService) Sync() {
	if s.syncRunning.Load() {
		slog.Error("library syncing already running")
		return
	}

	s.syncRunning.Store(true)
	defer s.syncRunning.Store(false)

	err := s.runSync()
	if err != nil {
		slog.Error("failed to run sync", "err", err)
		return
	}
}

func fixArr(arr []string) []string {
	seen := map[string]bool{}
	res := make([]string, 0, len(arr))

	for _, value := range arr {
		value = anvil.String(value)
		if value == "" {
			continue
		}

		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

const (
	UNKNOWN_ARTIST_ID   = "unknown"
	UNKNOWN_ARTIST_NAME = "UNKNOWN"
)

// TODO(patrik): Add testing for this
// TODO(patrik): This needs some work
func FixMetadata(metadata *Metadata) error {
	album := &metadata.Album

	album.Name = anvil.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	if len(album.Artists) == 0 {
		// TODO(patrik): Instead of this just validate the metadata
		// and reject it
		album.Artists = []string{UNKNOWN_ARTIST_NAME}
	}

	album.Artists = fixArr(album.Artists)

	album.Tags = append(album.Tags, metadata.General.Tags...)
	for i, tag := range album.Tags {
		album.Tags[i] = utils.Slug(strings.TrimSpace(tag))
	}

	album.Tags = fixArr(album.Tags)

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = anvil.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)
		for i, tag := range t.Tags {
			t.Tags[i] = utils.Slug(strings.TrimSpace(tag))
		}

		t.Tags = fixArr(t.Tags)

		if len(t.Artists) == 0 {
			// TODO(patrik): Instead of this just validate the metadata
			// and reject it
			t.Artists = []string{UNKNOWN_ARTIST_NAME}
		}

		t.Artists = fixArr(t.Artists)
	}

	// err := validate.ValidateStruct(&metadata.Album,
	// 	validate.Field(&metadata.Album.Name, validate.Required),
	// 	validate.Field(&metadata.Album.Artists, validate.Length(1, 0)),
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// for _, track := range metadata.Tracks {
	// 	err := validate.ValidateStruct(&track,
	// 		validate.Field(&track.File, validate.Required),
	// 		validate.Field(&track.Name, validate.Required),
	// 		validate.Field(&track.Artists, validate.Length(1, 0)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SyncHelper struct {
	artists map[string]string

	albums map[string]struct{}
	tracks map[string]struct{}
}

func (helper *SyncHelper) getOrCreateArtist(ctx context.Context, db *database.Database, name string) (string, error) {
	slug := utils.Slug(name)

	if artist, exists := helper.artists[slug]; exists {
		return artist, nil
	}

	dbArtist, err := db.GetArtistBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbArtist, err = db.CreateArtist(ctx, database.CreateArtistParams{
				Slug: slug,
				Name: name,
			})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	helper.artists[slug] = dbArtist.Id
	return dbArtist.Id, nil
}

func (helper *SyncHelper) setAlbumFeaturingArtists(ctx context.Context, db *database.Database, albumId string, artists []string) error {
	err := db.RemoveAllAlbumFeaturingArtists(ctx, albumId)
	if err != nil {
		return err
	}

	for _, artistName := range artists {
		artistId, err := helper.getOrCreateArtist(ctx, db, artistName)
		if err != nil {
			return err
		}

		err = db.AddFeaturingArtistToAlbum(ctx, albumId, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setAlbumTags(ctx context.Context, db *database.Database, albumId string, tags []string) error {
	err := db.RemoveAllTagsFromAlbum(ctx, albumId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToAlbum(ctx, slug, albumId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setTrackFeaturingArtists(ctx context.Context, db *database.Database, trackId string, artists []string) error {
	err := db.RemoveAllTrackFeaturingArtists(ctx, trackId)
	if err != nil {
		return err
	}

	for _, artistName := range artists {
		artistId, err := helper.getOrCreateArtist(ctx, db, artistName)
		if err != nil {
			return err
		}

		err = db.AddFeaturingArtistToTrack(ctx, trackId, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setTrackTags(ctx context.Context, db *database.Database, trackId string, tags []string) error {
	err := db.RemoveAllTagsFromTrack(ctx, trackId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToTrack(ctx, slug, trackId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

// TODO(patrik): Update the errors for album
func (helper *SyncHelper) syncAlbum(ctx context.Context, metadata *Metadata, db *database.Database) error {
	err := FixMetadata(metadata)
	if err != nil {
		return err
	}

	dbAlbum, err := db.GetAlbumById(ctx, metadata.Album.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			artist, err := helper.getOrCreateArtist(ctx, db, metadata.Album.Artists[0])
			if err != nil {
				return fmt.Errorf("failed to create artist for album: %w", err)
			}

			dbAlbum, err = db.CreateAlbum(ctx, database.CreateAlbumParams{
				Id:       metadata.Album.Id,
				Name:     metadata.Album.Name,
				ArtistId: artist,
			})
			if err != nil {
				return fmt.Errorf("failed to create album: %w", err)
			}
		} else {
			return err
		}
	}

	helper.albums[dbAlbum.Id] = struct{}{}

	changes := database.AlbumChanges{}

	// TODO(patrik): More updates

	changes.Name = types.Change[string]{
		Value:   metadata.Album.Name,
		Changed: metadata.Album.Name != dbAlbum.Name,
	}

	artist, err := helper.getOrCreateArtist(ctx, db, metadata.Album.Artists[0])
	if err != nil {
		return fmt.Errorf("failed to create artist for album: %w", err)
	}

	changes.ArtistId = types.Change[string]{
		Value:   artist,
		Changed: artist != dbAlbum.ArtistId,
	}

	changes.CoverArt = types.Change[sql.NullString]{
		Value: sql.NullString{
			String: metadata.General.Cover,
			Valid:  metadata.General.Cover != "",
		},
		Changed: metadata.General.Cover != dbAlbum.CoverArt.String,
	}

	changes.Year = types.Change[sql.NullInt64]{
		Value: sql.NullInt64{
			Int64: metadata.Album.Year,
			Valid: metadata.Album.Year != 0,
		},
		Changed: metadata.Album.Year != dbAlbum.Year.Int64,
	}

	err = db.UpdateAlbum(ctx, dbAlbum.Id, changes)
	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}

	err = helper.setAlbumFeaturingArtists(
		ctx,
		db,
		dbAlbum.Id,
		metadata.Album.Artists[1:],
	)
	if err != nil {
		return fmt.Errorf("failed to set album featuring artists: %w", err)
	}

	err = helper.setAlbumTags(
		ctx,
		db,
		dbAlbum.Id,
		metadata.Album.Tags,
	)
	if err != nil {
		return fmt.Errorf("failed to set album tags: %w", err)
	}

	for i, track := range metadata.Tracks {
		stat, err := os.Stat(track.File)
		if err != nil {
			return fmt.Errorf("failed to stat track[%d] file (%s): %w", i, track.File, err)
		}

		modifiedTime := stat.ModTime().UnixMilli()

		artist, err := helper.getOrCreateArtist(ctx, db, track.Artists[0])
		if err != nil {
			return fmt.Errorf("failed to set create artist for track[%d]: %w", i, err)
		}

		dbTrack, err := db.GetTrackById(ctx, track.Id)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				probeResult, err := utils.ProbeTrack(track.File)
				if err != nil {
					return fmt.Errorf("failed to probe track[%d] file (%s): %w", i, track.File, err)
				}

				trackId, err := db.CreateTrack(ctx, database.CreateTrackParams{
					Id:           track.Id,
					Filename:     track.File,
					ModifiedTime: modifiedTime,
					MediaType:    probeResult.MediaType,
					Name:         track.Name,
					OtherName:    sql.NullString{},
					AlbumId:      dbAlbum.Id,
					ArtistId:     artist,
					Duration:     int64(probeResult.Duration),
					Number: sql.NullInt64{
						Int64: track.Number,
						Valid: track.Number != 0,
					},
					Year: sql.NullInt64{
						Int64: track.Year,
						Valid: track.Year != 0,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create track[%d]: %w", i, err)
				}

				helper.tracks[trackId] = struct{}{}

				err = helper.setTrackFeaturingArtists(
					ctx,
					db,
					trackId,
					track.Artists[1:],
				)
				if err != nil {
					return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
				}

				err = helper.setTrackTags(ctx, db, trackId, track.Tags)
				if err != nil {
					return fmt.Errorf("failed to set track[%d] tags: %w", i, err)
				}

				continue
			}
		}

		helper.tracks[dbTrack.Id] = struct{}{}

		err = helper.setTrackFeaturingArtists(
			ctx,
			db,
			dbTrack.Id,
			track.Artists[1:],
		)
		if err != nil {
			return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
		}

		err = helper.setTrackTags(ctx, db, dbTrack.Id, track.Tags)
		if err != nil {
			return fmt.Errorf("failed to set track[%d] tags: %w", i, err)
		}

		// TODO(patrik): Check modified time and probe again
		// TODO(patrik): Update track

		changes := database.TrackChanges{}

		if modifiedTime > dbTrack.ModifiedTime {
			probeResult, err := utils.ProbeTrack(track.File)
			if err != nil {
				return fmt.Errorf("failed to probe track[%d] file (%s): %w", i, track.File, err)
			}

			dur := int64(probeResult.Duration)
			changes.Duration = types.Change[int64]{
				Value:   dur,
				Changed: dur != dbTrack.Duration,
			}

			changes.MediaType = types.Change[types.MediaType]{
				Value:   probeResult.MediaType,
				Changed: probeResult.MediaType != dbTrack.MediaType,
			}

			changes.ModifiedTime = types.Change[int64]{
				Value:   modifiedTime,
				Changed: modifiedTime != dbTrack.ModifiedTime,
			}
		}

		// TODO(patrik): Implement all the changes here

		changes.Filename = types.Change[string]{
			Value:   track.File,
			Changed: track.File != dbTrack.Filename,
		}

		changes.Name = types.Change[string]{
			Value:   track.Name,
			Changed: track.Name != dbTrack.Name,
		}

		changes.ArtistId = types.Change[string]{
			Value:   artist,
			Changed: artist != dbTrack.ArtistId,
		}

		changes.Number = types.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: track.Number,
				Valid: track.Number != 0,
			},
			Changed: track.Number != dbTrack.Number.Int64,
		}

		changes.Year = types.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: track.Year,
				Valid: track.Year != 0,
			},
			Changed: track.Year != dbTrack.Year.Int64,
		}

		err = db.UpdateTrack(ctx, dbTrack.Id, changes)
		if err != nil {
			return fmt.Errorf("failed to update track[%d]: %w", i, err)
		}
	}

	return nil
}
