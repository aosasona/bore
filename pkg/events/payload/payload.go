package payload

import (
	"context"
	"encoding/json"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
)

type ProjectionOptions struct {
	Sequence  int64               // The sequence number of the event being projected.
	Aggregate aggregate.Aggregate // The aggregate associated with the event.
}

type Payload interface {
	json.Marshaler
	Type() action.Action
	ApplyProjection(ctx context.Context, tx bun.Tx, options *ProjectionOptions) error
}

type RawPayload interface {
	[]byte | json.RawMessage
}

// Decode decodes the given JSON data into the provided payload type.
func Decode[T Payload, P RawPayload](data P, payload T) (T, error) {
	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, err
	}
	return payload, nil
}
