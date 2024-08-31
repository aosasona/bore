package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// DataDir is the path to the directory where the application stores its data
	DataDir string `toml:"data_dir"`

	// EnableNativeClipboard enables passing clipboard data to the native clipboard on copy
	EnableNativeClipboard bool `toml:"enable_native_clipboard"`
}

// DefaultDataDir returns the default data directory path
func DefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	return home + "/.bore"
}

// DefaultConfigFilePath returns the default configuration file path
func DefaultConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	return filepath.Join(home, ".config", "bore", "config.toml")
}

func DefaultConfig() *Config {
	return &Config{
		DataDir:               DefaultDataDir(),
		EnableNativeClipboard: false,
	}
}

func WriteConfigToFile(config *Config, path string) error {
	// Ensure that the parent directory exists or create it
	dir := filepath.Dir(path)

	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("Failed to create parent directory: %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("Failed to check parent directory: %s", err)
	} else if !stat.IsDir() {
		return fmt.Errorf("Parent directory is not a directory: %s", dir)
	}

	// Create the config file
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
// NOTE: probably switch to https://github.com/knadh/koanf if we need to load other formats like JSON and YAML in the future
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
