package core

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// setupTestRepo creates a temporary directory, initializes a Zark repository,
// and makes an initial commit so that HEAD exists. It also changes the
// current working directory to the new repo and returns a function to clean up.
func setupTestRepo(t *testing.T) (*Repository, string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)
	if err := repo.Init(); err != nil {
		t.Fatalf("Failed to initialize test repository: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(repo.Path); err != nil {
		t.Fatalf("Failed to change directory to repo path: %v", err)
	}

	cleanup := func() {
		os.Chdir(originalWD)
	}

	dummyFilePath := "test.txt"
	if err := os.WriteFile(dummyFilePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	if err := AddFiles(repo, []string{dummyFilePath}); err != nil {
		t.Fatalf("Failed to add files for initial commit: %v", err)
	}

	// FIX: Pass the entire repo object to NewStorage.
	storage := NewStorage(repo)
	index, _ := LoadIndex(repo.IndexPath)
	var treeEntries []TreeEntry
	for _, entry := range index.Entries {
		treeEntries = append(treeEntries, TreeEntry{Name: entry.Path, Hash: entry.Hash, Type: "blob"})
	}
	tree := NewTree(treeEntries)
	storage.Store(tree)

	commit := NewCommit(tree.Hash(), "", "tester", "tester@example.com", "initial commit")
	storage.Store(commit)

	mainRefPath := filepath.Join(repo.RefsDir, "heads", "main")
	os.WriteFile(mainRefPath, []byte(commit.Hash()+"\n"), 0644)

	return repo, commit.Hash(), cleanup
}

func TestBranching(t *testing.T) {
	t.Run("Create a new branch successfully", func(t *testing.T) {
		repo, headHash, cleanup := setupTestRepo(t)
		defer cleanup()

		branchName := "new-feature"
		err := CreateBranch(repo, branchName)
		if err != nil {
			t.Fatalf("CreateBranch failed: expected no error, got %v", err)
		}

		branchPath := filepath.Join(repo.RefsDir, "heads", branchName)
		content, err := os.ReadFile(branchPath)
		if err != nil {
			t.Fatalf("Could not read new branch file '%s': %v", branchPath, err)
		}

		if strings.TrimSpace(string(content)) != headHash {
			t.Errorf("Branch content is incorrect: got %s, want %s", string(content), headHash)
		}
	})

	t.Run("Fail to create a branch that already exists", func(t *testing.T) {
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		err := CreateBranch(repo, "main")
		if err == nil {
			t.Fatalf("CreateBranch succeeded unexpectedly for an existing branch")
		}
	})

	t.Run("List branches and identify the current branch", func(t *testing.T) {
		repo, _, cleanup := setupTestRepo(t)
		defer cleanup()

		CreateBranch(repo, "feature-a")
		CreateBranch(repo, "feature-b")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := ListBranches(repo)
		if err != nil {
			t.Fatalf("ListBranches failed: %v", err)
		}

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		cleanOutput := re.ReplaceAllString(output, "")

		lines := strings.Split(strings.TrimSpace(cleanOutput), "\n")
		foundCurrent := false
		foundA := false
		foundB := false
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "* main" {
				foundCurrent = true
			} else if trimmedLine == "feature-a" {
				foundA = true
			} else if trimmedLine == "feature-b" {
				foundB = true
			}
		}

		if !foundCurrent {
			t.Errorf("Output should have marked 'main' as the current branch. Got:\n%s", output)
		}
		if !foundA {
			t.Errorf("Output should have contained branch 'feature-a'. Got:\n%s", output)
		}
		if !foundB {
			t.Errorf("Output should have contained branch 'feature-b'. Got:\n%s", output)
		}
	})
}