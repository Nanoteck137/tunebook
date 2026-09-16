package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"path"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var playlistErr = NewServiceErrCreator("playlist")

var (
	ErrPlaylistServicePlaylistNotFound       = playlistErr.New("playlist not found")
	ErrPlaylistServiceTrackNotFound          = playlistErr.New("track not found")
	ErrPlaylistServiceTrackAlreadyAdded      = playlistErr.New("track already added")
	ErrPlaylistServiceItemNotFound           = playlistErr.New("item not found")
	ErrPlaylistServiceFilterNotFound         = playlistErr.New("filter not found")
	ErrPlaylistServiceAnchorTrackNotFound    = playlistErr.New("anchor track not found")
	ErrPlaylistServiceAnchorPlaylistNotFound = playlistErr.New("anchor playlist not found")
	ErrPlaylistServiceNotAuthorized          = playlistErr.New("not authorized")
)

type PlaylistService struct {
	logger *slog.Logger

	db *database.Database

	filesystem   *FilesystemService
	imageService *ImageService

	emitter broker.EventEmitter
}

func NewPlaylistService(
	logger *slog.Logger,
	db *database.Database,
	filesystem *FilesystemService,
	imageService *ImageService,
	emitter broker.EventEmitter,
) *PlaylistService {
	return &PlaylistService{
		logger:       logger,
		db:           db,
		filesystem:   filesystem,
		imageService: imageService,
		emitter:      emitter,
	}
}

type GetPlaylistsParams struct {
	Page  types.PageParams
	Query types.QueryParams
}

func (s *PlaylistService) GetPlaylists(
	ctx context.Context,
	params GetPlaylistsParams,
) ([]database.Playlist, types.Page, error) {
	playlists, page, err := s.db.GetPlaylists(ctx, database.GetPlaylistsParams{
		Page:  params.Page,
		Query: params.Query,
	})
	if err != nil {
		return nil, types.Page{}, playlistErr.Wrap(
			"get playlists: db get", err)
	}

	return playlists, page, nil
}

