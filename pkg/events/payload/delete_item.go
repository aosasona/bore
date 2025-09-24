package payload

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/events/action"
)

type DeleteItem struct {
	ID string `json:"item_id"`
}

// ApplyProjection implements Payload.
func (d *DeleteItem) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options *ProjectionOptions,
) error {
	panic("unimplemented")
}

// Type implements Payload.
func (d *DeleteItem) Type() action.Action {
	panic("unimplemented")
}

var _ Payload = (*DeleteItem)(nil)
