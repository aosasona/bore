package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/daos"
	"go.trulyao.dev/bore/pkg/system"
)

const (
	FormatBase64    = "base64"
	FormatPlainText = "plain"
)

type CopyOpts struct {
	// Collection ID to associate the copied content with
	CollectionId string

	// Format of the content to copy
	Format string
}

type HandlerInterface interface {
	Copy(r io.Reader, opts CopyOpts) (string, error)

	PasteLastCopied(io.Writer) error

	// DecodeToFormat decodes the content to the specified format
	DecodeToFormat([]byte, string) ([]byte, error)

	// TODO: add a PasteManyIdx method that returns a list of artifacts with their numeric index from the bottom (which is then mapped to their UUID ids) with 0 being most recent
}

type Handler struct {
	dao             *daos.Queries
	config          *config.Config
	nativeClipboard system.NativeClipboardInterface
}

func New(
	dao *daos.Queries,
	config *config.Config,
	nativeClipboard system.NativeClipboardInterface,
) *Handler {
	return &Handler{dao: dao, config: config, nativeClipboard: nativeClipboard}
}

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(r io.Reader, opts CopyOpts) (string, error) {
	if !ValidateFormat(opts.Format) {
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	if content, err = h.DecodeToFormat(content, opts.Format); err != nil {
		return "", err
	}

	// Check if the content already exists, if it does, just update the last modified time
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	createArtifactParams := daos.UpsertArtifactParams{Content: content}
	if opts.CollectionId != "" {
		createArtifactParams.CollectionID = sql.NullString{String: opts.CollectionId, Valid: true}
	}

	// Persist to main store
	artifact, err := h.dao.UpsertArtifact(ctx, createArtifactParams)
	if err != nil {
		return "", fmt.Errorf("Failed to write to bore store: %s", err.Error())
	}

	// Write to native clipboard if enabled and present
	if h.config.EnableNativeClipboard {
		if !h.nativeClipboard.IsAvailable() {
			fmt.Fprint(
				os.Stderr,
				"[WARNING] `EnableNativeClipboard` is set to true in your config but no native clipboard was found on this machine",
			)
			return artifact.ID, nil
		}

		if err = h.nativeClipboard.Copy(content); err != nil {
			return artifact.ID, fmt.Errorf(
				"Copied to bore store but failed to write to native clipboard: %s",
				err.Error(),
			)
		}
	}

	return artifact.ID, nil
}

// PasteLastCopied returns the content of the last artifact (sorted by last modified time) from the database
func (h *Handler) PasteLastCopied(w io.Writer) error {
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	artifact, err := h.dao.GetMostRecentArtifact(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	_, err = w.Write(artifact.Content)
	return err
}

// DecodeToFormat decodes the content to the specified (and supported) format
func (h *Handler) DecodeToFormat(content []byte, format string) ([]byte, error) {
	switch format {
	case FormatBase64:
		destination := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
		if _, err := base64.StdEncoding.Decode(destination, content); err != nil {
			return nil, err
		}
		return destination, nil
	case FormatPlainText:
		return content, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func ValidateFormat(format string) bool {
	return format == FormatBase64 || format == FormatPlainText
}

var _ HandlerInterface = (*Handler)(nil)
