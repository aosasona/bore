package main

import (
	"os"

	"github.com/urfave/cli/v2"
	boreapp "go.trulyao.dev/bore/cmd/bore/app"
	"go.trulyao.dev/bore/pkg/config"
	"go.trulyao.dev/bore/pkg/handler"
)

var (
	app *boreapp.App

	version = "latest"
)

func Execute() error {
	var err error
	app, err = boreapp.New(config.DefaultConfigFilePath())
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
				Value:   config.DefaultConfigFilePath(),
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
			return app.UpdateConfigPath(c.String("config"))
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
