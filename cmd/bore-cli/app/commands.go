package app

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/commands"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

func (a *App) createRootCmd() (*cli.App, error) {
	var commandManager *commands.Manager

	// nolint:exhaustruct
	rootCmd := &cli.App{
		Name:                 "bore",
		Usage:                "An SQLite-backed clipboard manager for headless environments",
		Version:              Version,
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{Name: "Ayodeji O.", Email: "ayodeji@trulyao.dev"},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() > 0 {
				return cli.ShowAppHelp(ctx)
			}

			// If the program was piped into, we need to read from stdin and copy that
			if commands.PipedIn() {
				return a.handler.Copy(ctx, handler.CliCopyOptions{Stdin: true})
			}

			return a.handler.Paste(ctx)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    handler.FlagConfig,
				Aliases: []string{"c"},
				Usage:   "Path to the configuration file",
				Value:   defaultConfigPath(),
			},
			&cli.StringFlag{
				Name:    handler.FlagDataDir,
				Aliases: []string{"d"},
				Usage:   "Path to the data directory where data is stored",
				Value:   defaultDataPath(),
			},
			&cli.StringFlag{
				Name:    handler.FlagFormat,
				Aliases: []string{"f"},
				Usage:   "Output format for the current command (e.g., json, base64, text)",
			},
			&cli.StringFlag{
				Name:    handler.FlagOutputFile,
				Aliases: []string{"o"},
				Usage:   "Path to a file where the pasted content will be saved. If not provided, content will be printed to stdout.",
				Value:   "",
			},
		},
		Before: func(ctx *cli.Context) error {
			var err error

			a.SetConfigPath(ctx.String("config"))
			a.SetDataDir(ctx.String("data-dir"))

			if err = a.Load(); err != nil {
				return cli.Exit("failed to load configuration: "+err.Error(), 1)
			}

			commandOptions := commands.NewCommandOptions{
				Handler: a.handler,
				Bore:    a.bore,
				TUI:     a.tuiManager,
				Config:  a.configManager,
				Version: Version,
			}

			commandManager, err = commands.NewManager(&commandOptions)
			if err != nil {
				return cli.Exit("failed to initialize commands: "+err.Error(), 1)
			}

			return nil
		},
	}

	rootCmd.Commands = commands.BuildAll(func() *commands.Manager {
		return commandManager
	})

	return rootCmd, nil
}
