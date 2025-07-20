package main

import (
	"fmt"
	"os"
)

func main() {
	cli, err := NewCli()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create CLI:", err)
		os.Exit(1)
	}

	if err := cli.execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
