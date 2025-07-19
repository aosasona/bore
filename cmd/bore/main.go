package main

import (
	"fmt"
	"os"
)

func main() {
	cli := NewCli()

	if err := cli.execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
