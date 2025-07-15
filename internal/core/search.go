package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// SearchCommits searches for commits based on query, author, and message.
// This version scans all objects to find commits, making it more thorough.
func SearchCommits(repo *Repository, query, author, message string) ([]*Commit, error) {
	var results []*Commit
	storage := NewStorage(repo)

	// Walk the objects directory to find all potential commit objects.
	err := filepath.Walk(repo.ObjectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		// A rough check: if it's a file in a 2-char directory, it's an object.
		if len(filepath.Base(filepath.Dir(path))) == 2 {
			hash := filepath.Base(filepath.Dir(path)) + info.Name()

			// Load the object data to inspect it.
			objData, err := storage.Load(hash)
			if err != nil {
				return nil // Skip objects that can't be loaded
			}

			// Check if it looks like a commit (contains a "tree" field in its JSON).
			if !strings.Contains(string(objData), `"tree"`) {
				return nil
			}

			var commit Commit
			if err := json.Unmarshal(objData, &commit); err != nil {
				return nil // Skip objects that don't unmarshal as commits
			}

			// It's a commit, now check if it matches the filter.
			commit.hash = hash // Set the hash since it's not in the JSON

			match := true
			if author != "" && !strings.Contains(commit.Author, author) {
				match = false
			}
			if message != "" && !strings.Contains(commit.Message, message) {
				match = false
			}
			if query != "" && !strings.Contains(commit.Message, query) && !strings.Contains(commit.Author, query) {
				match = false
			}

			if match {
				results = append(results, &commit)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}