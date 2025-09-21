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
	Create(ctx context.Context, item *models.Item) error
	FindLatest(ctx context.Context, collectionID string) (*models.Item, error)
	FindById(ctx context.Context, identifier string) (*models.Item, error)
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
