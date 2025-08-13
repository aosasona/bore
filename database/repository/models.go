package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
)

type Action string

const (
	ActionCopyV1           Action = "copy_v1"
	ActionCreateCollection Action = "create_collection"
)

type Relay struct {
	bun.BaseModel `bun:"table:relays,alias:r"`

	ID            string    `bun:"id,pk"`
	Alias         string    `bun:"alias,notnull"`
	Address       string    `bun:"address,notnull"`
	Metadata      []byte    `bun:"metadata"`
	AddedAt       time.Time `bun:"added_at,notnull,default:(unixepoch())"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull,default:(unixepoch())"`
}

type Peer struct {
	bun.BaseModel `bun:"table:peers,alias:p"`

	ID         string    `bun:"id,pk"`
	Name       string    `bun:"name,notnull"`
	Metadata   []byte    `bun:"metadata"`
	AddedAt    time.Time `bun:"added_at,notnull,default:(unixepoch())"`
	LastSeenAt string    `bun:"last_seen_at,notnull"`

	RelayID string `bun:",notnull"`

	Relay *Relay `bun:"rel:belongs-to,join:relay_id=id"`
}

type Collection struct {
	bun.BaseModel `bun:"table:collections,alias:co"`

	ID        string       `bun:"id,pk"`
	Name      string       `bun:"name,notnull"`
	Hash      string       `bun:"hash,notnull"`
	PinnedAt  bun.NullTime `bun:"pinned_at"`
	CreatedAt time.Time    `bun:"created_at,notnull,default:(unixepoch())"`
	UpdatedAt string       `bun:"updated_at,notnull,default:(unixepoch())"`

	PeerID string `bun:",notnull"`

	Peer *Peer `bun:"rel:belongs-to,join:peer_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (c *Collection) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	c.ID = ulid.Make().String()
	return nil
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID          string         `bun:"id,pk"`
	AggregateID sql.NullString `bun:"aggregate_id"`
	Action      Action         `bun:"action,notnull"`
	Version     int            `bun:"version,notnull,default:1"`
	Payload     []byte         `bun:"payload,notnull"`

	IngestedAt time.Time    `bun:"ingested_at,notnull,default:(unixepoch())"`
	LoggedAt   bun.NullTime `bun:"logged_at,default:(unixepoch())"`
	AppliedAt  bun.NullTime `bun:"applied_at,default:(unixepoch())"`

	PeerID string `bun:",notnull"`

	Peer *Peer `bun:"rel:belongs-to,join:peer_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (e *Event) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	e.ID = ulid.Make().String()
	return nil
}

type Clip struct {
	bun.BaseModel `bun:"table:clips,alias:cl"`

	ID        string    `bun:"id,pk"`
	Content   []byte    `bun:"contnet,notnull"`
	Hash      string    `bun:"hash,notnull"`
	Mimetype  string    `bun:"mimetype,notnull"`
	Size      int64     `bun:"size,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:(unixepoch())"`
	UpdatedAt string    `bun:"updated_at,notnull,default:(unixepoch())"`

	CollectionID string `bun:",notnull"`
	PeerID       string `bun:",notnull"`

	Collection *Collection `bun:"rel:belongs-to,join:collection_id=id"`
	Peer       *Peer       `bun:"rel:belongs-to,join:peer_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (c *Clip) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	c.ID = ulid.Make().String()
	return nil
}

var (
	_ bun.BeforeInsertHook = (*Clip)(nil)
	_ bun.BeforeInsertHook = (*Collection)(nil)
	_ bun.BeforeInsertHook = (*Event)(nil)
)
