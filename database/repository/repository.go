package repository

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
)

type ClipRepository interface {
	GetLastClip(context.Context) (Clip, error)
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

func withLock[T any](r *repo, fn func(*repo) T) T {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fn(r)
}

var _ Repository = (*repo)(nil)
