package bore

import (
	"errors"
	"strings"

	"github.com/uptrace/bun"
)

type (
	Bore struct {
		// connection is the database connection used by this bore instance
		db *bun.DB
	}

	Args struct {
		// DataPath is the path to the storage directory.
		DataPath string

		// Config is the configuration for the bore instance.
		Config Config // optional, can be nil
	}
)

var (
	ErrInvalidArgs         = errors.New("invalid arguments provided to New function")
	ErrStoragePathRequired = errors.New("storage path is required")
)

// New creates a new Bore instance with the provided configuration.
// TODO: create data dir if it does not exist
func New(config *Config) (*Bore, error) {
	if config == nil {
		return nil, ErrInvalidArgs
	}

	if strings.TrimSpace(config.DataDir) == "" {
		return nil, ErrStoragePathRequired
	}

	panic("todo")
}
