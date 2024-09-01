package main

import (
	"os"

	"github.com/urfave/cli/v2"
	boreapp "go.trulyao.dev/bore/cmd/bore/app"
	"go.trulyao.dev/bore/pkg/config"
)

var (
	// TODD: find a way to init this before anything is setup
	app *boreapp.App
	err error

	version = "0.1.0"
)

func Execute() error {
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
			cli.ShowAppHelp(c)
			return nil
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
