//go:build windows

package system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.trulyao.dev/bore/v2/clipboard/internal/win32"
)

type windowsClipboard struct {
	mu sync.Mutex
}

func NewNativeClipboard() (NativeClipboard, error) {
	return &windowsClipboard{}, nil
}

// Available implements NativeClipboard.
func (w *windowsClipboard) Available() bool {
	return true
}

// Binaries implements NativeClipboard.
func (w *windowsClipboard) Binaries() Binaries {
	return Binaries{}
}

// Clear implements NativeClipboard.
func (w *windowsClipboard) Clear(ctx context.Context) error {
	return w.Write(ctx, []byte{})
}

// Read implements NativeClipboard.
func (w *windowsClipboard) Read(ctx context.Context) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// NOTE: The Win32 calls don't take a context, so we have to manually handle the timeout
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		s, err := win32.ReadAll()
		ch <- result{data: []byte(s), err: err}
	}()

	select {
	case <-ctxWithTimeout.Done():
		return nil, ctxWithTimeout.Err()
	case result := <-ch:
		if result.err != nil {
			return nil, fmt.Errorf("failed to read from clipboard: %w", result.err)
		}
		return result.data, nil
	}
}

// Write implements NativeClipboard.
func (w *windowsClipboard) Write(ctx context.Context, data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// NOTE: The Win32 calls don't take a context, so we have to manually handle the timeout
	ch := make(chan error, 1)
	go func() {
		ch <- win32.WriteAll(string(data))
	}()

	select {
	case <-ctxWithTimeout.Done():
		return ctxWithTimeout.Err()
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("failed to write to clipboard: %w", err)
		}
		return nil
	}
}

var _ NativeClipboard = (*windowsClipboard)(nil)
