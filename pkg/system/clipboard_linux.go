//go:build linux

package system

type nativeClipboard struct {
	binName string
}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	panic("not implemented")
}

// IsAvailable checks if a clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	panic("not implemented")
}

// Copy copies the content to the system clipboard
func (n *nativeClipboard) Copy(content []byte) error {
	panic("not implemented")
}

// Paste returns the last copied content from the system clipboard
func (n *nativeClipboard) Paste() ([]byte, error) {
	panic("not implemented")
}

var _ NativeClipboardInterface = (*nativeClipboard)(nil)
