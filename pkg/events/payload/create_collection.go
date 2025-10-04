package payload

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/events/action"
)

type CreateCollection struct{}

// ApplyProjection implements Payload.
func (c *CreateCollection) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options ProjectionOptions,
) error {
	if !options.Aggregate.IsValid() {
		return errors.New("invalid aggregate")
	}

	// TODO: implement
	panic("unimplemented")
}

// Type implements Payload.
func (c *CreateCollection) Type() action.Action {
	return action.ActionCreateCollection
}

var _ Payload = (*CreateCollection)(nil)
