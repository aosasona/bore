package handler

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"
)

// Paste returns the content of the last artifact (sorted by last modified time) from the database
func (h *Handler) Paste(source Source, w io.Writer) (string, error) {
	switch source {
	case SourceBore:
		return h.PasteFromBore(w)

	case SourceSystem:
		return h.PasteFromSystemClipboard(w)

	default:
		return "", fmt.Errorf("unsupported source: %s", source)
	}
}

// Paste the last copied content from the bore clipboard
func (h *Handler) PasteFromBore(w io.Writer) (string, error) {
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	artifact, err := h.dao.GetLatestArtifact(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	_, err = w.Write(artifact.Content)
	return "", err
}

// Paste the last copied content from the system clipboard
func (h *Handler) PasteFromSystemClipboard(w io.Writer) (string, error) {
	if !h.nativeClipboard.IsAvailable() {
		return "", fmt.Errorf("no native clipboard found")
	}

	content, err := h.nativeClipboard.Paste()
	if err != nil {
		return "", fmt.Errorf("failed to paste from system clipboard: %s", err.Error())
	}

	_, err = w.Write(content)
	return "", err
}
