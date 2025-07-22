package repository

import (
	"sync"

	"github.com/uptrace/bun"
)

type baseRepository struct {
	mu sync.Mutex
	db *bun.DB
}

func New(db *bun.DB) *baseRepository {
	return &baseRepository{
		db: db,
	}
}

// withLock is a helper function to execute a function with a lock on the repository.
func (b *baseRepository) withLock(f func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	f()
}
