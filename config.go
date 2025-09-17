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

	// DefaultCollection is the name of the default collection to use when none is specified
	DefaultCollection string `toml:"default_collection" json:"default_collection"`
}

// DefaultConfig returns the default configuration for the bore application.
func DefaultConfig() Config {
	return Config{
		DataDir:              ".",
		ClipboardPassthrough: true,
		DefaultCollection:    "",
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
