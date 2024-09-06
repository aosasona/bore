package app

import (
	"database/sql"

	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/daos"
	"go.trulyao.dev/bore/pkg/db"
	"go.trulyao.dev/bore/pkg/handler"
	"go.trulyao.dev/bore/pkg/system"
)

type App struct {
	config          *config.Config
	db              *sql.DB
	nativeClipboard system.NativeClipboardInterface
}

func New(configPath string) (*App, error) {
	a := new(App)

	// Load the configuration
	conf, err := config.Load(configPath)
	if err != nil {
		return a, err
	}
	a.config = conf

	// Initialize the database connection
	db, err := db.Connect(conf.DataDir)
	if err != nil {
		return a, err
	}
	a.db = db

	a.nativeClipboard, _ = system.NewNativeClipboard()

	return a, nil
}

func (a *App) Daos() *daos.Queries {
	return daos.New(a.db)
}

func (a *App) Handler() handler.HandlerInterface {
	return handler.New(a.Daos(), a.config, a.nativeClipboard)
}

func (a *App) loadNativeClipboard() {
}

func (a *App) UpdateConfigPath(configPath string) error {
	_ = a.db.Close()

	var err error

	if a.config, err = config.Load(configPath); err != nil {
		return err
	}

	if a.db, err = db.Connect(a.config.DataDir); err != nil {
		return err
	}

	a.nativeClipboard, _ = system.NewNativeClipboard()

	return nil
}
