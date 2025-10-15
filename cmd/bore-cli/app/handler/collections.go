package handler

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
)

func (h *Handler) ListCollections(c *cli.Context) error {
	format := c.String(FlagFormat)
	if format == "" {
		format = string(PasteFormatText)
	}

	_ = format
	panic("not implemented")
}

func (h *Handler) CreateCollection(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("collection name is required", 1)
	} else if c.NArg() > 1 {
		return cli.Exit("too many arguments", 1)
	}

	name := c.Args().First()
	force := c.Bool(FlagForce)

	result, err := h.bore.Collections().Create(
		c.Context,
		bore.CreateCollectionOptions{Name: name, AppendSuffixIfExists: force},
	)
	if err != nil {
		return err
	}

	_, _ = c.App.Writer.Write([]byte(result.Name + " (" + result.ID + ")\n"))
	return nil
}

func (h *Handler) DeleteCollection(c *cli.Context) error {
	panic("not implemented")
}

func (h *Handler) SetDefaultCollection(c *cli.Context) error {
	panic("not implemented")
}
