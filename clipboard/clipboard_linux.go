//go:build linux

package clipboard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type Program int

const (
	ProgramXsel Program = iota
	ProgramXclip
	ProgramWlClipboard
)

// linuxClipboard implements the NativeClipboard interface for Linux systems, with support for xsel, xclip, or wl-clipboard binaries.
type linuxClipboard struct {
	// Program specifies which clipboard program to use.
	program Program

	// binaries holds the paths to the binaries used for clipboard operations.
	binaries Binaries
}

type clipboardCandidate map[Program][]string

var clipboardCandidatesMap = clipboardCandidate{
	ProgramXsel:        []string{"xsel"},
	ProgramXclip:       []string{"xclip"},
	ProgramWlClipboard: []string{"wl-copy", "wl-paste"},
}

func NewNativeClipboard() (NativeClipboard, error) {
	for program, bins := range clipboardCandidatesMap {
		var binaries Binaries

		switch program {
		case ProgramXsel, ProgramXclip:
			if path, err := exec.LookPath(bins[0]); err == nil {
				binaries = Binaries{copy: path, paste: path}
			}

		case ProgramWlClipboard:
			copyPath, copyErr := exec.LookPath(bins[0])
			pastePath, pasteErr := exec.LookPath(bins[1])
			if copyErr == nil && pasteErr == nil {
				binaries = Binaries{copy: copyPath, paste: pastePath}
			}
		}

		if !binaries.Empty() {
			return &linuxClipboard{
				program:  program,
				binaries: binaries,
			}, nil
		}
	}

	return nil, errors.New(
		"no supported clipboard found for linux platform, the following are currently supported: `xclip`, `xsel` and `wl-clipboard`",
	)
}

// Available implements NativeClipboard.
func (l *linuxClipboard) Available() bool {
	return l.binaries.copy != "" && l.binaries.paste != ""
}

// Binaries implements NativeClipboard.
func (l *linuxClipboard) Binaries() Binaries {
	return l.binaries
}

// Clear implements NativeClipboard.
func (l *linuxClipboard) Clear(ctx context.Context) error {
	return l.Write(ctx, []byte{})
}

// Read implements NativeClipboard.
func (l *linuxClipboard) Read(ctx context.Context) ([]byte, error) {
	if !l.Available() {
		return nil, errors.New("clipboard binaries not available")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	args := []string{}
	switch l.program {
	case ProgramXclip:
		args = []string{"-selection", "clipboard", "-o"}
	case ProgramXsel:
		args = []string{"--clipboard", "--output"}
	}

	cmd := exec.CommandContext(ctxWithTimeout, l.binaries.paste, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("(`%s`) failed to read from clipboard: %w", l.binaries.paste, err)
	}

	return output, nil
}

// Write implements NativeClipboard.
func (l *linuxClipboard) Write(ctx context.Context, data []byte) error {
	if !l.Available() {
		return errors.New("clipboard binaries not available")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	args := []string{}
	switch l.program {
	case ProgramXclip:
		args = []string{"-selection", "clipboard"}
	case ProgramXsel:
		args = []string{"--clipboard", "--input"}
	}

	cmd := exec.CommandContext(ctxWithTimeout, l.binaries.copy, args...)
	cmd.Stdin = bytes.NewReader(data)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("(`%s`) failed to write to clipboard: %w", l.binaries.copy, err)
	}

	return nil
}

var _ NativeClipboard = (*linuxClipboard)(nil)
