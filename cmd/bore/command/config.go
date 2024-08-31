package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/pkg/config"
)

var (
	initCommand = &cli.Command{
		Name:  "init",
		Usage: "Create the default config file in the provided path",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Path to the configuration file",
				Value:   config.DefaultConfigFilePath(),
			},
		},
		Action: handleInitAction,
	}

	dumpCurrentConfigCommand = &cli.Command{
		Name:   "dump",
		Usage:  "Dump the current configuration",
		Action: handleDumpCurrentConfigAction,
	}

	Config = &cli.Command{
		Name:  "config",
		Usage: "Manage configuration",
		Subcommands: []*cli.Command{
			initCommand,
			dumpCurrentConfigCommand,
		},
	}
)

func handleDumpCurrentConfigAction(c *cli.Context) error {
	config := config.Get()
	fmt.Fprintf(c.App.Writer, "%#v", config)
	return nil
}

func handleInitAction(c *cli.Context) error {
	path := c.String("path")
	if path == "" {
		return fmt.Errorf("No path provided")
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

	fmt.Fprintf(c.App.Writer, "Config file created at `%s`\n", path)
	return nil
}
