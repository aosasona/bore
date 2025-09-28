package payload

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/events/action"
)

type BumpItem struct{}

// ApplyProjection implements Payload.
func (b *BumpItem) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options ProjectionOptions,
) error {
	if !options.Aggregate.IsValid() {
		return nil
	}

	return repo.Items().Bump(ctx, tx, options.Aggregate.ID(), options.Sequence)
}

// Type implements Payload.
func (b *BumpItem) Type() action.Action {
	return action.ActionBumpItem
}

var _ Payload = (*BumpItem)(nil)
