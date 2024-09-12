package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/cmd/bore/app/models"
	"go.trulyao.dev/bore/cmd/bore/app/styles"
	"go.trulyao.dev/bore/pkg/config"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (a *App) ConfigCommand() *cli.Command {
	var (
		initCommand = &cli.Command{
			Name:  "init",
			Usage: "Create the default config file in the provided path",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "Path to the configuration file e.g. ./config.toml",
				},
			},
			Action: a.CreateConfigFile,
		}

		dumpCurrentConfigCommand = &cli.Command{
			Name:   "dump",
			Usage:  "Dump the current configuration",
			Action: a.DumpCurrentConfig,
		}
	)

	return &cli.Command{
		Name:  "config",
		Usage: "Manage configuration",
		Subcommands: []*cli.Command{
			initCommand,
			dumpCurrentConfigCommand,
		},
	}
}

func (a *App) DumpCurrentConfig(ctx *cli.Context) error {
	if a.config == nil {
		return fmt.Errorf("Config is not loaded")
	}

	if ctx.Bool("json") {
		jsonConfig, err := json.Marshal(a.config)
		if err != nil {
			return fmt.Errorf("Failed to marshal config to JSON: %s", err)
		}

		fmt.Fprintf(ctx.App.Writer, "%s", jsonConfig)
		return nil
	}

	columns := []table.Column{
		{Title: "Property", Width: 30},
		{Title: "Value", Width: 12},
		{Title: "Description", Width: 16},
	}

	rows := []table.Row{
		{"config.path", a.config.Path, "Path to the configuration file"},
		{
			"config.data_dir",
			a.config.DataDir,
			"Path to the data directory i.e. where the clipboard history is stored",
		},
		{
			"config.enable_native_clipboard",
			fmt.Sprintf("%t", a.config.EnableNativeClipboard),
			"Enable native clipboard passthrough (if not available, warnings will be shown)",
		},
		{
			"config.show_id_on_copy",
			fmt.Sprintf("%t", a.config.ShowIdOnCopy),
			"Show the ID of the copied content after a successful copy",
		},
		{
			"native_clipboard.is_available",
			fmt.Sprintf("%t", a.nativeClipboard.IsAvailable()),
			"Whether a native clipboard is available",
		},
		{
			"native_clipboard.paths.copy_bin",
			a.nativeClipboard.Paths().CopyBinPath,
			"Path to the system binary that copies data to the clipboard",
		},
		{
			"native_clipboard.paths.paste_bin",
			a.nativeClipboard.Paths().PasteBinPath,
			"Path to the system binary that pastes data from the clipboard",
		},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(8),
		table.WithFocused(true),
	)
	t.SetStyles(styles.TableStyle())

	m := models.NewConfigDumpModel(t)
	_, err := tea.NewProgram(m).Run()
	return err
}

func (a *App) CreateConfigFile(ctx *cli.Context) error {
	path := ctx.String("path")
	if path == "" {
		path = ctx.Args().First()
	}

	if path == "" {
		path = config.DefaultConfigFilePath()
	}

	// If the path is just a directory, append the default config file name
	s, err := os.Stat(path)
	if !strings.HasSuffix(path, ".toml") && (err == nil && s.IsDir()) {
		path = filepath.Join(path, "config.toml")
	}

	// Expand the path to handle ~, . etc
	path, err = filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path: %s", err)
	}

	defaultConfig := config.DefaultConfig()

	if err := config.WriteConfigToFile(defaultConfig, path); err != nil {
		return err
	}

	fmt.Fprintf(ctx.App.Writer, "Config file created at `%s`\n", path)
	return nil
}
