package events

import "github.com/uptrace/bun"

// copyEvent is emitted when a copy operation is performed.
type copyEvent struct {
	// The content that was copied to the clipboard
	Content []byte `json:"content"`

	// Hash is the SHA256 hash of the content that was copied.
	Hash string `json:"hash"`

	// MimeType is the MIME type of the content that was copied.
	MimeType MimeType `json:"mime_type"`

	// CollectionID is the identifier of the collection to which this event belongs.
	CollectionID string `json:"collection_id"`
}

// Type implements Event.
func (c *copyEvent) Type() string {
	return "copy_v1"
}

// Apply implements Event.
func (c *copyEvent) Apply(db *bun.DB) error {
	panic("unimplemented")
}

// MarshalJSON implements Event.
func (c *copyEvent) MarshalJSON() ([]byte, error) {
	panic("unimplemented")
}

// UnmarshalJSON implements Event.
func (c *copyEvent) UnmarshalJSON([]byte) error {
	panic("unimplemented")
}

var _ Event = (*copyEvent)(nil)
