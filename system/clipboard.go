package system

type BinPath struct {
	CopyBin  string // Path to the binary for copying to clipboard
	PasteBin string // Path to the binary for pasting from clipboard
}

// NativeClipboard is an interface that defines methods for common clipboard operations.
type NativeClipboard interface {
	// IsAvailable checks if a clipboard implementation is available.
	IsAvailable() bool

	// Write writes bytes to the clipboard.
	Write(data []byte) error

	// Read reads bytes from the clipboard.
	Read() ([]byte, error)

	// Clear clears the clipboard content.
	Clear() error

	// BinPaths returns the paths to the binaries used for clipboard operations (for example, on MacOS, these might be pbcopy and pbpaste).
	BinPaths() []BinPath
}
