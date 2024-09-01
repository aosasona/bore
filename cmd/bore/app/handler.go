package app

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/pkg/handler"
)

func (a *App) CopyCommand() *cli.Command {
	return &cli.Command{
		Name:  "copy",
		Usage: "Copy the content of the provided file or STDIN to the clipboard",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Path to the file to copy",
			},
		},
		Action: a.Copy,
	}
}

func (a *App) PasteCommand() *cli.Command {
	return &cli.Command{
		Name:   "paste",
		Usage:  "Paste the last copied content",
		Action: a.Paste,
		Flags:  []cli.Flag{}, // Add flags for: last ranges of copied content, file to write to, collection to paste from etc.
	}
}

func (a *App) Copy(ctx *cli.Context) error {
	if strings.TrimSpace(ctx.String("file")) != "" {
		return a.CopyFromFile(ctx)
	}

	return a.CopyFromStdin(ctx)
}

// CopyFromStdin copies the content from the STDIN to the database
// TODO: implement collections and artifact types
func (a *App) CopyFromStdin(ctx *cli.Context) error {
	content := bufio.NewReader(ctx.App.Reader)
	id, err := a.Handler().Copy(content, handler.CopyOpts{})
	if err != nil {
		return err
	}

	if a.config.ShowIdOnCopy {
		fmt.Fprintln(ctx.App.Writer, fmt.Sprintf("Copied content with ID: %s\n", id))
	}

	return nil
}

func (a *App) CopyFromFile(ctx *cli.Context) error {
	return fmt.Errorf("not implemented")
}

func (a *App) Paste(ctx *cli.Context) error {
	err := a.Handler().PasteLastCopied(ctx.App.Writer)
	if err != nil {
		return err
	}

	return nil
}
