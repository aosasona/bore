package commands

import (
	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
)

type CollectionsCommand struct{}

type ListCollectionsCommand struct{}

type CreateCollectionCommand struct{}

type DeleteCollectionCommand struct{}

type RenameCollectionCommand struct{}

type ShowDefaultCollectionCommand struct{}

type SetDefaultCollectionCommand struct{}

type UnsetDefaultCollectionCommand struct{}

func (c CollectionsCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "collections",
		Usage: "Manage clipboard collections",
		Subcommands: []*cli.Command{
			ListCollectionsCommand{}.Build(run),
			CreateCollectionCommand{}.Build(run),
			DeleteCollectionCommand{}.Build(run),
			RenameCollectionCommand{}.Build(run),
			ShowDefaultCollectionCommand{}.Build(run),
		},
	}
}

func (c CollectionsCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return cli.ShowSubcommandHelp(ctx)
}

func (l ListCollectionsCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "list",
		Usage: "List all collections",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    handler.FlagFormat,
				Aliases: []string{"f"},
				Usage:   "Output format (text, json)",
			},
		},
		Action: func(ctx *cli.Context) error {
			return run(ctx, l)
		},
	}
}

func (l ListCollectionsCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.ListCollections(ctx)
}

func (c CreateCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a new collection",
		ArgsUsage: "[collection name]",
		Args:      true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    handler.FlagForce,
				Aliases: []string{"f"},
				Usage:   "Force creation even if a collection with the same name exists (uses a random suffix)",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			return run(ctx, c)
		},
	}
}

func (c CreateCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.CreateCollection(ctx)
}

func (d DeleteCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a collection",
		Args:      true,
		ArgsUsage: "[collection id]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    handler.FlagForce,
				Aliases: []string{"f"},
				Usage:   "Force deletion without confirmation",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			return run(ctx, d)
		},
	}
}

func (d DeleteCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.DeleteCollection(ctx)
}

func (r RenameCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:      "rename",
		Usage:     "Rename a collection",
		ArgsUsage: "[collection id] [new name]",
		Args:      true,
		Action: func(ctx *cli.Context) error {
			return run(ctx, r)
		},
	}
}

func (r RenameCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.RenameCollection(ctx)
}

func (s ShowDefaultCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name: "default",
		Action: func(ctx *cli.Context) error {
			return run(ctx, s)
		},
		Subcommands: []*cli.Command{
			SetDefaultCollectionCommand{}.Build(run),
			UnsetDefaultCollectionCommand{}.Build(run),
		},
		Usage: "Manage the default collection",
	}
}

func (s ShowDefaultCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.ShowDefaultCollection(ctx)
}

func (s SetDefaultCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:      "set",
		Usage:     "Set the default collection",
		ArgsUsage: "[collection id]",
		Action: func(ctx *cli.Context) error {
			return run(ctx, s)
		},
	}
}

func (s SetDefaultCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.SetDefaultCollection(ctx)
}

func (u UnsetDefaultCollectionCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "unset",
		Usage: "Unset the default collection",
		Action: func(ctx *cli.Context) error {
			return run(ctx, u)
		},
	}
}

func (u UnsetDefaultCollectionCommand) Execute(ctx *cli.Context, manager *Manager) error {
	return manager.handler.UnsetDefaultCollection(ctx)
}
