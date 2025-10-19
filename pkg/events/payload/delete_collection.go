package payload

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events/action"
)

type DeleteCollection struct{}

// ApplyProjection implements Payload.
func (d *DeleteCollection) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options ProjectionOptions,
) error {
	if !options.Aggregate.IsValid() {
		return errs.New("invalid aggregate")
	}

	return repo.Collections().DeleteById(ctx, tx, options.Aggregate.ID())
}

// Type implements Payload.
func (d *DeleteCollection) Type() action.Action {
	return action.ActionDeleteCollection
}

var _ Payload = (*DeleteCollection)(nil)
