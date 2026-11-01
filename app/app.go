package app

import (
	"context"
	"log/slog"

	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/database/migrations"
	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tasks"
	"github.com/nanoteck137/tunebook/tools/broker"
)

type App struct {
	db     *database.Database
	config *config.Config

	authService         *service.AuthService
	userService         *service.UserService
	taskService         *service.TaskService
	notificationService *service.NotificationService
	searchService       *service.SearchService
	libraryService      *service.LibraryService
	imageService        *service.ImageService
	mediaService        *service.MediaService
	filesystemService   *service.FilesystemService

	artistService   *service.ArtistService
	albumService    *service.AlbumService
	trackService    *service.TrackService
	playlistService *service.PlaylistService
	historyService  *service.HistoryService
	queueService    *service.QueueService
	jobService      *service.JobService

	broker *broker.Broker
}

func (app *App) UserService() *service.UserService {
	return app.userService
}

func (app *App) ArtistService() *service.ArtistService {
	return app.artistService
}

func (app *App) AlbumService() *service.AlbumService {
	return app.albumService
}

func (app *App) TrackService() *service.TrackService {
	return app.trackService
}

func (app *App) PlaylistService() *service.PlaylistService {
	return app.playlistService
}

func (app *App) HistoryService() *service.HistoryService {
	return app.historyService
}

func (app *App) QueueService() *service.QueueService {
	return app.queueService
}

func (app *App) JobService() *service.JobService {
	return app.jobService
}

func (app *App) TaskService() *service.TaskService {
	return app.taskService
}

func (app *App) NotificationService() *service.NotificationService {
	return app.notificationService
}

func (app *App) MediaService() *service.MediaService {
	return app.mediaService
}

func (app *App) Broker() *broker.Broker {
	return app.broker
}

func (app *App) ImageService() *service.ImageService {
	return app.imageService
}

func (app *App) FilesystemService() *service.FilesystemService {
	return app.filesystemService
}

func (app *App) LibraryService() *service.LibraryService {
	return app.libraryService
}

func (app *App) SearchService() *service.SearchService {
	return app.searchService
}

func (app *App) AuthService() *service.AuthService {
	return app.authService
}

func (app *App) DB() *database.Database {
	return app.db
}

func (app *App) Config() *config.Config {
	return app.config
}

func (app *App) Bootstrap() error {
	var err error

	app.filesystemService = service.NewFilesystemService(
		slog.With(slog.String("service", "filesystem")),
		app.config.DataDir,
	)

	err = app.filesystemService.EnsureBaseDirs()
	if err != nil {
		return err
	}

	app.db, err = database.Open(app.filesystemService.DatabaseFile())
	if err != nil {
		return err
	}

	if app.config.RunMigrations {
		err = migrations.RunMigrateUp(context.Background(), app.db)
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
		newServiceLogger("notification"),
		app.config,
	)

	if !app.notificationService.IsEnabled() {
		slog.Warn("notification service disabled, ntfy_base_url or ntfy_topic not set")
	}

	app.broker = broker.NewBroker()

	app.taskService = service.NewTaskService(
		newServiceLogger("task"),
		app.broker,
	)

	app.jobService = service.NewJobService(
		newServiceLogger("job"),
		app.db,
		app.broker,
	)

	app.imageService = service.NewImageService(
		newServiceLogger("image"),
		app.db,
		app.filesystemService,
	)

	app.authService = service.NewAuthService(
		newServiceLogger("auth"),
		app.db,
		app.config,
		app.imageService,
	)

	app.userService = service.NewUserService(
		newServiceLogger("user"),
		app.db,
		app.filesystemService,
		app.imageService,
		app.broker,
	)

	app.searchService = service.NewSearchService(
		newServiceLogger("search"),
		app.db,
		app.config,
	)

	app.mediaService = service.NewMediaService(
		newServiceLogger("media"),
		app.db,
		app.filesystemService,
		app.config,
	)

	app.libraryService = service.NewLibraryService(
		newServiceLogger("library"),
		app.db,
		app.config,
		app.filesystemService,
		app.notificationService,
		app.mediaService,
		app.broker,
	)

	app.artistService = service.NewArtistService(
		newServiceLogger("artist"),
		app.db,
		app.imageService,
		app.filesystemService,
	)

	app.albumService = service.NewAlbumService(
		newServiceLogger("album"),
		app.db,
		app.imageService,
		app.filesystemService,
	)

	app.trackService = service.NewTrackService(
		newServiceLogger("track"),
		app.db,
		app.broker,
	)

	app.playlistService = service.NewPlaylistService(
		newServiceLogger("playlist"),
		app.db,
		app.filesystemService,
		app.imageService,
		app.broker,
	)

	app.historyService = service.NewHistoryService(
		newServiceLogger("track_history"),
		app.db,
	)

	app.queueService = service.NewQueueService(
		newServiceLogger("queue"),
		app.db,
	)

	app.broker.RegisterProducer(app.libraryService)
	app.broker.RegisterProducer(app.taskService)
	app.broker.RegisterProducer(app.jobService)

	taskList := []service.Task{
		tasks.NewLibrarySyncTask(app.jobService),
		tasks.NewSearchIndexTask(app.jobService),
		tasks.NewUserStatsRecalculateTask(app.userService, app.jobService),
		tasks.NewUserTrackStatsRebuildTask(app.userService, app.jobService),
		tasks.NewAuthCleanupTask(app.jobService),
		tasks.NewCacheCleanupTask(app.jobService),
		tasks.NewLibraryCleanupTask(app.jobService),
		tasks.NewJobsCleanupTask(app.jobService),
		tasks.NewUserReviewRebuildTask(app.userService, app.jobService),
		tasks.NewTestTask(app.jobService),
	}

	for _, task := range taskList {
		err = app.taskService.AddTask(task)
		if err != nil {
			return err
		}
	}

	jobList := []service.Job{
		jobs.NewLibrarySyncJob(app.libraryService, app.jobService),
		jobs.NewSearchIndexJob(app.searchService),
		jobs.NewAuthCleanupJob(app.authService),
		jobs.NewCacheCleanupJob(app.filesystemService),
		jobs.NewLibraryCleanupJob(app.libraryService, app.jobService),
		jobs.NewJobsCleanupJob(app.jobService),
		jobs.NewGeneratePlaylistImageJob(app.playlistService),
		jobs.NewUserStatsUpdateJob(app.userService),
		jobs.NewUserTrackStatsRebuildJob(app.userService),
		jobs.NewUserReviewGenerateJob(app.userService),
		jobs.NewUserReviewGenerateTotalJob(app.userService),
		jobs.NewUserReviewRebuildJob(app.userService),
		jobs.NewTestJob(),
	}

	for _, job := range jobList {
		err = app.jobService.AddJob(job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Start() {
	app.jobService.Start()
	app.taskService.Start()
	app.broker.Start()
}

func (app *App) Shutdown() error {
	app.jobService.Stop()
	app.taskService.Stop()

	err := app.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}
