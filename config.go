package bore

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// DataDir is the path to the storage directory.
	DataDir string `toml:"data_dir" json:"data_dir"`

	// ClipboardPassthrough enables passing the bore clipboard data to the native clipboard on copy
	ClipboardPassthrough bool `toml:"clipboard_passthrough" json:"clipboard_passthrough"`

	// DeleteOnPaste deletes the content after it has been pasted
	DeleteOnPaste bool `toml:"delete_on_paste" json:"delete_on_paste"`
}

// DefaultConfig returns the default configuration for the bore application.
func DefaultConfig() Config {
	return Config{
		DataDir:              ".",
		ClipboardPassthrough: true,
		DeleteOnPaste:        false,
	}
}

// FromBytes reads the configuration from a byte slice and populates the Config struct.
func (c *Config) FromBytes(data []byte) (*Config, error) {
	decoder := toml.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}

// TOML exports the configuration to a byte slice in TOML format.
func (c *Config) TOML() ([]byte, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
