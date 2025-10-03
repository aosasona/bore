package bore

import (
	"os"
	"strings"
	"sync"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/clipboard"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events"
)

type Bore struct {
	namespaceMutex sync.Mutex

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

	// MARK: Namespaces
	items *clipboardNamespace
}

// New creates a new Bore instance with the provided configuration.
func New(config *Config) (*Bore, error) {
	if config == nil {
		return nil, errs.ErrInvalidConstructorArg
	}

	if strings.TrimSpace(config.DataDir) == "" {
		return nil, errs.ErrStoragePathRequired
	}

	if err := os.MkdirAll(config.DataDir, 0o755); err != nil {
		return nil, errs.Wrap(err, "failed to create data directory")
	}

	conn, err := database.Connect(config.DataDir)
	if err != nil {
		return nil, errs.ErrFailedToConnectToDB.WithError(err)
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
		return nil, errs.New("repository is not initialized")
	}

	return b.repository, nil
}

// SystemClipboard returns the native clipboard interface for the current platform.
func (b *Bore) SystemClipboard() (clipboard.NativeClipboard, error) {
	if b.clipboard == nil {
		return nil, errs.New("clipboard is not initialized")
	}

	return b.clipboard, nil
}

func (b *Bore) DB() (*bun.DB, error) {
	if b.db == nil {
		return nil, errs.New("database connection is not initialized")
	}

	return b.db, nil
}

// Config returns the configuration of the Bore instance.
func (b *Bore) Config() (*Config, error) {
	if b.config == nil {
		return nil, errs.New("configuration is not initialized")
	}

	return b.config, nil
}

// Clipboard returns the items namespace for managing clipboard items.
func (b *Bore) Clipboard() *clipboardNamespace {
	b.namespaceMutex.Lock()
	defer b.namespaceMutex.Unlock()

	if b.items == nil {
		b.items = &clipboardNamespace{b}
	}

	return b.items
}

func (b *Bore) Close() error {
	if err := b.db.Close(); err != nil {
		return errs.ErrFailedToCloseDB.WithError(err)
	}

	return nil
}

func (b *Bore) Reset() error {
	if err := os.RemoveAll(b.config.DataDir); err != nil {
		return errs.ErrFailedToRemoveDataDir.WithError(err)
	}

	return nil
}
