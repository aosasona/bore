package bore

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/clipboard"
	"go.trulyao.dev/bore/v2/pkg/events"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type (
	Bore struct {
		// connection is the database connection used by this bore instance
		db *bun.DB

		// config holds the configuration for this bore instance
		config *Config

		// clipboard is the native clipboard interface for the current platform
		clipboard clipboard.NativeClipboard

		// events is the event manager for this bore instance
		manager *events.Manager

		// repository is the interface for accessing database operations
		repository repository.Repository
	}
)

var (
	ErrInvalidArgs         = errors.New("invalid arguments provided to New function")
	ErrStoragePathRequired = errors.New("storage path is required")
)

// New creates a new Bore instance with the provided configuration.
func New(config *Config) (*Bore, error) {
	if config == nil {
		return nil, ErrInvalidArgs
	}

	if strings.TrimSpace(config.DataDir) == "" {
		return nil, ErrStoragePathRequired
	}

	if err := os.MkdirAll(config.DataDir, 0o755); err != nil {
		return nil, errors.New("failed to create data directory: " + err.Error())
	}

	conn, err := database.Connect(config.DataDir)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	clipboard, _ := clipboard.NewNativeClipboard()
	repository := repository.NewRepository(conn)

	return &Bore{
		db:         conn,
		config:     config,
		clipboard:  clipboard,
		repository: repository,
		manager:    events.NewManager(conn, repository),
	}, nil
}

// Repository returns the current repository implementation for the current Bore instance.
func (b *Bore) Repository() (repository.Repository, error) {
	if b.repository == nil {
		return nil, errors.New("repository is not initialized")
	}

	return b.repository, nil
}

// SystemClipboard returns the native clipboard interface for the current platform.
func (b *Bore) SystemClipboard() (clipboard.NativeClipboard, error) {
	if b.clipboard == nil {
		return nil, errors.New("clipboard is not initialized")
	}

	return b.clipboard, nil
}

func (b *Bore) DB() (*bun.DB, error) {
	if b.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	return b.db, nil
}

// Config returns the configuration of the Bore instance.
func (b *Bore) Config() (*Config, error) {
	if b.config == nil {
		return nil, errors.New("configuration is not initialized")
	}

	return b.config, nil
}

type CopyOptions struct {
	CollectionID string // Optional collection ID to associate with the copied item.
	Mimetype     mimetype.MimeType
}

// Copy copies the provided data to the Bore instance.
// TODO: implement database op and optionally use system clipbpard
func (b *Bore) Copy(ctx context.Context, data []byte, opts *CopyOptions) error {
	if b.clipboard.Available() && b.config.ClipboardPassthrough {
		if err := b.clipboard.Write(ctx, data); err != nil {
			return err
		}
	}

	// TODO: fill this properly
	e, err := events.NewWithGeneratedID(
		aggregate.AggregateTypeItem,
		&payload.CreateItem{
			Content:      data,
			Mimetype:     opts.Mimetype,
			CollectionID: opts.CollectionID,
		},
	)
	if err != nil {
		return errors.New("failed to create copy event: " + err.Error())
	}

	_ = e

	// b.manager.Apply(ctx, )
	panic("not implemented")
}

// Paste retrieves the last copied data from the Bore instance.
// TODO: implement database op and optionally use system clipbpard
func (b *Bore) Paste(ctx context.Context) ([]byte, error) {
	if !b.clipboard.Available() {
		return nil, errors.New("clipboard is not available on this platform")
	}

	return b.clipboard.Read(ctx)
}

func (b *Bore) Close() error {
	if err := b.db.Close(); err != nil {
		return errors.New("failed to close database connection: " + err.Error())
	}

	return nil
}

func (b *Bore) Reset() error {
	if err := os.RemoveAll(b.config.DataDir); err != nil {
		return errors.New("failed to remove data directory: " + err.Error())
	}

	return nil
}
