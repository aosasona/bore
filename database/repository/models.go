package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type Action string

// TODO: validation with https://github.com/go-playground/validator
// TODO: create default collection
const (
	ActionCopyV1           Action = "copy_v1"
	ActionCreateCollection Action = "create_collection"
	ActionDeleteClip       Action = "delete_clip"
	ActionDeleteCollection Action = "delete_collection"
)

type Relay struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:relays,alias:r"`

	ID            string    `bun:"id,pk"`
	Alias         string    `bun:"alias,notnull"                                 validate:"min=3,max=32"`
	Address       string    `bun:"address,notnull"`
	Metadata      []byte    `bun:"metadata"`
	AddedAt       time.Time `bun:"added_at,notnull,default:(unixepoch())"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull,default:(unixepoch())"`
}

type Peer struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:peers,alias:p"`

	ID         string    `bun:"id,pk"`
	Name       string    `bun:"name,notnull"                           validate:"alphanumunicode,min=3,max=32"`
	Metadata   []byte    `bun:"metadata"`
	AddedAt    time.Time `bun:"added_at,notnull,default:(unixepoch())"`
	LastSeenAt string    `bun:"last_seen_at,notnull"`

	RelayID string `bun:",notnull"`

	Relay *Relay `bun:"rel:belongs-to,join:relay_id=id"`
}

type Collection struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:collections,alias:co"`

	ID        string       `bun:"id,pk"`
	Name      string       `bun:"name,notnull"`
	Hash      string       `bun:"hash,notnull"`
	PinnedAt  bun.NullTime `bun:"pinned_at"`
	CreatedAt time.Time    `bun:"created_at,notnull,default:(unixepoch())"`
	UpdatedAt string       `bun:"updated_at,notnull,default:(unixepoch())"`

	PeerID string `bun:"peer_id,notnull"`

	Peer *Peer `bun:"rel:belongs-to,join:peer_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (collection *Collection) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if collection.ID == "" {
		collection.ID = ulid.Make().String()
	}

	return nil
}

type Event struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:events,alias:e"`

	ID          string         `bun:"id,pk"`
	AggregateID sql.NullString `bun:"aggregate_id"`
	Action      Action         `bun:"action,notnull"`
	Version     int            `bun:"version,notnull,default:1" validate:"gte=1"`
	Payload     []byte         `bun:"payload,notnull"`

	IngestedAt time.Time    `bun:"ingested_at,notnull,default:(unixepoch())"`
	LoggedAt   bun.NullTime `bun:"logged_at,default:(unixepoch())"`
	AppliedAt  bun.NullTime `bun:"applied_at,default:(unixepoch())"`

	PeerID string `bun:"peer_id,notnull"`

	Peer *Peer `bun:"rel:belongs-to,join:peer_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (event *Event) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if event.ID == "" {
		event.ID = ulid.Make().String()
	}

	return nil
}

type Clip struct {
	validation.ValidateStructMixin
	bun.BaseModel `bun:"table:clips,alias:cl"`

	ID      string `bun:"id,pk"                                    validate:"required,ulid"`
	Content []byte `bun:"contnet,notnull"                          validate:"required"`
	Hash    string `bun:"hash,notnull"                             validate:"required,sha256"`
	// TODO: implement custom validator for mimetype
	Mimetype  string    `bun:"mimetype,notnull"                         validate:"required,mimetype"`
	Size      int64     `bun:"size,notnull" 						   validate:"required,gt=0"`
	CreatedAt time.Time `bun:"created_at,notnull,default:(unixepoch())"`
	UpdatedAt string    `bun:"updated_at,notnull,default:(unixepoch())"`

	CollectionID string `bun:"collection_id" validate:"required,ulid"`

	Collection *Collection `bun:"rel:belongs-to,join:collection_id=id"`
}

// BeforeInsert implements bun.BeforeInsertHook.
func (clip *Clip) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	if clip.ID == "" {
		clip.ID = ulid.Make().String()
	}

	return nil
}

var (
	_ bun.BeforeInsertHook = (*Clip)(nil)
	_ bun.BeforeInsertHook = (*Collection)(nil)
	_ bun.BeforeInsertHook = (*Event)(nil)
)
