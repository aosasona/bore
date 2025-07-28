package device

import "sync"

var (
	identifier string
	mu         sync.Mutex
)

// GetIdentifier reads the device identifier from either the cache or the identity file.
func GetIdentifier() string {
	mu.Lock()
	defer mu.Unlock()

	panic("not implemented")
}

// CreateIdentifier creates a new device identifier and writes it to the identity file.
// The identity file is stored in the provided data directory.
//
// NOTE: The data is not written to the database to prevent accidental duplication of device identifiers when or if the database is reset, copied, or moved.
func CreateIdentifier(dataDir string) string {
	mu.Lock()
	defer mu.Unlock()

	panic("not implemented")
}

func readIdentityFile() (string, error) {
	panic("not implemented")
}
