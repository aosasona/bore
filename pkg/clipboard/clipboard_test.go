package clipboard_test

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	"go.trulyao.dev/bore/v2/pkg/clipboard"
)

func skipOnWSL(t *testing.T) {
	t.Helper()

	if runtime.GOOS != "linux" {
		return
	}

	if os.Getenv("WSL_INTEROP") != "" || os.Getenv("WSL_DISTRO_NAME") != "" {
		t.Skip("skipping native clipboard test on WSL")
	}

	version, err := os.ReadFile("/proc/version")
	if err != nil {
		return
	}

	if strings.Contains(strings.ToLower(string(version)), "microsoft") {
		t.Skip("skipping native clipboard test on WSL")
	}
}

func Test_IsAvailable(t *testing.T) {
	skipOnWSL(t)

	c, err := clipboard.NewNativeClipboard()
	if err != nil {
		t.Fatalf("Failed to create clipboard: %v", err)
	}

	available := c.Available()
	if !available {
		t.Error("Expected clipboard to be available, but it is not")
	}
}

func Test_SetAndGet(t *testing.T) {
	skipOnWSL(t)

	c, err := clipboard.NewNativeClipboard()
	if err != nil {
		t.Fatalf("Failed to create clipboard: %v", err)
	}

	testText := "Hello, Clipboard!"
	err = c.Write(context.TODO(), []byte(testText))
	if err != nil {
		t.Fatalf("Failed to set clipboard text: %v", err)
	}

	gotText, err := c.Read(context.TODO())
	if err != nil {
		t.Fatalf("Failed to get clipboard text: %v", err)
	}

	if string(gotText) != testText {
		t.Errorf("Expected clipboard text '%s', got '%s'", testText, gotText)
	}
}
