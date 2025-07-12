package system

type Binaries struct {
	// Path to the binary for copying to clipboard
	Copy string

	// Path to the binary for pasting from clipboard
	Paste string
}

// NativeClipboard is an interface that defines methods for common clipboard operations.
type NativeClipboard interface {
	// Available checks if a clipboard implementation is available for the current platform
	Available() bool

	// Write writes bytes to the clipboard.
	Write(data []byte) error

	// Read reads bytes from the clipboard.
	Read() ([]byte, error)

	// Clear clears the clipboard content.
	Clear() error

	// Binaries returns the paths to the (often) individual binaries used for clipboard operations
	// For example; on MacOS, these might be `pbcopy` and `pbpaste`, while on Linux, they might be `xclip` or `xsel`.
	Binaries() Binaries
}
