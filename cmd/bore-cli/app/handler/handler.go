package handler

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/events"
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
	FlagDataDir    = "data-dir"
	FlagDelete     = "delete"
	FlagFormat     = "format"
	FlagSystem     = "system"
	FlagInputFile  = "input-file"
	FlagMimeType   = "mime-type"
	FlagOutputFile = "output-file"
)

type (
	PasteOptions struct {
		Collection    string
		DeleteOnPaste bool
		Format        PasteFormat
		OutputFile    string
	}

	CopyOptions struct {
		Collection string
		InputFile  string
		MimeType   events.MimeType
		System     bool
	}
)

type Handler struct {
	bore *bore.Bore
}

func New(bore *bore.Bore) *Handler {
	return &Handler{bore: bore}
}

func (h *Handler) Copy(ctx *cli.Context) error {
	panic("implement me")
}

func (h *Handler) Paste(ctx *cli.Context) error {
	options := &PasteOptions{
		Collection:    ctx.String(FlagCollection),
		Format:        PasteFormat(ctx.String(FlagFormat)),
		DeleteOnPaste: ctx.Bool(FlagDelete),
		OutputFile:    ctx.String(FlagOutputFile),
	}

	if ctx.Bool(FlagSystem) {
		return h.pasteFromSystem(ctx, options)
	}

	return h.pasteFromDatabase(ctx, options)
}

func (h *Handler) pasteFromSystem(ctx *cli.Context, options *PasteOptions) error {
	clipboard, err := h.bore.SystemClipboard()
	if err != nil {
		return err
	}

	if !clipboard.Available() {
		return ErrClipboardNotAvailable
	}

	content, err := clipboard.Read(ctx.Context)
	if err != nil {
		return err
	}

	if options.DeleteOnPaste {
		if err := clipboard.Clear(ctx.Context); err != nil {
			return err
		}
	}

	switch options.Format {
	case PasteFormatText:
		break

	case PasteFormatBase64:
		base64Content := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
		base64.StdEncoding.Encode(base64Content, content)
		content = base64Content

	case PasteFormatJSON:
		jsonContent, err := json.Marshal(map[string]string{"content": string(content)})
		if err != nil {
			return err
		}
		content = jsonContent

	default:
		return cli.Exit("unsupported format: "+string(options.Format), 1)
	}

	if options.OutputFile != "" {
		return h.writeToFile(ctx, options.OutputFile, content)
	}

	return h.writeToStdout(ctx, content)
}

func (h *Handler) pasteFromDatabase(ctx *cli.Context, options *PasteOptions) error {
	panic("implement me")
}

func (h *Handler) writeToFile(ctx *cli.Context, filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

func (h *Handler) writeToStdout(ctx *cli.Context, content []byte) error {
	_, err := ctx.App.Writer.Write(content)
	return err
}
