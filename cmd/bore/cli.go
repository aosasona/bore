package main

import (
	"go.trulyao.dev/bore/v2"
)

type Cli struct {
	// configPath is the path to the configuration file.
	configPath string

	// dataDir is the path to the data directory where data is stored.
	dataDir string

	bore *bore.Bore
}

func NewCli() (*Cli, error) {
	c := &Cli{}

	c.dataDir = defaultDataPath()
	c.configPath = defaultConfigPath()
	c.Load()

	return c, nil
}

func (c *Cli) SetConfigPath(path string) {
	c.configPath = path
}

func (c *Cli) SetDataDir(path string) {
	c.dataDir = path
}

func (c *Cli) Load() error {
	panic("not implemented yet")
}

// func loadConfig() (bore.Config, error) {
// 	configPath := defaultConfigPath()
// 	if configPath == "" {
// 		return bore.Config{}, errors.New("unable to determine default config path")
// 	}
//
// 	if err := createDefaultDirectories(); err != nil {
// 		return bore.Config{}, errors.New("failed to create default directories: " + err.Error())
// 	}
//
// 	if _, err := os.Stat(configPath); err != nil && !os.IsNotExist(err) {
// 		if err := createConfigFile(configPath); err != nil {
// 			return bore.Config{}, errors.New("failed to create config file: " + err.Error())
// 		}
// 	}
//
// 	return readConfigFile(configPath)
// }
//
// func readConfigFile(filePath string) (bore.Config, error) {
// 	configStr, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return bore.Config{}, errors.New("failed to read config file: " + err.Error())
// 	}
//
// 	var config bore.Config
// 	if _, err := config.FromBytes(configStr); err != nil {
// 		return bore.Config{}, errors.New("failed to read config file: " + err.Error())
// 	}
// 	config.DataDir = defaultDataPath()
//
// 	return config, nil
// }
//
// func createConfigFile(filePath string) error {
// 	config := bore.DefaultConfig()
// 	configStr, err := config.TOML()
// 	if err != nil {
// 		return errors.New("failed to convert config to string: " + err.Error())
// 	}
//
// 	if err := os.WriteFile(filePath, configStr, 0644); err != nil {
// 		return errors.New("failed to write config file: " + err.Error())
// 	}
//
// 	return nil
// }
//
// func createDefaultDirectories() error {
// 	if err := os.MkdirAll(defaultDataPath(), 0755); err != nil {
// 		return errors.New("failed to create data directory: " + err.Error())
// 	}
//
// 	configDir := path.Dir(defaultConfigPath())
// 	if err := os.MkdirAll(configDir, 0755); err != nil {
// 		return errors.New("failed to create storage directory: " + err.Error())
// 	}
//
// 	return nil
// }
