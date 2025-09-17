package repository

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

// TODO: move to event sourcing package
type Action string

const (
	ActionCopyV1           Action = "copy_v1"
	ActionCreateCollection Action = "create_collection"
	ActionDeleteClip       Action = "delete_clip"
	ActionDeleteCollection Action = "delete_collection"
)

type Collection struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:collections,alias:co"`

	ID string `bun:"id,pk"`
	// TODO: fill
}

// BeforeInsert implements bun.BeforeInsertHook.
func (collection *Collection) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if collection.ID == "" {
		collection.ID = ulid.Make().String()
	}

	if err := collection.Validate(); err != nil {
		return err
	}

	return nil
}

type Item struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:items,alias:cl"`

	ID      string `bun:"id,pk"`
	Content []byte `bun:"content,notnull"`
}
