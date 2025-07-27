package migrations

import (
	"embed"
	"log/slog"

	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	slog.Debug("Initializing migrations")
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic("Failed to discover migrations: " + err.Error())
	}
	slog.Debug("Migrations ran successfully")
}
