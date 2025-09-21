package models

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type Collection struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:collections,alias:co"`

	ID        string    `bun:"id,pk"`
	Name      string    `bun:"name,notnull"                                          validate:"required,alphanumunicode"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (collection *Collection) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if collection.ID == "" {
		return errors.New("ID is required")
	}

	if err := collection.Validate(); err != nil {
		return err
	}

	return nil
}
