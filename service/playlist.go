package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
)

var playlistErr = NewServiceErrCreator("playlist")

var (
	ErrPlaylistServicePlaylistNotFound    = playlistErr.New("playlist not found")
	ErrPlaylistServiceTrackNotFound       = playlistErr.New("track not found")
	ErrPlaylistServiceTrackAlreadyAdded   = playlistErr.New("track already added")
	ErrPlaylistServiceItemNotFound        = playlistErr.New("item not found")
	ErrPlaylistServiceFilterNotFound      = playlistErr.New("filter not found")
	ErrPlaylistServiceAnchorTrackNotFound = playlistErr.New("anchor track not found")
	ErrPlaylistServiceNotAuthorized       = playlistErr.New("not authorized")
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

type GetPlaylistsParams struct {
	Page   types.PageParams
	Filter types.FilterParams
}

func (s *PlaylistService) GetPlaylists(
	ctx context.Context,
	params GetPlaylistsParams,
) ([]database.Playlist, types.Page, error) {
	playlists, page, err := s.db.GetPlaylists(ctx, database.GetPlaylistsParams{
		Page:   params.Page,
		Filter: params.Filter,
	})
	if err != nil {
		return nil, types.Page{}, playlistErr.Wrap("get playlists: db get", err)
	}

	return playlists, page, nil
}

func (s *PlaylistService) checkOwnership(playlist database.Playlist, userId string) error {
	if playlist.OwnerId != userId {
		return ErrPlaylistServiceNotAuthorized
	}

	return nil
}

type GetPlaylistByIdParams struct {
	PlaylistId string
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

		return database.Playlist{}, playlistErr.Wrap("get playlist by id", err)
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
		return "", playlistErr.Wrap("create: db create", err)
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
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
	})
	if err != nil {
		return playlistErr.Wrap("edit: get playlist", err)
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	changes := database.PlaylistChanges{}

	if params.Name != nil {
		changes.Name = types.Change[string]{
			Value:   *params.Name,
			Changed: *params.Name != playlist.Name,
		}
	}

	if params.CoverUrl != nil {
		url := *params.CoverUrl

		cover, err := s.imageService.DownloadCoverForPlaylist(
			ctx,
			DownloadCoverForPlaylistParams{
				PlaylistId: playlist.Id,
				Url:        url,
			},
		)
		if err != nil {
			return err
		}

		changes.CoverArt = types.Change[sql.NullString]{
			Value: sql.NullString{
				String: cover,
				Valid:  cover != "",
			},
			Changed: cover != playlist.CoverArt.String,
		}
	}

	err = s.db.UpdatePlaylist(ctx, playlist.Id, changes)
	if err != nil {
		return playlistErr.Wrap("edit: db update", err)
	}

	// TODO(patrik): Make this better
	err = os.RemoveAll(s.dataDir.CacheImages().Playlist(playlist.Id))
	if err != nil {
		return playlistErr.Wrap("edit: remove cache", err)
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	err = s.db.DeletePlaylist(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("delete: db delete", err)
	}

	err = os.RemoveAll(s.dataDir.Playlist(playlist.Id))
	if err != nil {
		return playlistErr.Wrap("delete: remove dir", err)
	}

	// TODO(patrik): Make this better
	err = os.RemoveAll(s.dataDir.CacheImages().Playlist(playlist.Id))
	if err != nil {
		return playlistErr.Wrap("delete: remove cache", err)
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	cover, err := s.imageService.UploadImageForPlaylist(
		ctx,
		UploadImageForPlaylistParams{
			PlaylistId: playlist.Id,
			File:       params.File,
		},
	)
	if err != nil {
		return playlistErr.Wrap("upload image: upload", err)
	}

	err = s.db.UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
		CoverArt: types.Change[sql.NullString]{
			Value: sql.NullString{
				String: cover,
				Valid:  cover != "",
			},
			Changed: cover != playlist.CoverArt.String,
		},
	})
	if err != nil {
		return playlistErr.Wrap("upload image: db update", err)
	}

	// TODO(patrik): Make this better
	err = os.RemoveAll(s.dataDir.CacheImages().Playlist(playlist.Id))
	if err != nil {
		return playlistErr.Wrap("upload image: remove cache", err)
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	cover, err := s.imageService.GenerateImageForPlaylist(
		ctx,
		GenerateImageForPlaylistParams{
			PlaylistId: playlist.Id,
		},
	)
	if err != nil {
		return playlistErr.Wrap("gen image: image gen", err)
	}

	err = s.db.UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
		CoverArt: types.Change[sql.NullString]{
			Value: sql.NullString{
				String: cover,
				Valid:  cover != "",
			},
			Changed: cover != playlist.CoverArt.String,
		},
	})
	if err != nil {
		return playlistErr.Wrap("gen image: db update", err)
	}

	// TODO(patrik): Make this better
	err = os.RemoveAll(s.dataDir.CacheImages().Playlist(playlist.Id))
	if err != nil {
		return playlistErr.Wrap("gen image: remove cache", err)
	}

	return nil
}

