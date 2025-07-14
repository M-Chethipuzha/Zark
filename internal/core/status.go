package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetStatus handles the core logic of determining and printing the repository status.
func GetStatus(repo *Repository) error {
	// 1. Load HEAD Tree entries
	storage := NewStorage(repo.ObjectsDir)
	headTree, err := GetHeadTreeEntries(repo, storage)
	if err != nil {
		// A non-fatal error, likely means no commits yet
		if !strings.Contains(err.Error(), "no commits yet") {
			return fmt.Errorf("could not get HEAD tree: %w", err)
		}
	}

	// 2. Load Index entries
	index, err := LoadIndex(repo.IndexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not load index: %w", err)
	}
	indexEntries := make(map[string]string)
	if index != nil {
		for _, entry := range index.Entries {
			indexEntries[entry.Path] = entry.Hash
		}
	}

	// --- Compare states and find differences ---
	stagedChanges := make(map[string]string)
	unstagedChanges := make(map[string]string)
	untrackedFiles := []string{}

	// Compare HEAD vs Index for staged changes
	for path, hash := range indexEntries {
		if headHash, ok := headTree[path]; !ok {
			stagedChanges[path] = "new file"
		} else if headHash != hash {
			stagedChanges[path] = "modified"
		}
	}
	for path := range headTree {
		if _, ok := indexEntries[path]; !ok {
			stagedChanges[path] = "deleted"
		}
	}

	// 3. Walk Working Directory and compare with Index
	err = filepath.Walk(repo.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, ".zark") || info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(repo.Path, path)

		if indexHash, ok := indexEntries[relPath]; ok {
			// File is in the index, check for modification
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			currentHashBytes := sha256.Sum256(content)
			currentHash := hex.EncodeToString(currentHashBytes[:])
			if currentHash != indexHash {
				unstagedChanges[relPath] = "modified"
			}
		} else {
			// File is not in the index, so it's untracked
			untrackedFiles = append(untrackedFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking working directory: %w", err)
	}

	// Check for files deleted from working tree but still in index
	for path := range indexEntries {
		fullPath := filepath.Join(repo.Path, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			unstagedChanges[path] = "deleted"
		}
	}

	// --- Print Status Report ---
	printStatus("Changes to be committed:", stagedChanges)
	printStatus("Changes not staged for commit:", unstagedChanges)

	if len(untrackedFiles) > 0 {
		fmt.Println("\nUntracked files:")
		fmt.Println("  (use \"zark add <file>...\" to include in what will be committed)")
		for _, path := range untrackedFiles {
			fmt.Printf("\t\033[31m%s\033[0m\n", path) // Red color for untracked files
		}
	}

	return nil
}

func printStatus(title string, changes map[string]string) {
	if len(changes) > 0 {
		fmt.Println(title)
		fmt.Println("  (use \"zark restore <file>...\" to unstage)")
		for path, status := range changes {
			color := "\033[32m" // Green for staged
			if title == "Changes not staged for commit:" {
				color = "\033[31m" // Red for unstaged
			}
			fmt.Printf("\t%s%s:   %s\033[0m\n", color, status, path)
		}
		fmt.Println()
	}
}
