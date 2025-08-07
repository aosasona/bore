package app

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func (a *App) createRootCmd() *cli.App {
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
			a.SetConfigPath(ctx.String("config"))
			a.SetDataDir(ctx.String("data-dir"))

			return a.Load()
		},
		Commands: []*cli.Command{
			a.infoCommand(),
		},
	}
}

func (a *App) infoCommand() *cli.Command {
	return &cli.Command{
		Name:  "info",
		Usage: "Display information about the current bore instance",
		Action: func(ctx *cli.Context) error {
			deviceID := a.bore.DeviceID()
			if deviceID == "" {
				deviceID = "<not set>"
			}

			fmt.Println("Bore Version:", Version)
			fmt.Println("Data Directory:", a.dataDir)
			fmt.Println("Config Path:", a.configPath)
			fmt.Println("Device ID:", deviceID)
			fmt.Println("Clipboard Passthrough:", a.bore.Config().ClipboardPassthrough)
			fmt.Println("Delete on Paste:", a.bore.Config().DeleteOnPaste)
			return nil
		},
	}
}

func (a *App) resetCommand() *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: "Reset the bore instance, clearing all data",
		Action: func(ctx *cli.Context) error {
			if err := a.bore.Reset(); err != nil {
				return cli.Exit("failed to reset bore: "+err.Error(), 1)
			}

			if err := os.Remove(a.configPath); err != nil {
				return cli.Exit("failed to remove configuration file: "+err.Error(), 1)
			}

			fmt.Println("Bore instance has been reset successfully.")
			return nil
		},
	}
}
