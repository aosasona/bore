package repository

import (
	"context"
	"sync"
	"time"

	"github.com/uptrace/bun"
)

var timeout = time.Second * 5

type ItemRepository interface {
	Create(ctx context.Context, clip *Item) error
	FindLatest(ctx context.Context, collectionID string) (*Item, error)
	FindById(ctx context.Context, identifier string) (*Item, error)
	DeleteById(ctx context.Context, identifier string) error
}

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
			r.items = &itemsRepository{db: r.db}
		}
		return r.items
	})
}

func withContext(ctx context.Context) (context.Context, func()) {
	return context.WithTimeout(ctx, timeout)
}

func withLock[T any](r *repo, fn func(*repo) T) T {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fn(r)
}

var _ Repository = (*repo)(nil)
