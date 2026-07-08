package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type ResetCommand struct{}

func (r ResetCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  "reset",
		Usage: "Reset the bore instance, clearing all data",
		Action: func(ctx *cli.Context) error {
			return run(ctx, r)
		},
	}
}

func (r ResetCommand) Execute(ctx *cli.Context, manager *Manager) error {
	if err := manager.bore.Reset(); err != nil {
		return cli.Exit("failed to reset bore: "+err.Error(), 1)
	}

	if err := os.Remove(manager.config.ConfigPath()); err != nil {
		return cli.Exit("failed to remove configuration file: "+err.Error(), 1)
	}

	fmt.Println("Bore instance has been reset successfully.")
	return nil
}

var _ SubCommand = (*ResetCommand)(nil)
