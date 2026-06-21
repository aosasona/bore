package commands

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

type CopyCommand struct{}

func (c CopyCommand) Build(run Runner) *cli.Command {
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
			return run(ctx, c)
		},
	}
}

func (c CopyCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.Copy(ctx, handler.CliCopyOptions{Stdin: PipedIn()})
}
