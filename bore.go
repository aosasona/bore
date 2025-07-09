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
	}
)

var (
	ErrInvalidArgs         = errors.New("invalid arguments provided to New function")
	ErrStoragePathRequired = errors.New("storage path is required")
)

func New(args *Args) (*Bore, error) {
	if args == nil {
		return nil, ErrInvalidArgs
	}

	if strings.TrimSpace(args.DataPath) == "" {
		return nil, ErrStoragePathRequired
	}

	panic("todo")
}
