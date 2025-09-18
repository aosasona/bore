package events

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
)

//go:generate go tool github.com/abice/go-enum --marshal

// ENUM(create_item,delete_item,create_collection,delete_collection)
type Type string

// An event represents any action or change that occurs within the system.
type Event struct {
	bun.BaseModel `bun:"table:events,alias:ev"`

	// ID is the unique identifier for the event.
	ID               string              `bun:"event_id,pk"                       json:"event_id"`          // ULID
	Sequence         int64               `bun:"sequence_id,autoincrement,notnull" json:"sequence"`          // The (auto-generated) sequential number of the event.
	Aggregate        aggregate.Aggregate `bun:",type:text"                        json:"aggregate"`         // The target entity of the event.
	AggregateVersion int64               `bun:"aggregate_version"                 json:"aggregate_version"` // The version of the aggregate after the event.
	Type             Type                `bun:",type:text"                        json:"type"`              // The type of event.
	Payload          json.RawMessage     `bun:",type:JSON"                        json:"payload"`           // The event payload, stored as JSON.
	OccurredAt       time.Time           `bun:"occurred_at,nullzero,notnull"      json:"occured_at"`        // The timestamp when the event occurred.
}
