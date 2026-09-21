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

var jobErr = NewServiceErrCreator("job")

const (
	// DefaultJobMaxAttempts is how many times a job is attempted before it is
	// permanently failed, unless overridden with WithMaxAttempts.
	DefaultJobMaxAttempts = 3

	jobRunTimeout = 5 * time.Minute
	stuckJobGrace = 5 * time.Minute

	jobBackoffBase = 30 * time.Second
	jobMaxBackoff  = 30 * time.Minute
)

type JobInfo struct {
	Name        string
	DisplayName string
}

type Job interface {
	Info() JobInfo
	Run(ctx context.Context, data string) error
}

type JobHandler func(ctx context.Context, data string) error

type jobEntry struct {
	job  Job
	info *JobInfo
}

type JobService struct {
	logger   *slog.Logger
	db       *database.Database
	handlers map[string]JobHandler
	jobs     map[string]*jobEntry
	mu       sync.RWMutex

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func NewJobService(logger *slog.Logger, db *database.Database) *JobService {
	return &JobService{
		logger:   logger,
		db:       db,
		handlers: make(map[string]JobHandler),
		jobs:     make(map[string]*jobEntry),
	}
}

func (s *JobService) Start() {
	s.stopCh = make(chan struct{})

	requeued, err := s.db.RequeueStuckJobs(
		context.Background(),
		time.Now().UnixMilli()-stuckJobGrace.Milliseconds(),
	)
	if err != nil {
		s.logger.Error("requeue stuck jobs", "err", err)
	} else if requeued > 0 {
		s.logger.Info("requeued stuck jobs", "count", requeued)
	}

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

func (s *JobService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("job queue worker stopped")
}

func (s *JobService) RegisterJob(name string, handler JobHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[name]; exists {
		s.logger.Warn(
			"job handler already registered, overwriting", "name", name)
	}

	s.handlers[name] = handler
	s.logger.Info("registered job handler", "name", name)
}

func (s *JobService) AddJob(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	info := job.Info()

	if info.Name == "" {
		return jobErr.New("job name must not be empty")
	}

	_, exists := s.jobs[info.Name]
	if exists {
		s.logger.Error("job with name already exists", "name", info.Name)
		return jobErr.Newf("job with name already exists: %s", info.Name)
	}

	s.jobs[info.Name] = &jobEntry{
		job:  job,
		info: &info,
	}

	s.handlers[info.Name] = job.Run

	s.logger.Info(
		"new job",
		"name", info.Name,
	)

	return nil
}

type PushJobOptions struct {
	MaxAttempts int
}

type PushJobOption func(*PushJobOptions)

func WithMaxAttempts(maxAttempts int) PushJobOption {
	return func(o *PushJobOptions) {
		o.MaxAttempts = maxAttempts
	}
}

func (s *JobService) PushJob(
	ctx context.Context,
	name string,
	data any,
	opts ...PushJobOption,
) error {
	options := PushJobOptions{
		MaxAttempts: DefaultJobMaxAttempts,
	}
	for _, opt := range opts {
		opt(&options)
	}

	s.mu.RLock()
	_, exists := s.handlers[name]
	s.mu.RUnlock()

	if !exists {
		return jobErr.Newf("no handler registered for job: %s", name)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return jobErr.Wrap("marshal job data", err)
	}

	id, err := s.db.CreateJob(ctx, database.CreateJobParams{
		Name:        name,
		Data:        string(raw),
		MaxAttempts: options.MaxAttempts,
	})
	if err != nil {
		return jobErr.Wrap("create job", err)
	}

	s.logger.Info("pushed job", "id", id, "name", name)

	return nil
}

func (s *JobService) ProcessPendingJobs(ctx context.Context) error {
	jobs, err := s.db.GetPendingJobs(ctx, 10)
	if err != nil {
		return jobErr.Wrap("get pending jobs", err)
	}

	if len(jobs) > 0 {
		s.logger.Info("processing pending jobs", "count", len(jobs))
	}

	for _, job := range jobs {
		err := s.processJob(ctx, job)
		if err != nil {
			s.logger.Error(
				"failed to process job",
				"id", job.Id,
				"name", job.Name,
				"err", err,
			)
		}
	}

	return nil
}

func (s *JobService) processJob(ctx context.Context, job database.Job) error {
	err := s.db.ClaimJob(ctx, job.Id)
	if err != nil {
		return jobErr.Wrap("claim job", err)
	}

	s.logger.Info(
		"running job",
		"id", job.Id,
		"name", job.Name,
		"attempt", job.Attempts+1,
	)

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
		return jobErr.New(errMsg)
	}

	runCtx, cancel := context.WithTimeout(ctx, jobRunTimeout)
	defer cancel()

	err = handler(runCtx, job.Data)
	if err != nil {
		shouldRequeue := job.Attempts+1 < job.MaxAttempts

		nextAttemptAt := int64(0)
		if shouldRequeue {
			delay := jobBackoff(job.Attempts + 1)
			nextAttemptAt = time.Now().Add(delay).UnixMilli()
		}

		s.logger.Error(
			"job failed",
			"id", job.Id,
			"name", job.Name,
			"attempt", job.Attempts+1,
			"maxAttempts", job.MaxAttempts,
			"requeue", shouldRequeue,
			"err", err,
		)

		s.db.FailJob(ctx, job.Id, database.FailJobParams{
			Requeue:       shouldRequeue,
			Error:         err.Error(),
			NextAttemptAt: nextAttemptAt,
		})
		return jobErr.Wrap("job handler failed", err)
	}

	err = s.db.CompleteJob(ctx, job.Id)
	if err != nil {
		return jobErr.Wrap("complete job", err)
	}

	s.logger.Info("job completed", "id", job.Id, "name", job.Name)

	return nil
}

func jobBackoff(attempt int) time.Duration {
	delay := jobBackoffBase
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= jobMaxBackoff {
			return jobMaxBackoff
		}
	}

	if delay > jobMaxBackoff {
		delay = jobMaxBackoff
	}

	return delay
}

// CleanupOldJobs deletes finished jobs older than olderThan and returns how
// many rows were removed.
func (s *JobService) CleanupOldJobs(
	ctx context.Context,
	olderThan time.Duration,
) (int64, error) {
	count, err := s.db.DeleteOldJobs(
		ctx,
		time.Now().UnixMilli()-olderThan.Milliseconds(),
	)
	if err != nil {
		return 0, jobErr.Wrap("delete old jobs", err)
	}

	return count, nil
}
