package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/pkg/config"
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
	output := fmt.Sprintf(`Loaded configuration from %s
DataDir: %s
EnableNativeClipboard: %t
`, a.config.Path, a.config.DataDir, a.config.EnableNativeClipboard)

	fmt.Fprintf(ctx.App.Writer, "%s", output)
	return nil
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
