package clipboard

import (
	"context"
	"fmt"
)

// NativeClipboard is an interface that defines methods for common clipboard operations.
type NativeClipboard interface {
	// Available checks if a clipboard implementation is available for the current platform
	Available() bool

	// Write writes bytes to the clipboard.
	Write(ctx context.Context, data []byte) error

	// Read reads bytes from the clipboard.
	Read(ctx context.Context) ([]byte, error)

	// Clear clears the clipboard content.
	Clear(ctx context.Context) error

	// Binaries returns the paths to the (often) individual binaries used for clipboard operations
	// For example; on MacOS, these might be `pbcopy` and `pbpaste`, while on Linux, they might be `xclip` or `xsel`.
	Binaries() Binaries
}

type Binaries struct {
	// Path to the binary for copying to clipboard
	copy string

	// Path to the binary for pasting from clipboard
	paste string
}

func (b Binaries) Copy() string {
	return b.copy
}

func (b Binaries) Paste() string {
	return b.paste
}

func (b Binaries) Empty() bool {
	return b.copy == "" && b.paste == ""
}

func (b Binaries) String() string {
	if b.Empty() {
		return "no binaries available"
	}
	return fmt.Sprintf("copy: %s, paste: %s", b.copy, b.paste)
}
