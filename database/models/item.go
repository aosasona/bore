package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type Item struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:items,alias:i"`

	ID                    string    `bun:"id,pk"                                                 validate:"required,ulid"`
	Content               []byte    `bun:"content,notnull"                                       validate:"required"`
	Hash                  string    `bun:"hash,notnull"                                          validate:"required,sha256"`
	Mimetype              string    `bun:"mimetype,notnull"                                      validate:"required,mimetype"`
	LastAppliedSequenceID int64     `bun:"last_applied_sequence_id,notnull"                      validate:"gte=0"`
	CreatedAt             time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt             time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	CollectionID sql.NullString `bun:"collection_id"`
	Collection   *Collection    `bun:"rel:belongs-to,join:collection_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (item *Item) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if item.ID == "" {
		return errors.New("ID is required")
	}

	if err := item.Validate(); err != nil {
		return err
	}

	return nil
}
