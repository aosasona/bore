package commands

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

type PasteCommand struct{}

func (p PasteCommand) Build(run Runner) *cli.Command {
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
			return run(ctx, p)
		},
	}
}

func (p PasteCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.Paste(ctx)
}
