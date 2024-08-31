package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/cmd/bore/command"
	"go.trulyao.dev/bore/pkg/config"
)

var root = &cli.App{
	Name:  "bore",
	Usage: "A minimal clipboard manager for terminal/headless environments",
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
	},
	Authors: []*cli.Author{
		{
			Name:  "Ayodeji O.",
			Email: "ayodeji@trulyao.dev",
		},
	},
	Before: func(c *cli.Context) error {
		configPath := c.String("config")
		if err := config.Load(configPath); err != nil {
			return err
		}

		return nil
	},
	Commands: []*cli.Command{
		command.Config,
	},
}

func Execute() error {
	return root.Run(os.Args)
}
