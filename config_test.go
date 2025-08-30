package bore_test

import (
	"testing"

	"go.trulyao.dev/bore/v2"
)

func Test_FromBytes(t *testing.T) {
	config := &bore.Config{}

	data := `
	data_dir = "~/.local/share/bore"
	clipboard_passthrough = false
	`

	_, err := config.FromBytes([]byte(data))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DataDir != "~/.local/share/bore" {
		t.Errorf("expected DataDir to be '~/.local/share/bore', got '%s'", config.DataDir)
	}

	if config.ClipboardPassthrough {
		t.Error("expected ClipboardPassthrough to be false, got true")
	}
}

func Test_TOML(t *testing.T) {
	config := &bore.Config{
		DataDir:              "~/.local/share/bore",
		ClipboardPassthrough: false,
	}

	tomlData, err := config.TOML()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `data_dir = "~/.local/share/bore"
clipboard_passthrough = false
delete_on_paste = true
`

	if string(tomlData) != expected {
		t.Errorf("expected TOML output to be:\n%s\nbut got:\n%s", expected, string(tomlData))
	}
}
