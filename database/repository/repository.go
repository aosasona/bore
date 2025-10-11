package repository

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

var ErrEmptyIdentifier = errs.New("identifier cannot be empty")

type ItemRepository interface {
	Create(ctx context.Context, tx bun.Tx, item *models.Item) error
	Bump(
		ctx context.Context,
		tx bun.Tx,
		identifier string,
		sequenceId int64,
	) error // Bump updates the sequence ID and updated_at timestamp of an item to move it to the top of the list.
	DeleteById(ctx context.Context, tx bun.Tx, identifier string) error

	FindLatest(ctx context.Context, collectionID string) (*models.Item, error)
	FindById(ctx context.Context, identifier string, collectionId string) (*models.Item, error)
	FindByHash(ctx context.Context, hash string, collectionId string) (*models.Item, error)
}

type CollectionLookupOptions struct {
	Identifier string
	Name       string
}

type OrderBy struct {
	Field     string
	Ascending bool
}

type Pagination struct {
	Limit  int
	Offset int
}

type FindAllOptions struct {
	OrderBy    []OrderBy
	Pagination *Pagination
}

type CollectionRepository interface {
	Create(ctx context.Context, tx bun.Tx, collection *models.Collection) error
	Rename(ctx context.Context, tx bun.Tx, identifier string, newName string) error
	DeleteById(ctx context.Context, tx bun.Tx, identifier string) error

	FindById(ctx context.Context, identifier string) (*models.Collection, error)
	FindByName(ctx context.Context, name string) (*models.Collection, error)
	// FindOne looks up a collection by either ID or name.
	FindOne(ctx context.Context, opts CollectionLookupOptions) (*models.Collection, error)
	FindAll(ctx context.Context, opts FindAllOptions) (models.Collections, error)
}

// Repository is the main interface for accessing all repositories.
// Some methods might require a transaction (bun.Tx) to be passed in if they modify data.
type Repository interface {
	Items() ItemRepository
	Collections() CollectionRepository
}

type repo struct {
	mu sync.Mutex
	db *bun.DB

	items       ItemRepository
	collections CollectionRepository
}

func NewRepository(db *bun.DB) Repository {
	return &repo{db: db}
}

// Items implements Repository.
func (r *repo) Items() ItemRepository {
	return withLock(r, func(r *repo) ItemRepository {
		if r.items == nil {
			r.items = &itemRepository{db: r.db}
		}
		return r.items
	})
}

// Collections implements Repository.
func (r *repo) Collections() CollectionRepository {
	return withLock(r, func(r *repo) CollectionRepository {
		if r.collections == nil {
			r.collections = &collectionRepository{db: r.db}
		}
		return r.collections
	})
}

func withLock[T any](r *repo, fn func(*repo) T) T {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fn(r)
}

var _ Repository = (*repo)(nil)
