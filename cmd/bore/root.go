package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "source"

// TODO: move path to the CLI itself (can be overriden by using the -c flag)
// TODO: automatically detect the data directory for each platform and remove this config option
func (c *Cli) execute() error {
	panic("not implemented")
}

func (c *Cli) createRootCmd() *cli.App {
	return &cli.App{
		Name:    "bore",
		Usage:   "A clipboard manager for the terminal",
		Version: Version,
		Authors: []*cli.Author{
			{Name: "Ayodeji O.", Email: "ayodeji@trulyao.dev"},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() > 0 {
				return cli.ShowAppHelp(ctx)
			}

			// If the program was piped into, we need to read from stdin and copy that
			fileinfo, err := os.Stdin.Stat()
			if err != nil {
				return cli.Exit("failed to read from stdin: "+err.Error(), 1)
			}

			if (fileinfo.Mode() & os.ModeCharDevice) == 0 {
				panic("copy not implemented yet")
			}

			panic("paste not implemented yet")
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the configuration file",
				Value:   defaultConfigPath(),
			},
			&cli.StringFlag{
				Name:    "data-dir",
				Aliases: []string{"d"},
				Usage:   "Path to the data directory where data is stored",
				Value:   defaultDataPath(),
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"o"},
				Usage:   "Output format for the current command (e.g., json, base64, text)",
			},
		},
		Before: func(ctx *cli.Context) error {
			c.SetConfigPath(ctx.String("config"))
			c.SetDataDir(ctx.String("data-dir"))

			return c.InstantiateBore()
		},
	}
}
