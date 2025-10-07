package payload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
)

type ProjectionOptions struct {
	Sequence  int64               // The sequence number of the event being projected.
	Aggregate aggregate.Aggregate // The aggregate associated with the event.
}

type Payload interface {
	Type() action.Action
	ApplyProjection(
		ctx context.Context,
		tx bun.Tx,
		repo repository.Repository,
		options ProjectionOptions,
	) error
}

type RawPayload interface {
	[]byte | json.RawMessage
}

// Decode decodes the given JSON data into an associated Payload struct based on the action type.
func Decode[P RawPayload](data P, a action.Action) (Payload, error) {
	var target Payload

	switch a {
	case action.ActionCreateCollection:
		target = new(CreateCollection)

	case action.ActionDeleteCollection:
		target = new(DeleteCollection)

	case action.ActionCreateItem:
		target = new(CreateItem)

	case action.ActionDeleteItem:
		target = new(DeleteItem)

	case action.ActionBumpItem:
		target = new(BumpItem)

	default:
		return nil, errs.New(fmt.Sprintf("unknown event action: %s", a))
	}

	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}

	return target, nil
}
