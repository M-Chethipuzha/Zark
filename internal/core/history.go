package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ShowHistory displays the commit log for the repository, walking the parent chain from HEAD.
func ShowHistory(repo *Repository) error {
	storage := NewStorage(repo)

	commitHash, err := ResolveRef(repo, "HEAD")
	if err != nil {
		if strings.Contains(err.Error(), "no commits yet") || strings.Contains(err.Error(), "broken HEAD reference") {
			fmt.Println("No commits yet.")
			return nil
		}
		return fmt.Errorf("failed to resolve HEAD: %w", err)
	}

	for commitHash != "" {
		commitData, err := storage.Load(commitHash)
		if err != nil {
			return fmt.Errorf("failed to load commit object %s: %w", commitHash, err)
		}

		var commit Commit
		if err := json.Unmarshal(commitData, &commit); err != nil {
			return fmt.Errorf("failed to unmarshal commit %s: %w", commitHash, err)
		}

		// FIX: Print the commitHash from the loop, as the unmarshaled commit struct
		// does not have its unexported 'hash' field populated.
		fmt.Printf("\033[33mcommit %s\033[0m\n", commitHash) // Yellow for commit hash
		fmt.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
		fmt.Printf("Date:   %s\n", commit.Timestamp.Format(time.RFC1123Z))
		fmt.Printf("\n\t%s\n\n", commit.Message)

		commitHash = commit.Parent
	}

	return nil
}