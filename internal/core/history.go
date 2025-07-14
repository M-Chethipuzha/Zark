package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ShowHistory displays the commit log for the repository, walking the parent chain from HEAD.
func ShowHistory(repo *Repository) error {
	storage := NewStorage(repo.ObjectsDir)

	// Resolve HEAD to get the hash of the most recent commit.
	commitHash, err := ResolveRef(repo, "HEAD")
	if err != nil {
		// This is a common case for a new repository with no commits.
		// We check for both possible error messages from ResolveRef.
		if strings.Contains(err.Error(), "no commits yet") || strings.Contains(err.Error(), "broken HEAD reference") {
			fmt.Println("No commits yet.")
			return nil
		}
		return fmt.Errorf("failed to resolve HEAD: %w", err)
	}

	// Walk through the commit history by following parent pointers.
	for commitHash != "" {
		commitData, err := storage.Load(commitHash)
		if err != nil {
			return fmt.Errorf("failed to load commit object %s: %w", commitHash, err)
		}

		var commit Commit
		if err := json.Unmarshal(commitData, &commit); err != nil {
			return fmt.Errorf("failed to unmarshal commit %s: %w", commitHash, err)
		}

		// Print commit details in a git-like format.
		fmt.Printf("\033[33mcommit %s\033[0m\n", commitHash) // Yellow for commit hash
		fmt.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
		fmt.Printf("Date:   %s\n", commit.Timestamp.Format(time.RFC1123Z))
		fmt.Printf("\n\t%s\n\n", commit.Message)

		// Move to the parent commit for the next iteration.
		commitHash = commit.Parent
	}

	return nil
}