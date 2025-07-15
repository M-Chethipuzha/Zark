package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLFS(t *testing.T) {
	repo, _, cleanup := setupTestRepo(t)
	defer cleanup()

	t.Run("Track a new LFS pattern", func(t *testing.T) {
		pattern := "*.zip"
		err := TrackLFS(repo, pattern)
		if err != nil {
			t.Fatalf("TrackLFS failed: %v", err)
		}

		// Verify that the .zarkattributes file was created and contains the pattern
		attributesPath := filepath.Join(repo.Path, ".zarkattributes")
		content, err := os.ReadFile(attributesPath)
		if err != nil {
			t.Fatalf("Could not read .zarkattributes file: %v", err)
		}

		expectedLine := "*.zip filter=lfs diff=lfs merge=lfs -text"
		if !strings.Contains(string(content), expectedLine) {
			t.Errorf("Expected .zarkattributes to contain '%s', but it did not. Got:\n%s", expectedLine, string(content))
		}
	})
}