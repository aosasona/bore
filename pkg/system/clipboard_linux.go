package system

type nativeClipboard struct {
	binName string
}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	// TODO: look for all common clipbaord managers on linux
	// xclip | xsel | (wl-copy & wl-paste)
	panic("not implemented")
}

// IsAvailable checks if a clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	return n.binName != ""
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
