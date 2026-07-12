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
	"time"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/library"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

const (
	maxSyncErrors = 5

	syncArtistsFile = "artists"
	syncAlbumsFile  = "albums"
	syncTracksFile  = "tracks"
)

var _ (broker.Event) = (*LibrarySyncStateEvent)(nil)
var _ broker.EventProducer = (*LibraryService)(nil)

type MissingItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LibrarySyncStateEvent struct {
	Errors []string `json:"errors"`

	NumArtists int `json:"numArtists"`
	NumAlbums  int `json:"numAlbums"`
	NumTracks  int `json:"numTracks"`

	MissingArtists []MissingItem `json:"missingArtists"`
	MissingAlbums  []MissingItem `json:"missingAlbums"`
	MissingTracks  []MissingItem `json:"missingTracks"`

	ArtistsSyncDurationMs int64 `json:"artistsSyncDurationMs"`
	AlbumsSyncDurationMs  int64 `json:"albumsSyncDurationMs"`
	TracksSyncDurationMs  int64 `json:"tracksSyncDurationMs"`
	TotalSyncDurationMs   int64 `json:"totalSyncDurationMs"`
}

func (e LibrarySyncStateEvent) GetEventType() string {
	return "library-sync-state"
}

var libraryErr = NewServiceErrCreator("library")

type LibraryService struct {
	logger *slog.Logger

	db      *database.Database
	dataDir types.DataDir
	config  *config.Config

	norificationService *NotificationService
	mediaService        *MediaService

	emitter broker.EventEmitter

	errors []error

	numArtists int
	numAlbums  int
	numTracks  int

	missingArtists []MissingItem
	missingAlbums  []MissingItem
	missingTracks  []MissingItem

	syncedArtistIds map[string]struct{}
	syncedAlbumIds  map[string]struct{}
	syncedTrackIds  map[string]struct{}

	artistsSyncDuration time.Duration
	albumsSyncDuration  time.Duration
	tracksSyncDuration  time.Duration
	totalSyncDuration   time.Duration
}

func NewLibraryService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
	config *config.Config,
	notificationService *NotificationService,
	mediaService *MediaService,
	emitter broker.EventEmitter,
) *LibraryService {
	return &LibraryService{
		logger:              logger,
		db:                  db,
		dataDir:             dataDir,
		config:              config,
		norificationService: notificationService,
		mediaService:        mediaService,
		emitter:             emitter,

		// NOTE(patrik): We need to initialize these or else the frontend
		// gets an error because these are nil
		missingArtists: []MissingItem{},
		missingAlbums:  []MissingItem{},
		missingTracks:  []MissingItem{},
	}
}

func (s *LibraryService) update() {
	if s.emitter != nil {
		s.emitter.EmitEvent(s.GetSyncStateEvent())
	}
}

func (s *LibraryService) GetInitEvents() []broker.Event {
	return []broker.Event{s.GetSyncStateEvent()}
}

func (s *LibraryService) GetSyncStateEvent() LibrarySyncStateEvent {
	errors := make([]string, len(s.errors))

	for i, err := range s.errors {
		errors[i] = err.Error()
	}

	return LibrarySyncStateEvent{
		Errors:                errors,
		NumArtists:            s.numArtists,
		NumAlbums:             s.numAlbums,
		NumTracks:             s.numTracks,
		MissingArtists:        s.missingArtists,
		MissingAlbums:         s.missingAlbums,
		MissingTracks:         s.missingTracks,
		ArtistsSyncDurationMs: s.artistsSyncDuration.Milliseconds(),
		AlbumsSyncDurationMs:  s.albumsSyncDuration.Milliseconds(),
		TracksSyncDurationMs:  s.tracksSyncDuration.Milliseconds(),
		TotalSyncDurationMs:   s.totalSyncDuration.Milliseconds(),
	}
}

func (s *LibraryService) addError(err error) bool {
	s.errors = append(s.errors, err)

	s.update()

	return len(s.errors) >= maxSyncErrors
}

