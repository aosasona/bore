package handler

import (
	"context"
	"database/sql"
	"io"
	"time"

	"go.trulyao.dev/bore/pkg/daos"
)

const (
	ArtifactTypeText = "text/plain"
)

type CopyOpts struct {
	CollectionId string
	ArtifactType string
}

type HandlerInterface interface {
	Copy(r io.Reader, opts CopyOpts) (string, error)

	Paste(artifactIds []string) (io.Reader, error)
}

type Handler struct {
	dao *daos.Queries
}

func New(dao *daos.Queries) *Handler {
	return &Handler{dao: dao}
}

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(r io.Reader, opts CopyOpts) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	// Check if the content already exists, if it does, just update the last modified time
	// TODO: write to the native clipboard regardless if enabled

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	createArtifactParams := daos.UpsertArtifactParams{Content: content, ArtifactType: "text"}
	if opts.CollectionId != "" {
		createArtifactParams.CollectionID = sql.NullString{String: opts.CollectionId, Valid: true}
	}

	if opts.ArtifactType != "" {
		createArtifactParams.ArtifactType = opts.ArtifactType
	}

	artifact, err := h.dao.UpsertArtifact(ctx, createArtifactParams)
	if err != nil {
		return "", err
	}

	return artifact.ID, nil
}

// Paste returns the content of the given IDs or the last content if no IDs are provided
func (h *Handler) Paste(ids []string) (io.Reader, error) {
	return nil, nil
}

var _ HandlerInterface = (*Handler)(nil)
