package bore

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/clipboard"
	"go.trulyao.dev/bore/v2/database"
)

type (
	Bore struct {
		// deviceId is the unique identifier for this device
		// TODO: init device ID on first run
		// TODO: load device ID from `.device_id` file
		deviceId string

		// connection is the database connection used by this bore instance
		db *bun.DB

		// config holds the configuration for this bore instance
		config *Config

		// clipboard is the native clipboard interface for the current platform
		clipboard clipboard.NativeClipboard
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

	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, errors.New("failed to create data directory: " + err.Error())
	}

	conn, err := database.Connect(config.DataDir)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	clipboard, err := clipboard.NewNativeClipboard()
	if err != nil {
		return nil, errors.New("failed to create native clipboard: " + err.Error())
	}

	return &Bore{
		db:        conn,
		config:    config,
		clipboard: clipboard,
	}, nil
}

func (b *Bore) DeviceID() string {
	return b.deviceId
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

// Copy copies the provided data to the Bore instance.
func (b *Bore) Copy(ctx context.Context, data []byte) error {
	if !b.clipboard.Available() {
		return errors.New("clipboard is not available on this platform")
	}

	return b.clipboard.Write(ctx, data)
}

// Paste retrieves the last copied data from the Bore instance.
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
