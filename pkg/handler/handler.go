package handler

import (
	"context"
	"database/sql"
	"io"
	"time"

	"go.trulyao.dev/bore/pkg/daos"
)

type HandlerInterface interface {
	Copy(string, io.Reader) (string, error)

	Paste([]string) (io.Reader, error)
}

type Handler struct {
	dao *daos.Queries
}

func New(dao *daos.Queries) *Handler {
	return &Handler{dao: dao}
}

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(id string, r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	// Check if the content already exists, if it does, just update the last modified time
	// TODO: write to the native clipboard regardless if enabled

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	article, err := h.dao.GetArtifactByContent(ctx, daos.GetArtifactByContentParams{
		Content:      string(content),
		CollectionID: sql.NullString{String: id, Valid: true},
	})
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	_ = article

	return "", nil
}

// Paste returns the content of the given IDs or the last content if no IDs are provided
func (h *Handler) Paste(ids []string) (io.Reader, error) {
	return nil, nil
}

var _ HandlerInterface = (*Handler)(nil)
