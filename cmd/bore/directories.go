package main

import (
	"path/filepath"
	"runtime"

	"github.com/apparentlymart/go-userdirs/userdirs"
	"github.com/mitchellh/go-homedir"
)

var directories = userdirs.ForApp("bore", "trulyao", "dev.trulyao.bore")

// defaultConfigPath returns the default path to the configuration file.
func defaultConfigPath() string {
	return directories.CachePath("config.toml")
}

// defaultDataDir returns the default data directory for the application.
// NOTE: The directory may or may not exist, depending on the platform and user configuration.
func defaultDataDir() string {
	if len(directories.DataDirs) > 0 {
		return directories.DataDirs[0]
	}

	home, err := homedir.Dir()
	if err != nil {
		panic("failed to get home directory: " + err.Error())
	}

	var dataDir string
	switch runtime.GOOS {
	case "windows":
		dataDir = filepath.Join(home, "AppData", "Local", "bore")
	case "darwin", "linux":
		// On macOS and Linux, we use the XDG Base Directory Specification to determine the data directory.
		// Although this is not ideal on macOS, it is a common practice by most CLI tools.
		dataDir = filepath.Join(home, ".local", "share", "bore")
	default:
		panic("unsupported operating system: " + runtime.GOOS)
	}

	return dataDir
}
