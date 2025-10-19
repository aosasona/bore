package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type CopyOptions struct {
	Collection string
	InputFile  string
	MimeType   mimetype.MimeType
	System     bool
}

type CliCopyOptions struct {
	Stdin bool
}

func (h *Handler) Copy(ctx *cli.Context, options CliCopyOptions) error {
	inputFile := ctx.String(FlagInputFile)

	format := PasteFormat(strings.TrimSpace(ctx.String(FlagFormat)))
	if format == "" {
		format = PasteFormatText
	}
	if !slices.Contains([]PasteFormat{PasteFormatText, PasteFormatBase64}, format) {
		return cli.Exit("invalid format for copy command: "+string(format), 1)
	}

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
		reader := bufio.NewReader(ctx.App.Reader)
		fmt.Println(
			"No input argument provided. Please enter the content to copy (end with Ctrl+D):",
		)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return cli.Exit("failed to read from stdin: "+err.Error(), 1)
			}

			content = append(content, line...)
			if err == io.EOF {
				break
			}
		}

		if len(content) == 0 {
			return cli.Exit("no content provided to copy", 1)
		}
		content = content[:len(content)-1] // Remove the last newline character

	case ctx.NArg() > 1:
		return cli.Exit("too many arguments provided. Only one argument is allowed", 1)
	default:
		content = []byte(ctx.Args().First())
	}

	if format == PasteFormatBase64 {
		if content, err = h.decodebase64Content(content); err != nil {
			return err
		}
	}

	return h.bore.Clipboard().Set(ctx.Context, content, bore.SetClipboardOptions{
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

	item, err := h.bore.Get(ctx.Context, bore.GetClipboardOptions{
		ItemID:              ctx.String(FlagIdentifier),
		CollectionID:        ctx.String(FlagCollection),
		FromSystemClipboard: ctx.Bool(FlagSystem),
		DeleteAfterPaste:    ctx.Bool(FlagDelete),
		SkipCollectionCheck: false,
	})
	if err != nil {
		return cli.Exit("failed to paste content: "+err.Error(), 1)
	}

	var content []byte
	if content, err = h.contentToFormat(item, format); err != nil {
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

func (h *Handler) contentToFormat(result bore.PasteResult, format PasteFormat) ([]byte, error) {
	switch format {
	case PasteFormatText:
		return result.Content, nil

	case PasteFormatBase64:
		base64Content := make([]byte, base64.StdEncoding.EncodedLen(len(result.Content)))
		base64.StdEncoding.Encode(base64Content, result.Content)
		return base64Content, nil

	case PasteFormatJSON:
		jsonContent, err := json.Marshal(map[string]string{
			"id":            result.Item.ID,
			"mimetype":      result.Item.Mimetype,
			"content":       string(result.Content),
			"collection_id": result.Item.CollectionID.String,
			"created_at":    result.Item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
		if err != nil {
			return nil, err
		}
		return jsonContent, nil

	default:
		return nil, cli.Exit("unsupported format: "+string(format), 1)
	}
}

func (h *Handler) decodebase64Content(content []byte) ([]byte, error) {
	decodedContent := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
	n, err := base64.StdEncoding.Decode(decodedContent, content)
	if err != nil {
		return nil, cli.Exit("failed to decode base64 content: "+err.Error(), 1)
	}
	return decodedContent[:n], nil
}
