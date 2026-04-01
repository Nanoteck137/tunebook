package core

import (
	"log/slog"

	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tasks"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
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
	})
	if err != nil {
		return err
	}

	app.db, err = database.Open(dataDir.DatabaseFile())
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
		newServiceLogger("notification"),
		app.config,
	)

	app.taskService = service.NewTaskService(newServiceLogger("task"))

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
	)

	app.albumService = service.NewAlbumService(
		newServiceLogger("album"),
		app.db,
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

	err = app.taskService.AddTask(tasks.NewAuthCleanupTask(app.authService))
	if err != nil {
		return err
	}

	err = app.taskService.AddTask(tasks.NewCacheCleanupTask(dataDir))
	if err != nil {
		return err
	}

	err = app.taskService.AddTask(tasks.NewLibrarySyncTask(app.libraryService))
	if err != nil {
		return err
	}

	err = app.taskService.AddTask(tasks.NewLibraryCleanupTask(app.libraryService))
	if err != nil {
		return err
	}

	err = app.taskService.AddTask(tasks.NewSearchIndexTask(app.searchService))
	if err != nil {
		return err
	}

	// TODO(patrik): This should not be in bootstrap
	app.taskService.DisplayTasks()
	app.taskService.Start()

	// TODO(patrik): This should not be in bootstrap
	go app.broker.Listen()

	return nil
}

func (app *BaseApp) Shutdown() error {
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
