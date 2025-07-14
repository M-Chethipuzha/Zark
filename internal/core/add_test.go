package core

import (
	"os"
	"testing"
)

func TestAddFiles(t *testing.T) {
	t.Run("Add a new file", func(t *testing.T) {
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		newFileName := "newfile.txt"
		if err := os.WriteFile(newFileName, []byte("new content"), 0644); err != nil {
			t.Fatalf("Failed to create new file: %v", err)
		}

		if err := AddFiles(repo, []string{newFileName}); err != nil {
			t.Fatalf("AddFiles failed: %v", err)
		}

		index, err := LoadIndex(repo.IndexPath)
		if err != nil {
			t.Fatalf("Failed to load index after add: %v", err)
		}

		found := false
		for _, entry := range index.Entries {
			if entry.Path == newFileName {
				found = true
				break
			}
		}

		if !found {
			t.Error("newfile.txt was not found in the index after being added")
		}
	})

	t.Run("Add a modified file", func(t *testing.T) {
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		testFileName := "test.txt"
		index, _ := LoadIndex(repo.IndexPath)
		var originalHash string
		for _, entry := range index.Entries {
			if entry.Path == testFileName {
				originalHash = entry.Hash
				break
			}
		}

		if err := os.WriteFile(testFileName, []byte("modified content"), 0644); err != nil {
			t.Fatalf("Failed to modify test.txt: %v", err)
		}

		if err := AddFiles(repo, []string{testFileName}); err != nil {
			t.Fatalf("AddFiles failed for modified file: %v", err)
		}

		newIndex, _ := LoadIndex(repo.IndexPath)
		var newHash string
		for _, entry := range newIndex.Entries {
			if entry.Path == testFileName {
				newHash = entry.Hash
				break
			}
		}

		if newHash == "" || newHash == originalHash {
			t.Error("Expected hash to change for modified file, but it did not")
		}
	})
}