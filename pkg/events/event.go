package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/events/action"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
)

var (
	ErrInvalidEventType = errors.New("invalid event type")
	ErrInvalidAggregate = errors.New("invalid aggregate")
)

// An event represents any action or change that occurs within the system.
type Event struct {
	bun.BaseModel `bun:"table:events,alias:ev"`

	// ID is the unique identifier for the event.
	ID               string               `bun:"event_id,pk"                       json:"event_id"`          // ULID
	Sequence         int64                `bun:"sequence_id,autoincrement,notnull" json:"sequence"`          // The (auto-generated) sequential number of the event.
	Aggregate        *aggregate.Aggregate `bun:"-"                                 json:"aggregate"`         // The target entity of the event.
	AggregateVersion int64                `bun:"aggregate_version"                 json:"aggregate_version"` // The version of the aggregate after the event.
	Type             action.Action        `bun:",type:text"                        json:"type"`              // The type of event.
	Payload          json.RawMessage      `bun:",type:JSON"                        json:"payload"`           // The event payload, stored as JSON.
	OccurredAt       time.Time            `bun:"occurred_at,nullzero,notnull"      json:"occured_at"`        // The timestamp when the event occurred.

	AggregateType string `bun:"aggregate_type,notnull" json:"aggregate_type"` // The type of the aggregate, stored for easier querying.
	AggregateID   string `bun:"aggregate_id,notnull"   json:"aggregate_id"`   // The ID of the aggregate, stored for easier querying.
}

// New creates a new event with the given aggregate, type, and payload.
func New(agg *aggregate.Aggregate, payload payload.Payload) (*Event, error) {
	if !agg.IsValid() {
		return nil, ErrInvalidAggregate
	}

	if !payload.Type().IsValid() {
		return nil, ErrInvalidEventType
	}

	data, err := payload.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:         ulid.Make().String(),
		OccurredAt: time.Now(),
		Aggregate:  agg,
		Type:       payload.Type(),
		Payload:    data,
	}, nil
}

func NewWithGeneratedID(aggType aggregate.AggregateType, payload payload.Payload) (*Event, error) {
	agg, err := aggregate.New(aggType)
	if err != nil {
		return nil, err
	}

	return New(agg, payload)
}

func (e *Event) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if !e.Aggregate.IsValid() {
		return ErrInvalidAggregate
	}

	if !e.Type.IsValid() {
		return ErrInvalidEventType
	}

	if e.ID == "" {
		e.ID = ulid.Make().String()
	}

	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}

	if err := e.SetAggregate(e.Aggregate); err != nil {
		return err
	}

	if e.Sequence != 0 {
		return errors.New("sequence must not be set manually")
	}

	return nil
}

func (e *Event) SetAggregate(agg *aggregate.Aggregate) error {
	if agg == nil {
		return ErrInvalidAggregate
	} else if !agg.IsValid() {
		return ErrInvalidAggregate
	}

	e.Aggregate = agg
	e.AggregateType = agg.Type()
	e.AggregateID = agg.ID()

	return nil
}

func (e *Event) Save(ctx context.Context, tx bun.Tx) error {
	_, err := tx.NewInsert().Model(e).Exec(ctx)
	return err
}
