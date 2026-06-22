package service

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/utils"
	"github.com/robfig/cron/v3"
)

var taskErr = NewServiceErrCreator("task")

type TaskInfo struct {
	Name        string
	DisplayName string
	Schedule    string
}

type Task interface {
	Info() TaskInfo
	Run(ctx context.Context) error
}

var _ (broker.Event) = (*TaskSyncStateEvent)(nil)

type TaskSyncStateEventTask struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsRunning   bool   `json:"isRunning"`

	order int
}

type TaskSyncStateEvent struct {
	Tasks []TaskSyncStateEventTask `json:"tasks"`
}

func (e TaskSyncStateEvent) GetEventType() string {
	return "task-sync-state"
}

type taskEntry struct {
	task    Task
	info    *TaskInfo
	entryId *cron.EntryID

	order int

	isRunning bool

	lastRunSuccess bool
	lastRunTime    *time.Duration
	lastRunError   error
}

type TaskService struct {
	logger *slog.Logger

	cron  *cron.Cron
	tasks map[string]*taskEntry

	wg sync.WaitGroup
	mu sync.RWMutex

	updateFunc UpdateFunc
}

func NewTaskService(logger *slog.Logger) *TaskService {
	return &TaskService{
		logger: logger,
		cron:   cron.New(),
		tasks:  map[string]*taskEntry{},
	}
}

func (s *TaskService) GetSyncStateEvent() TaskSyncStateEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := TaskSyncStateEvent{
		Tasks: []TaskSyncStateEventTask{},
	}

	for _, task := range s.tasks {
		displayName := task.info.DisplayName
		if displayName == "" {
			displayName = task.info.Name
		}

		res.Tasks = append(res.Tasks, TaskSyncStateEventTask{
			Name:        task.info.Name,
			DisplayName: displayName,
			IsRunning:   task.isRunning,

			order: task.order,
		})
	}

	sort.SliceStable(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].order < res.Tasks[j].order
	})

	return res
}

func (s *TaskService) SetUpdateFunc(f UpdateFunc) {
	s.updateFunc = f
}

func (s *TaskService) update() {
	if s.updateFunc != nil {
		s.updateFunc()
	}
}

func (s *TaskService) AddTask(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	info := task.Info()

	// TODO(patrik): Maybe add some validation for TaskInfo

	_, exists := s.tasks[info.Name]
	if exists {
		s.logger.Error("task with name already exists", "name", info.Name)

		return taskErr.Newf("task with name already exists: %s", info.Name)
	}

	var entryId *cron.EntryID

	if info.Schedule != "" {
		id, err := s.cron.AddFunc(info.Schedule, func() {
			s.RunTask(context.Background(), info.Name)
		})
		if err != nil {
			s.logger.Error("add task to cron instance", "err", err)
			return taskErr.Wrap("add cron func", err)
		}

		entryId = &id
	}

	order := len(s.tasks)
	s.tasks[info.Name] = &taskEntry{
		task:    task,
		entryId: entryId,
		info:    &info,
		order:   order,
	}

	s.logger.Info(
		"new task",
		"name", info.Name,
		"schedule", info.Schedule,
		"order", order,
	)

	return nil
}

func (s *TaskService) Start() {
	s.cron.Start()
}

func (s *TaskService) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()

	s.Wait()
}

func (s *TaskService) Wait() {
	s.wg.Wait()
}

func (s *TaskService) RunTask(ctx context.Context, name string) {
	s.mu.Lock()
	entry, exists := s.tasks[name]
	if !exists {
		s.mu.Unlock()
		s.logger.Error("no task with name", "name", name)
		return
	}

	if entry.isRunning {
		s.mu.Unlock()
		s.logger.Error("task already running", "name", name)
		return
	}

	entry.isRunning = true
	s.mu.Unlock()

	s.update()

	s.wg.Add(1)
	defer func() {
		s.wg.Done()

		s.mu.Lock()
		entry.isRunning = false
		s.mu.Unlock()

		s.update()
	}()

	s.logger.Debug("running task", "task", name)

	timer := utils.SimpleTimer{}
	timer.Start()

	err := entry.task.Run(ctx)

	s.mu.Lock()
	dur := timer.Stop()
	entry.lastRunTime = &dur

	if err != nil {
		s.logger.Error("task error", "task", name, "err", err, "duration", dur)
		entry.lastRunError = err
		entry.lastRunSuccess = false
	} else {
		s.logger.Info("task success", "task", name, "duration", dur)
		entry.lastRunError = nil
		entry.lastRunSuccess = true
	}
	s.mu.Unlock()
}
