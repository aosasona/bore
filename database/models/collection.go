package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type Collection struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:collections"`

	ID        string    `bun:"id,pk"                                                 json:"id"`
	Name      string    `bun:"name,notnull"                                          json:"name"       validate:"required,collection_name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Items []*Item `bun:"rel:has-many,join:id=collection_id" json:"omitzero"`
}

type CollectionWithItemsCount struct {
	Collection
	ItemsCount int `bun:"items_count"`
}

type Collections []*CollectionWithItemsCount

func (collection *Collection) Aggregate() (aggregate.Aggregate, error) {
	return aggregate.WithID(aggregate.AggregateTypeCollection, collection.ID)
}

// BeforeAppendModel implements schema.BeforeAppendModelHook.
func (collection *Collection) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	if err := collection.Validate(); err != nil {
		return err
	}

	return nil
}

func (collection *Collection) MarshalJSON() ([]byte, error) {
	aggregate, err := aggregate.WithID(aggregate.AggregateTypeCollection, collection.ID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"id":         collection.ID,
		"aggregate":  aggregate,
		"name":       collection.Name,
		"created_at": collection.CreatedAt,
		"updated_at": collection.UpdatedAt,
	})
}

var _ bun.BeforeAppendModelHook = (*Collection)(nil)
