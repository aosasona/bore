package system

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

type ProgramType int

const (
	ProgXclip = iota
	ProgXsel
	ProgWlClipboard
)

type nativeClipboard struct {
	// Clipboard program type
	program ProgramType

	// Path to the binary responsible for copying
	copyBinPath string

	// Path to the binary responsible for pasting
	pasteBinPath string
}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	n := new(nativeClipboard)
	var (
		path string
		err  error
	)

	// Check for `xclip`
	if path, err = exec.LookPath("xclip"); err == nil {
		n.copyBinPath = path
		n.pasteBinPath = path
		n.program = ProgXclip
		return n, nil
	}

	// Check for `xsel`
	if path, err = exec.LookPath("xsel"); err == nil {
		n.copyBinPath = path
		n.pasteBinPath = path
		n.program = ProgXsel
		return n, nil
	}

	// Check for `wl-clipboard`
	// For this, we need to check for the presence of both binaries: wl-copy and wc-paste
	wlCopyPath, wcerr := exec.LookPath("wl-copy")
	wlPastePath, wperr := exec.LookPath("wl-paste")
	if (wlCopyPath != "" && wcerr == nil) && (wlPastePath != "" && wperr == nil) {
		n.copyBinPath = wlCopyPath
		n.pasteBinPath = wlPastePath
		n.program = ProgWlClipboard
		return n, nil
	}

	return n, errors.New(
		"no supported clipboard found, the following are currently supported: `xclip`, `xsel` and `wl-clipboard`",
	)
}

// IsAvailable checks if a clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	return (n.copyBinPath != "" && n.pasteBinPath != "")
}

// Copy copies the content to the system clipboard
func (n *nativeClipboard) Copy(content []byte) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, n.copyBinPath)
	cmd.Stdin = bytes.NewReader(content)

	return cmd.Run()
}

// Paste returns the last copied content from the system clipboard
func (n *nativeClipboard) Paste() ([]byte, error) {
	opts := []string{}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	if n.program == ProgXclip {
		opts = []string{"-o"}
	}

	cmd := exec.CommandContext(ctx, n.pasteBinPath, opts...)
	return cmd.Output()
}

func (n *nativeClipboard) Paths() ProgramPaths {
	return ProgramPaths{
		CopyBinPath:  n.copyBinPath,
		PasteBinPath: n.pasteBinPath,
	}
}

var _ NativeClipboardInterface = (*nativeClipboard)(nil)
