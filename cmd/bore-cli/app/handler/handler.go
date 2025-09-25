package handler

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
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
		MimeType   mimetype.MimeType
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
	format := PasteFormat(ctx.String(FlagFormat))
	outputFile := ctx.String(FlagOutputFile)

	content, err := h.bore.Paste(ctx.Context, bore.PasteOptions{
		ItemID:              ctx.String(FlagIdentifier),
		CollectionID:        ctx.String(FlagCollection),
		FromSystemClipboard: ctx.Bool(FlagSystem),
		DeleteAfterPaste:    ctx.Bool(FlagDelete),
		SkipCollectionCheck: false,
	})

	if content, err = h.contentToFormat(content, format); err != nil {
		return err
	}

	if outputFile != "" {
		return h.writeToFile(ctx, outputFile, content)
	}

	return h.writeToStdout(ctx, content)
}

func (h *Handler) writeToFile(ctx *cli.Context, filename string, content []byte) error {
	return os.WriteFile(filename, content, 0o644)
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
