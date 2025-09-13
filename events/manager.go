package events

import (
	"crypto/sha256"

	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/device"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type Manager struct {
	identity   *device.Identity
	repository repository.Repository
}

func NewManager(repository repository.Repository, identity *device.Identity) *Manager {
	return &Manager{
		identity:   identity,
		repository: repository,
	}
}

func (m *Manager) Copy(
	content []byte,
	mimeType mimetype.MimeType,
	collectionID string,
) *copyEvent {
	return &copyEvent{
		identity:     m.identity,
		Content:      content,
		Hash:         hash(content),
		MimeType:     mimeType,
		CollectionID: collectionID,
	}
}

func (m *Manager) DeleteClip(clipId string) *deleteClipEvent {
	return &deleteClipEvent{
		identity: m.identity,
		ClipId:   clipId,
	}
}

func (m *Manager) Log(event Event) (Log, error) {
	return event.Apply(m.repository)
}

func hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return string(h.Sum(nil))
}
