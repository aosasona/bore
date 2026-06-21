package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	InfoName  = "info"
	InfoUsage = "Display information about the current bore instance"
)

type InfoCommand struct{}

func (i InfoCommand) Name() string {
	return InfoName
}

func (i InfoCommand) Usage() string {
	return InfoUsage
}

func (i InfoCommand) Build(run Runner) *cli.Command {
	// nolint:exhaustruct
	return &cli.Command{
		Name:  i.Name(),
		Usage: i.Usage(),
		Action: func(ctx *cli.Context) error {
			return run(ctx, i)
		},
	}
}

func (i InfoCommand) Execute(ctx *cli.Context, manager *Manager) error {
	config, err := manager.config.Read()
	if err != nil {
		return cli.Exit("failed to get bore configuration: "+err.Error(), 1)
	}

	fmt.Println("Bore Version:", manager.config.Version())
	fmt.Println("Data Directory:", config.DataDir)
	fmt.Println("Config Path:", manager.config.ConfigPath())
	fmt.Println("Default Collection:", config.DefaultCollection)
	fmt.Println("Clipboard Passthrough:", config.ClipboardPassthrough)
	return nil
}
