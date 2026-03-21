package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/parser"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

var (
	ErrPlaylistServicePlaylistNotFound = errors.New("playlist service: playlist not found")
)

type PlaylistService struct {
	logger *slog.Logger

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
		logger:       logger,
		db:           db,
		dataDir:      dataDir,
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

type DeletePlaylistParams struct {
	PlaylistId string
	UserId     string
}

func (s *PlaylistService) DeletePlaylist(
	ctx context.Context,
	params DeletePlaylistParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	err = s.db.DeletePlaylist(ctx, playlist.Id)
	if err != nil {
		return err
	}

	err = os.RemoveAll(s.dataDir.Playlist(playlist.Id))
	if err != nil {
		return err
	}

	err = os.RemoveAll(s.dataDir.Cache().Playlist(playlist.Id))
	if err != nil {
		return err
	}

	return nil
}

type UploadPlaylistImageParams struct {
	PlaylistId string
	UserId     string

	File *multipart.FileHeader
}

func (s *PlaylistService) UploadPlaylistImage(
	ctx context.Context,
	params UploadPlaylistImageParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	ext := path.Ext(params.File.Filename)

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

	srcImage, err := params.File.Open()
	if err != nil {
		return err
	}

	_, err = io.Copy(tmp, srcImage)
	if err != nil {
		return err
	}

	tmp.Close()

	imageType, err := s.imageService.ValidateImage(tmpPath)
	if err != nil {
		return err
	}

	playlistDir := s.dataDir.Playlist(playlist.Id)

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

	coverArt := "uploaded" + imageExt
	output := path.Join(playlistDir, coverArt)
	err = os.Rename(tmpPath, output)
	if err != nil {
		return fmt.Errorf("failed to promote temp file: %w", err)
	}

	err = s.db.UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
		CoverArt: types.Change[sql.NullString]{
			Value: sql.NullString{
				String: coverArt,
				Valid:  coverArt != "",
			},
			Changed: coverArt != playlist.CoverArt.String,
		},
	})
	if err != nil {
		return err
	}

	err = os.RemoveAll(s.dataDir.Cache().Playlist(playlist.Id))
	if err != nil {
		return err
	}

	return nil
}

type GeneratePlaylistImageParams struct {
	PlaylistId string
	UserId     string
}

func (s *PlaylistService) GeneratePlaylistImage(
	ctx context.Context,
	params GeneratePlaylistImageParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	images, err := s.db.GetPlaylistTrackImages(ctx, playlist.Id, 4)
	if err != nil {
		return err
	}

	imgs := [4]string{}

	for i, img := range images {
		if !img.Valid {
			continue
		}

		imgs[i] = img.String
	}

	playlistDir := s.dataDir.Playlist(playlist.Id)

	err = utils.CreateDirectories([]string{
		playlistDir,
	})
	if err != nil {
		return err
	}

	coverArt := "generated.png"
	out := path.Join(playlistDir, coverArt)
	err = utils.GeneratePlaylistCover(imgs, out, 512)
	if err != nil {
		return err
	}

	err = s.db.UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
		CoverArt: types.Change[sql.NullString]{
			Value: sql.NullString{
				String: coverArt,
				Valid:  coverArt != "",
			},
			Changed: coverArt != playlist.CoverArt.String,
		},
	})
	if err != nil {
		return err
	}

	err = os.RemoveAll(s.dataDir.Cache().Playlist(playlist.Id))
	if err != nil {
		return err
	}

	return nil
}

type GetPlaylistItemsParams struct {
	PlaylistId string
	UserId     string

	Page   types.PageParams
	Filter types.FilterParams

	FilterId string
}

