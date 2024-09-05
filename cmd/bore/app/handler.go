package app

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/pkg/handler"
)

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
				Value:   handler.FormatPlainText,
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
				Value:   handler.FormatPlainText,
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
	format := ctx.String("format")
	if !handler.ValidateFormat(format) {
		return fmt.Errorf("unsupported format: %s", format)
	}

	var fn func(*cli.Context) (io.Reader, error) = a.CopyFromPrompt

	// If the program was piped to, read directly from STDIN
	fi, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("failed to get STDIN file info: %w", err)
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fn = a.CopyFromStdin
	}

	// Sanity check
	if fn == nil {
		return errors.New("unsupported input method")
	}

	reader, err := fn(ctx)
	if err != nil {
		return err
	}

	id, err := a.Handler().Copy(reader, handler.CopyOpts{Format: format})
	if err != nil {
		return err
	}

	if a.config.ShowIdOnCopy {
		fmt.Fprintln(ctx.App.Writer, fmt.Sprintf("Copied content with ID: %s", id))
	}

	return nil
}

// CopyFromPrompt reads the content from the terminal prompt
func (a *App) CopyFromPrompt(ctx *cli.Context) (io.Reader, error) {
	content := []byte{}

	reader := bufio.NewReader(ctx.App.Reader)
	fmt.Print("Enter content to copy (press esc + enter to finish):\n")

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == 27 {
			break
		}

		content = append(content, b)
	}

	return bytes.NewReader(content), nil
}

// CopyFromStdin copies the content from the STDIN
func (a *App) CopyFromStdin(ctx *cli.Context) (io.Reader, error) {
	return bufio.NewReader(ctx.App.Reader), nil
}

func (a *App) Paste(ctx *cli.Context) error {
	err := a.Handler().PasteLastCopied(ctx.App.Writer)
	return err
}
