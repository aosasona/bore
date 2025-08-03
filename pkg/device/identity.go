package device

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
)

const IdentityFileName = ".identity"

var instance *identity

type identity struct {
	// identifier is the unique device identifier.
	// It is initialized to an empty string and will be set when the identifier is created or read from the identity file.
	identifier string

	// dataDir is the directory where the application's data is stored; same as the one provided on initialization.
	dataDir string

	sync.Mutex
}

// Identity returns a shared instance of the identity manager.
func Identity(dataDir string) *identity {
	if instance == nil {
		instance = &identity{dataDir: dataDir}
	}

	return instance
}

// GetIdentifier reads the device identifier from either the cache or the identity file.
func (i *identity) GetIdentifier() string {
	i.Lock()
	defer i.Unlock()

	if i.identifier != "" {
		return i.identifier
	}

	f, err := i.readIdentityFile()
	if err != nil {
		if os.IsNotExist(err) {
			return CreateIdentifier(dataDir)
		}

		panic("failed to read identity file: " + err.Error())
	}

	identifier = f
	return identifier
}

// CreateIdentifier creates a new device identifier and writes it to the identity file.
// The identity file is stored in the provided data directory.
//
// NOTE: The data is not written to the database to prevent accidental duplication of device identifiers when or if the database is reset, copied, or moved.
func (i *identity) CreateIdentifier() string {
	mu.Lock()
	defer mu.Unlock()

	// We should NOT overwrite the identifier if it already exists and is valid.
	content, err := readIdentityFile(dataDir)
	if err != nil && !os.IsNotExist(err) {
		panic("failed to read identity file: " + err.Error())
	}

	content = strings.TrimSpace(content)
	if isValidIdentifier(content) {
		identifier = content
		return content
	}

	id := generateNewIdentifier()
	if err := writeIdentityFile(dataDir, id); err != nil {
		panic("failed to write identity file: " + err.Error())
	}

	identifier = id
	return id
}

// readIdentityFile reads the device identifier from the identity file in the specified data directory
func (i *identity) readIdentityFile(dataDir string) (string, error) {
	filePath := filepath.Join(dataDir, IdentityFileName)
	f, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

// writeIdentityFile writes the device identifier to the identity file in the specified data directory
func (i *identity) writeIdentityFile(dataDir, identifier string) error {
	filePath := filepath.Join(dataDir, IdentityFileName)
	if err := os.WriteFile(filePath, []byte(identifier), 0o600); err != nil {
		return err
	}

	return nil
}

// generateNewIdentifier generates a new ULID identifier
func (i *identity) generateNewIdentifier() string {
	return ulid.Make().String()
}

// isValidIdentifier checks if the provided identifier is a valid ULID
func (i *identity) isValidIdentifier(id string) bool {
	parsedUlid, err := ulid.Parse(id)
	return err == nil && !parsedUlid.IsZero()
}
