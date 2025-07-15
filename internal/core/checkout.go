package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Checkout handles the core logic of checking out a branch or commit.
func Checkout(repo *Repository, ref string) error {
	// FIX: Pass the entire repo object to NewStorage.
	storage := NewStorage(repo)

	isBranch := false
	branchRefPath := filepath.Join(repo.RefsDir, "heads", ref)
	if _, err := os.Stat(branchRefPath); err == nil {
		isBranch = true
	}

	commitHash, err := ResolveRef(repo, ref)
	if err != nil {
		return fmt.Errorf("failed to resolve ref '%s': %w", ref, err)
	}

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

	var treeEntries []TreeEntry
	if err := json.Unmarshal(treeData, &treeEntries); err != nil {
		return fmt.Errorf("failed to unmarshal tree %s: %w", commit.TreeHash, err)
	}
	tree := Tree{Entries: treeEntries}

	index, err := LoadIndex(repo.IndexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load index: %w", err)
	}

	if index != nil {
		for _, entry := range index.Entries {
			os.Remove(filepath.Join(repo.Path, entry.Path))
		}
	}

	newIndex := NewIndex()

	for _, entry := range tree.Entries {
		filePath := filepath.Join(repo.Path, entry.Name)
		blobData, err := storage.Load(entry.Hash)
		if err != nil {
			return fmt.Errorf("failed to load blob %s for file %s: %w", entry.Hash, entry.Name, err)
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		if err := os.WriteFile(filePath, blobData, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		info, _ := os.Stat(filePath)
		newIndex.Add(entry.Name, entry.Hash, "100644", info.Size(), info.ModTime())
	}

	if err := newIndex.Save(repo.IndexPath); err != nil {
		return fmt.Errorf("failed to save new index after checkout: %w", err)
	}

	var headContent string
	if isBranch {
		headContent = fmt.Sprintf("ref: refs/heads/%s", ref)
	} else {
		headContent = commitHash
	}

	if err := os.WriteFile(repo.HeadPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	return nil
}