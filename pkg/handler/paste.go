package handler

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"go.trulyao.dev/bore/pkg/daos"
	"golang.org/x/sync/errgroup"
)

type PasteOpts struct {
	// Delete the content from the clipboard after pasting
	DeleteOnPaste bool

	// Source to paste from
	Source Source
}

// PasteLast returns the content of the last artifact (sorted by last modified time) from the database
func (h *Handler) PasteLast(source Source, w io.Writer, opts PasteOpts) (string, error) {
	switch source {
	case SourceBore:
		return h.PasteFromBore(w, opts)

	case SourceSystem:
		return h.PasteFromSystemClipboard(w, opts)

	default:
		return "", fmt.Errorf("unsupported source: %s", source)
	}
}

// Paste the last copied content from the bore clipboard
func (h *Handler) PasteFromBore(w io.Writer, opts PasteOpts) (string, error) {
	var (
		artifact daos.Artifact
		err      error
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Get the latest artifact from the database (and remove it if the delete flag is set)
	if opts.DeleteOnPaste {
		artifact, err = h.dao.DeleteAndReturnLatestArtifact(ctx)
	} else {
		artifact, err = h.dao.GetLatestArtifact(ctx)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	_, err = w.Write(artifact.Content)
	if err != nil {
		return "", fmt.Errorf("failed to write content to writer: %s", err.Error())
	}

	return artifact.ID, nil
}

// Paste the last copied content from the system clipboard
func (h *Handler) PasteFromSystemClipboard(w io.Writer, opts PasteOpts) (string, error) {
	if !h.nativeClipboard.IsAvailable() {
		return "", fmt.Errorf("no native clipboard found")
	}

	content, err := h.nativeClipboard.Paste()
	if err != nil {
		return "", fmt.Errorf("failed to paste from system clipboard: %s", err.Error())
	}

	group, _ := errgroup.WithContext(context.TODO())

	// Delete the content from the clipboard if the delete flag is set
	group.Go(func() error {
		if opts.DeleteOnPaste {
			return h.nativeClipboard.Clear()
		}
		return nil
	})

	// Dump the content to the writer
	group.Go(func() error {
		_, err := w.Write(content)
		return err
	})

	err = group.Wait()
	return "", err
}
