package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckout(t *testing.T) {
	repo, initialCommitHash, cleanup := setupTestRepo(t)
	defer cleanup()

	CreateBranch(repo, "feature-branch")
	Checkout(repo, "feature-branch")

	featureFilePath := "feature.txt"
	os.WriteFile(featureFilePath, []byte("feature content"), 0644)
	AddFiles(repo, []string{featureFilePath})

	// FIX: Pass the entire repo object to NewStorage.
	storage := NewStorage(repo)
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

	featureRefPath := filepath.Join(repo.RefsDir, "heads", "feature-branch")
	os.WriteFile(featureRefPath, []byte(featureCommitHash+"\n"), 0644)

	t.Run("Checkout a branch", func(t *testing.T) {
		if err := Checkout(repo, "main"); err != nil {
			t.Fatalf("Checkout to main branch failed: %v", err)
		}

		headContent, _ := os.ReadFile(repo.HeadPath)
		if strings.TrimSpace(string(headContent)) != "ref: refs/heads/main" {
			t.Errorf("HEAD should point to 'ref: refs/heads/main', but got '%s'", string(headContent))
		}

		if _, err := os.Stat(featureFilePath); !os.IsNotExist(err) {
			t.Error("feature.txt should not exist on main branch, but it does")
		}

		testFilePath := "test.txt"
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Error("test.txt should exist on main branch, but it does not")
		}
	})

	t.Run("Checkout a commit (detached HEAD)", func(t *testing.T) {
		if err := Checkout(repo, initialCommitHash); err != nil {
			t.Fatalf("Checkout to commit hash failed: %v", err)
		}

		headContent, _ := os.ReadFile(repo.HeadPath)
		if strings.TrimSpace(string(headContent)) != initialCommitHash {
			t.Errorf("HEAD should be detached and point to commit hash %s, but got '%s'", initialCommitHash, string(headContent))
		}
	})
}