package events

import (
	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
)

type deleteEvent struct {
	identity *string `json:"-" mapstructure:"-"`

	// Identifier is the unique identifier of the clip to be deleted.
	Identifier string `json:"identifier" mapstructure:"identifier"`
}

// Action implements Event.
func (d *deleteEvent) Action() repository.Action {
	return repository.ActionDeleteClip
}

// Apply implements Event.
func (d *deleteEvent) Apply(db *bun.DB) (Log, error) {
	panic("unimplemented")
}

// Replay implements Event.
func (d *deleteEvent) Replay(db *bun.DB) error {
	panic("unimplemented")
}

var _ Event = (*deleteEvent)(nil)
