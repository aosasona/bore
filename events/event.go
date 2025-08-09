package events

import (
	"github.com/uptrace/bun"
)

// TODO: add device registration event

// Metadata hold data about the event such as the device ID, aggregate ID, version, and timestamp.
type Metadata struct {
	// DeviceID is the unique identifier for the device that generated the event.
	DeviceID string `json:"device_id"`

	// AggregateID is the unique identifier for the aggregate/entity that the event belongs to.
	AggregateID string `json:"aggregate_id"`

	// Version is the version of the event schema.
	Version int `json:"version"`

	// Timestamp is the time when the event was created.
	Timestamp int64 `json:"timestamp"`
}

type Event interface {
	// Type returns the type of the event as a string.
	Type() string

	// Apply replays an event against the provided database connection.
	Apply(db *bun.DB) error
}
