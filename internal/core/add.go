package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddFiles handles the core logic of adding files to the index.
func AddFiles(repo *Repository, paths []string) error {
	// Load the index, or create a new one if it doesn't exist.
	index, err := LoadIndex(repo.IndexPath)
	if err != nil {
		if os.IsNotExist(err) {
			index = NewIndex()
		} else {
			return fmt.Errorf("failed to load index: %w", err)
		}
	}

	// FIX: Pass the entire repo object to NewStorage.
	storage := NewStorage(repo)

	for _, path := range paths {
		// Walk the file path. If it's a file, it will be visited once.
		// If it's a directory, it will visit all files within it.
		err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip the .zark directory itself
			if strings.Contains(currentPath, ".zark") {
				return nil
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			absPath, err := filepath.Abs(currentPath)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", currentPath, err)
			}

			// Get the relative path from the repository root
			relPath, err := filepath.Rel(repo.Path, absPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", absPath, err)
			}

			// Read file content
			content, err := os.ReadFile(currentPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", currentPath, err)
			}

			// Create a blob, store it, and add it to the index
			blob := NewBlob(content)
			if err := storage.Store(blob); err != nil {
				return fmt.Errorf("failed to store blob for %s: %w", relPath, err)
			}

			index.Add(relPath, blob.Hash(), "100644", info.Size(), info.ModTime())
			fmt.Printf("added '%s'\n", relPath)

			return nil
		})

		if err != nil {
			return fmt.Errorf("error processing path %s: %w", path, err)
		}
	}

	// Save the updated index to disk
	return index.Save(repo.IndexPath)
}