package app

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/view"
)

var Version string

func init() {
	if Version != "" {
		return // was set by -ldflags
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	if output, err := cmd.Output(); err == nil {
		latestHash := strings.TrimSpace(string(output))
		if len(latestHash) > 7 {
			Version = latestHash[:7]
		}
	}

	if Version == "" {
		Version = "dev"
	}
}

type App struct {
	bore *bore.Bore

	// configPath is the path to the configuration file.
	configPath string

	// dataDir is the path to the data directory where data is stored.
	dataDir string

	handler *handler.Handler

	views *view.ViewManager
}

func New() (*App, error) {
	app := &App{
		configPath: defaultConfigPath(),
		dataDir:    defaultDataPath(),
		views:      view.NewViewManager(),
	}

	return app, nil
}

func (a *App) Execute() error {
	app := a.createRootCmd()
	if err := app.Run(os.Args); err != nil {
		return cli.Exit("error: "+err.Error(), 1)
	}
	return nil
}

func (a *App) SetConfigPath(path string) {
	a.configPath = path
}

func (a *App) SetDataDir(path string) {
	a.dataDir = path
}

func (a *App) Load() error {
	if err := a.createDirectories(); err != nil {
		return errors.New("failed to create directories: " + err.Error())
	}

	configExists, err := a.configFileExists()
	if err != nil {
		return err
	}

	if !configExists {
		if err := a.createDefaultConfigFile(); err != nil {
			return err
		}
	}

	config, err := a.readConfig()
	if err != nil {
		return err
	}

	bore, err := bore.New(config)
	if err != nil {
		return errors.New("failed to create bore instance: " + err.Error())
	}

	a.bore = bore
	a.handler = handler.New(bore, a.views)

	return nil
}

func (a *App) createDirectories() error {
	configDir := path.Dir(a.configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}

	if err := os.MkdirAll(a.dataDir, 0o755); err != nil {
		return errors.New("failed to create data directory: " + err.Error())
	}

	return nil
}

func (a *App) configFileExists() (bool, error) {
	if _, err := os.Stat(a.configPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.New("failed to check config file: " + err.Error())
	}
	return true, nil
}

func (a *App) createDefaultConfigFile() error {
	config := bore.DefaultConfig()
	configStr, err := config.TOML()
	if err != nil {
		return errors.New("failed to convert config to string: " + err.Error())
	}

	if err := os.WriteFile(a.configPath, configStr, 0o644); err != nil {
		return errors.New("failed to write config file: " + err.Error())
	}

	return nil
}

func (a *App) readConfig() (*bore.Config, error) {
	configStr, err := os.ReadFile(a.configPath)
	if err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	config := new(bore.Config)
	if _, err := config.FromBytes(configStr); err != nil {
		return nil, errors.New("failed to parse config file: " + err.Error())
	}
	config.DataDir = a.dataDir

	return config, nil
}
