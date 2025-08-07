package main

import (
	"fmt"
	"log/slog"
	"os"

	"go.trulyao.dev/bore/v2/cmd/bore-cli/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := a.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
