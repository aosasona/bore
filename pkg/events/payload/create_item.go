package payload

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/lib"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type CreateItem struct {
	Content      []byte            `json:"content"`
	Mimetype     mimetype.MimeType `json:"mimetype"`
	CollectionID string            `json:"collection_id"`
}

// ApplyProjection implements Payload.
func (c *CreateItem) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options *ProjectionOptions,
) error {
	if options == nil {
		return errors.New("options cannot be nil")
	} else if !options.Aggregate.IsValid() {
		return errors.New("invalid aggregate")
	}

	hash := lib.ComputeChecksum(c.Content)

	row := models.Item{
		ID:                    options.Aggregate.ID(),
		Content:               c.Content,
		Hash:                  hash,
		Mimetype:              c.Mimetype.String(),
		LastAppliedSequenceID: options.Sequence,
		CollectionID:          sql.NullString{String: c.CollectionID, Valid: c.CollectionID != ""},
	}

	return repo.Items().Create(ctx, tx, &row)
}

// Type implements Payload.
func (c *CreateItem) Type() action.Action {
	return action.ActionCreateItem
}

var _ Payload = (*CreateItem)(nil)
