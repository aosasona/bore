package app

import (
	"database/sql"

	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/daos"
	"go.trulyao.dev/bore/pkg/db"
	"go.trulyao.dev/bore/pkg/handler"
)

type App struct {
	config *config.Config
	db     *sql.DB
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

	return a, nil
}

func (a *App) Daos() *daos.Queries {
	return daos.New(a.db)
}

func (a *App) Handler() handler.HandlerInterface {
	return handler.New(a.Daos())
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

	return nil
}
