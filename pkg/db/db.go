package db

import (
	csha256 "crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	guuid "github.com/google/uuid"
	"github.com/mattes/migrate"
	"github.com/mattn/go-sqlite3"
	"go.trulyao.dev/bore/sql/migrations"
)

func uuid() string {
	return guuid.New().String()
}

func sha256(s string) string {
	c := csha256.New()
	c.Write([]byte(s))
	return fmt.Sprintf("%x", c.Sum(nil))
}

func GetDSN(dataDir string) string {
	path := filepath.Join(dataDir, "data.db")
	return "file:" + path + "?_foreign_keys=on&mode=rwc&_journal_mode=WAL&cache=shared"
}

func Connect(dataDir string) (*sql.DB, error) {
	if !slices.Contains(sql.Drivers(), "sqlite3_custom") {
		// Add UUID function
		sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// UUID function
				if err := conn.RegisterFunc("uuid", uuid, false); err != nil {
					return err
				}

				// SHA256 function
				if err := conn.RegisterFunc("sha256", sha256, false); err != nil {
					return err
				}

				return nil
			},
		})
	}

	// Run migrations
	err := migrations.Migrate(filepath.Join(dataDir, "data.db"))
	if err != nil && err.Error() != migrate.ErrNoChange.Error() {
		fmt.Fprintln(
			os.Stderr,
			"Failed to perform database migration, your database may be corrupt, please report if you think this is a bug.",
		)

		return nil, err
	}

	// Establish connection
	conn, err := sql.Open("sqlite3_custom", GetDSN(dataDir))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
