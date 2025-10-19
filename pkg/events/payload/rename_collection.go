package payload

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
)

type RenameCollection struct {
	NewName string `json:"new_name"`
}

// ApplyProjection implements Payload.
func (r *RenameCollection) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options ProjectionOptions,
) error {
	if options.Aggregate.Type() != aggregate.AggregateTypeCollection.String() {
		return errs.New("invalid aggregate type for RenameCollection projection")
	}

	return repo.Collections().Rename(ctx, tx, options.Aggregate.ID(), r.NewName)
}

// Type implements Payload.
func (r *RenameCollection) Type() action.Action {
	return action.ActionRenameCollection
}

var _ Payload = (*RenameCollection)(nil)
