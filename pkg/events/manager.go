package events

import "github.com/uptrace/bun"

// Manager handles event sourcing operations.
type Manager struct {
	db *bun.DB
}

func NewManager(db *bun.DB) *Manager {
	return &Manager{db: db}
}
