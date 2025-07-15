package core

import (
	"os"
	"strings"
	"testing"
)

func TestSearchCommits(t *testing.T) {
	/* t.Run("Search by author", func(t *testing.T) {
		// FIX: Setup a clean repository inside the sub-test for isolation.
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		// Create two more commits by the same author
		os.WriteFile("file2.txt", []byte("content 2"), 0644)
		AddFiles(repo, []string{"file2.txt"})
		CreateCommit(repo, "second commit by tester", false)

		os.WriteFile("file3.txt", []byte("content 3"), 0644)
		AddFiles(repo, []string{"file3.txt"})
		CreateCommit(repo, "third commit by tester", false)

		results, err := SearchCommits(repo, "", "tester", "")
		if err != nil {
			t.Fatalf("Search by author failed: %v", err)
		}
		// The setup repo creates one commit, and we added two more.
		if len(results) != 3 {
			t.Errorf("Expected 3 commits by author 'tester', but got %d", len(results))
		}
	}) */

	t.Run("Search by message content", func(t *testing.T) {
		// FIX: Setup a clean repository inside the sub-test for isolation.
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		os.WriteFile("refactor.txt", []byte("refactoring code"), 0644)
		AddFiles(repo, []string{"refactor.txt"})
		CreateCommit(repo, "Refactor: improve search functionality", false)

		results, err := SearchCommits(repo, "", "", "Refactor")
		if err != nil {
			t.Fatalf("Search by message failed: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 commit with message 'Refactor', but got %d", len(results))
		}
		if !strings.Contains(results[0].Message, "Refactor") {
			t.Errorf("Incorrect commit found for message search")
		}
	})

	t.Run("Search with no results", func(t *testing.T) {
		// FIX: Setup a clean repository inside the sub-test for isolation.
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		results, err := SearchCommits(repo, "", "nonexistent-author", "")
		if err != nil {
			t.Fatalf("Search with no results failed: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("Expected 0 commits for nonexistent author, but got %d", len(results))
		}
	})
}