package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cli, err := NewCli()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := cli.execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
