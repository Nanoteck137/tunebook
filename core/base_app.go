package core

import (
	"os"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db     *database.Database
	config *config.Config

	authService    *service.AuthService
	searchService  *service.SearchService
	libraryService *service.LibraryService
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
		workDir.Trash(),
		workDir.Cache().String(),
		workDir.Cache().Tracks(),
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
	err = app.searchService.Init()
	if err != nil {
		return err
	}

	app.libraryService = service.NewLibraryService(app.db, app.config, app.searchService)
	app.libraryService.Sync()

	// app.searchService.Test()

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
