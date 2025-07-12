//go:build darwin

package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type macClipboard struct {
	// binaries holds the paths to the binaries used for clipboard operations.
	binaries Binaries
}

// NewNativeClipboard creates a new instance of NativeClipboard for macOS using the pbcopy and pbpaste binaries.
func NewNativeClipboard(binaries Binaries) (NativeClipboard, error) {
	var (
		err error

		copyBinary, pasteBinary string
	)

	if copyBinary, err = exec.LookPath("pbcopy"); err != nil {
		return nil, fmt.Errorf("pbcopy binary not found in PATH: %w", err)
	}

	if pasteBinary, err = exec.LookPath("pbpaste"); err != nil {
		return nil, fmt.Errorf("pbpaste binary not found in PATH: %w", err)
	}

	binaries = Binaries{
		copy:  copyBinary,
		paste: pasteBinary,
	}

	if binaries.Empty() {
		return nil, fmt.Errorf("required clipboard binaries not found in PATH: %v", binaries)
	}

	return &macClipboard{binaries}, nil
}

// Available implements NativeClipboard.
func (n *macClipboard) Available() bool {
	return n.binaries.copy != "" && n.binaries.paste != ""
}

// Binaries implements NativeClipboard.
func (n *macClipboard) Binaries() Binaries {
	return n.binaries
}

// Clear implements NativeClipboard.
func (n *macClipboard) Clear(ctx context.Context) error {
	return n.Write(ctx, []byte{})
}

// Read implements NativeClipboard.
func (n *macClipboard) Read(ctx context.Context) ([]byte, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxWithTimeout, n.binaries.paste)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read from clipboard: %w", err)
	}

	return output, nil
}

// Write implements NativeClipboard.
func (n *macClipboard) Write(ctx context.Context, data []byte) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxWithTimeout, n.binaries.copy)
	cmd.Stdin = bytes.NewReader(data)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return nil
}

var _ NativeClipboard = (*macClipboard)(nil)
