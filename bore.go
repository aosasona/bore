package bore

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/clipboard"
	"go.trulyao.dev/bore/v2/pkg/events"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
	"go.trulyao.dev/bore/v2/pkg/lib"
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
	Passthrough  bool   // Whether to also copy to the system clipboard if available.
	CollectionID string // Optional collection ID to associate with the copied item.
	Mimetype     mimetype.MimeType
}

// Copy copies the provided data to the Bore instance.
func (b *Bore) Copy(ctx context.Context, data []byte, opts CopyOptions) error {
	forwardToSystemClipboard := b.config.ClipboardPassthrough || opts.Passthrough
	if b.clipboard.Available() && forwardToSystemClipboard {
		if err := b.clipboard.Write(ctx, data); err != nil {
			return err
		}
	}

	hash := lib.ComputeChecksum(data)
	existingItem, err := b.repository.Items().FindByHash(ctx, hash, opts.CollectionID)
	if err != nil {
		return errors.New("failed to check for existing item: " + err.Error())
	}

	var e *events.Event
	if existingItem != nil {
		var existingAgg aggregate.Aggregate
		existingAgg, err = aggregate.NewWithID(aggregate.AggregateTypeItem, existingItem.ID)
		if err != nil {
			return errors.New("failed to create aggregate for existing item: " + err.Error())
		}

		e, err = events.New(existingAgg, &payload.BumpItem{})
	} else {
		e, err = events.NewWithGeneratedID(
			aggregate.AggregateTypeItem,
			&payload.CreateItem{
				Content:      data,
				Mimetype:     opts.Mimetype,
				CollectionID: opts.CollectionID,
			},
		)
	}
	if err != nil {
		return errors.New("failed to create copy event: " + err.Error())
	}

	if _, _, err = b.manager.Apply(ctx, e, events.AppendOptions{ExpectedVersion: -1}); err != nil {
		return errors.New("failed to apply copy event: " + err.Error())
	}

	return nil
}

type PasteOptions struct {
	ItemID              string // Optional item identifier to filter pasted items.
	CollectionID        string // Optional collection ID to filter pasted items.
	FromSystemClipboard bool   // Whether to paste from the system clipboard if available.
	DeleteAfterPaste    bool   // Whether to delete the pasted item after pasting.
	SkipCollectionCheck bool   // Whether to skip checking if the collection exists.
}

type PasteResult struct {
	Content []byte
	Item    *models.Item
}

// Paste retrieves the last copied data from the Bore instance.
func (b *Bore) Paste(ctx context.Context, options PasteOptions) (PasteResult, error) {
	if b.clipboard.Available() && options.FromSystemClipboard {
		rawContent, err := b.clipboard.Read(ctx)
		if err != nil {
			return PasteResult{}, errors.New("failed to read from system clipboard: " + err.Error())
		}

		return PasteResult{Content: rawContent, Item: nil}, nil
	}

	options.CollectionID = strings.TrimSpace(options.CollectionID)
	if options.CollectionID != "" && !options.SkipCollectionCheck {
		exists, err := b.repository.Collections().Exists(ctx, options.CollectionID)
		if err != nil {
			return PasteResult{}, errors.New("failed to check collection existence: " + err.Error())
		} else if !exists {
			return PasteResult{}, errors.New("requested collection does not exist")
		}
	}

	var (
		item *models.Item
		err  error
	)

	identifier := strings.TrimSpace(options.ItemID)
	if identifier == "" {
		item, err = b.repository.Items().FindLatest(ctx, options.CollectionID)
	} else {
		item, err = b.repository.Items().FindById(ctx, identifier, options.CollectionID)
	}

	if err != nil {
		return PasteResult{}, errors.New("failed to find latest item: " + err.Error())
	}

	if item == nil {
		return PasteResult{}, nil
	}

	if options.DeleteAfterPaste {
		agg, err := aggregate.NewWithID(aggregate.AggregateTypeItem, item.ID)
		if err != nil {
			return PasteResult{}, errors.New(
				"failed to create aggregate for deletion: " + err.Error(),
			)
		}

		e, err := events.New(agg, &payload.DeleteItem{})
		if err != nil {
			return PasteResult{}, errors.New("failed to create delete event: " + err.Error())
		}

		if _, _, err = b.manager.Apply(ctx, e); err != nil {
			return PasteResult{}, errors.New("failed to apply delete event: " + err.Error())
		}
	}

	return PasteResult{Content: item.Content, Item: item}, nil
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
