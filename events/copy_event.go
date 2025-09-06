package events

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/device"
)

// copyEvent is emitted when a copy operation is performed.
type copyEvent struct {
	identity *device.Identity `json:"-" mapstructure:"-"`

	// The content that was copied to the clipboard
	Content []byte `json:"content" mapstructure:"content"`

	// Hash is the SHA256 hash of the content that was copied.
	Hash string `json:"hash" mapstructure:"hash"`

	// MimeType is the MIME type of the content that was copied.
	MimeType MimeType `json:"mime_type" mapstructure:"mime_type"`

	// CollectionID is the identifier of the collection to which this event belongs.
	CollectionID string `json:"collection_id" mapstructure:"collection_id"`
}

// Action implements Event.
func (c *copyEvent) Action() repository.Action {
	return repository.ActionCopyV1
}

// Apply implements Event.
func (c *copyEvent) Apply(repository repository.Repository) (Log, error) {
	panic("unimplemented")
}

// MarshalJSON implements Event.
func (c *copyEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"content":       c.Content,
		"hash":          c.Hash,
		"mime_type":     c.MimeType.String(),
		"collection_id": c.CollectionID,
	})
}

// UnmarshalJSON implements Event.
func (c *copyEvent) UnmarshalJSON(raw []byte) error {
	data := map[string]any{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	if err := mapstructure.Decode(data, c); err != nil {
		return fmt.Errorf("failed to decode copy event: %w", err)
	}

	if len(c.Content) == 0 {
		return errors.New("content cannot be empty")
	}

	if c.Hash == "" {
		return errors.New("hash cannot be empty")
	}

	if c.MimeType == nil {
		return errors.New("mime type cannot be nil")
	}

	if c.CollectionID == "" {
		return errors.New("collection ID cannot be empty")
	}

	return nil
}

// Play implements Event.
func (c *copyEvent) Play(repository.Repository, Log) error {
	panic("unimplemented")
}

var _ Event = (*copyEvent)(nil)
