package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveRef resolves a reference (like "main" or a commit hash) to a full commit hash.
func ResolveRef(repo *Repository, ref string) (string, error) {
	// Check for SHA-256 hash length (64 characters)
	if len(ref) == 64 {
		return ref, nil
	}

	// Check if it's a branch
	refPath := filepath.Join(repo.RefsDir, "heads", ref)
	if data, err := os.ReadFile(refPath); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	// Check if it's HEAD
	if ref == "HEAD" {
		headData, err := os.ReadFile(repo.HeadPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("no commits yet, HEAD does not exist")
			}
			return "", fmt.Errorf("failed to read HEAD: %w", err)
		}
		headRef := strings.TrimSpace(string(headData))
		if strings.HasPrefix(headRef, "ref: ") {
			branchPath := filepath.Join(repo.ZarkDir, headRef[5:])
			if data, err := os.ReadFile(branchPath); err == nil {
				return strings.TrimSpace(string(data)), nil
			}
			return "", fmt.Errorf("broken HEAD reference: %s", headRef[5:])
		}
		return headRef, nil
	}

	return "", fmt.Errorf("could not resolve reference: %s", ref)
}

// GetHeadTreeEntries returns a map of file paths to blob hashes for the current HEAD commit.
func GetHeadTreeEntries(repo *Repository, storage *Storage) (map[string]string, error) {
	headHash, err := ResolveRef(repo, "HEAD")
	if err != nil {
		if strings.Contains(err.Error(), "no commits yet") {
			return make(map[string]string), nil
		}
		return nil, err
	}

	commitData, err := storage.Load(headHash)
	if err != nil {
		return nil, fmt.Errorf("failed to load HEAD commit object: %w", err)
	}

	var commit Commit
	if err := json.Unmarshal(commitData, &commit); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HEAD commit: %w", err)
	}

	treeData, err := storage.Load(commit.TreeHash)
	if err != nil {
		return nil, fmt.Errorf("failed to load tree object %s: %w", commit.TreeHash, err)
	}

	var treeEntries []TreeEntry
	if err := json.Unmarshal(treeData, &treeEntries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree: %w", err)
	}

	entries := make(map[string]string)
	for _, entry := range treeEntries {
		entries[entry.Name] = entry.Hash
	}

	return entries, nil
}