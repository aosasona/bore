package app

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

func pipedIn() bool {
	// If the program was piped into, we need to read from stdin and copy that
	fileinfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read from stdin:", err)
		return false
	}

	return (fileinfo.Mode() & os.ModeCharDevice) == 0
}

func (a *App) createRootCmd() *cli.App {
	// nolint:exhaustruct
	return &cli.App{
		Name:                 "bore",
		Usage:                "A clipboard manager for the terminal",
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
			if pipedIn() {
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
			a.SetConfigPath(ctx.String("config"))
			a.SetDataDir(ctx.String("data-dir"))

			return a.Load()
		},
		Commands: []*cli.Command{
			a.infoCommand(),
			a.resetCommand(),
			a.copyCommand(),
			a.pasteCommand(),
			a.collectionsCommand(),
		},
	}
}

func (a *App) infoCommand() *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "info",
		Usage: "Display information about the current bore instance",
		Action: func(ctx *cli.Context) error {
			config, err := a.bore.Config()
			if err != nil {
				return cli.Exit("failed to get bore configuration: "+err.Error(), 1)
			}

			fmt.Println("Bore Version:", Version)
			fmt.Println("Data Directory:", a.dataDir)
			fmt.Println("Config Path:", a.configPath)
			fmt.Println("Default Collection:", config.DefaultCollection)
			fmt.Println("Clipboard Passthrough:", config.ClipboardPassthrough)
			return nil
		},
	}
}

func (a *App) resetCommand() *cli.Command {
	// nolint:exhaustruct
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

func (a *App) copyCommand() *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "copy",
		Usage: "Copy content to the clipboard",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    handler.FlagCollection,
				Aliases: []string{"c"},
				Usage:   "Collection ID to associate with the copied content",
			},
			&cli.StringFlag{
				Name:    handler.FlagMimeType,
				Aliases: []string{"m"},
				Usage:   "MIME type of the content being copied (e.g., text/plain, image/png)",
				Value:   "text/plain",
			},
			&cli.StringFlag{
				Name:    handler.FlagInputFile,
				Aliases: []string{"i"},
				Usage:   "Path to a file to read content from. If not provided, content will be read from stdin.",
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    handler.FlagSystem,
				Aliases: []string{"s"},
				Usage:   "Copy content to the system clipboard ONLY",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    handler.FlagFormat,
				Aliases: []string{"f"},
				Usage:   "Input format of the content being copied (text, base64)",
			},
		},
		Args:      true,
		ArgsUsage: "[content]",
		Action: func(ctx *cli.Context) error {
			return a.handler.Copy(ctx, handler.CliCopyOptions{Stdin: pipedIn()})
		},
	}
}

func (a *App) pasteCommand() *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "paste",
		Usage: "Paste content from the clipboard",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    handler.FlagCollection,
				Aliases: []string{"c"},
				Usage:   "Collection ID to paste content from",
			},
			&cli.StringFlag{
				Name:        handler.FlagFormat,
				Aliases:     []string{"f"},
				Usage:       "Format to output the pasted content (text, json, base64)",
				Value:       string(handler.PasteFormatText),
				DefaultText: string(handler.PasteFormatText),
			},
			&cli.BoolFlag{
				Name:    handler.FlagSystem,
				Aliases: []string{"s"},
				Usage:   "Paste content from the system clipboard instead",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    handler.FlagDelete,
				Aliases: []string{"d"},
				Usage:   "Delete the content from the clipboard after pasting",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    handler.FlagOutputFile,
				Aliases: []string{"o"},
				Usage:   "Path to a file where the pasted content will be saved. If not provided, content will be printed to stdout.",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    handler.FlagIdentifier,
				Aliases: []string{"id"},
				Usage:   "Identifier of the specific clipboard entry to paste. If not provided, the most recent entry will be used.",
				Value:   "",
			},
		},
		Action: func(ctx *cli.Context) error {
			return a.handler.Paste(ctx)
		},
	}
}

func (a *App) collectionsCommand() *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "collections",
		Usage: "Manage clipboard collections",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all collections",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    handler.FlagFormat,
						Aliases: []string{"f"},
						Usage:   "Output format (text, json)",
					},
				},
				Action: func(ctx *cli.Context) error {
					return a.handler.ListCollections(ctx)
				},
			},
			{
				Name:      "create",
				Usage:     "Create a new collection",
				ArgsUsage: "[collection name]",
				Args:      true,
				Action: func(ctx *cli.Context) error {
					return a.handler.CreateCollection(ctx)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a collection",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     handler.FlagCollection,
						Aliases:  []string{"c"},
						Usage:    "ID of the collection to delete",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    handler.FlagForce,
						Aliases: []string{"f"},
						Usage:   "Force deletion without confirmation",
						Value:   false,
					},
				},
				Action: func(ctx *cli.Context) error {
					return a.handler.DeleteCollection(ctx)
				},
			},
			{
				Name:  "set-default",
				Usage: "Set the default collection",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     handler.FlagCollection,
						Aliases:  []string{"c"},
						Usage:    "ID of the collection to set as default",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return a.handler.SetDefaultCollection(ctx)
				},
			},
		},
	}
}
