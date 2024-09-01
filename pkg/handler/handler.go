package handler

import (
	"io"

	"go.trulyao.dev/bore/pkg/daos"
)

type HandlerInterface interface {
	Copy(io.Reader) (string, error)

	Paste([]int) (io.Reader, error)
}

type Handler struct {
	dao *daos.Queries
}

func New(dao *daos.Queries) *Handler {
	return &Handler{dao: dao}
}

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	// TODO: remove
	_ = content

	return "", nil
}

// Paste returns the content of the given IDs or the last content if no IDs are provided
func (h *Handler) Paste(ids []int) (io.Reader, error) {
	return nil, nil
}

var _ HandlerInterface = (*Handler)(nil)
