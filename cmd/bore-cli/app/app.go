package app

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/config"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/tui"
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
	// configPath is the path to the configuration file.
	configPath string

	// dataDir is the path to the data directory where data is stored.
	dataDir string

	bore          *bore.Bore
	handler       *handler.Handler
	tuiManager    *tui.Manager
	configManager *config.Manager
}

func New() (*App, error) {
	app := &App{
		configPath: defaultConfigPath(),
		dataDir:    defaultDataPath(),
		tuiManager: tui.NewViewManager(),
	}

	return app, nil
}

func (a *App) Execute() error {
	app := a.createRootCmd()
	return app.Run(os.Args)
}

func (a *App) SetConfigPath(path string) {
	a.configPath = path
}

func (a *App) SetDataDir(path string) {
	a.dataDir = path
}

func (a *App) Load() error {
	configManager, err := config.NewManager(config.Options{
		ConfigPath: a.configPath,
		DataDir:    a.dataDir,
	})
	if err != nil {
		return err
	}

	a.configManager = configManager

	config, err := a.configManager.Read()
	if err != nil {
		return err
	}

	bore, err := bore.New(config)
	if err != nil {
		return errors.New("failed to create bore instance: " + err.Error())
	}

	a.bore = bore
	a.handler = handler.New(bore, a.tuiManager, configManager)

	return nil
}
