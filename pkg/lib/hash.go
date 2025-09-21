package lib

import (
	"crypto/sha256"
	"fmt"
)

// Generate a SHA256 hash of a string and return the hex representation
func Checksum(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
