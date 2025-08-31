package handler

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/database/repository"
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
	FlagIdentifier = "identifier"
	FlagSystem     = "system"
	FlagInputFile  = "input-file"
	FlagMimeType   = "mime-type"
	FlagOutputFile = "output-file"
)

type (
	PasteOptions struct {
		Identifier string
		Collection string

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
		Collection: ctx.String(FlagCollection),
		Identifier: ctx.String(FlagIdentifier),

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

	if content, err = h.contentToFormat(content, options.Format); err != nil {
		return err
	}

	if options.OutputFile != "" {
		return h.writeToFile(ctx, options.OutputFile, content)
	}

	return h.writeToStdout(ctx, content)
}

func (h *Handler) pasteFromDatabase(ctx *cli.Context, options *PasteOptions) error {
	repo, err := h.bore.Repository()
	if err != nil {
		return err
	}

	var clip *repository.Clip

	if strings.TrimSpace(options.Identifier) != "" {
		if clip, err = repo.Clips().FindById(ctx.Context, options.Identifier); err != nil {
			return err
		}
	} else {
		if clip, err = repo.Clips().FindLatestClip(ctx.Context, options.Collection); err != nil {
			return err
		}
	}

	if clip == nil {
		return cli.Exit("no clip found", 1)
	}

	content := clip.Content
	if content, err = h.contentToFormat(content, options.Format); err != nil {
		return err
	}

	// TODO: handle delete on paste
	// TODO: handle output to file

	panic("not implemented")
}

func (h *Handler) writeToFile(ctx *cli.Context, filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

func (h *Handler) writeToStdout(ctx *cli.Context, content []byte) error {
	_, err := ctx.App.Writer.Write(content)
	return err
}

func (h *Handler) contentToFormat(content []byte, format PasteFormat) ([]byte, error) {
	switch format {
	case PasteFormatText:
		return content, nil

	case PasteFormatBase64:
		base64Content := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
		base64.StdEncoding.Encode(base64Content, content)
		return base64Content, nil

	case PasteFormatJSON:
		jsonContent, err := json.Marshal(map[string]string{"content": string(content)})
		if err != nil {
			return nil, err
		}
		return jsonContent, nil

	default:
		return nil, cli.Exit("unsupported format: "+string(format), 1)
	}
}
