package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/broker"
)

var jobErr = NewServiceErrCreator("job")

const (
	// DefaultJobMaxAttempts is how many times a job is attempted before it is
	// permanently failed, unless overridden with WithMaxAttempts.
	DefaultJobMaxAttempts = 3

	// DefaultJobsLimit is how many jobs are returned/sent when listing or
	// syncing jobs.
	DefaultJobsLimit = 50

	jobRunTimeout = 5 * time.Minute

	jobBackoffBase = 30 * time.Second
	jobMaxBackoff  = 30 * time.Minute
)

type JobInfo struct {
	Name        string
	DisplayName string

	// FailOnRestart marks the job as failed instead of requeueing it if the
	// server restarts while the job is running. Defaults to false, meaning the
	// job is requeued (and retried) after a restart. Jobs like a full library
	// sync should set this so they don't resume from a half-finished run.
	FailOnRestart bool
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

	emitter broker.EventEmitter

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func NewJobService(
	logger *slog.Logger,
	db *database.Database,
	emitter broker.EventEmitter,
) *JobService {
	return &JobService{
		logger:   logger,
		db:       db,
		handlers: make(map[string]JobHandler),
		jobs:     make(map[string]*jobEntry),
		emitter:  emitter,
	}
}

var _ broker.Event = (*GetJobsResponse)(nil)
var _ broker.EventProducer = (*JobService)(nil)

// GetEventType implements broker.Event.
func (e GetJobsResponse) GetEventType() string {
	return "job-sync-state"
}

func (s *JobService) update() {
	if s.emitter == nil {
		return
	}

	jobs, err := s.GetJobs(context.Background(), DefaultJobsLimit)
	if err != nil {
		s.logger.Error("get jobs for sync event", "err", err)
		return
	}

	s.emitter.EmitEvent(jobs)
}

func (s *JobService) GetInitEvents() []broker.Event {
	jobs, err := s.GetJobs(context.Background(), DefaultJobsLimit)
	if err != nil {
		s.logger.Error("get jobs for init event", "err", err)
		return nil
	}

	return []broker.Event{jobs}
}

func (s *JobService) Start() {
	s.stopCh = make(chan struct{})

	s.recoverStuckJobs()

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

// recoverStuckJobs handles jobs that were left in the "running" state, meaning
// they were interrupted by a shutdown, crash or restart. Restartable jobs are
// requeued so they are retried, while jobs marked FailOnRestart are marked
// failed instead.
func (s *JobService) recoverStuckJobs() {
	stuck, err := s.db.GetStuckJobs(context.Background())
	if err != nil {
		s.logger.Error("get stuck jobs", "err", err)
		return
	}

	if len(stuck) == 0 {
		return
	}

	requeued := 0
	failed := 0
	for _, job := range stuck {
		if s.jobRestartable(job.Name) {
			err := s.db.RequeueJob(
				context.Background(),
				job.Id,
				time.Now().UnixMilli(),
			)
			if err != nil {
				s.logger.Error("requeue stuck job", "id", job.Id, "err", err)
				continue
			}

			requeued++
			continue
		}

		err := s.db.FailJob(
			context.Background(),
			job.Id,
			database.FailJobParams{
				Requeue: false,
				Error:   "job interrupted by server restart",
			},
		)
		if err != nil {
			s.logger.Error("fail stuck job", "id", job.Id, "err", err)
			continue
		}

		failed++
	}

	s.logger.Info(
		"recovered stuck jobs",
		"requeued", requeued,
		"failed", failed,
	)
}

func (s *JobService) jobRestartable(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.jobs[name]
	if !exists {
		return true
	}

	return !entry.info.FailOnRestart
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
	UniqueKey   string
}

type PushJobOption func(*PushJobOptions)

func WithMaxAttempts(maxAttempts int) PushJobOption {
	return func(o *PushJobOptions) {
		o.MaxAttempts = maxAttempts
	}
}

// WithUniqueKey sets a unique key for the job. Only one job with the same
// unique key may be pending or running at a time; pushing again while one is
// active is a no-op. Keys are used to dedupe jobs that must not be queued
// multiple times, e.g. a full library sync or image generation for the same
// playlist.
func WithUniqueKey(uniqueKey string) PushJobOption {
	return func(o *PushJobOptions) {
		o.UniqueKey = uniqueKey
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

	if options.UniqueKey != "" {
		active, err := s.db.HasActiveJobWithUniqueKey(ctx, options.UniqueKey)
		if err != nil {
			return jobErr.Wrap("check active job by unique key", err)
		}

		if active {
			s.logger.Info(
				"job already active, skipping push",
				"name", name,
				"uniqueKey", options.UniqueKey,
			)
			return nil
		}
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return jobErr.Wrap("marshal job data", err)
	}

	id, err := s.db.CreateJob(ctx, database.CreateJobParams{
		Name:        name,
		Data:        string(raw),
		UniqueKey:   options.UniqueKey,
		MaxAttempts: options.MaxAttempts,
	})
	if err != nil {
		return jobErr.Wrap("create job", err)
	}

	s.logger.Info("pushed job", "id", id, "name", name)

	s.update()

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

	s.update()

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
		s.update()
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
		s.update()
		return jobErr.Wrap("job handler failed", err)
	}

	err = s.db.CompleteJob(ctx, job.Id)
	if err != nil {
		return jobErr.Wrap("complete job", err)
	}

	s.update()

	s.logger.Info("job completed", "id", job.Id, "name", job.Name)

	return nil
}

// GetJobsResponseItem is a single job as returned to API clients.
type JobItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Error       string `json:"error"`
	Attempts    int    `json:"attempts"`
	MaxAttempts int    `json:"maxAttempts"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type GetJobsResponse struct {
	Jobs []JobItem `json:"jobs"`
}

type JobsCleanupParams struct {
	OlderThanMs int64 `json:"olderThanMs"`
}

func (s *JobService) GetJobs(ctx context.Context, limit int) (GetJobsResponse, error) {
	jobs, err := s.db.GetJobs(ctx, limit)
	if err != nil {
		return GetJobsResponse{}, jobErr.Wrap("get jobs", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	res := GetJobsResponse{
		Jobs: make([]JobItem, 0, len(jobs)),
	}

	for _, job := range jobs {
		displayName := job.Name
		if entry, exists := s.jobs[job.Name]; exists && entry.info.DisplayName != "" {
			displayName = entry.info.DisplayName
		}

		res.Jobs = append(res.Jobs, JobItem{
			Id:          job.Id,
			Name:        job.Name,
			DisplayName: displayName,
			Status:      job.Status,
			Error:       job.Error,
			Attempts:    job.Attempts,
			MaxAttempts: job.MaxAttempts,
			Created:     job.Created,
			Updated:     job.Updated,
		})
	}

	return res, nil
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
