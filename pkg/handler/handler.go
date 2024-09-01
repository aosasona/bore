package handler

import (
	"database/sql"
	"io"
)

type HandlerInterface interface {
	Copy(io.Reader) (string, error)

	Paste([]int) (io.Reader, error)
}

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(r io.Reader) (string, error) {
	return "", nil
}

// Paste returns the content of the given IDs or the last content if no IDs are provided
func (h *Handler) Paste(ids []int) (io.Reader, error) {
	return nil, nil
}

var _ HandlerInterface = (*Handler)(nil)
