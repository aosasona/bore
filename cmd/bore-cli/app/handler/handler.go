package handler

import (
	"io"

	"go.trulyao.dev/bore/v2"
)

type Handler struct {
	bore *bore.Bore
}

func New(bore *bore.Bore) *Handler {
	return &Handler{bore: bore}
}

func (h *Handler) CopyFromStdin(reader *io.Reader) {
	panic("implement me")
}
