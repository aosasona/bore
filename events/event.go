package events

import (
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/models"
)

// TODO: add device registration event

// Timestamp holds the timestamps for different stages of an event's lifecycle.
type Timestamp struct {
	// LoggedAt is the time when the event was logged on the source device.
	LoggedAt time.Time `json:"logged_at"`

	// IngestedAt is the time when the event was ingested into the database on the current device.
	IngestedAt time.Time `json:"ingested_at"`

	// AppliedAt is the time when the event was applied to the database, e.g. when the actual database operation like creating a clip was performed.
	AppliedAt time.Time `json:"applied_at"`
}

// Metadata hold data about the event such as the device ID, aggregate ID, version, and timestamp.
type Metadata struct {
	// DeviceID is the unique identifier for the device that generated the event.
	DeviceID string `json:"device_id"`

	// AggregateID is the unique identifier for the aggregate/entity that the event belongs to.
	AggregateID string `json:"aggregate_id"`

	// Version is the version of the event schema.
	Version int `json:"version"`
}

type Log struct {
	// Metadata holds the metadata about the event itself
	Metadata Metadata `json:"metadata"`

	// Timestamp holds the timestamps for different stages of the event's lifecycle.
	Timestamp Timestamp `json:"timestamp"`
}

type Event interface {
	// Action returns the type of the event as a string.
	Action() models.Action

	// Apply executes the appropriate database operations for this event and returns a proper Log that can be saved to the events log and replayed later, or streamed to other devices.
	Apply(db *bun.DB) (Log, error)

	// Replay replays a log entry from a previously saved event log.
	// TODO: this should be on the event manager itself
	Replay(db *bun.DB) error
}
