package core

import (
	"os"
	"testing"
)

func TestSecretScanning(t *testing.T) {
	t.Run("Commit fails when a secret is detected", func(t *testing.T) {
		// FIX: Setup a clean repository inside the sub-test for isolation.
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		secretFile := "config.yaml"
		// This key is long enough to trigger the regex
		secretContent := "api_key: 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0'"
		os.WriteFile(secretFile, []byte(secretContent), 0644)

		AddFiles(repo, []string{secretFile})

		// We expect this to fail
		err := ScanForSecrets(repo)
		if err == nil {
			t.Fatal("Expected secret scan to fail, but it succeeded")
		}
	})

	t.Run("Commit succeeds when no secret is present", func(t *testing.T) {
		// FIX: Setup a clean repository inside the sub-test for isolation.
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		safeFile := "safe_config.yaml"
		safeContent := "setting: 'some_value'"
		os.WriteFile(safeFile, []byte(safeContent), 0644)

		AddFiles(repo, []string{safeFile})

		// We expect this to succeed
		err := ScanForSecrets(repo)
		if err != nil {
			t.Fatalf("Expected secret scan to succeed, but it failed: %v", err)
		}
	})
}