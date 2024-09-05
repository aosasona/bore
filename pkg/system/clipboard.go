package system

// NativeClipboard is the interface for interacting with the underlying clipboard
type NativeClipboardInterface interface {
	// IsAvailable checks if a native system clipboard is available
	IsAvailable() bool

	// Copy copies the content to the system clipboard
	Copy(content []byte) error

	// Paste returns the last copied content from the system clipboard
	Paste() ([]byte, error)
}
