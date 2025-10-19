package handler

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/config"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/tui"
)

var ErrClipboardNotAvailable = cli.Exit("system clipboard is not available", 1)

type PasteFormat string

const (
	PasteFormatText   PasteFormat = "text"
	PasteFormatJSON   PasteFormat = "json"
	PasteFormatBase64 PasteFormat = "base64"
)

const (
	FlagConfig     = "config"
	FlagCollection = "collection"
	FlagName       = "name"
	FlagDataDir    = "data-dir"
	FlagDelete     = "delete"
	FlagFormat     = "format"
	FlagForce      = "force"
	FlagIdentifier = "identifier"
	FlagSystem     = "system"
	FlagInputFile  = "input-file"
	FlagMimeType   = "mime-type"
	FlagOutputFile = "output-file"
)

type Handler struct {
	tuiManager    *tui.Manager
	configManager *config.Manager
	bore          *bore.Bore
}

func New(bore *bore.Bore, tuiManager *tui.Manager, configManager *config.Manager) *Handler {
	return &Handler{
		bore:          bore,
		tuiManager:    tuiManager,
		configManager: configManager,
	}
}
