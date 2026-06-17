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

const defaultQueuePerPage = 50

type GetQueueResult struct {
	Items        []database.QueueItemTrack
	CurrentIndex int
	Page         types.Page
}

type GetQueueIdsResult struct {
	Entries      []database.QueueItemEntry
	CurrentIndex int
}

func (s *QueueService) GetQueue(ctx context.Context, userId string, page, perPage int) (*GetQueueResult, error) {
	queue, err := s.getOrCreateQueue(ctx, userId)
	if err != nil {
		return nil, err
	}

	if perPage <= 0 {
		perPage = defaultQueuePerPage
	}

	count, err := s.db.GetQueueItemCount(ctx, queue.Id)
	if err != nil {
		return nil, err
	}

	items, err := s.db.GetQueueItems(ctx, database.GetQueueItemsParams{
		QueueId: queue.Id,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	pageInfo := types.Page{
		Page:       page,
		PerPage:    perPage,
		TotalItems: count,
		TotalPages: (count + perPage - 1) / perPage,
	}

	return &GetQueueResult{
		Items:        items,
		CurrentIndex: queue.CurrentIndex,
		Page:         pageInfo,
	}, nil
}

func (s *QueueService) GetQueueIds(ctx context.Context, userId string) (*GetQueueIdsResult, error) {
	queue, err := s.getOrCreateQueue(ctx, userId)
	if err != nil {
		return nil, err
	}

	entries, err := s.db.GetQueueItemEntries(ctx, queue.Id)
	if err != nil {
		return nil, err
	}

	return &GetQueueIdsResult{
		Entries:      entries,
		CurrentIndex: queue.CurrentIndex,
	}, nil
}

func (s *QueueService) GetQueueItemAtIndex(ctx context.Context, userId string, index int) (database.QueueItemTrack, error) {
	queue, err := s.getOrCreateQueue(ctx, userId)
	if err != nil {
		return database.QueueItemTrack{}, err
	}

	return s.db.GetQueueItemAtPosition(ctx, queue.Id, index)
}

func (s *QueueService) getOrCreateQueue(ctx context.Context, userId string) (database.Queue, error) {
	queue, err := s.db.GetQueueByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			err = s.db.CreateQueue(ctx, userId, 0)
			if err != nil {
				return database.Queue{}, err
			}

			queue, err = s.db.GetQueueByUserId(ctx, userId)
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

func (s *QueueService) ReplaceQueue(ctx context.Context, params ReplaceQueueParams) error {
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
		CurrentIndex: database.Change[int]{Value: currentIndex, Changed: true},
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

type AddItemsParams struct {
	UserId   string
	TrackIds []string
	Position string // "next" or "end"
}

func (s *QueueService) AddItems(ctx context.Context, params AddItemsParams) error {
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
		insertIndex := queue.CurrentIndex + 1
		if insertIndex > len(existingIds) {
			insertIndex = len(existingIds)
		}

		newIds = make([]string, 0, len(existingIds)+len(params.TrackIds))
		newIds = append(newIds, existingIds[:insertIndex]...)
		newIds = append(newIds, params.TrackIds...)
		newIds = append(newIds, existingIds[insertIndex:]...)
	} else {
		newIds = append(existingIds, params.TrackIds...)
	}

	return s.ReplaceQueue(ctx, ReplaceQueueParams{
		UserId:       params.UserId,
		TrackIds:     newIds,
		CurrentIndex: queue.CurrentIndex,
	})
}

type RemoveItemParams struct {
	UserId string
	ItemId string
}

func (s *QueueService) RemoveItem(ctx context.Context, params RemoveItemParams) error {
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

	return s.ReplaceQueue(ctx, ReplaceQueueParams{
		UserId:       params.UserId,
		TrackIds:     newIds,
		CurrentIndex: newIndex,
	})
}

type SetPositionParams struct {
	UserId string
	Index  int
}

func (s *QueueService) SetPosition(ctx context.Context, params SetPositionParams) error {
	queue, err := s.getOrCreateQueue(ctx, params.UserId)
	if err != nil {
		return err
	}

	return s.db.UpdateQueue(ctx, queue.Id, database.QueueChanges{
		CurrentIndex: database.Change[int]{Value: params.Index, Changed: true},
	})
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
		CurrentIndex: database.Change[int]{Value: 0, Changed: true},
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}
