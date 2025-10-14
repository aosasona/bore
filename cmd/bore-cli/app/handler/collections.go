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
	name := c.Args().First()
	force := c.Bool(FlagForce)

	result, err := h.bore.Collections().Create(
		c.Context,
		bore.CreateCollectionOptions{Name: name, AppendSuffixIfExists: force},
	)
	if err != nil {
		return cli.Exit("failed to create collection: "+err.Error(), 1)
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
