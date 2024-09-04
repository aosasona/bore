package app

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/pkg/handler"
)

const (
	FormatBase64    = "base64"
	FormatPlainText = "plain-text"
)

// TODO: implement formats support
// TODO: implement native clipboard support
func (a *App) CopyCommand() *cli.Command {
	return &cli.Command{
		Name:  "copy",
		Usage: "Copy content from STDIN or prompt",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "The format of the content to copy. Available formats: base64, plain-text",
				Value:   FormatPlainText,
			},
		},
		Action: a.Copy,
	}
}

func (a *App) PasteCommand() *cli.Command {
	return &cli.Command{
		Name:   "paste",
		Usage:  "Paste from the bore database",
		Action: a.Paste,

		// TODO: Add flags for: last ranges of copied content, file to write to, collection to paste from etc.
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "The format of the content to copy. Available formats: base64, plain-text",
				Value:   FormatPlainText,
			},
			&cli.BoolFlag{
				Name:    "from-system",
				Aliases: []string{"s"},
				Usage:   "Paste from the system clipboard instead of the bore clipboard (also imports the content into the bore clipboard)",
			},
		},
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
