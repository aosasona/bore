package payload

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
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
	panic("unimplemented")
}

// Type implements Payload.
func (d *DeleteCollection) Type() action.Action {
	panic("unimplemented")
}

var _ Payload = (*DeleteCollection)(nil)
