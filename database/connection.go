package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"go.trulyao.dev/bore/v2/migrations"
)

const DriverCustom = "sqlite3_custom"

func isDev() bool {
	return os.Getenv("APP_ENV") == "development" || strings.HasPrefix(os.Args[0], "go")
}

func Connect(dataDir string) (*bun.DB, error) {
	connection, err := createDbConnection(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	db := bun.NewDB(connection, sqlitedialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithEnabled(isDev()),
		),
	)

	ctx := context.Background()
	if err = runPragmas(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to run PRAGMAs: %w", err)
	}

	if err = runMigration(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func dsn(dataDir string) string {
	path := filepath.Join(dataDir, "data.db")
	return "file:" + path + "?_foreign_keys=on&mode=rwc&_journal_mode=WAL&cache=shared"
}

func createDbConnection(dataDir string) (*sql.DB, error) {
	connection, err := sql.Open(sqliteshim.ShimName, dsn(dataDir))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return connection, nil
}

func runMigration(db *bun.DB) error {
	ctx := context.Background()
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Init(ctx); err != nil {
		return err
	}

	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx)

	groups, err := migrator.Migrate(ctx)
	if err != nil {
		fmt.Printf("Applied migrations:\n")
		for _, group := range groups.Migrations.Applied() {
			fmt.Printf("- %s\n", group.String())
		}
		return err
	}

	if groups.IsZero() {
		return nil
	}

	fmt.Printf("migrated %d groups\n", len(groups.Migrations.Applied()))
	return nil
}

func runPragmas(ctx context.Context, tx bun.IDB) error {
	if _, err := tx.NewRaw("PRAGMA busy_timeout = 10000;").Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewRaw("PRAGMA foreign_keys = ON;").Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewRaw("PRAGMA journal_mode = WAL;").Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewRaw("PRAGMA synchronous = NORMAL;").Exec(ctx); err != nil {
		return err
	}

	return nil
}
