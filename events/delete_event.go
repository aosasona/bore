package events

import (
	"context"
	"encoding/json"
	"time"

	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/device"
)

type deleteClipEvent struct {
	identity *device.Identity `json:"-" mapstructure:"-"`

	// ClipId is the unique identifier of the clip to be deleted.
	ClipId string `json:"clipId" mapstructure:"clipId"`
}

// MarshalJSON implements Event.
func (d *deleteClipEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"clipId": d.ClipId,
	})
}

// UnmarshalJSON implements Event.
func (d *deleteClipEvent) UnmarshalJSON(raw []byte) error {
	data := map[string]any{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	if v, ok := data["clipId"].(string); ok {
		d.ClipId = v
	}

	return nil
}

// Action implements Event.
func (d *deleteClipEvent) Action() repository.Action {
	return repository.ActionDeleteClip
}

// Apply implements Event.
func (d *deleteClipEvent) Apply(repository repository.Repository) (Log, error) {
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	err := repository.Clips().DeleteById(ctx, d.ClipId)
	if err != nil {
		return Log{}, err
	}

	identity, err := d.identity.GetIdentifier()
	if err != nil {
		return Log{}, err
	}

	return Log{
		Action: d.Action(),
		Metadata: Metadata{
			DeviceID:    identity,
			AggregateID: d.ClipId,
			Version:     1,
		},
		Timestamp: Timestamp{
			IngestedAt: time.Now(),
			AppliedAt:  time.Now(),
		},
	}, nil
}

// Play implements Event.
func (d *deleteClipEvent) Play(repository.Repository, Log) error {
	panic("unimplemented")
}

var _ Event = (*deleteClipEvent)(nil)
