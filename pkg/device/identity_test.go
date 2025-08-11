package device_test

import (
	"os"
	"testing"

	"go.trulyao.dev/bore/v2/pkg/device"
)

const testDir = "/tmp/bore_device"

func init() {
	if err := os.MkdirAll(testDir, 0755); err != nil {
		panic("Failed to create test directory: " + err.Error())
	}
}

func Test_GetIdentifier(t *testing.T) {
	i := device.NewIdentity(testDir)

	identifier, err := i.GetIdentifier()
	if err != nil {
		t.Fatalf("GetIdentifier failed: %v", err)
	}

	if identifier == "" {
		t.Fatal("GetIdentifier returned an empty identifier")
	}

	if !i.IsValidIdentifier(identifier) {
		t.Fatalf("GetIdentifier returned an invalid identifier: %s", identifier)
	}

	t.Logf("Device identifier: %s", identifier)
}
