package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var queueErr = NewServiceErrCreator("queue")

var (
	ErrQueueServiceQueueNotFound = queueErr.New("queue not found")
	ErrQueueServiceItemNotFound  = queueErr.New("item not found")
)

type QueueService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewQueueService(
	logger *slog.Logger,
	db *database.Database,
) *QueueService {
	return &QueueService{
		logger: logger,
		db:     db,
	}
}

type GetQueueParams struct {
	Page types.PageParams

	QueueId string
	UserId  string
}

type GetQueueResult struct {
	Items        []database.QueueItemTrack
	CurrentIndex int
	Page         types.Page
}

func (s *QueueService) GetQueue(
	ctx context.Context,
	params GetQueueParams,
) (GetQueueResult, error) {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return GetQueueResult{}, err
	}

	items, page, err := s.db.GetQueueItems(ctx, database.GetQueueItemsParams{
		Page:    params.Page,
		QueueId: queue.Id,
		UserId:  queue.UserId,
	})
	if err != nil {
		return GetQueueResult{}, queueErr.Wrap("get queue", err)
	}

	return GetQueueResult{
		Items:        items,
		CurrentIndex: queue.CurrentIndex,
		Page:         page,
	}, nil
}

type QueueIdsResult struct {
	Entries      []database.QueueItemEntry
	CurrentIndex int
}

func (s *QueueService) GetQueueIds(
	ctx context.Context,
	queueId string,
	userId string,
) (QueueIdsResult, error) {
	queue, err := s.getOrCreateQueue(ctx, queueId, userId)
	if err != nil {
		return QueueIdsResult{}, err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id, queue.UserId)
	if err != nil {
		return QueueIdsResult{}, queueErr.Wrap("get queue ids", err)
	}

	return QueueIdsResult{
		Entries:      entries,
		CurrentIndex: queue.CurrentIndex,
	}, nil
}

type GetQueueItemAtIndexParams struct {
	QueueId string
	UserId  string
	Index   int
}

func (s *QueueService) GetQueueItemAtIndex(
	ctx context.Context,
	params GetQueueItemAtIndexParams,
) (database.QueueItemTrack, error) {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return database.QueueItemTrack{}, err
	}

	item, err := s.db.GetQueueItemAtPosition(ctx, queue.Id, queue.UserId, params.Index)
	if err != nil {
		return database.QueueItemTrack{}, queueErr.Wrap(
			"get queue item at index", err)
	}

	return item, nil
}

func (s *QueueService) getOrCreateQueue(
	ctx context.Context,
	queueId string,
	userId string,
) (database.Queue, error) {
	queue, err := s.db.GetQueueById(ctx, queueId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			_, err = s.db.CreateQueue(ctx, database.CreateQueueParams{
				Id:           queueId,
				UserId:       userId,
				CurrentIndex: 0,
			})
			if err != nil {
				return database.Queue{}, queueErr.Wrap("create queue", err)
			}

			queue, err = s.db.GetQueueById(ctx, queueId)
			if err != nil {
				return database.Queue{}, queueErr.Wrap("get queue by id", err)
			}

			return queue, nil
		}

		return database.Queue{}, queueErr.Wrap("get or create queue", err)
	}

	return queue, nil
}

type ReplaceQueueParams struct {
	QueueId      string
	UserId       string
	TrackIds     []string
	CurrentIndex int
	Shuffle      bool
}

func (s *QueueService) ReplaceQueue(
	ctx context.Context,
	params ReplaceQueueParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return err
	}

	trackIds := params.TrackIds
	if params.Shuffle {
		shuffled := make([]string, len(trackIds))
		copy(shuffled, trackIds)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		trackIds = shuffled
	}

	tx, err := s.db.Begin()
	if err != nil {
		return queueErr.Wrap("replace queue: begin", err)
	}
	defer tx.Rollback()

	err = tx.ClearQueueItems(ctx, queue.Id, queue.UserId)
	if err != nil {
		return queueErr.Wrap("replace queue: clear items", err)
	}

	items := make([]database.CreateQueueItemParams, len(trackIds))
	for i, trackId := range trackIds {
		items[i] = database.CreateQueueItemParams{
			TrackId:  trackId,
			Position: i,
		}
	}

	err = tx.CreateQueueItems(ctx, database.CreateQueueItemsParams{
		QueueId: queue.Id,
		UserId:  queue.UserId,
		Items:   items,
	})
	if err != nil {
		return queueErr.Wrap("replace queue: create items", err)
	}

	currentIndex := params.CurrentIndex
	if currentIndex < 0 || currentIndex >= len(trackIds) {
		currentIndex = 0
	}

	err = tx.UpdateQueue(ctx, queue.Id, database.QueueChanges{
		CurrentIndex: database.Change[int]{
			Value:   currentIndex,
			Changed: true,
		},
	})
	if err != nil {
		return queueErr.Wrap("replace queue: update", err)
	}

	err = tx.Commit()
	if err != nil {
		return queueErr.Wrap("replace queue: commit", err)
	}

	return nil
}

