package events

import (
	"crypto/sha256"

	"github.com/uptrace/bun"
)

type Manager struct {
	db *bun.DB
}

func NewManager(db *bun.DB) *Manager {
	return &Manager{db}
}

func (m *Manager) NewCopyV1Event(
	content []byte,
	mimeType MimeType,
	collectionID string,
) *copyEvent {
	return &copyEvent{
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
