package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Checkout handles the core logic of checking out a branch or commit.
func Checkout(repo *Repository, ref string) error {
	storage := NewStorage(repo.ObjectsDir)

	// First, check if the ref is a branch name.
	isBranch := false
	branchRefPath := filepath.Join(repo.RefsDir, "heads", ref)
	if _, err := os.Stat(branchRefPath); err == nil {
		isBranch = true
	}

	// Resolve the ref to a commit hash.
	commitHash, err := ResolveRef(repo, ref)
	if err != nil {
		return fmt.Errorf("failed to resolve ref '%s': %w", ref, err)
	}

	// Load the target commit and its tree.
	commitData, err := storage.Load(commitHash)
	if err != nil {
		return fmt.Errorf("failed to load commit object %s: %w", commitHash, err)
	}
	var commit Commit
	if err := json.Unmarshal(commitData, &commit); err != nil {
		return fmt.Errorf("failed to unmarshal commit %s: %w", commitHash, err)
	}

	treeData, err := storage.Load(commit.TreeHash)
	if err != nil {
		return fmt.Errorf("failed to load tree object %s: %w", commit.TreeHash, err)
	}

	// FIX: Unmarshal into a slice of entries first, then create the Tree struct.
	var treeEntries []TreeEntry
	if err := json.Unmarshal(treeData, &treeEntries); err != nil {
		return fmt.Errorf("failed to unmarshal tree %s: %w", commit.TreeHash, err)
	}
	tree := Tree{Entries: treeEntries}

	// Get the current index to know which files are currently tracked.
	index, err := LoadIndex(repo.IndexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// Clean the working directory of all currently tracked files.
	if index != nil {
		for _, entry := range index.Entries {
			os.Remove(filepath.Join(repo.Path, entry.Path))
		}
	}

	// Create a new index for the checkout state.
	newIndex := NewIndex()

	// Write the files from the target tree to the working directory.
	for _, entry := range tree.Entries {
		filePath := filepath.Join(repo.Path, entry.Name)
		blobData, err := storage.Load(entry.Hash)
		if err != nil {
			return fmt.Errorf("failed to load blob %s for file %s: %w", entry.Hash, entry.Name, err)
		}

		// Ensure the directory for the file exists.
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		// Write the file content.
		if err := os.WriteFile(filePath, blobData, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		// We need file info to add to the index, so we stat the file we just wrote.
		info, _ := os.Stat(filePath)
		newIndex.Add(entry.Name, entry.Hash, "100644", info.Size(), info.ModTime())
	}

	// Save the new index.
	if err := newIndex.Save(repo.IndexPath); err != nil {
		return fmt.Errorf("failed to save new index after checkout: %w", err)
	}

	// Update HEAD.
	var headContent string
	if isBranch {
		// If it's a branch, HEAD should point to the branch ref.
		headContent = fmt.Sprintf("ref: refs/heads/%s", ref)
	} else {
		// Otherwise, it's a detached HEAD pointing to a commit.
		headContent = commitHash
	}

	if err := os.WriteFile(repo.HeadPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	return nil
}