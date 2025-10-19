package config

import (
	"os"
	"path"

	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

type (
	Manager struct {
		config     *bore.Config
		dataDir    string
		configPath string
	}

	Options struct {
		DataDir    string
		ConfigPath string
	}
)

func NewManager(opts Options) (*Manager, error) {
	m := &Manager{
		config:     nil,
		dataDir:    opts.DataDir,
		configPath: opts.ConfigPath,
	}

	if err := m.createDirectories(); err != nil {
		return nil, err
	}

	exists, err := m.ConfigExists()
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := m.WriteDefault(); err != nil {
			return nil, errs.Wrap(err, "failed to write default config")
		}
	}

	return m, nil
}

func (m *Manager) ConfigPath() string { return m.configPath }

func (m *Manager) DataDir() string   { return m.dataDir }
func (m *Manager) ConfigDir() string { return path.Dir(m.configPath) }

func (m *Manager) createDirectories() error {
	configDir := path.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return errs.Wrap(err, "failed to create config directory")
	}

	if err := os.MkdirAll(m.dataDir, 0o700); err != nil {
		return errs.Wrap(err, "failed to create data directory")
	}

	return nil
}

func (m *Manager) ConfigExists() (bool, error) {
	if _, err := os.Stat(m.configPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errs.Wrap(err, "failed to stat config file")
	}
	return true, nil
}

func (m *Manager) WriteDefault() error {
	defaultConfig := bore.DefaultConfig()
	defaultConfig.DataDir = m.dataDir
	return m.Write(&defaultConfig)
}

func (m *Manager) Write(config *bore.Config) error {
	configBytes, err := config.TOML()
	if err != nil {
		return errs.New("failed to serialize config: " + err.Error())
	}

	if err := os.WriteFile(m.configPath, configBytes, 0o600); err != nil {
		return errs.New("failed to write config file: " + err.Error())
	}

	return nil
}

func (m *Manager) Read() (*bore.Config, error) {
	if m.config != nil {
		return m.config, nil
	}

	configStr, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, errs.New("failed to read config file: " + err.Error())
	}

	config := new(bore.Config)
	if _, err := config.FromBytes(configStr); err != nil {
		return nil, errs.New("failed to parse config file: " + err.Error())
	}
	config.DataDir = m.dataDir
	m.config = config

	return m.config, nil
}

func (m *Manager) SetDefaultCollectionID(identifier string) error {
	config, err := m.Read()
	if err != nil {
		return errs.Wrap(err, "failed to read config")
	}

	if err := config.SetDefaultCollection(identifier); err != nil {
		return errs.Wrap(err, "failed to set default collection")
	}

	if err := m.Write(config); err != nil {
		return errs.Wrap(err, "failed to write config")
	}

	return nil
}

func (m *Manager) UnsetDefaultCollectionID() error {
	config, err := m.Read()
	if err != nil {
		return errs.Wrap(err, "failed to read config")
	}

	config.DefaultCollection = ""

	if err := m.Write(config); err != nil {
		return errs.Wrap(err, "failed to write config")
	}

	return nil
}
