package repository

import (
	"context"
	"sync"
	"time"

	"github.com/uptrace/bun"
)

var timeout = time.Second * 5

type ClipRepository interface {
	Create(ctx context.Context, clip *Clip) error
	FindLatest(ctx context.Context, collectionID string) (*Clip, error)
	FindById(ctx context.Context, identifier string) (*Clip, error)
	DeleteById(ctx context.Context, identifier string) error
}

type Repository interface {
	Clips() ClipRepository
}

type repo struct {
	mu sync.Mutex
	db *bun.DB

	clips ClipRepository
}

func NewRepository(db *bun.DB) Repository {
	return &repo{db: db}
}

// Clips implements Repository.
func (r *repo) Clips() ClipRepository {
	return withLock(r, func(r *repo) ClipRepository {
		if r.clips == nil {
			r.clips = &clipRepository{db: r.db}
		}
		return r.clips
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
