package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type nativeClipboard struct {
	// Path to the binary responsible for copying
	copyBinPath string

	// Path to the binary responsible for pasting
	pasteBinPath string
}

func NewNativeClipboard() (NativeClipboardInterface, error) {
	var err error
	n := new(nativeClipboard)

	n.copyBinPath, err = exec.LookPath("clip.exe")
	if n.copyBinPath == "" || err != nil {
		return n, fmt.Errorf("failed to find clip.exe: %s", err.Error())
	}

	n.pasteBinPath, err = exec.LookPath("powershell.exe")
	if n.pasteBinPath == "" || err != nil {
		return n, fmt.Errorf("failed to find powershell.exe: %s", err.Error())
	}

	return &nativeClipboard{}, nil
}

// IsAvailable checks if a clipboard is available on the current system
func (n *nativeClipboard) IsAvailable() bool {
	return n.copyBinPath != "" && n.pasteBinPath != ""
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
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, n.pasteBinPath, "Get-Clipboard")
	return cmd.Output()
}

// Clear clears the system clipboard
func (n *nativeClipboard) Clear() error {
	return n.Copy([]byte{})
}

func (n *nativeClipboard) Paths() ProgramPaths {
	return ProgramPaths{
		CopyBinPath:  n.copyBinPath,
		PasteBinPath: n.pasteBinPath,
	}
}

var _ NativeClipboardInterface = (*nativeClipboard)(nil)
