package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateCommit creates a new commit object from the index, signs it if requested,
// and updates the current branch reference.
func CreateCommit(repo *Repository, message string, sign bool) (string, error) {
	index, err := LoadIndex(repo.IndexPath)
	if err != nil || len(index.Entries) == 0 {
		return "", fmt.Errorf("nothing to commit, index is empty")
	}

	storage := NewStorage(repo)
	var treeEntries []TreeEntry
	for _, entry := range index.Entries {
		treeEntries = append(treeEntries, TreeEntry{
			Mode: entry.Mode,
			Name: entry.Path,
			Hash: entry.Hash,
			Type: "blob",
		})
	}
	tree := NewTree(treeEntries)
	if err := storage.Store(tree); err != nil {
		return "", fmt.Errorf("failed to store tree object: %w", err)
	}

	parent, _ := ResolveRef(repo, "HEAD")
	config, err := repo.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	commit := NewCommit(tree.Hash(), parent, config.User.Name, config.User.Email, message)

	if sign {
		// Pass the commit object itself to be signed
		signedCommit, err := SignCommit(commit)
		if err != nil {
			return "", fmt.Errorf("failed to sign commit: %w", err)
		}
		// The signing process modifies the commit (e.g., adds a signature to the message),
		// so we need to re-calculate its hash.
		commit = signedCommit
		commit.rehash()
	}

	if err := storage.Store(commit); err != nil {
		return "", fmt.Errorf("failed to store commit object: %w", err)
	}

	// Update the current branch to point to the new commit
	headData, err := os.ReadFile(repo.HeadPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	headRefStr := strings.TrimSpace(string(headData))
	if !strings.HasPrefix(headRefStr, "ref: ") {
		return "", fmt.Errorf("cannot commit in detached HEAD state")
	}
	branchPath := filepath.Join(repo.ZarkDir, headRefStr[5:])

	if err := os.WriteFile(branchPath, []byte(commit.Hash()+"\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to update branch reference: %w", err)
	}

	return commit.Hash(), nil
}