package core

import (
	"context"
	"os"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/broker"
	"github.com/nanoteck137/dwebble/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db     *database.Database
	config *config.Config

	authService    *service.AuthService
	searchService  *service.SearchService
	libraryService *service.LibraryService
	imageService   *service.ImageService
	mediaService   *service.MediaService

	broker *broker.Broker
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

	dirs := []string{
		workDir.Artists(),
		workDir.Albums(),
		workDir.Tracks(),
		workDir.Playlists(),
		workDir.Trash(),
		workDir.Cache().String(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
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

	app.authService = service.NewAuthService(app.db, app.config)
	// TODO(patrik): This should be a worker
	go app.authService.CleanRoutine()

	app.searchService = service.NewSearchService(app.db, app.config)

	// TODO(patrik): Do this lazily
	err = app.searchService.Init()
	if err != nil {
		return err
	}

	app.libraryService = service.NewLibraryService(app.db, app.config, app.searchService)

	app.imageService = service.NewImageService(app.db, app.config.WorkDir())

	app.mediaService = service.NewMediaService(app.db, app.config.WorkDir())

	app.broker = broker.NewBroker(func() []broker.Event {
		return []broker.Event{
			app.libraryService.GetSyncStateEvent(),
		}
	})

	// TODO(patrik): Move to worker?
	go app.broker.Listen()

	app.libraryService.SetUpdateFunc(func() {
		app.broker.EmitEvent(app.libraryService.GetSyncStateEvent())
	})

	// TODO(patrik): Remove test code
	playlists, err := app.DB().GetAllPlaylists(context.TODO())
	if err != nil {
		return err
	}

	for _, playlist := range playlists {
		err = app.DB().DeletePlaylist(context.TODO(), playlist.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *BaseApp) Shutdown() error {
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
