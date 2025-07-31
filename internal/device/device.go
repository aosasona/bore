package device

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
)

var (
	identifier string
	mu         sync.Mutex
)

const IdentityFileName = ".identity"

// GetIdentifier reads the device identifier from either the cache or the identity file.
func GetIdentifier(dataDir string) string {
	mu.Lock()
	defer mu.Unlock()

	if identifier != "" {
		return identifier
	}

	f, err := readIdentityFile(dataDir)
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
func CreateIdentifier(dataDir string) string {
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
func readIdentityFile(dataDir string) (string, error) {
	filePath := filepath.Join(dataDir, IdentityFileName)
	f, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

// writeIdentityFile writes the device identifier to the identity file in the specified data directory
func writeIdentityFile(dataDir, identifier string) error {
	filePath := filepath.Join(dataDir, IdentityFileName)
	if err := os.WriteFile(filePath, []byte(identifier), 0o600); err != nil {
		return err
	}

	return nil
}

// generateNewIdentifier generates a new ULID identifier
func generateNewIdentifier() string {
	return ulid.Make().String()
}

// isValidIdentifier checks if the provided identifier is a valid ULID
func isValidIdentifier(id string) bool {
	parsedUlid, err := ulid.Parse(id)
	return err == nil && !parsedUlid.IsZero()
}
