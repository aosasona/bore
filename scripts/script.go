package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DirMigrations = "migrations"

var (
	createDownMigration = false
	migrationName       = ""
)

// Usage: go run script.go [-down] create-migration [migration_name]
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
	migrationName = regexp.MustCompile("[^a-zA-Z0-9_]").
		ReplaceAllString(strings.ToLower(migrationName), "_")

	println("Creating new migration...")

	lastMigrationIndex, err := getLastMigrationIndex()
	if err != nil {
		log.Fatalf("Error getting last migration index: %v", err)
	}

	newMigrationIndex := lastMigrationIndex + 1
	upFilename := fmt.Sprintf(
		"%06d_%s.up.sql",
		newMigrationIndex,
		migrationName,
	)
	downFilename := fmt.Sprintf(
		"%06d_%s.down.sql",
		newMigrationIndex,
		migrationName,
	)

	upFilePath := fmt.Sprintf("%s/%s", DirMigrations, upFilename)
	downFilePath := fmt.Sprintf("%s/%s", DirMigrations, downFilename)

	upFile, err := os.Create(upFilePath)
	if err != nil {
		log.Fatalf("Error creating up migration file: %v", err)
	}
	defer upFile.Close()

	_, err = upFile.WriteString("-- Up migration SQL goes here\n")
	if err != nil {
		log.Fatalf("Error writing to up migration file: %v", err)
	}

	if createDownMigration {
		downFile, err := os.Create(downFilePath)
		if err != nil {
			log.Fatalf("Error creating down migration file: %v", err)
		}
		defer downFile.Close()

		_, err = downFile.WriteString("-- Down migration SQL goes here\n")
		if err != nil {
			log.Fatalf("Error writing to down migration file: %v", err)
		}
	}

	logMessage := fmt.Sprintf(
		"Migration created successfully:\nUp: %s",
		upFilePath,
	)

	if createDownMigration {
		logMessage += fmt.Sprintf("\nDown: %s", downFilePath)
	}
	log.Println(logMessage)
}

func getLastMigrationIndex() (int, error) {
	files, err := os.ReadDir(DirMigrations)
	if err != nil {
		return -1, err
	}

	lastMigrationIndex := -1
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		underscoreIndex := strings.Index(file.Name(), "_")
		if underscoreIndex == -1 {
			log.Printf("Skipping file %s: no underscore found", file.Name())
			continue
		}

		migrationIndex, err := strconv.Atoi(file.Name()[:underscoreIndex])
		if err != nil {
			log.Printf("Skipping file %s: invalid migration index", file.Name())
			continue
		}

		if migrationIndex > lastMigrationIndex {
			lastMigrationIndex = migrationIndex
		}
	}

	return lastMigrationIndex, nil
}
