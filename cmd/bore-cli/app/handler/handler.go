package handler

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
)

type PasteFormat string

const (
	PasteFormatText   PasteFormat = "text"
	PasteFormatJSON   PasteFormat = "json"
	PasteFormatBase64 PasteFormat = "base64"
)

const (
	FlagConfig     = "config"
	FlagCollection = "collection"
	FlagDataDir    = "data-dir"
	FlagDelete     = "delete"
	FlagFormat     = "format"
	FlagFromSystem = "system"
	FlagInputFile  = "input-file"
	FlagMimeType   = "mime-type"
	FlagOutputFile = "output-file"
)

type Handler struct {
	bore *bore.Bore
}

func New(bore *bore.Bore) *Handler {
	return &Handler{bore: bore}
}

func (h *Handler) Paste(ctx *cli.Context) {
	panic("implement me")
}
