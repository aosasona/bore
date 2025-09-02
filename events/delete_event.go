package events

import (
	"encoding/json"

	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/device"
)

type deleteClipEvent struct {
	identity *device.Identity `json:"-" mapstructure:"-"`

	// Identifier is the unique identifier of the clip to be deleted.
	Identifier string `json:"identifier" mapstructure:"identifier"`
}

// MarshalJSON implements Event.
func (d *deleteClipEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"identifier": d.Identifier,
	})
}

// UnmarshalJSON implements Event.
func (d *deleteClipEvent) UnmarshalJSON(raw []byte) error {
	data := map[string]any{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	if v, ok := data["identifier"].(string); ok {
		d.Identifier = v
	}

	return nil
}

// Action implements Event.
func (d *deleteClipEvent) Action() repository.Action {
	return repository.ActionDeleteClip
}

// Apply implements Event.
func (d *deleteClipEvent) Apply(repository repository.Repository) (Log, error) {
	panic("unimplemented")
}

// Replay implements Event.
func (d *deleteClipEvent) Replay(repository repository.Repository) error {
	panic("unimplemented")
}

var _ Event = (*deleteClipEvent)(nil)
