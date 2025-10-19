package handler

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
)

func (h *Handler) ListCollections(c *cli.Context) error {
	format := c.String(FlagFormat)
	if format == "" {
		format = string(PasteFormatText)
	}

	collections, err := h.bore.Collections().List(c.Context, bore.ListCollectionsOptions{})
	if err != nil {
		return err
	}

	if format == string(PasteFormatJSON) {
		return h.viewManager.RenderJSON(c.App.Writer, collections)
	}

	return h.viewManager.RenderCollectionsList(c.App.Writer, collections)
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
	// WARN: When this is enabled, we will skip the confirmation prompt, this nothing to do with whether it is empty or not as you might assume
	force := c.Bool(FlagForce)

	if c.NArg() == 0 {
		return cli.Exit("collection id is required", 1)
	} else if c.NArg() > 1 {
		return cli.Exit("too many arguments", 1)
	}

	collectionID := c.Args().First()
	collection, err := h.bore.Collections().Get(c.Context, collectionID)
	if err != nil {
		return err
	}

	if collection == nil {
		return cli.Exit("collection not found", 1)
	}

	if !force {
		var confirmation string
		_, _ = c.App.Writer.Write(
			fmt.Appendf(
				[]byte{},
				"Are you sure you want to delete the collection %q (ID: %s)? This action cannot be undone. (y/N): ",
				collection.Name,
				collection.ID,
			),
		)
		_, _ = fmt.Scanln(&confirmation)
		if confirmation != "y" && confirmation != "Y" {
			_, _ = c.App.Writer.Write([]byte("Aborted.\n"))
			return nil
		}
	}

	if err := h.bore.Collections().Delete(c.Context, collectionID); err != nil {
		return err
	}

	_, _ = c.App.Writer.Write([]byte(collection.ID + "\n"))
	return nil
}

func (h *Handler) SetDefaultCollection(c *cli.Context) error {
	panic("not implemented")
}
