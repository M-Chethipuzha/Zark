package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateBranch creates a new branch pointing to the current HEAD commit.
func CreateBranch(repo *Repository, branchName string) error {
	branchPath := filepath.Join(repo.RefsDir, "heads", branchName)
	if _, err := os.Stat(branchPath); err == nil {
		return fmt.Errorf("a branch named '%s' already exists", branchName)
	}

	// Get the current HEAD commit hash.
	headHash, err := ResolveRef(repo, "HEAD")
	if err != nil {
		return fmt.Errorf("cannot create branch: %w. Make a commit first", err)
	}

	// Write the commit hash to the new branch file.
	if err := os.WriteFile(branchPath, []byte(headHash+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to create branch file: %w", err)
	}

	fmt.Printf("Branch '%s' created at %s\n", branchName, headHash[:8])
	return nil
}

// ListBranches lists all local branches.
func ListBranches(repo *Repository) error {
	headsDir := filepath.Join(repo.RefsDir, "heads")

	// Get current branch from HEAD
	headData, err := os.ReadFile(repo.HeadPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}
	currentRef := strings.TrimSpace(string(headData))

	err = filepath.Walk(headsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			branchName := filepath.Base(path)

			// Check if this is the current branch
			if strings.HasSuffix(currentRef, branchName) {
				fmt.Printf("* \033[32m%s\033[0m\n", branchName) // Green for current branch
			} else {
				fmt.Printf("  %s\n", branchName)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}
	return nil
}