func (s *PlaylistService) GetPlaylistItems(
	ctx context.Context,
	params GetPlaylistItemsParams,
) ([]database.OrderedTrack, types.Page, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	if params.FilterId != "" {
		// TODO(patrik): Maybe log the error, or check for NotFound
		filter, err := s.db.GetPlaylistFilterById(ctx, params.FilterId, playlist.Id)
		if err == nil {
			params.Filter.Filter = filter.Filter
		}
	}

	tracks, page, err := s.db.GetPlaylistTracks(ctx, database.GetPlaylistTracksParams{
		PlaylistId: playlist.Id,
		Page:       params.Page,
		Filter:     params.Filter,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	for i, track := range tracks {
		tracks[i].Track.Order = utils.Pointer(track.Order + 1)
	}

	return tracks, page, err
}

type AddItemToPlaylistParams struct {
	PlaylistId string
	UserId     string

	TrackId string
}

func (s *PlaylistService) AddItemToPlaylist(
	ctx context.Context,
	params AddItemToPlaylistParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	// TODO(patrik): Replace with a simpler track query, we only need to
	// know if the track exists and not all the data it has
	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		// TODO(patrik): Handle
		// if errors.Is(err, database.ErrItemNotFound) {
		// 	return nil, TrackNotFound()
		// }

		return err
	}

	index, err := s.db.GetNextPlaylistItemIndex(ctx, playlist.Id)
	if err != nil {
		return err
	}

	err = s.db.CreatePlaylistItem(ctx, database.CreatePlaylistItemParams{
		PlaylistId: playlist.Id,
		TrackId:    track.Id,
		Order:      index,
	})
	if err != nil {
		// TODO(patrik): Handle
		// if errors.Is(err, database.ErrItemAlreadyExists) {
		// 	return nil, PlaylistAlreadyHasTrack()
		// }

		return err
	}

	return nil
}

type RemovePlaylistItemParams struct {
	PlaylistId string
	UserId     string

	TrackId string
}

func (s *PlaylistService) RemovePlaylistItem(
	ctx context.Context,
	params RemovePlaylistItemParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	// TODO(patrik): Check for trackId exists?
	err = s.db.DeletePlaylistItem(ctx, playlist.Id, params.TrackId)
	if err != nil {
		return err
	}

	return nil
}

type ReorderPlaylistItemsParams struct {
	PlaylistId string
	UserId     string

	Before        bool
	AnchorTrackId string
	TrackIds      []string
}

func (s *PlaylistService) ReorderPlaylistItems(
	ctx context.Context,
	params ReorderPlaylistItemsParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	current, err := tx.GetPlaylistItems(ctx, playlist.Id)
	if err != nil {
		return err
	}

	// Index current items by ID for O(1) lookup.
	index := make(map[string]database.PlaylistItem, len(current))
	for _, item := range current {
		index[item.TrackId] = item
	}

	// Validate that all supplied trackIDs exist in the playlist and resolve them to PlaylistItems.
	items := make([]database.PlaylistItem, 0, len(params.TrackIds))
	for _, id := range params.TrackIds {
		item, ok := index[id]
		if !ok {
			// TODO(patrik): Handle error
			return fmt.Errorf("track %q not found in playlist %q", id, playlist.Id)
		}

		items = append(items, item)
	}

	// Validate that anchorTrackID exists in the playlist (unless it's empty).
	if params.AnchorTrackId != "" {
		if _, ok := index[params.AnchorTrackId]; !ok {
			// TODO(patrik): Handle error
			return fmt.Errorf("anchor track %q not found in playlist %q", params.AnchorTrackId, playlist.Id)
		}
	}

	// Build a set of IDs to move for O(1) lookup.
	moveSet := make(map[string]bool, len(items))
	for _, item := range items {
		moveSet[item.TrackId] = true
	}

	// Collect all items that are NOT being moved, preserving their order.
	stationary := make([]database.PlaylistItem, 0, len(current))
	for _, item := range current {
		if !moveSet[item.TrackId] {
			stationary = append(stationary, item)
		}
	}

	// Find the insertion index within the stationary slice.
	// Defaults to 0 so that an empty anchorTrackID prepends the moved items.
	insertAt := 0
	if params.AnchorTrackId != "" {
		for i, item := range stationary {
			if item.TrackId == params.AnchorTrackId {
				insertAt = i + 1
				break
			}
		}
	}

	// Splice: stationary[:insertAt] + items + stationary[insertAt:]
	spliced := make([]database.PlaylistItem, 0, len(current))
	spliced = append(spliced, stationary[:insertAt]...)
	spliced = append(spliced, items...)
	spliced = append(spliced, stationary[insertAt:]...)

	for i, item := range spliced {
		err := tx.UpdatePlaylistItem(
			ctx,
			item.PlaylistId,
			item.TrackId,
			database.PlaylistItemChanges{
				Order: types.Change[int]{
					Value:   i,
					Changed: i != item.Order,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type GetPlaylistFiltersParams struct {
	PlaylistId string
	UserId     string
}

func (s *PlaylistService) GetPlaylistFilters(
	ctx context.Context,
	params GetPlaylistFiltersParams,
) ([]database.PlaylistFilter, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return nil, err
	}

	filters, err := s.db.GetPlaylistFiltersByPlaylistId(ctx, playlist.Id)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

// TODO(patrik): Handle errors
// TODO(patrik): Move this, this needs to be used in the validation
func (s *PlaylistService) testFilter(filterStr string) error {
	ast, err := parser.ParseExpr(filterStr)
	if err != nil {
		return err
		// return InvalidFilter(err)
	}

	a := adapter.TrackResolverAdapter{}
	r := filter.New(&a)
	_, err = r.Resolve(ast)
	if err != nil {
		return err
		// return InvalidFilter(err)
	}

	return nil
}

type CreatePlaylistFilterParams struct {
	PlaylistId string
	UserId     string

	Name   string
	Filter string
}

func (s *PlaylistService) CreatePlaylistFilter(
	ctx context.Context,
	params CreatePlaylistFilterParams,
) (string, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return "", err
	}

	err = s.testFilter(params.Filter)
	if err != nil {
		// TODO(patrik): Better error
		return "", err
	}

	filterId, err := s.db.CreatePlaylistFilter(
		ctx,
		database.CreatePlaylistFilterParams{
			PlaylistId: playlist.Id,
			Name:       params.Name,
			Filter:     params.Filter,
		},
	)
	if err != nil {
		return "", err
	}

	return filterId, nil
}

type EditPlaylistFilterParams struct {
	PlaylistId string
	UserId     string
	FilterId   string

	Name   *string
	Filter *string
}

func (s *PlaylistService) EditPlaylistFilter(
	ctx context.Context,
	params EditPlaylistFilterParams,
) error {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
		UserId:     params.UserId,
	})
	if err != nil {
		return err
	}

	filter, err := s.db.GetPlaylistFilterById(ctx, params.FilterId, playlist.Id)
	if err != nil {
		// TODO(patrik): Handle error
		// if errors.Is(err, database.ErrItemNotFound) {
		// 	return nil, PlaylistFilter()
		// }

		return err
	}

	changes := database.PlaylistFilterChanges{}

	if params.Name != nil {
		changes.Name = types.Change[string]{
			Value:   *params.Name,
			Changed: *params.Name != filter.Name,
		}
	}

	if params.Filter != nil {
		// TODO(patrik): Test filter?
		// s.testFilter(*params.Filter)

		changes.Filter = types.Change[string]{
			Value:   *params.Filter,
			Changed: *params.Filter != filter.Filter,
		}
	}

	err = s.db.UpdatePlaylistFilter(ctx, filter.Id, playlist.Id, changes)
	if err != nil {
		return err
	}

	return nil
}
