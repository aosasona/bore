package migrations

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed *.sql
var migrations embed.FS

type Migration struct {
	Name    string
	Content string
}

// GetMigrations returns the list of migrations
func GetMigrations() ([]Migration, error) {
	files, err := migrations.ReadDir(".")
	if err != nil {
		return nil, err
	}

	allMigrations := make([]Migration, 0, len(files))
	for _, file := range files {
		content, err := migrations.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		allMigrations = append(allMigrations, Migration{
			Name:    file.Name(),
			Content: string(content),
		})
	}

	return allMigrations, nil
}

func Migrate(databasePath string) error {
	sourceDriver, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	dsn := "sqlite3:" + strings.Repeat(string(os.PathSeparator), 2) + databasePath

	migrator, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	return migrator.Up()
}
