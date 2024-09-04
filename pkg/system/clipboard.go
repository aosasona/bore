package system

// NativeClipboard is the interface for interacting with the underlying clipboard
type NativeClipboardInterface interface {
	IsAvailable() bool
	Paste(content []byte) error
	Copy() ([]byte, error)
}
