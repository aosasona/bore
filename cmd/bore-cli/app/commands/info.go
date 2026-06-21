package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type InfoCommand struct {
	*Manager
}

func (c *Manager) Info() *InfoCommand {
	return &InfoCommand{Manager: c}
}

func (i InfoCommand) Name() string {
	return "info"
}

func (i InfoCommand) Usage() string {
	return "Display information about the current bore instance"
}

func (i *InfoCommand) Execute(ctx *cli.Context) error {
	config, err := i.config.Read()
	if err != nil {
		return cli.Exit("failed to get bore configuration: "+err.Error(), 1)
	}

	fmt.Println("Bore Version:", i.config.Version())
	fmt.Println("Data Directory:", config.DataDir)
	fmt.Println("Config Path:", i.config.ConfigPath())
	fmt.Println("Default Collection:", config.DefaultCollection)
	fmt.Println("Clipboard Passthrough:", config.ClipboardPassthrough)
	return nil
}
