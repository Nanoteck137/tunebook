package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

var (
	ErrPlaylistServicePlaylistNotFound = errors.New("playlist service: playlist not found")
)

type PlaylistService struct {
	logger  *slog.Logger

	db      *database.Database
	dataDir types.DataDir

	imageService *ImageService
}

func NewPlaylistService(
	logger *slog.Logger, 
	db *database.Database, 
	dataDir types.DataDir, 
	imageService *ImageService,
) *PlaylistService {
	return &PlaylistService{
		logger: logger,
		db:     db,
		dataDir: dataDir,
		imageService: imageService,
	}
}

func (s *PlaylistService) GetPlaylistsByUser(
	ctx context.Context,
	userId string,
) ([]database.Playlist, error) {
	return s.db.GetPlaylistsByUser(ctx, userId)
}

type GetPlaylistByIdParams struct {
	PlaylistId string
	UserId     string
}

func (s *PlaylistService) GetPlaylistById(
	ctx context.Context,
	params GetPlaylistByIdParams,
) (database.Playlist, error) {
	playlist, err := s.db.GetPlaylistById(ctx, params.PlaylistId)
	if err != nil {
		if errors.Is(database.ErrItemNotFound, err) {
			return database.Playlist{}, ErrPlaylistServicePlaylistNotFound
		}

		return database.Playlist{}, fmt.Errorf("playlist-service: get playlist by id: %w", err)
	}

	if playlist.OwnerId != params.UserId {
		return database.Playlist{}, ErrPlaylistServicePlaylistNotFound
	}

	return playlist, nil
}

type CreatePlaylistParams struct {
	Name    string
	OwnerId string
}

func (s *PlaylistService) CreatePlaylist(
	ctx context.Context,
	params CreatePlaylistParams,
) (string, error) {
	playlistId, err := s.db.CreatePlaylist(ctx, database.CreatePlaylistParams{
		Name:    params.Name,
		OwnerId: params.OwnerId,
	})
	if err != nil {
		return "", fmt.Errorf("playlist-service: create playlist: %w", err)
	}

	return playlistId, nil
}

type EditPlaylistParams struct {
	PlaylistId string
	UserId     string

	Name     *string
	CoverUrl *string
}

func (s *PlaylistService) EditPlaylist(
	ctx context.Context,
	params EditPlaylistParams,
) error {
	dbPlaylist, err := s.db.GetPlaylistById(ctx, params.PlaylistId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrPlaylistServicePlaylistNotFound
		}

		return fmt.Errorf("playlist-service: edit playlist: get playlist by id: %w", err)
	}

	if dbPlaylist.OwnerId != params.UserId {
		return ErrPlaylistServicePlaylistNotFound
	}

	changes := database.PlaylistChanges{}

	if params.Name != nil {
		changes.Name = types.Change[string]{
			Value:   *params.Name,
			Changed: *params.Name != dbPlaylist.Name,
		}
	}

	if params.CoverUrl != nil {
		url := *params.CoverUrl

		// TODO(patrik): Cleanup, move to utils
		getImageExtFromContentType := func(contentType string) (string, error) {
			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				return "", fmt.Errorf("failed to parse content type: %w", err)
			}

			// TODO(patrik): Add support for more exts
			switch mediaType {
			case "image/png":
				return ".png", nil
			case "image/jpeg":
				return ".jpeg", nil
			default:
				return "", fmt.Errorf("unsupported media type: %s", mediaType)
			}
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		ext, err := getImageExtFromContentType(contentType)
		if err != nil {
			return err
		}

		// TODO(patrik): The tmp dir should be inside the work dir
		tmp, err := os.CreateTemp("", "tmp-image-*"+ext)
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpPath := tmp.Name()
		defer tmp.Close()

		// always clean up temp file if something goes wrong
		defer func() {
			_, err := os.Stat(tmpPath)
			if err == nil {
				os.Remove(tmpPath)
			}
		}()

		_, err = io.Copy(tmp, resp.Body)
		if err != nil {
			return err
		}

		tmp.Close()

		imageType, err := s.imageService.ValidateImage(tmpPath)
		if err != nil {
			return err
		}

		// TODO(patrik): I hate this
		playlistDir := s.dataDir.Playlist(dbPlaylist.Id)

		err = utils.CreateDirectories([]string{
			playlistDir,
		})
		if err != nil {
			return err
		}

		imageExt, ok := imageType.ToExt()
		if !ok {
			return errors.New("invalid image type")
		}

		cover := "downloaded" + imageExt
		output := filepath.Join(playlistDir, cover)
		err = os.Rename(tmpPath, output)
		if err != nil {
			return fmt.Errorf("failed to promote temp file: %w", err)
		}

		changes.CoverArt = types.Change[sql.NullString]{
			Value: sql.NullString{
				String: cover,
				Valid:  cover != "",
			},
			Changed: cover != dbPlaylist.CoverArt.String,
		}
	}

	err = s.db.UpdatePlaylist(ctx, dbPlaylist.Id, changes)
	if err != nil {
		return fmt.Errorf("playlist-service: edit playlist: update playlist: %w", err)
	}

	return nil
}
