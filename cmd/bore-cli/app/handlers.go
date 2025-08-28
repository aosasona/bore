package app

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

func (a *App) copyCommand() *cli.Command {
	return &cli.Command{
		Name:  "copy",
		Usage: "Copy content to the clipboard",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "collection",
				Aliases: []string{"c"},
				Usage:   "Collection ID to associate with the copied content",
			},
			&cli.StringFlag{
				Name:    "mime-type",
				Aliases: []string{"m"},
				Usage:   "MIME type of the content being copied (e.g., text/plain, image/png)",
				Value:   "text/plain",
			},
			&cli.StringFlag{
				Name:    "input-file",
				Aliases: []string{"i"},
				Usage:   "Path to a file to read content from. If not provided, content will be read from stdin.",
				Value:   "",
			},
		},
		Action: func(ctx *cli.Context) error {
			panic("not implemented")
		},
	}
}

func (a *App) pasteCommand() *cli.Command {
	return &cli.Command{
		Name:  "paste",
		Usage: "Paste content from the clipboard",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "collection",
				Aliases: []string{"c"},
				Usage:   "Collection ID to paste content from",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Format to output the pasted content (text, json, base64)",
				Value:   string(handler.PasteFormatText),
			},
			&cli.BoolFlag{
				Name:    "from-system",
				Aliases: []string{"s"},
				Usage:   "Paste content from the system clipboard instead",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete the content from the clipboard after pasting",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "output-file",
				Aliases: []string{"o"},
				Usage:   "Path to a file where the pasted content will be saved. If not provided, content will be printed to stdout.",
				Value:   "",
			},
		},
		Action: func(ctx *cli.Context) error {
			panic("not implemented")
		},
	}
}
