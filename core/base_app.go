package core

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/broker"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

const (
	jobAuthCleanup  = "auth-cleanup"
	jobCacheCleanup = "cache-cleanup"
)

var _ service.Job = (*AuthCleanupJob)(nil)

type AuthCleanupJob struct {
	authService *service.AuthService
}

func (j *AuthCleanupJob) Name() string {
	return jobAuthCleanup
}

func (j *AuthCleanupJob) Schedule() string {
	return "@every 30m"
}

func (j *AuthCleanupJob) Run(ctx context.Context) error {
	time.Sleep(4 * time.Second)
	j.authService.RunCleanup()
	return nil
}

var _ service.Job = (*CacheCleanupJob)(nil)

type CacheCleanupJob struct {
	workDir types.WorkDir
}

func (j *CacheCleanupJob) Name() string {
	return jobCacheCleanup
}

func (j *CacheCleanupJob) Schedule() string {
	return ""
}

func (j *CacheCleanupJob) Run(ctx context.Context) error {
	cacheDir := j.workDir.Cache()

	err := os.RemoveAll(cacheDir.String())
	if err != nil {
		return err
	}

	err = utils.CreateDirectories([]string{
		cacheDir.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db     *database.Database
	config *config.Config

	authService         *service.AuthService
	jobService          *service.JobService
	notificationService *service.NotificationService
	searchService       *service.SearchService
	libraryService      *service.LibraryService
	imageService        *service.ImageService
	mediaService        *service.MediaService

	broker *broker.Broker
}

func (app *BaseApp) JobService() *service.JobService {
	return app.jobService
}

func (app *BaseApp) NotificationService() *service.NotificationService {
	return app.notificationService
}

func (app *BaseApp) MediaService() *service.MediaService {
	return app.mediaService
}

func (app *BaseApp) Broker() *broker.Broker {
	return app.broker
}

func (app *BaseApp) ImageService() *service.ImageService {
	return app.imageService
}

func (app *BaseApp) LibraryService() *service.LibraryService {
	return app.libraryService
}

func (app *BaseApp) SearchService() *service.SearchService {
	return app.searchService
}

func (app *BaseApp) AuthService() *service.AuthService {
	return app.authService
}

func (app *BaseApp) DB() *database.Database {
	return app.db
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	err = utils.CreateDirectories([]string{
		workDir.Artists(),
		workDir.Albums(),
		workDir.Tracks(),
		workDir.Playlists(),
		workDir.Trash(),
		workDir.Cache().String(),
	})
	if err != nil {
		return err
	}

	app.db, err = database.Open(workDir.DatabaseFile())
	if err != nil {
		return err
	}

	if app.config.RunMigrations {
		err = app.db.RunMigrateUp()
		if err != nil {
			return err
		}
	}

	newServiceLogger := func(name string) *slog.Logger {
		return slog.With(
			slog.String("service", name),
		)
	}

	app.notificationService = service.NewNotificationService(
		newServiceLogger("notification-service"),
		app.config,
	)

	app.jobService = service.NewJobService(newServiceLogger("job-service"))
	err = app.jobService.Init()
	if err != nil {
		return err
	}

	app.authService = service.NewAuthService(app.db, app.config)
	// TODO(patrik): This should be a worker
	// go app.authService.CleanRoutine()

	app.searchService = service.NewSearchService(app.db, app.config)

	// TODO(patrik): Do this lazily
	err = app.searchService.Init()
	if err != nil {
		return err
	}

	app.mediaService = service.NewMediaService(app.db, app.config.WorkDir())

	app.libraryService = service.NewLibraryService(
		app.db,
		app.config,
		app.notificationService,
		app.mediaService,
		app.searchService,
	)

	app.imageService = service.NewImageService(
		newServiceLogger("image-service"),
		app.db,
		app.config.WorkDir(),
	)

	app.broker = broker.NewBroker(func() []broker.Event {
		return []broker.Event{
			app.libraryService.GetSyncStateEvent(),
			app.jobService.GetSyncStateEvent(),
		}
	})

	app.libraryService.SetUpdateFunc(func() {
		app.broker.EmitEvent(app.libraryService.GetSyncStateEvent())
	})

	app.jobService.SetUpdateFunc(func() {
		app.broker.EmitEvent(app.jobService.GetSyncStateEvent())
	})

	err = app.jobService.AddJob(&AuthCleanupJob{
		authService: app.authService,
	})
	if err != nil {
		return err
	}

	err = app.jobService.AddJob(&CacheCleanupJob{
		workDir: workDir,
	})
	if err != nil {
		return err
	}

	// TODO(patrik): This should not be in bootstrap
	app.jobService.DisplayJobs()
	app.jobService.Start()

	// TODO(patrik): This should not be in bootstrap
	go app.broker.Listen()

	app.jobService.RunJob(context.Background(), jobAuthCleanup)

	return nil
}

func (app *BaseApp) Shutdown() error {
	app.jobService.Stop()

	err := app.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
