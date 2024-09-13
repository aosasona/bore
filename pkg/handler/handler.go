package handler

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/daos"
	"go.trulyao.dev/bore/pkg/system"
)

// TODO: add PasteIdx method
type HandlerInterface interface {
	// Copy the content from the reader to the clipboard
	Copy(r io.Reader, opts CopyOpts) (string, error)

	// PasteLast the last copied content from the specified source to the writer
	PasteLast(Source, io.Writer, PasteOpts) (string, error)

	// RemoveLast the clipboard's content (last copied content)
	RemoveLast(Source) error

	// Delete the content at the specified index
	// NOTE: this is only applicable to the bore clipboard
	RemoveIdx(source Source, id string) error

	// Decodes the content from the specified format
	Decode(content []byte, from Format) ([]byte, error)

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

// Remove the last copied content from the clipboard
func (h *Handler) RemoveLast(source Source) error {
	switch source {
	case SourceBore:
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()

		if err := h.dao.DeleteLatestArtifact(ctx); err != nil {
			return fmt.Errorf("failed to delete latest artifact: %s", err.Error())
		}

		return nil

	case SourceSystem:
		if !h.nativeClipboard.IsAvailable() {
			return fmt.Errorf("no native clipboard found")
		}

		return h.nativeClipboard.Clear()

	default:
		return fmt.Errorf("unsupported source: %s", source)
	}
}

// RemoveIdx removes the content at the specified index from the clipboard
// NOTE: this is only applicable to the bore clipboard
func (h *Handler) RemoveIdx(source Source, id string) error {
	if source != SourceBore {
		return fmt.Errorf("unsupported source: %s", source)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	if err := h.dao.DeleteArtifactById(ctx, id); err != nil {
		return fmt.Errorf("failed to delete artifact: %s", err.Error())
	}

	return nil
}

var _ HandlerInterface = (*Handler)(nil)
