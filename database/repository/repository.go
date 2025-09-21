package repository

import (
	"context"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
)

var timeout = time.Second * 5

type ItemRepository interface {
	Create(ctx context.Context, tx bun.Tx, item *models.Item) error
	DeleteById(ctx context.Context, tx bun.Tx, identifier string) error
	FindLatest(ctx context.Context, collectionID string) (*models.Item, error)
	FindById(ctx context.Context, identifier string) (*models.Item, error)
}

// Repository is the main interface for accessing all repositories.
// Some methods might require a transaction (bun.Tx) to be passed in if they modify data.
type Repository interface {
	Items() ItemRepository
}

type repo struct {
	mu sync.Mutex
	db *bun.DB

	items ItemRepository
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

func withLock[T any](r *repo, fn func(*repo) T) T {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fn(r)
}

var _ Repository = (*repo)(nil)