func setArtistTags(
	ctx context.Context,
	db database.DB,
	artistId string,
	tags []string,
) error {
	err := db.RemoveAllArtistTags(ctx, artistId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := slug.Make(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddArtistTag(ctx, slug, artistId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncSingleArtist(
	ctx context.Context,
	entry *library.ArtistEntry,
) error {
	dbArtist, err := s.db.GetArtistById(ctx, entry.Id)
	if err != nil {
		if !errors.Is(err, database.ErrItemNotFound) {
			return fmt.Errorf("get artist: %w", err)
		}

		return s.createArtist(ctx, entry)
	}

	return s.updateArtist(ctx, entry, dbArtist)
}

func (s *LibraryService) createArtist(
	ctx context.Context,
	entry *library.ArtistEntry,
) error {
	coverArt := entry.GetCoverArt()

	_, err := s.db.CreateArtist(ctx, database.CreateArtistParams{
		Id:   entry.Id,
		Name: entry.Name,
		CoverArt: sql.NullString{
			String: coverArt,
			Valid:  coverArt != "",
		},
	})
	if err != nil {
		return fmt.Errorf("create artist: %w", err)
	}

	return s.setArtistRelations(ctx, entry)
}

func (s *LibraryService) updateArtist(
	ctx context.Context,
	entry *library.ArtistEntry,
	dbArtist database.Artist,
) error {
	changes := buildArtistChanges(entry, dbArtist)

	err := s.db.UpdateArtist(ctx, dbArtist.Id, changes)
	if err != nil {
		return fmt.Errorf("update artist: %w", err)
	}

	return s.setArtistRelations(ctx, entry)
}

func buildArtistChanges(
	entry *library.ArtistEntry,
	dbArtist database.Artist,
) database.ArtistChanges {
	coverArt := entry.GetCoverArt()

	changes := database.ArtistChanges{}

	changes.Name = database.Change[string]{
		Value:   entry.Name,
		Changed: entry.Name != dbArtist.Name,
	}

	changes.CoverArt = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: coverArt,
			Valid:  coverArt != "",
		},
		Changed: coverArt != dbArtist.CoverArt.String,
	}

	return changes
}

func (s *LibraryService) setArtistRelations(
	ctx context.Context,
	entry *library.ArtistEntry,
) error {
	err := setArtistTags(
		ctx,
		s.db.DB,
		entry.Id,
		entry.Tags,
	)
	if err != nil {
		return fmt.Errorf("set artist tags: %w", err)
	}

	return nil
}

func (s *LibraryService) syncArtists(
	ctx context.Context,
	libraryDir string,
) error {
	p := path.Join(libraryDir, syncArtistsFile)
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for idx := 0; decoder.More(); idx++ {
		var entry library.ArtistEntry
		err := decoder.Decode(&entry)
		if err != nil {
			stop := s.addError(
				fmt.Errorf(
					"decode next artist entry[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		err = entry.Validate()
		if err != nil {
			stop := s.addError(
				fmt.Errorf("entry validation[%d]: %w", idx, err),
			)
			if stop {
				break
			}

			continue
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		err = s.syncSingleArtist(ctx, &entry)
		if err != nil {
			stop := s.addError(
				fmt.Errorf("sync artist[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		s.syncedArtistIds[entry.Id] = struct{}{}
		s.numArtists += 1
		s.update()
	}

	return nil
}

func setAlbumFeaturingArtists(
	ctx context.Context,
	db database.DB,
	albumId string,
	artistIds []string,
) error {
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

func setAlbumTags(
	ctx context.Context,
	db database.DB,
	albumId string,
	tags []string,
) error {
	err := db.RemoveAllAlbumTags(ctx, albumId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := slug.Make(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddAlbumTag(ctx, slug, albumId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncSingleAlbum(
	ctx context.Context,
	entry *library.AlbumEntry,
) error {
	dbAlbum, err := s.db.GetAlbumById(ctx, entry.Id)
	if err != nil {
		if !errors.Is(err, database.ErrItemNotFound) {
			return fmt.Errorf("get album: %w", err)
		}

		return s.createAlbum(ctx, entry)
	}

	return s.updateAlbum(ctx, entry, dbAlbum)
}

func (s *LibraryService) createAlbum(
	ctx context.Context,
	entry *library.AlbumEntry,
) error {
	coverArt := entry.GetCoverArt()

	_, err := s.db.CreateAlbum(ctx, database.CreateAlbumParams{
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
		AlbumType: entry.AlbumType,
	})
	if err != nil {
		return fmt.Errorf("create album: %w", err)
	}

	return s.setAlbumRelations(ctx, entry)
}

func (s *LibraryService) updateAlbum(
	ctx context.Context,
	entry *library.AlbumEntry,
	dbAlbum database.Album,
) error {
	changes := buildAlbumChanges(entry, dbAlbum)

	err := s.db.UpdateAlbum(ctx, dbAlbum.Id, changes)
	if err != nil {
		return fmt.Errorf("update album: %w", err)
	}

	return s.setAlbumRelations(ctx, entry)
}

func buildAlbumChanges(
	entry *library.AlbumEntry,
	dbAlbum database.Album,
) database.AlbumChanges {
	coverArt := entry.GetCoverArt()

	changes := database.AlbumChanges{}

	changes.Name = database.Change[string]{
		Value:   entry.Name,
		Changed: entry.Name != dbAlbum.Name,
	}

	changes.ArtistId = database.Change[string]{
		Value:   entry.ArtistId,
		Changed: entry.ArtistId != dbAlbum.ArtistId,
	}

	changes.CoverArt = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: coverArt,
			Valid:  coverArt != "",
		},
		Changed: coverArt != dbAlbum.CoverArt.String,
	}

	changes.Year = database.Change[sql.NullInt64]{
		Value: sql.NullInt64{
			Int64: entry.Year,
			Valid: entry.Year != 0,
		},
		Changed: entry.Year != dbAlbum.Year.Int64,
	}

	changes.AlbumType = database.Change[library.AlbumType]{
		Value:   entry.AlbumType,
		Changed: entry.AlbumType != dbAlbum.AlbumType,
	}

	return changes
}

func (s *LibraryService) setAlbumRelations(
	ctx context.Context,
	entry *library.AlbumEntry,
) error {
	err := setAlbumFeaturingArtists(
		ctx,
		s.db.DB,
		entry.Id,
		entry.FeaturingArtistIds,
	)
	if err != nil {
		return fmt.Errorf("set album featuring artists: %w", err)
	}

	err = setAlbumTags(
		ctx,
		s.db.DB,
		entry.Id,
		entry.Tags,
	)
	if err != nil {
		return fmt.Errorf("set album tags: %w", err)
	}

	return nil
}

func (s *LibraryService) syncAlbums(
	ctx context.Context,
	libraryDir string,
) error {
	p := path.Join(libraryDir, syncAlbumsFile)
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for idx := 0; decoder.More(); idx++ {
		var entry library.AlbumEntry
		err := decoder.Decode(&entry)
		if err != nil {
			stop := s.addError(
				fmt.Errorf(
					"decode next album entry[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		err = entry.Validate()
		if err != nil {
			stop := s.addError(
				fmt.Errorf("entry validation[%d]: %w", idx, err),
			)
			if stop {
				break
			}

			continue
		}

		if entry.Path != "" {
			entry.Path = path.Join(libraryDir, entry.Path)
		}

		err = s.syncSingleAlbum(ctx, &entry)
		if err != nil {
			stop := s.addError(
				fmt.Errorf("sync album[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		s.syncedAlbumIds[entry.Id] = struct{}{}
		s.numAlbums += 1
		s.update()
	}

	return nil
}

func setTrackFeaturingArtists(
	ctx context.Context,
	db database.DB,
	trackId string,
	artistIds []string,
) error {
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

func setTrackTags(
	ctx context.Context,
	db database.DB,
	trackId string,
	tags []string,
) error {
	err := db.RemoveAllTrackTags(ctx, trackId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := slug.Make(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTrackTag(ctx, slug, trackId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (s *LibraryService) syncSingleTrack(
	ctx context.Context,
	entry *library.TrackEntry,
) error {
	trackFile := entry.GetTrackFile()

	stat, err := os.Stat(trackFile)
	if err != nil {
		return fmt.Errorf(
			"stat track file %q (%s): %w", entry.Name, trackFile, err)
	}

	modifiedTime := stat.ModTime().UnixMilli()

	dbTrack, err := s.db.GetTrackById(ctx, entry.Id)
	if err != nil {
		if !errors.Is(err, database.ErrItemNotFound) {
			return fmt.Errorf("get track: %w", err)
		}

		return s.createTrack(ctx, entry, trackFile, modifiedTime)
	}

	return s.updateTrack(ctx, entry, dbTrack, trackFile, modifiedTime)
}

func (s *LibraryService) createTrack(
	ctx context.Context,
	entry *library.TrackEntry,
	trackFile string,
	modifiedTime int64,
) error {
	probeResult, err := s.mediaService.ProbeMedia(ctx, trackFile)
	if err != nil {
		return fmt.Errorf(
			"probe new track %q (%s): %w", entry.Name, trackFile, err)
	}

	_, err = s.db.CreateTrack(ctx, database.CreateTrackParams{
		Id:           entry.Id,
		Filename:     trackFile,
		ModifiedTime: modifiedTime,
		MediaFormat:  probeResult.MediaFormat,
		Name:         entry.Name,
		AlbumId:      entry.AlbumId,
		ArtistId:     entry.ArtistId,
		Duration:     int64(probeResult.Duration.Seconds()),
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
		return fmt.Errorf("create track: %w", err)
	}

	return s.setTrackRelations(ctx, entry)
}

func (s *LibraryService) updateTrack(
	ctx context.Context,
	entry *library.TrackEntry,
	dbTrack database.Track,
	trackFile string,
	modifiedTime int64,
) error {
	changes, err := s.buildTrackChanges(ctx, entry, dbTrack, trackFile, modifiedTime)
	if err != nil {
		return err
	}

	err = s.db.UpdateTrack(ctx, dbTrack.Id, changes)
	if err != nil {
		return fmt.Errorf("update track: %w", err)
	}

	return s.setTrackRelations(ctx, entry)
}

func (s *LibraryService) buildTrackChanges(
	ctx context.Context,
	entry *library.TrackEntry,
	dbTrack database.Track,
	trackFile string,
	modifiedTime int64,
) (database.TrackChanges, error) {
	changes := database.TrackChanges{}

	if modifiedTime > dbTrack.ModifiedTime || dbTrack.Filename != trackFile {
		probeResult, err := s.mediaService.ProbeMedia(ctx, trackFile)
		if err != nil {
			return database.TrackChanges{}, fmt.Errorf(
				"probe updated track %q (%s): %w",
				entry.Name, trackFile, err)
		}

		dur := int64(probeResult.Duration.Seconds())
		changes.Duration = database.Change[int64]{
			Value:   dur,
			Changed: dur != dbTrack.Duration,
		}

		changes.MediaFormat = database.Change[types.MediaFormat]{
			Value:   probeResult.MediaFormat,
			Changed: probeResult.MediaFormat != dbTrack.MediaFormat,
		}

		changes.ModifiedTime = database.Change[int64]{
			Value:   modifiedTime,
			Changed: modifiedTime != dbTrack.ModifiedTime,
		}
	}

	changes.Filename = database.Change[string]{
		Value:   trackFile,
		Changed: trackFile != dbTrack.Filename,
	}

	changes.Name = database.Change[string]{
		Value:   entry.Name,
		Changed: entry.Name != dbTrack.Name,
	}

	changes.AlbumId = database.Change[string]{
		Value:   entry.AlbumId,
		Changed: entry.AlbumId != dbTrack.AlbumId,
	}

	changes.ArtistId = database.Change[string]{
		Value:   entry.ArtistId,
		Changed: entry.ArtistId != dbTrack.ArtistId,
	}

	changes.Number = database.Change[sql.NullInt64]{
		Value: sql.NullInt64{
			Int64: entry.Number,
			Valid: entry.Number != 0,
		},
		Changed: entry.Number != dbTrack.Number.Int64,
	}

	changes.Year = database.Change[sql.NullInt64]{
		Value: sql.NullInt64{
			Int64: entry.Year,
			Valid: entry.Year != 0,
		},
		Changed: entry.Year != dbTrack.Year.Int64,
	}

	return changes, nil
}

func (s *LibraryService) setTrackRelations(
	ctx context.Context,
	entry *library.TrackEntry,
) error {
	err := setTrackFeaturingArtists(
		ctx,
		s.db.DB,
		entry.Id,
		entry.FeaturingArtistIds,
	)
	if err != nil {
		return fmt.Errorf("set track featuring artists: %w", err)
	}

	err = setTrackTags(
		ctx,
		s.db.DB,
		entry.Id,
		entry.Tags,
	)
	if err != nil {
		return fmt.Errorf("set track tags: %w", err)
	}

	return nil
}

func (s *LibraryService) syncTracks(
	ctx context.Context,
	libraryDir string,
) error {
	p := path.Join(libraryDir, syncTracksFile)
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for idx := 0; decoder.More(); idx++ {
		var entry library.TrackEntry
		err := decoder.Decode(&entry)
		if err != nil {
			stop := s.addError(
				fmt.Errorf(
					"decode next track entry[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		err = entry.Validate()
		if err != nil {
			stop := s.addError(
				fmt.Errorf("entry validation[%d]: %w", idx, err),
			)
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
			stop := s.addError(
				fmt.Errorf("sync track[%d]: %w", idx, err))
			if stop {
				break
			}

			continue
		}

		s.syncedTrackIds[entry.Id] = struct{}{}
		s.numTracks += 1
		s.update()
	}

	return nil
}

func (s *LibraryService) Sync(ctx context.Context) error {
	s.logger.Info("starting library sync...")

	s.norificationService.SendSimple(SendSimpleParams{
		Title:    "Starting library sync",
		Message:  "Starting to sync the library",
		Tags:     []string{"library", "syncing"},
		Priority: 0,
	})

	defer s.finishSync()

	s.initSyncState()
	s.update()

	p := s.config.LibraryDir

	artistTimer := utils.SimpleTimer{}
	artistTimer.Start()

	err := s.syncArtists(ctx, p)
	artistTimer.Stop()
	if err != nil {
		return libraryErr.Wrap("sync artists", err)
	}

	albumTimer := utils.SimpleTimer{}
	albumTimer.Start()

	err = s.syncAlbums(ctx, p)
	albumTimer.Stop()
	if err != nil {
		return libraryErr.Wrap("sync albums", err)
	}

	trackTimer := utils.SimpleTimer{}
	trackTimer.Start()

	err = s.syncTracks(ctx, p)
	trackTimer.Stop()
	if err != nil {
		return libraryErr.Wrap("sync tracks", err)
	}

	existingArtistIds, err := s.db.GetAllArtistIds(ctx)
	if err != nil {
		return libraryErr.Wrap("get all artist ids", err)
	}

	existingAlbumIds, err := s.db.GetAllAlbumIds(ctx)
	if err != nil {
		return libraryErr.Wrap("get all album ids", err)
	}

	existingTrackIds, err := s.db.GetAllTrackIds(ctx)
	if err != nil {
		return libraryErr.Wrap("get all track ids", err)
	}

	s.findMissingArtists(ctx, existingArtistIds)
	s.findMissingAlbums(ctx, existingAlbumIds)
	s.findMissingTracks(ctx, existingTrackIds)

	s.artistsSyncDuration = artistTimer.Duration()
	s.albumsSyncDuration = albumTimer.Duration()
	s.tracksSyncDuration = trackTimer.Duration()
	s.totalSyncDuration =
		s.artistsSyncDuration + s.albumsSyncDuration + s.tracksSyncDuration

	return s.clearCache()
}

func (s *LibraryService) initSyncState() {
	s.errors = []error{}
	s.numArtists = 0
	s.numAlbums = 0
	s.numTracks = 0

	s.artistsSyncDuration = 0
	s.albumsSyncDuration = 0
	s.tracksSyncDuration = 0
	s.totalSyncDuration = 0

	s.missingArtists = []MissingItem{}
	s.missingAlbums = []MissingItem{}
	s.missingTracks = []MissingItem{}

	s.syncedArtistIds = make(map[string]struct{})
	s.syncedAlbumIds = make(map[string]struct{})
	s.syncedTrackIds = make(map[string]struct{})
}

func (s *LibraryService) finishSync() {
	s.update()

	s.logger.Info("stopped library sync")

	message := fmt.Sprintf(
		"%d Error(s)\nTotal sync time %s",
		len(s.errors),
		utils.PrettyDuration(s.totalSyncDuration),
	)

	tags := []string{"library", "syncing"}

	if len(s.errors) > 0 {
		tags = append(tags, "warning")
	}

	s.norificationService.SendSimple(SendSimpleParams{
		Title:   "Stopped library sync",
		Message: message,
		Tags:    tags,
	})
}

func (s *LibraryService) findMissingArtists(
	ctx context.Context,
	existingIds []string,
) {
	for _, id := range existingIds {
		if _, exists := s.syncedArtistIds[id]; exists {
			continue
		}

		artist, err := s.db.GetArtistById(ctx, id)
		if err != nil {
			continue
		}

		s.logger.Warn("missing artist",
			"id", id,
			"name", artist.Name,
		)

		s.missingArtists = append(s.missingArtists, MissingItem{
			Id:   id,
			Name: artist.Name,
		})
	}
}

func (s *LibraryService) findMissingAlbums(
	ctx context.Context,
	existingIds []string,
) {
	for _, id := range existingIds {
		if _, exists := s.syncedAlbumIds[id]; exists {
			continue
		}

		album, err := s.db.GetAlbumById(ctx, id)
		if err != nil {
			continue
		}

		s.logger.Warn("missing album",
			"id", id,
			"name", album.Name,
			"artist", album.ArtistName,
		)

		displayName := fmt.Sprintf(
			"%s (%s)", album.Name, album.ArtistName)

		s.missingAlbums = append(s.missingAlbums, MissingItem{
			Id:   id,
			Name: displayName,
		})
	}
}

func (s *LibraryService) findMissingTracks(
	ctx context.Context,
	existingIds []string,
) {
	for _, id := range existingIds {
		if _, exists := s.syncedTrackIds[id]; exists {
			continue
		}

		track, err := s.db.GetTrackById(ctx, id)
		if err != nil {
			continue
		}

		s.logger.Warn("missing track",
			"id", id,
			"name", track.Name,
			"album", track.AlbumName,
			"artist", track.ArtistName,
		)

		displayName := fmt.Sprintf(
			"%s (%s) (%s)",
			track.Name, track.AlbumName, track.ArtistName)

		s.missingTracks = append(s.missingTracks, MissingItem{
			Id:   id,
			Name: displayName,
		})
	}
}

func (s *LibraryService) clearCache() error {
	dir := s.dataDir.Cache()
	s.logger.Info("clearing the cache", "path", dir)

	err := os.RemoveAll(dir)
	if err != nil {
		return libraryErr.Wrap("clear cache", err)
	}

	err = utils.CreateDirectories([]string{dir})
	if err != nil {
		return libraryErr.Wrap("create cache dir", err)
	}

	return nil
}

func (s *LibraryService) Cleanup(ctx context.Context) error {
	deletedArtists := 0
	deletedAlbums := 0
	deletedTracks := 0

	for _, item := range s.missingTracks {
		err := s.db.DeleteTrack(ctx, item.Id)
		if err != nil {
			s.logger.Warn(
				"delete track",
				"id", item.Id,
				"name", item.Name,
				"err", err,
			)

			continue
		}

		deletedTracks++
	}

	for _, item := range s.missingAlbums {
		err := s.db.DeleteAlbum(ctx, item.Id)
		if err != nil {
			s.logger.Warn(
				"delete album",
				"id", item.Id,
				"name", item.Name,
				"err", err,
			)

			continue
		}

		deletedAlbums++
	}

	for _, item := range s.missingArtists {
		err := s.db.DeleteArtist(ctx, item.Id)
		if err != nil {
			s.logger.Warn(
				"delete artist",
				"id", item.Id,
				"name", item.Name,
				"err", err,
			)

			continue
		}

		deletedArtists++
	}

	s.logger.Info("library cleanup completed",
		"deletedArtists", deletedArtists,
		"deletedAlbums", deletedAlbums,
		"deletedTracks", deletedTracks,
	)

	return nil
}
