package service

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/nrednav/cuid2"
)

var filesystemErr = NewServiceErrCreator("filesystem")

var newFsTempFileId, _ = cuid2.Init(cuid2.WithLength(16))

type DataDir string

func (d DataDir) String() string {
	return string(d)
}

func (d DataDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d DataDir) Temp() string {
	return path.Join(d.String(), "temp")
}

func (d DataDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d DataDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d DataDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d DataDir) User(id string) string {
	return path.Join(d.Users(), id)
}

func (d DataDir) Cache() string {
	return path.Join(d.String(), "cache")
}

func (d DataDir) CacheImages() ImageCacheDir {
	return ImageCacheDir(path.Join(d.Cache(), "images"))
}

func (d DataDir) CacheTranscoding() string {
	return path.Join(d.Cache(), "transcoding")
}

type ImageCacheDir string

func (d ImageCacheDir) String() string {
	return string(d)
}

func (d ImageCacheDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d ImageCacheDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d ImageCacheDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d ImageCacheDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d ImageCacheDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d ImageCacheDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

func (d ImageCacheDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d ImageCacheDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d ImageCacheDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d ImageCacheDir) User(id string) string {
	return path.Join(d.Users(), id)
}

type FilesystemService struct {
	logger  *slog.Logger
	dataDir DataDir
}

func NewFilesystemService(
	logger *slog.Logger,
	dataDir string,
) *FilesystemService {
	return &FilesystemService{
		logger:  logger,
		dataDir: DataDir(dataDir),
	}
}

// DataDir returns the underlying DataDir path
func (s *FilesystemService) DataDir() DataDir {
	return s.dataDir
}

// DatabaseFile returns the path to the database file
func (s *FilesystemService) DatabaseFile() string {
	return s.dataDir.DatabaseFile()
}

// EnsureBaseDirs creates the base data directories
func (s *FilesystemService) EnsureBaseDirs() error {
	return s.EnsureDirs([]string{
		s.dataDir.Users(),
		s.dataDir.Playlists(),
		s.dataDir.Cache(),
		s.dataDir.Temp(),
	})
}

// Generic filesystem operations

func (s *FilesystemService) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (s *FilesystemService) EnsureDir(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return nil
}

func (s *FilesystemService) EnsureDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := s.EnsureDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesystemService) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *FilesystemService) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (s *FilesystemService) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (s *FilesystemService) RemoveFile(path string) error {
	return os.Remove(path)
}

// Temp file operations

func (s *FilesystemService) TempDir() string {
	return s.dataDir.Temp()
}

func (s *FilesystemService) CreateTempFile(ext string) (string, error) {
	dir := s.dataDir.Temp()
	if err := s.EnsureDir(dir); err != nil {
		return "", filesystemErr.Wrap("create temp file: ensure dir", err)
	}

	name := fmt.Sprintf("%s%s", newFsTempFileId(), ext)
	return path.Join(dir, name), nil
}

// Cache operations

func (s *FilesystemService) CacheDir() string {
	return s.dataDir.Cache()
}

func (s *FilesystemService) ClearCache() error {
	dir := s.dataDir.Cache()
	s.logger.Info("clearing the cache", "path", dir)

	err := os.RemoveAll(dir)
	if err != nil {
		return filesystemErr.Wrap("clear cache", err)
	}

	err = s.EnsureDir(dir)
	if err != nil {
		return filesystemErr.Wrap("clear cache: ensure dir", err)
	}

	return nil
}

func (s *FilesystemService) EnsureCacheDir() error {
	return s.EnsureDir(s.dataDir.Cache())
}

// Playlist operations

func (s *FilesystemService) PlaylistDir(playlistId string) string {
	return s.dataDir.Playlist(playlistId)
}

func (s *FilesystemService) EnsurePlaylistDir(playlistId string) error {
	return s.EnsureDir(s.dataDir.Playlist(playlistId))
}

func (s *FilesystemService) RemovePlaylistDir(playlistId string) error {
	return os.RemoveAll(s.dataDir.Playlist(playlistId))
}

func (s *FilesystemService) PlaylistImagePath(playlistId string) string {
	return s.dataDir.CacheImages().Playlist(playlistId)
}

func (s *FilesystemService) EnsurePlaylistImageCacheDirs(playlistId string) error {
	cacheDir := s.dataDir.CacheImages()
	return s.EnsureDirs([]string{
		cacheDir.String(),
		cacheDir.Playlists(),
		cacheDir.Playlist(playlistId),
	})
}

func (s *FilesystemService) ClearPlaylistImageCache(playlistId string) error {
	return os.RemoveAll(s.dataDir.CacheImages().Playlist(playlistId))
}

// User operations

func (s *FilesystemService) UserDir(userId string) string {
	return s.dataDir.User(userId)
}

func (s *FilesystemService) EnsureUserDir(userId string) error {
	return s.EnsureDir(s.dataDir.User(userId))
}

func (s *FilesystemService) UserImagePath(userId string) string {
	return s.dataDir.CacheImages().User(userId)
}

func (s *FilesystemService) EnsureUserImageCacheDirs(userId string) error {
	cacheDir := s.dataDir.CacheImages()
	return s.EnsureDirs([]string{
		cacheDir.String(),
		cacheDir.Users(),
		cacheDir.User(userId),
	})
}

func (s *FilesystemService) ClearUserImageCache(userId string) error {
	return os.RemoveAll(s.dataDir.CacheImages().User(userId))
}

// Artist operations

func (s *FilesystemService) ArtistImagePath(artistId string) string {
	return s.dataDir.CacheImages().Artist(artistId)
}

func (s *FilesystemService) EnsureArtistImageCacheDirs(artistId string) error {
	cacheDir := s.dataDir.CacheImages()
	return s.EnsureDirs([]string{
		cacheDir.String(),
		cacheDir.Artists(),
		cacheDir.Artist(artistId),
	})
}

func (s *FilesystemService) ClearArtistImageCache(artistId string) error {
	return os.RemoveAll(s.dataDir.CacheImages().Artist(artistId))
}

// Album operations

func (s *FilesystemService) AlbumImagePath(albumId string) string {
	return s.dataDir.CacheImages().Album(albumId)
}

func (s *FilesystemService) EnsureAlbumImageCacheDirs(albumId string) error {
	cacheDir := s.dataDir.CacheImages()
	return s.EnsureDirs([]string{
		cacheDir.String(),
		cacheDir.Albums(),
		cacheDir.Album(albumId),
	})
}

func (s *FilesystemService) ClearAlbumImageCache(albumId string) error {
	return os.RemoveAll(s.dataDir.CacheImages().Album(albumId))
}

// Transcoding cache operations

func (s *FilesystemService) TranscodingCacheDir() string {
	return s.dataDir.CacheTranscoding()
}

func (s *FilesystemService) TranscodingTrackDir(trackId string) string {
	return path.Join(s.dataDir.CacheTranscoding(), "tracks", trackId)
}

func (s *FilesystemService) EnsureTranscodingCacheDirs() error {
	cacheDir := s.dataDir.CacheTranscoding()
	tracksDir := path.Join(cacheDir, "tracks")
	return s.EnsureDirs([]string{
		cacheDir,
		tracksDir,
	})
}

func (s *FilesystemService) EnsureTranscodingTrackDir(trackId string) error {
	trackDir := s.TranscodingTrackDir(trackId)
	return s.EnsureDir(trackDir)
}

func (s *FilesystemService) TranscodingPath(trackId, filename string) string {
	return path.Join(s.TranscodingTrackDir(trackId), filename)
}
