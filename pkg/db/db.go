package db

import (
	"database/sql"
	"path/filepath"

	guuid "github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"go.trulyao.dev/bore/pkg/config"
)

func uuid() string {
	return guuid.New().String()
}

// TODO: stop using default data directory
func GetDSN() string {
	path := filepath.Join(config.DefaultDataDir(), "data.db")
	return "file:" + path + "?_foreign_keys=on&mode=rwc&_journal_mode=WAL&cache=shared"
}

func Connect() (*sql.DB, error) {
	// Add UUID function
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("uuid", uuid, false); err != nil {
				return err
			}

			return nil
		},
	})

	// Establish connection
	conn, err := sql.Open("sqlite3_custom", GetDSN())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
