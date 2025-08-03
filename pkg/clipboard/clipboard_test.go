package clipboard_test

import (
	"context"
	"testing"

	"go.trulyao.dev/bore/v2/pkg/clipboard"
)

func Test_IsAvailable(t *testing.T) {
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
