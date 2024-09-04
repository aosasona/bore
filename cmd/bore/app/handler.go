package app

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

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
	fi, _ := os.Stdin.Stat()
	isPipe := (fi.Mode() & os.ModeCharDevice) == 0

	if isPipe {
		return a.CopyFromStdin(ctx)
	}

	return a.CopyFromPrompt(ctx)
}

// CopyFromPrompt reads the content from the terminal prompt
func (a *App) CopyFromPrompt(ctx *cli.Context) error {
	content := []byte{}

	reader := bufio.NewReader(ctx.App.Reader)
	fmt.Print("Enter content to copy (press esc + enter to finish):\n")

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if b == 27 {
			break
		}

		content = append(content, b)
	}

	contentReader := bytes.NewReader(content)
	id, err := a.Handler().Copy(contentReader, handler.CopyOpts{})
	if err != nil {
		return err
	}

	if a.config.ShowIdOnCopy {
		fmt.Fprintln(ctx.App.Writer, fmt.Sprintf("Copied content with ID: %s", id))
	}

	return nil
}

// CopyFromStdin copies the content from the STDIN
func (a *App) CopyFromStdin(ctx *cli.Context) error {
	content := bufio.NewReader(ctx.App.Reader)
	id, err := a.Handler().Copy(content, handler.CopyOpts{})
	if err != nil {
		return err
	}

	if a.config.ShowIdOnCopy {
		fmt.Fprintln(ctx.App.Writer, fmt.Sprintf("Copied content with ID: %s", id))
	}

	return nil
}

func (a *App) Paste(ctx *cli.Context) error {
	err := a.Handler().PasteLastCopied(ctx.App.Writer)
	return err
}
