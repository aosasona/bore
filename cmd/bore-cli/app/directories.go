package app

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/mitchellh/go-homedir"
)

var directories = userdirs.ForApp("bore", "trulyao", "dev.trulyao.bore")

// defaultConfigPath returns the default path to the configuration file.
func defaultConfigPath() string {
	switch runtime.GOOS {
	// For UNIX-like systems, we will stick with the more familiar XDG Base Directory Specification since users are more likely to look at something like ~/.config/bore/config.toml even on MacOS
	case "darwin", "linux":
		xdgConfigDir := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigDir != "" {
			return filepath.Join(xdgConfigDir, "bore", "config.toml")
		}

		// If the XDG_CONFIG_HOME environment variable is not set, we fall back to the default
		// location of ~/.config/bore/config.toml
		home, err := homedir.Dir()
		if err != nil {
			home = "~"
		}
		return filepath.Join(home, ".config", "bore", "config.toml")

	default:
		return directories.CachePath("config.toml")
	}
}

// defaultDataPath returns the default data directory for the application.
// NOTE: The directory may or may not exist, depending on the platform and user configuration.
func defaultDataPath() string {
	switch runtime.GOOS {
	case "darwin", "linux":
		xdgDataDir := os.Getenv("XDG_DATA_HOME")
		if xdgDataDir != "" {
			return filepath.Join(xdgDataDir, "bore")
		}

		// If the XDG_DATA_HOME environment variable is not set, we fall back to the default location of ~/.local/share/bore
		home, err := homedir.Dir()
		if err != nil {
			home = "~"
		}

		return filepath.Join(home, ".local", "share", "bore")

	case "windows":
		if len(directories.DataDirs) > 0 {
			return directories.DataDirs[0]
		}

		home, err := homedir.Dir()
		if err != nil {
			panic("failed to get home directory: " + err.Error())
		}

		return filepath.Join(home, "AppData", "Local", "bore")

	default:
		if len(directories.DataDirs) > 0 {
			return directories.DataDirs[0]
		}

		panic("unsupported operating system: " + runtime.GOOS)
	}
}
