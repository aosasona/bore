package device

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
)

const IdentityFileName = ".identity"

type identity struct {
	// identifier is the unique device identifier.
	// It is initialized to an empty string and will be set when the identifier is created or read from the identity file.
	identifier string

	// path is the path to the identity file.
	path string

	sync.Mutex
}

// Identity returns a shared instance of the identity manager.
func Identity(dataDir string) *identity {
	return &identity{
		path:       filepath.Join(dataDir, IdentityFileName),
		identifier: "",
	}
}

// GetIdentifier reads the device identifier from either the cache or the identity file.
func (i *identity) GetIdentifier() (string, error) {
	i.Lock()
	defer i.Unlock()

	if i.identifier != "" {
		return i.identifier, nil
	}

	f, err := i.readIdentityFile()
	if err != nil {
		if os.IsNotExist(err) {
			return i.CreateIdentifier()
		}

		return "", fmt.Errorf("[GetIdentifier] failed to read identity file: %w", err)
	}

	i.identifier = f
	return i.identifier, nil
}

// CreateIdentifier creates a new device identifier and writes it to the identity file.
// The identity file is stored in the provided data directory.
//
// NOTE: The data is not written to the database to prevent accidental duplication of device identifiers when or if the database is reset, copied, or moved.
func (i *identity) CreateIdentifier() (string, error) {
	i.Lock()
	defer i.Unlock()

	// We should NOT overwrite the identifier if it already exists and is valid.
	content, err := i.readIdentityFile()
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("[CreateIdentifier] failed to read identity file: %w", err)
	}

	content = strings.TrimSpace(content)
	if i.isValidIdentifier(content) {
		i.identifier = content
		return content, nil
	}

	id := i.generateNewIdentifier()
	if err := i.writeIdentityFile(id); err != nil {
		return "", fmt.Errorf("[CreateIdentifier] failed to write identity file: %w", err)
	}

	i.identifier = id
	return id, nil
}

func (i *identity) RsetIdentifier() error {
	i.Lock()
	defer i.Unlock()

	// Reset the identifier by removing the identity file
	if err := os.Remove(i.path); err != nil && !os.IsNotExist(err) {
		return err
	}

	i.identifier = ""
	return nil
}

// readIdentityFile reads the device identifier from the identity file in the specified data directory
func (i *identity) readIdentityFile() (string, error) {
	f, err := os.ReadFile(i.path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

// writeIdentityFile writes the device identifier to the identity file in the specified data directory
func (i *identity) writeIdentityFile(identifier string) error {
	if err := os.WriteFile(i.path, []byte(identifier), 0o600); err != nil {
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
