//go:build !windows && !darwin && !linux

package clipboard

import "context"

type notSupportedClipboard struct{}

// Available implements NativeClipboard.
func (n *notSupportedClipboard) Available() bool {
	return false
}

// Binaries implements NativeClipboard.
func (n *notSupportedClipboard) Binaries() Binaries {
	panic("not supported on this platform")
}

// Clear implements NativeClipboard.
func (n *notSupportedClipboard) Clear(ctx context.Context) error {
	panic("not supported on this platform")
}

// Read implements NativeClipboard.
func (n *notSupportedClipboard) Read(ctx context.Context) ([]byte, error) {
	panic("not supported on this platform")
}

// Write implements NativeClipboard.
func (n *notSupportedClipboard) Write(ctx context.Context, data []byte) error {
	panic("not supported on this platform")
}

var _ NativeClipboard = (*notSupportedClipboard)(nil)
