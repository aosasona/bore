package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// DataPath is the path to the directory where the application stores its data
	DataPath string `toml:"data_path"`

	// EnableNativeClipboard enables passing clipboard data to the native clipboard on copy
	EnableNativeClipboard bool `toml:"enable_native_clipboard"`
}

func DefaultConfig() *Config {
	return &Config{
		DataPath:              "",
		EnableNativeClipboard: false,
	}
}

func WriteConfigToFile(config *Config, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to create config file: %s", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(config); err != nil {
		return fmt.Errorf("Failed to write config file: %s", err)
	}

	return nil
}

// Parse the configuration file and return a Config object
func ParseConfig(path string) (*Config, error) {
	config := new(Config)

	if path == "" {
		return DefaultConfig(), nil
	}

	if path[0] == '~' {
		// Expand the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to expand tilde to home directory: %s", err)
		}
		path = home + path[1:]
	}

	// We will automatically fallback to the default configuration if the file is not found
	s, err := os.Stat(path)
	if err != nil {
		return DefaultConfig(), nil
	}

	if s.IsDir() {
		return config, fmt.Errorf("Config file path is a directory")
	}

	_, err = toml.DecodeFile(path, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %s", err)
	}

	return config, nil
}