func (s *PlaylistService) checkOwnership(
	playlist database.Playlist,
	userId string,
) error {
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

type GetPlaylistImageParams struct {
	PlaylistId  string
	Size        int
	ImageFormat types.ImageFormat
}

func (s *PlaylistService) GetPlaylistImage(
	ctx context.Context,
	params GetPlaylistImageParams,
) (string, error) {
	playlist, err := s.GetPlaylistById(ctx, GetPlaylistByIdParams{
		PlaylistId: params.PlaylistId,
	})
	if err != nil {
		return "", err
	}

	err = s.filesystem.EnsurePlaylistImageCacheDirs(playlist.Id)
	if err != nil {
		return "", playlistErr.Wrap("get playlist image", err)
	}

	input := ""
	if playlist.CoverArt.Valid {
		playlistDir := s.filesystem.PlaylistDir(playlist.Id)
		input = path.Join(playlistDir, playlist.CoverArt.String)
	}

	p, err := s.imageService.ProcessImage(ProcessImageParams{
		Input:       input,
		Default:     "default_album.png",
		OutputDir:   s.filesystem.PlaylistImagePath(playlist.Id),
		Size:        params.Size,
		ImageFormat: params.ImageFormat,
	})
	if err != nil {
		return "", playlistErr.Wrap("get playlist image", err)
	}

	return p, nil
}

type CreatePlaylistParams struct {
	Name    string
	OwnerId string
}

func (s *PlaylistService) CreatePlaylist(
	ctx context.Context,
	params CreatePlaylistParams,
) (string, error) {
	position, err := s.db.GetNextPlaylistPosition(ctx, params.OwnerId)
	if err != nil {
		return "", playlistErr.Wrap("create: db get next position", err)
	}

	playlistId, err := s.db.CreatePlaylist(ctx, database.CreatePlaylistParams{
		Name:     params.Name,
		OwnerId:  params.OwnerId,
		Position: position,
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
		changes.Name = database.Change[string]{
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
			return playlistErr.Wrap("edit: download cover", err)
		}

		changes.CoverArt = database.Change[sql.NullString]{
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

	err = s.filesystem.ClearPlaylistImageCache(playlist.Id)
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

	tx, err := s.db.Begin()
	if err != nil {
		return playlistErr.Wrap("delete: db begin", err)
	}
	defer tx.Rollback()

	err = tx.DeletePlaylist(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("delete: db delete", err)
	}

	err = tx.ReorderPlaylistsAfterDelete(ctx, playlist.OwnerId, playlist.Position)
	if err != nil {
		return playlistErr.Wrap("delete: db reorder playlists", err)
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("delete: db commit", err)
	}

	err = s.filesystem.RemovePlaylistDir(playlist.Id)
	if err != nil {
		return playlistErr.Wrap("delete: remove dir", err)
	}

	err = s.filesystem.ClearPlaylistImageCache(playlist.Id)
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
		CoverArt: database.Change[sql.NullString]{
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

	err = s.filesystem.ClearPlaylistImageCache(playlist.Id)
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
		CoverArt: database.Change[sql.NullString]{
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

	err = s.filesystem.ClearPlaylistImageCache(playlist.Id)
	if err != nil {
		return playlistErr.Wrap("gen image: remove cache", err)
	}

	return nil
}

type GetPlaylistItemsParams struct {
	PlaylistId string

	Page  types.PageParams
	Query types.QueryParams

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

			return nil, types.Page{}, playlistErr.Wrap(
				"get items: db get filter", err)
		}

		params.Query.Filter = filter.Filter
	}

	tracks, page, err := s.db.GetPlaylistTracks(
		ctx,
		database.GetPlaylistTracksParams{
			PlaylistId: playlist.Id,
			Page:       params.Page,
			Query:      params.Query,
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

	err = s.db.MarkPlaylistUpdated(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("add item: db mark playlist updated", err)
	}

	s.emitQuickPlaylistChanged(ctx, params.UserId, playlist.Id)

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

	item, err := s.db.GetPlaylistItemByTrackId(
		ctx, playlist.Id, params.TrackId)
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

	err = tx.MarkPlaylistUpdated(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("remove item: db mark playlist updated", err)
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("remove item: db commit", err)
	}

	s.emitQuickPlaylistChanged(ctx, params.UserId, playlist.Id)

	return nil
}

func (s *PlaylistService) emitQuickPlaylistChanged(
	ctx context.Context,
	userId string,
	playlistId string,
) {
	settings, err := s.db.GetUserSettingsById(ctx, userId)
	if err != nil || !settings.QuickPlaylist.Valid {
		return
	}

	if settings.QuickPlaylist.String != playlistId {
		return
	}

	s.emitter.EmitEventToUser(
		userId,
		broker.QuickPlaylistChangedEvent{
			PlaylistId: playlistId,
		},
	)
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
				Position: database.Change[int]{
					Value:   i,
					Changed: i != item.Position,
				},
			},
		)
		if err != nil {
			return playlistErr.Wrap("reorder items: db update item", err)
		}
	}

	err = tx.MarkPlaylistUpdated(ctx, playlist.Id)
	if err != nil {
		return playlistErr.Wrap("reorder items: db mark playlist updated", err)
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("reorder items: db commit", err)
	}

	return nil
}

type ReorderPlaylistsParams struct {
	UserId string

	Before           bool
	AnchorPlaylistId string
	PlaylistIds      []string
}

func (s *PlaylistService) ReorderPlaylists(
	ctx context.Context,
	params ReorderPlaylistsParams,
) error {
	current, err := s.db.GetUserPlaylists(ctx, params.UserId)
	if err != nil {
		return playlistErr.Wrap("reorder playlists: db get playlists", err)
	}

	index := make(map[string]database.Playlist, len(current))
	for _, playlist := range current {
		index[playlist.Id] = playlist
	}

	playlists := make([]database.Playlist, 0, len(params.PlaylistIds))
	for _, id := range params.PlaylistIds {
		playlist, ok := index[id]
		if !ok {
			continue
		}

		playlists = append(playlists, playlist)
	}

	if len(playlists) == 0 {
		return nil
	}

	if params.AnchorPlaylistId != "" {
		if _, ok := index[params.AnchorPlaylistId]; !ok {
			return ErrPlaylistServiceAnchorPlaylistNotFound
		}
	}

	moveSet := make(map[string]bool, len(playlists))
	for _, playlist := range playlists {
		moveSet[playlist.Id] = true
	}

	stationary := make([]database.Playlist, 0, len(current))
	for _, playlist := range current {
		if !moveSet[playlist.Id] {
			stationary = append(stationary, playlist)
		}
	}

	insertAt := 0
	if params.AnchorPlaylistId != "" {
		for i, playlist := range stationary {
			if playlist.Id == params.AnchorPlaylistId {
				insertAt = i + 1
				break
			}
		}
	}

	spliced := make([]database.Playlist, 0, len(current))
	spliced = append(spliced, stationary[:insertAt]...)
	spliced = append(spliced, playlists...)
	spliced = append(spliced, stationary[insertAt:]...)

	tx, err := s.db.Begin()
	if err != nil {
		return playlistErr.Wrap("reorder playlists: db begin", err)
	}
	defer tx.Rollback()

	for i, playlist := range spliced {
		err := tx.UpdatePlaylist(
			ctx,
			playlist.Id,
			database.PlaylistChanges{
				Position: database.Change[int]{
					Value:   i,
					Changed: i != playlist.Position,
				},
			},
		)
		if err != nil {
			return playlistErr.Wrap("reorder playlists: db update playlist", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return playlistErr.Wrap("reorder playlists: db commit", err)
	}

	return nil
}
