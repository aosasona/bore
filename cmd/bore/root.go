package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/cmd/bore/commands"
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
	Commands: []*cli.Command{
		commands.Config,
	},
}

func Execute() error {
	return root.Run(os.Args)
}
