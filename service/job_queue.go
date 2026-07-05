package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nanoteck137/tunebook/database"
)

var jobQueueErr = NewServiceErrCreator("job_queue")

type JobHandler func(ctx context.Context, data string) error

type JobQueueService struct {
	logger   *slog.Logger
	db       *database.Database
	handlers map[string]JobHandler
	mu       sync.RWMutex

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func NewJobQueueService(logger *slog.Logger, db *database.Database) *JobQueueService {
	return &JobQueueService{
		logger:   logger,
		db:       db,
		handlers: make(map[string]JobHandler),
	}
}

func (s *JobQueueService) Start() {
	s.stopCh = make(chan struct{})

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := s.ProcessPendingJobs(context.Background())
				if err != nil {
					s.logger.Error("process pending jobs", "err", err)
				}
			case <-s.stopCh:
				return
			}
		}
	}()

	s.logger.Info("job queue worker started")
}

func (s *JobQueueService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("job queue worker stopped")
}

func (s *JobQueueService) RegisterJob(name string, handler JobHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[name]; exists {
		s.logger.Warn("job handler already registered, overwriting", "name", name)
	}

	s.handlers[name] = handler
	s.logger.Info("registered job handler", "name", name)
}

func (s *JobQueueService) PushJob(ctx context.Context, name string, data any) error {
	s.mu.RLock()
	_, exists := s.handlers[name]
	s.mu.RUnlock()

	if !exists {
		return jobQueueErr.Newf("no handler registered for job: %s", name)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return jobQueueErr.Wrap("marshal job data", err)
	}

	id, err := s.db.CreateJob(ctx, database.CreateJobParams{
		Name: name,
		Data: string(raw),
	})
	if err != nil {
		return jobQueueErr.Wrap("create job", err)
	}

	s.logger.Info("pushed job", "id", id, "name", name)

	return nil
}

func (s *JobQueueService) ProcessPendingJobs(ctx context.Context) error {
	jobs, err := s.db.GetPendingJobs(ctx, 10)
	if err != nil {
		return jobQueueErr.Wrap("get pending jobs", err)
	}

	if len(jobs) > 0 {
		s.logger.Info("processing pending jobs", "count", len(jobs))
	}

	for _, job := range jobs {
		err := s.processJob(ctx, job)
		if err != nil {
			s.logger.Error("failed to process job", "id", job.Id, "name", job.Name, "err", err)
		}
	}

	return nil
}

func (s *JobQueueService) processJob(ctx context.Context, job database.Job) error {
	err := s.db.ClaimJob(ctx, job.Id)
	if err != nil {
		return jobQueueErr.Wrap("claim job", err)
	}

	s.logger.Info("running job", "id", job.Id, "name", job.Name, "attempt", job.Attempts+1)

	s.mu.RLock()
	handler, exists := s.handlers[job.Name]
	s.mu.RUnlock()

	if !exists {
		errMsg := fmt.Sprintf("no handler registered for job: %s", job.Name)
		s.logger.Error("no handler for job", "id", job.Id, "name", job.Name)
		s.db.FailJob(ctx, job.Id, database.FailJobParams{
			Requeue: false,
			Error:   errMsg,
		})
		return jobQueueErr.New(errMsg)
	}

	err = handler(ctx, job.Data)
	if err != nil {
		shouldRequeue := job.Attempts+1 < job.MaxAttempts

		s.logger.Error("job failed", "id", job.Id, "name", job.Name, "attempt", job.Attempts+1, "maxAttempts", job.MaxAttempts, "requeue", shouldRequeue, "err", err)

		s.db.FailJob(ctx, job.Id, database.FailJobParams{
			Requeue: shouldRequeue,
			Error:   err.Error(),
		})
		return jobQueueErr.Wrap("job handler failed", err)
	}

	err = s.db.CompleteJob(ctx, job.Id)
	if err != nil {
		return jobQueueErr.Wrap("complete job", err)
	}

	s.logger.Info("job completed", "id", job.Id, "name", job.Name)

	return nil
}
