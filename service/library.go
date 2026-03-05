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
	"sync/atomic"
	"time"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/broker"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
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

var _ (broker.Event) = (*LibrarySyncStateEvent)(nil)

type LibrarySyncStateEvent struct {
	IsRunning bool     `json:"isRunning"`
	Errors    []string `json:"errors"`

	NumArtists int `json:"numArtists"`
	NumAlbums  int `json:"numAlbums"`
	NumTracks  int `json:"numTracks"`

	ArtistsSyncDurationMs int64 `json:"artistsSyncDurationMs"`
	AlbumsSyncDurationMs  int64 `json:"albumsSyncDurationMs"`
	TracksSyncDurationMs  int64 `json:"tracksSyncDurationMs"`
	TotalSyncDurationMs   int64 `json:"totalSyncDurationMs"`
}

func (e LibrarySyncStateEvent) GetEventType() string {
	return "library-sync-state"
}

type UpdateFunc func()

type LibraryService struct {
	db     *database.Database
	config *config.Config

	searchService *SearchService

	errors      []error
	syncRunning atomic.Bool

	numArtists int
	numAlbums  int
	numTracks  int

	artistsSyncDuration time.Duration
	albumsSyncDuration  time.Duration
	tracksSyncDuration  time.Duration
	totalSyncDuration   time.Duration

	updateFunc UpdateFunc
}

func NewLibraryService(db *database.Database, config *config.Config, searchService *SearchService) *LibraryService {
	return &LibraryService{
		db:            db,
		config:        config,
		searchService: searchService,
		syncRunning:   atomic.Bool{},
	}
}

func (s *LibraryService) SetUpdateFunc(f UpdateFunc) {
	s.updateFunc = f
}

func (s *LibraryService) update() {
	if s.updateFunc != nil {
		s.updateFunc()
	}
}

func (s *LibraryService) GetSyncStateEvent() LibrarySyncStateEvent {
	errors := make([]string, len(s.errors))

	isRunning := s.syncRunning.Load()

	for i, err := range s.errors {
		errors[i] = err.Error()
	}

	return LibrarySyncStateEvent{
		IsRunning:           isRunning,
		Errors:              errors,
		NumArtists:          s.numArtists,
		NumAlbums:           s.numAlbums,
		NumTracks:           s.numTracks,
		ArtistsSyncDurationMs: s.artistsSyncDuration.Milliseconds(),
		AlbumsSyncDurationMs:  s.albumsSyncDuration.Milliseconds(),
		TracksSyncDurationMs:  s.tracksSyncDuration.Milliseconds(),
		TotalSyncDurationMs:   s.totalSyncDuration.Milliseconds(),
	}
}

func (s *LibraryService) addError(err error) bool {
	s.errors = append(s.errors, err)

	s.update()

	// TODO(patrik): Make constant
	return len(s.errors) >= 5
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

func (s *LibraryService) syncSingleArtist(ctx context.Context, entry *ArtistEntry) error {
	coverArt := entry.GetCoverArt()

	dbArtist, err := s.db.GetArtistById(ctx, entry.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			_, err = s.db.CreateArtist(ctx, database.CreateArtistParams{
				Id:   entry.Id,
				Slug: entry.Slug,
				Name: entry.Name,
				CoverArt: sql.NullString{
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

		changes.CoverArt = types.Change[sql.NullString]{
			Value: sql.NullString{
				String: coverArt,
				Valid:  coverArt != "",
			},
			Changed: coverArt != dbArtist.CoverArt.String,
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

		err = s.syncSingleArtist(ctx, &entry)
		if err != nil {
			return err
		}

		s.numArtists += 1
		s.update()
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

func (s *LibraryService) syncSingleAlbum(ctx context.Context, entry *AlbumEntry) error {
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

		err = s.syncSingleAlbum(ctx, &entry)
		if err != nil {
			return err
		}

		s.numAlbums += 1
		s.update()
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

func (s *LibraryService) syncSingleTrack(ctx context.Context, entry *TrackEntry) error {
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

	for idx := 0; decoder.More(); idx++ {
		var entry TrackEntry
		err := decoder.Decode(&entry)
		if err != nil {
			stop := s.addError(fmt.Errorf("failed to decode next track entry[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		err = s.syncSingleTrack(ctx, &entry)
		if err != nil {
			stop := s.addError(fmt.Errorf("failed to sync track[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		s.numTracks += 1
		s.update()
	}

	return nil
}

func (s *LibraryService) runSync() error {
	p := s.config.LibraryDir

	slog.Info("Starting library sync...")
	defer slog.Info("Stopped library sync")

	s.errors = []error{}
	s.numArtists = 0
	s.numAlbums = 0
	s.numTracks = 0

	s.update()

	ctx := context.TODO()

	artistTimer := utils.SimpleTimer{}
	artistTimer.Start()

	// TODO(patrik): Check for deleted artists
	err := s.syncArtists(ctx, p)
	if err != nil {
		return err
	}

	artistTimer.Stop()

	albumTimer := utils.SimpleTimer{}
	albumTimer.Start()

	// TODO(patrik): Check for deleted albums
	err = s.syncAlbums(ctx, p)
	if err != nil {
		return err
	}

	albumTimer.Stop()

	trackTimer := utils.SimpleTimer{}
	trackTimer.Start()

	// TODO(patrik): Check for deleted tracks
	err = s.syncTracks(ctx, p)
	if err != nil {
		return err
	}

	trackTimer.Stop()

	s.artistsSyncDuration = artistTimer.Duration()
	s.albumsSyncDuration = albumTimer.Duration()
	s.tracksSyncDuration = trackTimer.Duration()
	s.totalSyncDuration = s.artistsSyncDuration + s.albumsSyncDuration + s.tracksSyncDuration

	return nil
}

func (s *LibraryService) Sync() {
	if s.syncRunning.Load() {
		slog.Error("library syncing already running")
		return
	}

	s.syncRunning.Store(true)
	defer func() {
		s.syncRunning.Store(false)
		s.update()
	}()

	err := s.runSync()
	if err != nil {
		slog.Error("failed to run sync", "err", err)
		return
	}
}
