package bore

type Config struct {
	// DataDir is the path to the storage directory.
	DataDir string `toml:"data_path" json:"data_path"`

	// SystemClipboardPassthrough enables passing the bore clipboard data to the native clipboard on copy
	SystemClipboardPassthrough bool `toml:"system_clipboard_passthrough" json:"system_clipboard_passthrough"`

	// DeleteOnPaste deletes the content after it has been pasted
	DeleteOnPaste bool `toml:"delete_on_paste" json:"delete_on_paste"`
}

// DefaultConfig returns the default configuration for the bore application.
func DefaultConfig() Config {
	return Config{
		DataDir:                    ".",
		SystemClipboardPassthrough: true,
		DeleteOnPaste:              false,
	}
}

// Load reads the configuration from a byte slice and populates the Config struct.
func (c *Config) Load(data []byte) error {
	panic("not implemented")
}

// TOML exports the configuration to a byte slice in TOML format.
func (c *Config) TOML() ([]byte, error) {
	// Export the configuration to a byte slice in TOML format.
	// This is a placeholder implementation.
	panic("not implemented")
}
