package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
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

type CopyOptions struct {
	Collection string
	InputFile  string
	MimeType   mimetype.MimeType
	System     bool
}

type Handler struct {
	bore *bore.Bore
}

func New(bore *bore.Bore) *Handler {
	return &Handler{bore: bore}
}

type CliCopyOptions struct {
	Stdin bool
}

func (h *Handler) Copy(ctx *cli.Context, options CliCopyOptions) error {
	inputFile := ctx.String(FlagInputFile)

	rawMimeType := ctx.String(FlagMimeType)
	if rawMimeType == "" {
		rawMimeType = "text/plain"
	}
	mimeType, err := mimetype.ParseMimeType(rawMimeType)
	if err != nil {
		return cli.Exit("invalid mime type: "+err.Error(), 1)
	}

	var content []byte

	switch {
	case inputFile != "":
		content, err = os.ReadFile(inputFile)
		if err != nil {
			return cli.Exit("failed to read input file: "+err.Error(), 1)
		}

	case options.Stdin:
		reader := bufio.NewReader(ctx.App.Reader)
		content, err = io.ReadAll(reader)
		if err != nil {
			return cli.Exit("failed to read from stdin: "+err.Error(), 1)
		}

	case ctx.NArg() == 0:
		return cli.Exit(
			"no input provided. Please provide input via argument or --input-file flag",
			1,
		)
	case ctx.NArg() > 1:
		return cli.Exit("too many arguments provided. Only one argument is allowed", 1)
	default:
		content = []byte(ctx.Args().First())
	}

	return h.bore.Copy(ctx.Context, content, bore.CopyOptions{
		Passthrough:  ctx.Bool(FlagSystem),
		CollectionID: ctx.String(FlagCollection),
		Mimetype:     mimeType,
	})
}

func (h *Handler) Paste(ctx *cli.Context) error {
	format := PasteFormat(ctx.String(FlagFormat))
	if format == "" {
		format = PasteFormatText
	}

	outputFile := ctx.String(FlagOutputFile)

	content, err := h.bore.Paste(ctx.Context, bore.PasteOptions{
		ItemID:              ctx.String(FlagIdentifier),
		CollectionID:        ctx.String(FlagCollection),
		FromSystemClipboard: ctx.Bool(FlagSystem),
		DeleteAfterPaste:    ctx.Bool(FlagDelete),
		SkipCollectionCheck: false,
	})
	if err != nil {
		return cli.Exit("failed to paste content: "+err.Error(), 1)
	}

	if content, err = h.contentToFormat(content, format); err != nil {
		return err
	}

	if outputFile != "" {
		return h.writeToFile(ctx, outputFile, content)
	}

	return h.writeToStdout(ctx, content)
}

func (h *Handler) writeToFile(_ *cli.Context, filename string, content []byte) error {
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
