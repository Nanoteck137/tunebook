package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/types"
)

var (
	ErrQueueServiceQueueNotFound = errors.New("queue service: queue not found")
	ErrQueueServiceItemNotFound  = errors.New("queue service: item not found")
)

type QueueService struct {
	logger *slog.Logger
	db     *database.Database
}

func NewQueueService(logger *slog.Logger, db *database.Database) *QueueService {
	return &QueueService{
		logger: logger,
		db:     db,
	}
}

type GetQueueParams struct {
	Page types.PageParams

	UserId string
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
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
	if err != nil {
		return GetQueueResult{}, err
	}

	items, page, err := s.db.GetQueueItems(ctx, database.GetQueueItemsParams{
		Page:    params.Page,
		QueueId: queue.Id,
	})
	if err != nil {
		return GetQueueResult{}, err
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
	userId string,
) (QueueIdsResult, error) {
	queue, err := s.getOrCreateQueue(ctx, userId)
	if err != nil {
		return QueueIdsResult{}, err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id)
	if err != nil {
		return QueueIdsResult{}, err
	}

	return QueueIdsResult{
		Entries:      entries,
		CurrentIndex: queue.CurrentIndex,
	}, nil
}

type GetQueueItemAtIndexParams struct {
	UserId string
	Index  int
}

func (s *QueueService) GetQueueItemAtIndex(
	ctx context.Context,
	params GetQueueItemAtIndexParams,
) (database.QueueItemTrack, error) {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
	if err != nil {
		return database.QueueItemTrack{}, err
	}

	item, err := s.db.GetQueueItemAtPosition(ctx, queue.Id, params.Index)
	if err != nil {
		return database.QueueItemTrack{}, err
	}

	return item, nil
}

func (s *QueueService) getOrCreateQueue(
	ctx context.Context,
	userId string,
) (database.Queue, error) {
	queue, err := s.db.GetQueueById(ctx, userId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			_, err = s.db.CreateQueue(ctx, database.CreateQueueParams{
				Id:           userId,
				CurrentIndex: 0,
			})
			if err != nil {
				return database.Queue{}, err
			}

			queue, err = s.db.GetQueueById(ctx, userId)
			if err != nil {
				return database.Queue{}, err
			}

			return queue, nil
		}

		return database.Queue{}, err
	}

	return queue, nil
}

type ReplaceQueueParams struct {
	UserId       string
	TrackIds     []string
	CurrentIndex int
	Shuffle      bool
}

func (s *QueueService) ReplaceQueue(
	ctx context.Context,
	params ReplaceQueueParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
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
		return err
	}
	defer tx.Rollback()

	err = tx.ClearQueueItems(ctx, queue.Id)
	if err != nil {
		return err
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
		Items:   items,
	})
	if err != nil {
		return err
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
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type AddItemsParams struct {
	UserId   string
	TrackIds []string
	Position string // "next" or "end"
}

func (s *QueueService) AddItems(
	ctx context.Context,
	params AddItemsParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
	if err != nil {
		return err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id)
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
		UserId:       params.UserId,
		TrackIds:     newIds,
		CurrentIndex: queue.CurrentIndex,
	})
	if err != nil {
		return err
	}

	return nil
}

type RemoveItemParams struct {
	UserId string
	ItemId string
}

func (s *QueueService) RemoveItem(
	ctx context.Context,
	params RemoveItemParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
	if err != nil {
		return err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id)
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
	UserId string
	Index  int
}

func (s *QueueService) SetPosition(
	ctx context.Context,
	params SetPositionParams,
) error {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
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
		return err
	}

	return nil
}

func (s *QueueService) ClearQueue(ctx context.Context, userId string) error {
	queue, err := s.getOrCreateQueue(ctx, userId)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.ClearQueueItems(ctx, queue.Id)
	if err != nil {
		return err
	}

	err = tx.UpdateQueue(ctx, queue.Id, database.QueueChanges{
		CurrentIndex: database.Change[int]{
			Value:   0,
			Changed: true,
		},
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}
