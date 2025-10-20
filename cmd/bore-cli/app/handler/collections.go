package handler

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
)

// ListCollections lists all collections in the Bore database.
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
		return h.tuiManager.RenderJSON(c.App.Writer, collections)
	}

	return h.tuiManager.RenderCollectionsList(c.App.Writer, collections)
}

// CreateCollection creates a new collection with the given name.
// If --force is used and a collection with the same name exists, a suffix will be appended to the name.
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

// DeleteCollection deletes a collection by ID, confirming with the user first unless --force is used.
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

// RenameCollection renames a collection.
func (h *Handler) RenameCollection(c *cli.Context) error {
	if c.NArg() < 2 {
		return cli.Exit("collection id and new name are required", 1)
	} else if c.NArg() > 2 {
		return cli.Exit("too many arguments", 1)
	}

	collectionID := c.Args().Get(0)
	newName := c.Args().Get(1)

	if collectionID == "" {
		return cli.Exit("collection id is required", 1)
	} else if newName == "" {
		return cli.Exit("new name is required", 1)
	}

	id, err := ulid.Parse(collectionID)
	if err != nil {
		return cli.Exit("invalid collection id", 1)
	}

	collection, err := h.bore.Collections().Get(c.Context, id.String())
	if err != nil {
		return err
	}

	if collection == nil {
		return cli.Exit("collection not found", 1)
	}

	if err = h.bore.Collections().Rename(c.Context, id.String(), newName); err != nil {
		return err
	}

	_, _ = c.App.Writer.Write([]byte(newName + "\n"))
	return nil
}

// SetDefaultCollection sets the default collection.
func (h *Handler) SetDefaultCollection(c *cli.Context) error {
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

	if err := h.configManager.SetDefaultCollectionID(collection.ID); err != nil {
		return err
	}

	_, _ = c.App.Writer.Write([]byte(collection.ID + "\n"))
	return nil
}

// UnsetDefaultCollection unsets the default collection.
func (h *Handler) UnsetDefaultCollection(c *cli.Context) error {
	if err := h.configManager.UnsetDefaultCollectionID(); err != nil {
		return err
	}

	_, _ = c.App.Writer.Write([]byte("Default collection unset.\n"))
	return nil
}

// ShowDefaultCollection shows the default collection.
func (h *Handler) ShowDefaultCollection(c *cli.Context) error {
	config, err := h.configManager.Read()
	if err != nil {
		return err
	}

	if config.DefaultCollection == "" {
		_, _ = c.App.Writer.Write([]byte("No default collection set.\n"))
		return nil
	}

	collection, err := h.bore.Collections().Get(c.Context, config.DefaultCollection)
	if err != nil {
		return err
	}

	if collection == nil {
		if err := h.configManager.UnsetDefaultCollectionID(); err != nil {
			return err
		}

		_, _ = c.App.Writer.Write([]byte("Default collection not found.\n"))
		return nil
	}

	_, _ = c.App.Writer.Write([]byte(collection.Name + " (" + collection.ID + ")\n"))
	return nil
}
