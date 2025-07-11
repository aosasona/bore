package main

import (
	"flag"
	"log"
	"strings"
)

var (
	createDownMigration = false
	migrationName       = ""
)

func main() {
	flag.BoolVar(&createDownMigration, "down", false, "Create a down migration file")
	flag.Parse()

	cmd := strings.TrimSpace(flag.Arg(0))

	switch cmd {
	case "create-migration", "cm":
		createMigration()

	default:
		log.Println("Unknown command:", cmd)
		log.Println("Available commands: create-migration (cm)")
		flag.Usage()
	}
}

func createMigration() {
	migrationName = flag.Arg(1)
	if migrationName == "" {
		log.Fatal("Migration name is required")
	}

	println("Creating new migration...")
}
