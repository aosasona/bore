package events

import (
	"encoding/json"
	"time"

	"go.trulyao.dev/bore/v2/database/repository"
)

// TODO: add device registration event

// Timestamp holds the timestamps for different stages of an event's lifecycle.
type Timestamp struct {
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
	// Action is the name of the event, e.g. "CreateClipEvent"
	Action repository.Action `json:"action"`

	// Metadata holds the metadata about the event itself
	Metadata Metadata `json:"metadata"`

	// Timestamp holds the timestamps for different stages of the event's lifecycle.
	Timestamp Timestamp `json:"timestamp"`
}

type Event interface {
	json.Marshaler
	json.Unmarshaler

	// Action returns the type of the event as a string.
	Action() repository.Action

	// Apply executes the appropriate database operations for this event and returns a proper Log that can be saved to the events log and replayed later, or streamed to other devices.
	Apply(repository.Repository) (Log, error)

	// Play replays the event from a log entry.
	// For example, if the event is a CreateClipEvent, it will create the clip in the database with the same ID and data as in the log instead of generating a new ID.
	Play(repository.Repository, Log) error
}
