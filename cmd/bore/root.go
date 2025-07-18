package main

import (
	"github.com/urfave/cli/v2"
)

// TODO: move path to the CLI itself (can be overriden by using the -c flag)
// TODO: automatically detect the data directory for each platform and remove this config option
func execute() error {
	panic("not implemented")
}

func createRootCmd() *cli.App {
	return &cli.App{}
}
