package core

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/database/migrations"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tasks"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
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

	artistService   *service.ArtistService
	albumService    *service.AlbumService
	trackService    *service.TrackService
	playlistService *service.PlaylistService
	historyService  *service.HistoryService
	queueService    *service.QueueService
	jobService      *service.JobService

	broker *broker.Broker
}

func (app *BaseApp) UserService() *service.UserService {
	return app.userService
}

func (app *BaseApp) ArtistService() *service.ArtistService {
	return app.artistService
}

func (app *BaseApp) AlbumService() *service.AlbumService {
	return app.albumService
}

func (app *BaseApp) TrackService() *service.TrackService {
	return app.trackService
}

func (app *BaseApp) PlaylistService() *service.PlaylistService {
	return app.playlistService
}

func (app *BaseApp) HistoryService() *service.HistoryService {
	return app.historyService
}

func (app *BaseApp) QueueService() *service.QueueService {
	return app.queueService
}

func (app *BaseApp) JobService() *service.JobService {
	return app.jobService
}

func (app *BaseApp) TaskService() *service.TaskService {
	return app.taskService
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

func (app *BaseApp) DataDir() types.DataDir {
	return types.DataDir(app.config.DataDir)
}

func (app *BaseApp) Bootstrap() error {
	var err error

	dataDir := types.DataDir(app.config.DataDir)

	err = utils.CreateDirectories([]string{
		dataDir.Users(),
		dataDir.Playlists(),
		dataDir.Cache(),
		dataDir.Temp(),
	})
	if err != nil {
		return err
	}

	app.db, err = database.Open(dataDir.DatabaseFile())
	if err != nil {
		return err
	}

	// TODO(patrik): Should this be in Bootstrap()?
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

	app.taskService = service.NewTaskService(newServiceLogger("task"))

	app.jobService = service.NewJobService(
		newServiceLogger("job"),
		app.db,
	)

	app.imageService = service.NewImageService(
		newServiceLogger("image"),
		app.db,
		dataDir,
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
		dataDir,
		app.imageService,
	)

	app.searchService = service.NewSearchService(
		newServiceLogger("search"),
		app.db,
		dataDir,
		app.config,
	)

	app.mediaService = service.NewMediaService(
		newServiceLogger("media"),
		app.db,
		dataDir,
		app.config,
	)

	app.libraryService = service.NewLibraryService(
		newServiceLogger("library"),
		app.db,
		dataDir,
		app.config,
		app.notificationService,
		app.mediaService,
	)

	app.artistService = service.NewArtistService(
		newServiceLogger("artist"),
		app.db,
		app.imageService,
		dataDir,
	)

	app.albumService = service.NewAlbumService(
		newServiceLogger("album"),
		app.db,
		app.imageService,
		dataDir,
	)

	app.trackService = service.NewTrackService(
		newServiceLogger("track"),
		app.db,
	)

	app.playlistService = service.NewPlaylistService(
		newServiceLogger("playlist"),
		app.db,
		dataDir,
		app.imageService,
	)

	app.historyService = service.NewHistoryService(
		newServiceLogger("track_history"),
		app.db,
	)

	app.queueService = service.NewQueueService(
		newServiceLogger("queue"),
		app.db,
	)

	app.broker = broker.NewBroker(func() []broker.Event {
		return []broker.Event{
			app.libraryService.GetSyncStateEvent(),
			app.taskService.GetSyncStateEvent(),
		}
	})

	app.libraryService.SetUpdateFunc(func() {
		app.broker.EmitEvent(app.libraryService.GetSyncStateEvent())
	})

	app.taskService.SetUpdateFunc(func() {
		app.broker.EmitEvent(app.taskService.GetSyncStateEvent())
	})

	taskList := []service.Task{
		tasks.NewLibrarySyncTask(app.libraryService),
		tasks.NewSearchIndexTask(app.searchService),
		tasks.NewUserStatsRecalculateTask(app.userService, app.jobService),
		tasks.NewAuthCleanupTask(app.authService),
		tasks.NewCacheCleanupTask(dataDir),
		tasks.NewLibraryCleanupTask(app.libraryService),
	}

	for _, task := range taskList {
		err = app.taskService.AddTask(task)
		if err != nil {
			return err
		}
	}

	app.jobService.RegisterJob(
		tasks.GeneratePlaylistImage,
		func(ctx context.Context, data string) error {
			var params service.GeneratePlaylistImageParams
			err := json.Unmarshal([]byte(data), &params)
			if err != nil {
				return err
			}

			return app.PlaylistService().GeneratePlaylistImage(ctx, params)
		},
	)

	app.jobService.RegisterJob(
		tasks.UserStatsUpdate,
		func(ctx context.Context, data string) error {
			var params service.UpdateUserStatsParams
			err := json.Unmarshal([]byte(data), &params)
			if err != nil {
				return err
			}

			return app.UserService().RecalculateUserStats(ctx, params.UserId)
		},
	)

	return nil
}

func (app *BaseApp) Start() {
	app.jobService.Start()
	app.taskService.Start()
	app.broker.Start()
}

func (app *BaseApp) Shutdown() error {
	app.jobService.Stop()
	app.taskService.Stop()

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
