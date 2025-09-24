package payload

import (
	"context"
	"errors"

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
	options ProjectionOptions,
) error {
	if !options.Aggregate.IsValid() {
		return errors.New("invalid aggregate")
	}

	panic("unimplemented")
}

// Type implements Payload.
func (d *DeleteItem) Type() action.Action {
	return action.ActionDeleteItem
}

var _ Payload = (*DeleteItem)(nil)
