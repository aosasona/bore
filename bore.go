package bore

import (
	"errors"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database"
)

type (
	Bore struct {
		// connection is the database connection used by this bore instance
		db *bun.DB

		// config holds the configuration for this bore instance
		config *Config
	}
)

var (
	ErrInvalidArgs         = errors.New("invalid arguments provided to New function")
	ErrStoragePathRequired = errors.New("storage path is required")
)

// New creates a new Bore instance with the provided configuration.
// TODO: add native clipboard in Bore struct
func New(config *Config) (*Bore, error) {
	if config == nil {
		return nil, ErrInvalidArgs
	}

	if strings.TrimSpace(config.DataDir) == "" {
		return nil, ErrStoragePathRequired
	}

	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, errors.New("failed to create data directory: " + err.Error())
	}

	conn, err := database.Connect(config.DataDir)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	return &Bore{
		db:     conn,
		config: config,
	}, nil
}

func (b *Bore) DB() *bun.DB {
	if b.db == nil {
		panic("database connection is not initialized")
	}

	return b.db
}

// Config returns the configuration of the Bore instance.
func (b *Bore) Config() *Config {
	if b.config == nil {
		panic("configuration is not initialized")
	}

	return b.config
}
