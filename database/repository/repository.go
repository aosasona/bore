package repository

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
)

type ItemRepository interface {
	Create(ctx context.Context, tx bun.Tx, item *models.Item) error
	DeleteById(ctx context.Context, tx bun.Tx, identifier string) error
	FindLatest(ctx context.Context, collectionID string) (*models.Item, error)
	FindById(ctx context.Context, identifier string) (*models.Item, error)
}

type CollectionRepository interface {
	Exists(ctx context.Context, identifier string) (bool, error)
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