type GetPlaylistItemsParams struct {
	PlaylistId string

	Page   types.PageParams
	Filter types.FilterParams

	FilterId string
}

func (s *PlaylistService) GetPlaylistItems(
	ctx context.Context,
	params GetPlaylistItemsParams,
) ([]database.PlaylistItemTrack, types.Page, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
	})
	if err != nil {
		return nil, types.Page{}, err
	}

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, types.Page{}, ErrPlaylistServiceFilterNotFound
			}

			return nil, types.Page{}, playlistErr.Wrap("get items: db get filter", err)
		}

		params.Filter.Filter = filter.Filter
	}

	tracks, page, err := s.db.GetPlaylistTracks(
		ctx,
		database.GetPlaylistTracksParams{
			PlaylistId: playlist.Id,
			Page:       params.Page,
			Filter:     params.Filter,
		},
	)
	if err != nil {
		return nil, types.Page{}, playlistErr.Wrap("get items: db get", err)
	}

	for i, track := range tracks {
		tracks[i].Track.Order = utils.Pointer(track.Position + 1)
	}

	return tracks, page, nil
}

type GetPlaylistItemIdsParams struct {
	PlaylistId string
}

func (s *PlaylistService) GetPlaylistItemIds(
	ctx context.Context,
	params GetPlaylistItemIdsParams,
) ([]string, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
	})
	if err != nil {
		return nil, err
	}

	ids, err := s.db.GetPlaylistItemIds(
		ctx,
		database.GetPlaylistItemIdsParams{
			PlaylistId: playlist.Id,
		},
	)
	if err != nil {
		return nil, playlistErr.Wrap("get item ids: db get", err)
	}

	return ids, nil
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	// TODO(patrik): Replace with a simpler track query, we only need to
	// know if the track exists and not all the data it has
	track, err := s.db.GetTrackById(ctx, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrPlaylistServiceTrackNotFound
		}

		return err
	}

	index, err := s.db.GetNextPlaylistItemIndex(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("add item: db get next index", err)
	}

	err = s.db.CreatePlaylistItem(ctx, database.CreatePlaylistItemParams{
		PlaylistId: playlist.Id,
		TrackId:    track.Id,
		Position:   index,
	})
	if err != nil {
		if errors.Is(err, database.ErrItemAlreadyExists) {
			return ErrPlaylistServiceTrackAlreadyAdded
		}

		return playlistErr.Wrap("add item: db create item", err)
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	item, err := s.db.GetPlaylistItemByTrackId(ctx, playlist.Id, params.TrackId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ErrPlaylistServiceItemNotFound
		}

		return playlistErr.Wrap("remove item: db get item", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return playlistErr.Wrap("remove item: db begin", err)
	}
	defer tx.Rollback()

	err = tx.DeletePlaylistItem(ctx, playlist.Id, params.TrackId)
	if err != nil {
		return playlistErr.Wrap("remove item: db delete item", err)
	}

	err = tx.ReorderPlaylistItemsAfterDelete(ctx, playlist.Id, item.Position)
	if err != nil {
		return playlistErr.Wrap("remove item: db reorder items", err)
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("remove item: db commit", err)
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
	})
	if err != nil {
		return err
	}

	err = s.checkOwnership(playlist, params.UserId)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return playlistErr.Wrap("reorder items: db begin", err)
	}
	defer tx.Rollback()

	current, err := tx.GetPlaylistItems(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("reorder items: db get items", err)
	}

	index := make(map[string]database.PlaylistItem, len(current))
	for _, item := range current {
		index[item.TrackId] = item
	}

	items := make([]database.PlaylistItem, 0, len(params.TrackIds))
	for _, id := range params.TrackIds {
		item, ok := index[id]
		if !ok {
			continue
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return nil
	}

	if params.AnchorTrackId != "" {
		if _, ok := index[params.AnchorTrackId]; !ok {
			return ErrPlaylistServiceAnchorTrackNotFound
		}
	}

	moveSet := make(map[string]bool, len(items))
	for _, item := range items {
		moveSet[item.TrackId] = true
	}

	stationary := make([]database.PlaylistItem, 0, len(current))
	for _, item := range current {
		if !moveSet[item.TrackId] {
			stationary = append(stationary, item)
		}
	}

	insertAt := 0
	if params.AnchorTrackId != "" {
		for i, item := range stationary {
			if item.TrackId == params.AnchorTrackId {
				insertAt = i + 1
				break
			}
		}
	}

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
				Position: types.Change[int]{
					Value:   i,
					Changed: i != item.Position,
				},
			},
		)
		if err != nil {
			return playlistErr.Wrap("reorder items: db update item", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("reorder items: db commit", err)
	}

	return nil
}
