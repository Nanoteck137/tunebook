package service

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"time"

	"github.com/maruel/natural"
	"github.com/nanoteck137/dwebble/tools/broker"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/robfig/cron/v3"
)

type Job interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

var _ (broker.Event) = (*LibrarySyncStateEvent)(nil)

type JobSyncStateEventJob struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
}

type JobSyncStateEvent struct {
	Jobs []JobSyncStateEventJob `json:"jobs"`
}

func (e JobSyncStateEvent) GetEventType() string {
	return "job-sync-state"
}

type jobEntry struct {
	job     Job
	entryId *cron.EntryID

	isRunning bool

	lastRunSuccess bool
	lastRunTime    *time.Duration
	lastRunError   error
}

type JobService struct {
	logger *slog.Logger

	cron *cron.Cron
	jobs map[string]*jobEntry

	updateFunc UpdateFunc
}

func NewJobService(logger *slog.Logger) *JobService {
	return &JobService{
		logger: logger,
		cron:   cron.New(),
		jobs:   map[string]*jobEntry{},
	}
}

func (s *JobService) Init() error {
	s.logger.Info("initializing job service")

	return nil
}

func (s *JobService) GetSyncStateEvent() JobSyncStateEvent {
	// TODO(patrik): Lock

	res := JobSyncStateEvent{
		Jobs: []JobSyncStateEventJob{},
	}

	for _, job := range s.jobs {
		res.Jobs = append(res.Jobs, JobSyncStateEventJob{
			Name:      job.job.Name(),
			IsRunning: job.isRunning,
		})
	}

	sort.SliceStable(res.Jobs, func(i, j int) bool {
		return natural.Less(res.Jobs[i].Name, res.Jobs[j].Name)
	})

	return res
}

func (s *JobService) SetUpdateFunc(f UpdateFunc) {
	s.updateFunc = f
}

func (s *JobService) update() {
	if s.updateFunc != nil {
		s.updateFunc()
	}
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

	s.jobs[job.Name()] = &jobEntry{
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
		s.logger.Info("registered job",
			"name", job.job.Name(),
			"schedule", job.job.Schedule(),
		)
	}
}

func (s *JobService) Start() {
	s.cron.Start()
}

func (s *JobService) Stop() {
	s.cron.Start()
}

func (s *JobService) RunJob(ctx context.Context, name string) {
	// TODO(patrik): Add lock
	entry, exists := s.jobs[name]
	if !exists {
		s.logger.Error("no job with name", "name", name)
		return
	}

	if entry.isRunning {
		s.logger.Error("job already running", "name", name)
		return
	}

	entry.isRunning = true
	s.update()

	defer func() {
		entry.isRunning = false
		s.update()
	}()

	s.logger.Info("running job", "job", entry.job.Name())

	timer := utils.SimpleTimer{}
	timer.Start()

	err := entry.job.Run(ctx)

	dur := timer.Stop()
	entry.lastRunTime = &dur

	if err != nil {
		s.logger.Error("job returned error", "err", err, "job", name, "duration", dur)
		entry.lastRunError = err
		entry.lastRunSuccess = false
	} else {
		s.logger.Info("job was successfully", "job", name, "duration", dur)
		entry.lastRunError = nil
		entry.lastRunSuccess = true
	}
}
