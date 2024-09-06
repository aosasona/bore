package system

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

type nativeClipboard struct {
	// Path to `pbcopy`
	copyBinPath string

	// Path to `pbpaste`
	pasteBinPath string
}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	var err error

	n := new(nativeClipboard)

	n.copyBinPath, err = exec.LookPath("pbcopy")
	if n.copyBinPath == "" || err != nil {
		return n, errors.New("`pbcopy` not available")
	}

	n.pasteBinPath, err = exec.LookPath("pbpaste")
	if n.pasteBinPath == "" || err != nil {
		return n, errors.New("`pbpaste` not available")
	}

	return n, nil
}

// IsAvailable checks if a clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	return (n.copyBinPath != "" && n.pasteBinPath != "")
}

// Copy copies the content to the system clipboard
func (n *nativeClipboard) Copy(content []byte) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, n.copyBinPath)
	cmd.Stdin = bytes.NewReader(content)

	return cmd.Run()
}

// Paste returns the last copied content from the system clipboard
func (n *nativeClipboard) Paste() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, n.pasteBinPath)
	return cmd.Output()
}

// Paths returns the programs used for copying and pasting
func (n *nativeClipboard) Paths() ProgramPaths {
	return ProgramPaths{
		CopyBinPath:  n.copyBinPath,
		PasteBinPath: n.pasteBinPath,
	}
}

var _ NativeClipboardInterface = (*nativeClipboard)(nil)
