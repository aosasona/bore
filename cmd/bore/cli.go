package main

import "go.trulyao.dev/bore/v2"

type Cli struct {
	bore *bore.Bore
}

func NewCli() *Cli {
	// TODO: read config
	return &Cli{}
}
