//go:build !windows && !darwin && !linux && !ios

package clipboard

import "context"

type stubClipboard struct{}

// Available implements NativeClipboard.
func (n *stubClipboard) Available() bool {
	return false
}

// Binaries implements NativeClipboard.
func (n *stubClipboard) Binaries() Binaries {
	panic("not supported on this platform")
}

// Clear implements NativeClipboard.
func (n *stubClipboard) Clear(ctx context.Context) error {
	panic("not supported on this platform")
}

// Read implements NativeClipboard.
func (n *stubClipboard) Read(ctx context.Context) ([]byte, error) {
	panic("not supported on this platform")
}

// Write implements NativeClipboard.
func (n *stubClipboard) Write(ctx context.Context, data []byte) error {
	panic("not supported on this platform")
}

var _ NativeClipboard = (*stubClipboard)(nil)
