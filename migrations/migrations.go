package migrations

import (
	"embed"
	"fmt"

	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	fmt.Println("Running migrations...")
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic("Failed to discover migrations: " + err.Error())
	}
	fmt.Println("Migrations ran successfully")
}
