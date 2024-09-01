package handler

import (
	"context"
	"database/sql"
	"fmt"
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

	PasteLastCopied(io.Writer) error

	// TODO: add a PasteManyIdx method that returns a list of artifacts with their numeric index from the bottom (which is then mapped to their UUID ids) with 0 being most recent
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

	createArtifactParams := daos.UpsertArtifactParams{
		Content:      content,
		ArtifactType: ArtifactTypeText,
	}
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

// PasteLastCopied returns the content of the last artifact (sorted by last modified time) from the database
func (h *Handler) PasteLastCopied(w io.Writer) error {
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	artifact, err := h.dao.GetMostRecentArtifact(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No content to paste")
		}

		return err
	}

	// TODO: escape text content

	_, err = w.Write(artifact.Content)
	return err
}

var _ HandlerInterface = (*Handler)(nil)
