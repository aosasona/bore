package commands

import "github.com/urfave/cli/v2"

type InfoCommand struct {
	*Command
}

func (c *Command) Info() *InfoCommand {
	return &InfoCommand{Command: c}
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) Usage() string {
	return "Display information about the current bore instance"
}

func (c *InfoCommand) Action() error {
	panic("not implemented")
}

func (c *InfoCommand) Build() *cli.Command {
	panic("not implemented")
}
