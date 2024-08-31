package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/sql/migrations"
)

func main() {
	config.CreateDirIfNotExists(config.DefaultDataDir())

	// Run migrations
	if err := migrations.Migrate(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintf(
			os.Stderr,
			"Failed to perform database migration, your database may be corrupt, please report if you think this is a bug.\nError: %s\n",
			err.Error(),
		)
		os.Exit(1)
	}

	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
