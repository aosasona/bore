package main

import (
	"errors"
	"os"

	"go.trulyao.dev/bore/v2"
)

type Cli struct {
	bore *bore.Bore
}

func NewCli() (*Cli, error) {
	// Load the configuration from the default config path.
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	config.DataDir = defaultDataPath()

	bore, err := bore.New(&config)
	if err != nil {
		return nil, err
	}

	return &Cli{bore: bore}, nil
}

func loadConfig() (bore.Config, error) {
	configPath := defaultConfigPath()
	if configPath == "" {
		return bore.Config{}, errors.New("unable to determine default config path")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := createConfigFile(configPath); err != nil {
				return bore.Config{}, errors.New("failed to create config file: " + err.Error())
			}
			return bore.Config{}, nil
		}
	}

	return readConfigFile(configPath)
}

func readConfigFile(path string) (bore.Config, error) {
	configStr, err := os.ReadFile(path)
	if err != nil {
		return bore.Config{}, errors.New("failed to read config file: " + err.Error())
	}

	var config bore.Config
	if err := config.Load(configStr); err != nil {
		return bore.Config{}, errors.New("failed to read config file: " + err.Error())
	}
	config.DataDir = defaultDataPath()

	return config, nil
}

func createConfigFile(path string) error {
	config := bore.DefaultConfig()
	configStr, err := config.TOML()
	if err != nil {
		return errors.New("failed to convert config to string: " + err.Error())
	}

	if err := os.WriteFile(path, configStr, 0644); err != nil {
		return errors.New("failed to write config file: " + err.Error())
	}

	return nil
}
