package events

import (
	"crypto/sha256"

	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/device"
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

func (m *Manager) NewCopyV1Event(
	content []byte,
	mimeType MimeType,
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

func hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return string(h.Sum(nil))
}