type AddItemsParams struct {
	QueueId  string
	UserId   string
	TrackIds []string
	Position string // "next" or "end"
}

func (s *QueueService) AddItems(
	ctx context.Context,
	params AddItemsParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id, queue.UserId)
	if err != nil {
		return err
	}

	existingIds := make([]string, len(entries))
	for i, entry := range entries {
		existingIds[i] = entry.TrackId
	}

	var newIds []string
	if params.Position == "next" {
		insertIndex := min(queue.CurrentIndex+1, len(existingIds))

		newIds = make([]string, 0, len(existingIds)+len(params.TrackIds))
		newIds = append(newIds, existingIds[:insertIndex]...)
		newIds = append(newIds, params.TrackIds...)
		newIds = append(newIds, existingIds[insertIndex:]...)
	} else {
		newIds = append(existingIds, params.TrackIds...)
	}

	err = s.ReplaceQueue(ctx, ReplaceQueueParams{
		QueueId:      params.QueueId,
		UserId:       params.UserId,
		TrackIds:     newIds,
		CurrentIndex: queue.CurrentIndex,
	})
	if err != nil {
		return err
	}

	return nil
}

type AddToQueueParams struct {
	QueueId             string
	UserId              string
	Source              string // "album", "playlist", "artist", "tracks"
	SourceId            string
	TrackIds            []string
	Position            string // "replace", "next", "end"
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

type AddAlbumToQueueParams struct {
	QueueId             string
	UserId              string
	AlbumId             string
	FilterId            string
	Position            string // "replace", "next", "end"
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

func (s *QueueService) AddToQueue(
	ctx context.Context,
	params AddToQueueParams,
) error {
	var trackIds []string
	var err error

	switch params.Source {
	case "album":
		trackIds, err = s.db.GetTrackIdsByAlbum(ctx, params.SourceId, "")
	case "playlist":
		trackIds, err = s.db.GetPlaylistItemIds(
			ctx,
			database.GetPlaylistItemIdsParams{
				PlaylistId: params.SourceId,
			},
		)
	case "artist":
		trackIds, err = s.db.GetTrackIdsByArtist(ctx, params.SourceId, "")
	case "tracks":
		trackIds = params.TrackIds
	case "filter":
		filter, err := s.db.GetTrackFilterById(ctx, params.SourceId)
		if err != nil {
			return queueErr.Wrap("get track filter", err)
		}

		trackIds, err = s.db.GetTrackIdsByFilter(ctx, filter.Filter)
		if err != nil {
			return queueErr.Wrap("get track ids by filter", err)
		}
	default:
		return queueErr.New("unknown source: " + params.Source)
	}

	if err != nil {
		return queueErr.Wrap("fetch track ids", err)
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

func (s *QueueService) AddAlbumToQueue(
	ctx context.Context,
	params AddAlbumToQueueParams,
) error {
	filterStr := ""

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			// TODO(patrik): Handle not found error
			return queueErr.Wrap("get track filter", err)
		}

		filterStr = filter.Filter
	}

	trackIds, err := s.db.GetTrackIdsByAlbum(ctx, params.AlbumId, filterStr)
	if err != nil {
		return queueErr.Wrap("get track ids by album", err)
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

type AddArtistToQueueParams struct {
	QueueId             string
	UserId              string
	ArtistId            string
	FilterId            string
	Position            string
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

func (s *QueueService) AddArtistToQueue(
	ctx context.Context,
	params AddArtistToQueueParams,
) error {
	filterStr := ""

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			// TODO(patrik): Handle not found error
			return queueErr.Wrap("get track filter", err)
		}

		filterStr = filter.Filter
	}

	trackIds, err := s.db.GetTrackIdsByArtist(ctx, params.ArtistId, filterStr)
	if err != nil {
		return queueErr.Wrap("get track ids by artist", err)
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

type AddPlaylistToQueueParams struct {
	QueueId             string
	UserId              string
	PlaylistId          string
	FilterId            string
	Position            string
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

func (s *QueueService) AddPlaylistToQueue(
	ctx context.Context,
	params AddPlaylistToQueueParams,
) error {
	filterStr := ""

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			// TODO(patrik): Handle not found error
			return queueErr.Wrap("get track filter", err)
		}

		filterStr = filter.Filter
	}

	trackIds, err := s.db.GetTrackIdsByPlaylist(
		ctx,
		database.GetTrackIdsByPlaylistParams{
			PlaylistId: params.PlaylistId,
			FilterStr:  filterStr,
		},
	)
	if err != nil {
		return queueErr.Wrap("get track ids by playlist", err)
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

type AddFavoritesToQueueParams struct {
	QueueId             string
	UserId              string
	FavoriteUserId      string
	FilterId            string
	Position            string
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

func (s *QueueService) AddFavoritesToQueue(
	ctx context.Context,
	params AddFavoritesToQueueParams,
) error {
	filterStr := ""

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			return queueErr.Wrap("get track filter", err)
		}

		filterStr = filter.Filter
	}

	trackIds, err := s.db.GetTrackIdsByUserFavorites(
		ctx,
		database.GetTrackIdsByUserFavoritesParams{
			UserId:    params.FavoriteUserId,
			FilterStr: filterStr,
		},
	)
	if err != nil {
		return queueErr.Wrap("get track ids by favorites", err)
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

type AddTracksToQueueParams struct {
	QueueId             string
	UserId              string
	TrackIds            []string
	FilterId            string
	Position            string
	Shuffle             bool
	CurrentIndex        int
	QueueIndexToTrackId string
}

func (s *QueueService) AddTracksToQueue(
	ctx context.Context,
	params AddTracksToQueueParams,
) error {
	trackIds := params.TrackIds

	if params.FilterId != "" {
		filter, err := s.db.GetTrackFilterById(ctx, params.FilterId)
		if err != nil {
			return queueErr.Wrap("get track filter", err)
		}

		trackIds, err = s.db.GetTrackIdsByFilter(ctx, filter.Filter)
		if err != nil {
			return queueErr.Wrap("get track ids by filter", err)
		}
	}

	if len(trackIds) == 0 {
		return nil
	}

	switch params.Position {
	case "next", "end":
		return s.AddItems(ctx, AddItemsParams{
			QueueId:  params.QueueId,
			UserId:   params.UserId,
			TrackIds: trackIds,
			Position: params.Position,
		})
	case "replace":
		currentIndex := params.CurrentIndex
		if params.QueueIndexToTrackId != "" {
			for i, id := range trackIds {
				if id == params.QueueIndexToTrackId {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex < 0 || currentIndex >= len(trackIds) {
			currentIndex = 0
		}

		return s.ReplaceQueue(ctx, ReplaceQueueParams{
			QueueId:      params.QueueId,
			UserId:       params.UserId,
			TrackIds:     trackIds,
			CurrentIndex: currentIndex,
			Shuffle:      params.Shuffle,
		})
	default:
		return queueErr.New("unknown position: " + params.Position)
	}
}

type RemoveItemParams struct {
	QueueId string
	UserId  string
	ItemId  string
}

func (s *QueueService) RemoveItem(
	ctx context.Context,
	params RemoveItemParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id, queue.UserId)
	if err != nil {
		return err
	}

	newIds := make([]string, 0, len(entries))
	newIndex := queue.CurrentIndex
	found := false

	for i, entry := range entries {
		if entry.QueueItemId == params.ItemId {
			found = true
			if i < queue.CurrentIndex {
				newIndex--
			}
			continue
		}
		newIds = append(newIds, entry.TrackId)
	}

	if !found {
		return ErrQueueServiceItemNotFound
	}

	if newIndex < 0 {
		newIndex = 0
	}

	err = s.ReplaceQueue(ctx, ReplaceQueueParams{
		QueueId:      params.QueueId,
		UserId:       params.UserId,
		TrackIds:     newIds,
		CurrentIndex: newIndex,
	})
	if err != nil {
		return err
	}

	return nil
}

type SetPositionParams struct {
	QueueId string
	UserId  string
	Index   int
}

func (s *QueueService) SetPosition(
	ctx context.Context,
	params SetPositionParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.QueueId, params.UserId)
	if err != nil {
		return err
	}

	err = s.db.UpdateQueue(ctx, queue.Id, database.QueueChanges{
		CurrentIndex: database.Change[int]{
			Value:   params.Index,
			Changed: true,
		},
	})
	if err != nil {
		return queueErr.Wrap("set position", err)
	}

	return nil
}

// TODO(patrik): Refactor: Use ClearQueueParams
func (s *QueueService) ClearQueue(
	ctx context.Context, 
	queueId string, 
	userId string,
) error {
	queue, err := s.getOrCreateQueue(ctx, queueId, userId)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return queueErr.Wrap("clear queue: begin", err)
	}
	defer tx.Rollback()

	err = tx.ClearQueueItems(ctx, queue.Id, queue.UserId)
	if err != nil {
		return queueErr.Wrap("clear queue: clear items", err)
	}

	err = tx.UpdateQueue(ctx, queue.Id, database.QueueChanges{
		CurrentIndex: database.Change[int]{
			Value:   0,
			Changed: true,
		},
	})
	if err != nil {
		return queueErr.Wrap("clear queue: update", err)
	}

	err = tx.Commit()
	if err != nil {
		return queueErr.Wrap("clear queue: commit", err)
	}

	return nil
}
