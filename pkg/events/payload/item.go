package payload

import (
	"context"
	"encoding/json"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

// TODO: implement fields
type CreateItem struct {
	Content      string            `json:"content"`
	Mimetype     mimetype.MimeType `json:"mimetype"`
	CollectionID string            `json:"collection_id"`
}

// ApplyProjection implements Payload.
func (c *CreateItem) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	options *ProjectionOptions,
) error {
	panic("unimplemented")
}

// Type implements Payload.
func (c *CreateItem) Type() action.Action {
	return action.ActionCreateItem
}

// MarshalJSON implements Payload.
func (c *CreateItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{})
}

var _ Payload = (*CreateItem)(nil)
