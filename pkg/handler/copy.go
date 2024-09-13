package handler

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"go.trulyao.dev/bore/pkg/daos"
)

// Copy copies the content of the reader to the database and returns the ID of the content
func (h *Handler) Copy(r io.Reader, opts CopyOpts) (string, error) {
	if !ValidateFormat(opts.Format) {
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	if content, err = h.Decode(content, opts.Format); err != nil {
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
			fmt.Fprintln(
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
