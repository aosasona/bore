package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Path is the path to the configuration file
	Path string `toml:"-" json:"config_path"`

	// DataDir is the path to the directory where the application stores its data
	DataDir string `toml:"data_dir" json:"data_dir"`

	// EnableNativeClipboard enables passing clipboard data to the native clipboard on copy
	EnableNativeClipboard bool `toml:"enable_native_clipboard" json:"enable_native_clipboard"`

	// ShowIdOnCopy shows the ID of the copied content after a successful copy
	ShowIdOnCopy bool `toml:"show_id_on_copy" json:"show_id_on_copy"`
}

func CreateDirIfNotExists(path string) error {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("Failed to create data directory: %s", err)
		}
	}

	if err != nil {
		return fmt.Errorf("Failed to check data directory: %s", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("Data directory is not a directory: %s", path)
	}

	return nil
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
	var (
		configRoot string
		err        error
	)

	if runtime.GOOS == "darwin" {
		configRoot, err = os.UserHomeDir()
		if err != nil {
			configRoot = "~"
		}
		configRoot = filepath.Join(configRoot, ".config")
	} else {
		configRoot, err = os.UserConfigDir()
	}

	return filepath.Join(configRoot, "bore", "config.toml")
}

func DefaultConfig() *Config {
	return &Config{
		Path:                  DefaultConfigFilePath(),
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

	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Config file does not exist: %s", path)
		}

		return nil, fmt.Errorf("Failed to check config file: %s", err)
	}

	if s.IsDir() {
		return config, fmt.Errorf("Config file path is a directory")
	}

	_, err = toml.DecodeFile(path, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %s", err)
	}

	config.Path = path

	return config, nil
}

// Load the configuration file and store it in the global variable
func Load(path string) (*Config, error) {
	config, err := ParseConfig(path)
	if err != nil {
		return nil, err
	}

	CreateDirIfNotExists(config.DataDir)
	return config, nil
}
