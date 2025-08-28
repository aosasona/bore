package app

import "github.com/urfave/cli/v2"

type PasteFormat string

const (
	PasteFormatText   PasteFormat = "text"
	PasteFormatJSON   PasteFormat = "json"
	PasteFormatBase64 PasteFormat = "base64"
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
				Value:   string(PasteFormatText),
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
		},
		Action: func(ctx *cli.Context) error {
			panic("not implemented")
		},
	}
}
