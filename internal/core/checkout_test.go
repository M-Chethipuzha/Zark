package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckout(t *testing.T) {
	// FIX: Correctly handle 3 return values from the setup function.
	repo, initialCommitHash, cleanup := setupTestRepo(t)
	// FIX: Add the deferred call to the cleanup function.
	defer cleanup()

	// Create a new branch and a new commit on it
	CreateBranch(repo, "feature-branch")
	Checkout(repo, "feature-branch")

	// Create a new file on the feature branch
	featureFilePath := "feature.txt" // Use relative path since we changed directory
	os.WriteFile(featureFilePath, []byte("feature content"), 0644)
	AddFiles(repo, []string{featureFilePath})

	// Manually create the second commit
	storage := NewStorage(repo.ObjectsDir)
	index, _ := LoadIndex(repo.IndexPath)
	var treeEntries []TreeEntry
	for _, entry := range index.Entries {
		treeEntries = append(treeEntries, TreeEntry{Name: entry.Path, Hash: entry.Hash, Type: "blob"})
	}
	tree := NewTree(treeEntries)
	storage.Store(tree)
	featureCommit := NewCommit(tree.Hash(), initialCommitHash, "tester", "tester@example.com", "add feature")
	storage.Store(featureCommit)
	featureCommitHash := featureCommit.Hash()

	// Update the feature branch to point to the new commit
	featureRefPath := filepath.Join(repo.RefsDir, "heads", "feature-branch")
	os.WriteFile(featureRefPath, []byte(featureCommitHash+"\n"), 0644)

	t.Run("Checkout a branch", func(t *testing.T) {
		// Switch back to the main branch
		if err := Checkout(repo, "main"); err != nil {
			t.Fatalf("Checkout to main branch failed: %v", err)
		}

		// Verify HEAD points to main
		headContent, _ := os.ReadFile(repo.HeadPath)
		if strings.TrimSpace(string(headContent)) != "ref: refs/heads/main" {
			t.Errorf("HEAD should point to 'ref: refs/heads/main', but got '%s'", string(headContent))
		}

		// Verify the feature file does not exist
		if _, err := os.Stat(featureFilePath); !os.IsNotExist(err) {
			t.Error("feature.txt should not exist on main branch, but it does")
		}

		// Verify the original test file does exist
		testFilePath := "test.txt"
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Error("test.txt should exist on main branch, but it does not")
		}
	})

	t.Run("Checkout a commit (detached HEAD)", func(t *testing.T) {
		// Checkout the initial commit directly by its hash
		if err := Checkout(repo, initialCommitHash); err != nil {
			t.Fatalf("Checkout to commit hash failed: %v", err)
		}

		// Verify HEAD contains the commit hash (detached HEAD)
		headContent, _ := os.ReadFile(repo.HeadPath)
		if strings.TrimSpace(string(headContent)) != initialCommitHash {
			t.Errorf("HEAD should be detached and point to commit hash %s, but got '%s'", initialCommitHash, string(headContent))
		}
	})
}