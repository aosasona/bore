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

	// LoggedAt is the time when the event was logged on the source device.
	LoggedAt int64 `json:"logged_at"`

	// IngestedAt is the time when the event was ingested into the database on the current device.
	IngestedAt int64 `json:"ingested_at"`

	// AppliedAt is the time when the event was applied to the database, e.g. when the actual database operation like creating a clip was performed.
	AppliedAt int64 `json:"applied_at"`
}

type Event interface {
	// Type returns the type of the event as a string.
	Type() string

	// Apply replays an event against the provided database connection.
	Apply(db *bun.DB) error
}
