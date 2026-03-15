package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/robfig/cron/v3"
)

type JobFunc func(ctx context.Context) error

type Job interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

type jobEntry struct {
	job     Job
	entryId *cron.EntryID

	// success
	// lastRun   *time.Time
	// lastError error

	runner JobFunc
}

type JobService struct {
	logger *slog.Logger

	cron *cron.Cron
	jobs map[string]jobEntry
}

func NewJobService(logger *slog.Logger) *JobService {
	return &JobService{
		logger: logger,
		cron:   cron.New(),
		jobs:   map[string]jobEntry{},
	}
}

func (s *JobService) Init() error {
	s.logger.Info("initializing job service")

	return nil
}

func (s *JobService) AddJob(job Job) error {
	_, exists := s.jobs[job.Name()]
	if exists {
		s.logger.Error("job with name already exists",
			slog.String("name", job.Name()),
		)

		return errors.New("job with name already exists: " + job.Name())
	}

	schedule := job.Schedule()

	var entryId *cron.EntryID

	if schedule != "" {
		id, err := s.cron.AddFunc(job.Schedule(), func() {

			s.RunJob(context.Background(), job.Name())
		})
		if err != nil {
			s.logger.Error("failed to add job to cron instance", "err", err)
			return err
		}

		entryId = &id
	}

	s.jobs[job.Name()] = jobEntry{
		job:     job,
		entryId: entryId,
	}

	s.logger.Info("added job",
		slog.String("name", job.Name()),
	)

	return nil
}

func (s *JobService) DisplayJobs() {
	for _, job := range s.jobs {
		s.logger.Info("registered job", "name", job.job.Name(), "schedule", job.job.Schedule())
	}
}

func (s *JobService) Start() {
	s.cron.Start()
}

func (s *JobService) Stop() {
	s.cron.Start()
}

func (s *JobService) getJobEntry(name string) (jobEntry, bool) {
	// TODO(patrik): Add lock
	res, exists := s.jobs[name]

	return res, exists
}

func (s *JobService) RunJob(ctx context.Context, name string) {
	entry, exists := s.getJobEntry(name)
	if !exists {
		s.logger.Error("no job with name", "name", name)
	}

	s.logger.Info("running job", "job", entry.job.Name())

	timer := utils.SimpleTimer{}
	timer.Start()

	err := entry.job.Run(ctx)
	if err != nil {
		s.logger.Error("job returned error", "err", err, "job", name)
		return
	}

	dur := timer.Stop()

	s.logger.Info("job was successfully", "job", entry.job.Name(), "duration", dur)
}
