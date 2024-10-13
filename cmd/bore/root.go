package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	boreapp "go.trulyao.dev/bore/cmd/bore/app"
	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/handler"
)

var (
	app           *boreapp.App
	dirConfigPath = ""

	version = "source"

	dirPaths = []string{"./bore.toml", "./.bore/config.toml"}
)

func Execute() error {
	configPath := config.DefaultConfigFilePath()

	// Check for the presence of the config files in the current directory (`./bore.toml` or `.bore/config.toml`)
	for _, path := range dirPaths {
		if s, err := os.Stat(path); err == nil && !s.IsDir() {
			dirConfigPath = path
			break
		}
	}

	if dirConfigPath != "" {
		configPath = dirConfigPath
	}

	var err error
	app, err = boreapp.New(configPath)
	if err != nil {
		return err
	}

	root := CreateRootCommand()
	return root.Run(os.Args)
}

func CreateRootCommand() *cli.App {
	return &cli.App{
		Name:    "bore",
		Usage:   "A minimal clipboard manager for terminal/headless environments",
		Version: version,
		Action: func(c *cli.Context) error {
			// If there was any argument passed, simply raise the help message
			if c.NArg() > 0 {
				return cli.ShowAppHelp(c)
			}

			// If the program was piped to, read directly from STDIN
			fi, _ := os.Stdin.Stat()
			if (fi.Mode() & os.ModeCharDevice) == 0 {
				return app.Copy(c)
			}

			return app.Paste(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output the result in JSON format",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "The format of the content to copy. Available formats: base64, plain-text",
				Value:   handler.FormatPlainText.String(),
			},
		},

		// Global flags have to be passed in BEFORE subcommands e.g `bore -c config.toml config dump`
		Before: func(c *cli.Context) error {
			configPath := strings.TrimSpace(c.String("config"))

			// Use the directory's config file if it exists and the user hasn't specified a config file
			if dirConfigPath != "" && configPath == "" {
				configPath = dirConfigPath
			}

			// If we still do not have a config file at this point, it is safe to use the default config file
			if configPath == "" {
				configPath = config.DefaultConfigFilePath()
			}

			return app.UpdateConfigPath(configPath)
		},
		Authors: []*cli.Author{
			{
				Name:  "Ayodeji O.",
				Email: "ayodeji@trulyao.dev",
			},
		},
		Commands: []*cli.Command{
			app.ConfigCommand(),
			app.CopyCommand(),
			app.PasteCommand(),
		},
	}
}
