package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// TrackLFS adds a pattern to the LFS tracking file, which is named .zarkattributes.
func TrackLFS(repo *Repository, pattern string) error {
	// The .zarkattributes file should be in the root of the working directory, not in .zark
	lfsFilePath := filepath.Join(repo.Path, ".zarkattributes")
	f, err := os.OpenFile(lfsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .zarkattributes: %w", err)
	}
	defer f.Close()

	// This line mimics Git's LFS configuration for a given file pattern.
	// It tells Zark to use the 'lfs' filter for this pattern.
	lfsLine := fmt.Sprintf("%s filter=lfs diff=lfs merge=lfs -text", pattern)
	if _, err := f.WriteString(lfsLine + "\n"); err != nil {
		return fmt.Errorf("failed to write to .zarkattributes: %w", err)
	}

	fmt.Printf("Tracking '%s' for large file storage.\n", pattern)
	return nil
}