package bore

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/oklog/ulid/v2"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

// NOTE: where and show the configuration is stored is determined by the application layer.
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

func (c *Config) SetDefaultCollection(identifier string) error {
	id, err := ulid.Parse(identifier)
	if err != nil {
		return errs.Wrap(err, "invalid collection identifier")
	}

	c.DefaultCollection = id.String()
	return nil
}

// FromBytes reads the configuration from a byte slice and populates the Config struct.
func (c *Config) FromBytes(data []byte) (*Config, error) {
	decoder := toml.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Decode(c); err != nil {
		return nil, errs.Wrap(err, "failed to decode configuration")
	}

	return c, nil
}

// TOML exports the configuration to a byte slice in TOML format.
func (c *Config) TOML() ([]byte, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(c); err != nil {
		return nil, errs.Wrap(err, "failed to encode configuration")
	}
	return buf.Bytes(), nil
}
