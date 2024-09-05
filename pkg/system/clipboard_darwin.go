//go:build darwin

package system

type nativeClipboard struct{}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	return &nativeClipboard{}, nil
}

// IsAvailable checks if the clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	panic("not implemented")
}

// PasteToNativeClipboard pastes the provided content to the system clipboard
func (n *nativeClipboard) Paste(content []byte) error {
	panic("not implemented")
}

// CopyFromNativeClipboard copies the content from the system clipboard
func (n *nativeClipboard) Copy() ([]byte, error) {
	panic("not implemented")
}

var _ NativeClipboardInterface = (*nativeClipboard)(nil)
