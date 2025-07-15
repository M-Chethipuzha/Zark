package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ScanForSecrets scans staged files for potential secrets before a commit.
func ScanForSecrets(repo *Repository) error {
	index, err := LoadIndex(repo.IndexPath)
	if err != nil {
		// If the index doesn't exist, there's nothing to scan.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// This regex is a simple example. Real-world scanners use more sophisticated patterns.
	// It looks for common keywords followed by what might be a secret string.
	re := regexp.MustCompile(`(?i)(api_key|secret|token|password)[\s=:]+['"]?([a-zA-Z0-9_.-]{20,})['"]?`)

	for _, entry := range index.Entries {
		filePath := filepath.Join(repo.Path, entry.Path)
		file, err := os.Open(filePath)
		if err != nil {
			// Silently skip files that can't be opened.
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			if re.MatchString(scanner.Text()) {
				// Found a potential secret. Abort the commit with an actionable error message.
				return fmt.Errorf(
					"error: potential secret found in %s on line %d.\n"+
						"To override, commit with the '--no-verify' flag (not implemented). Aborting commit",
					entry.Path, lineNumber,
				)
			}
		}
	}

	return nil
}

// SignCommit is a placeholder for GPG signing logic.
func SignCommit(commit *Commit) (*Commit, error) {
	// In a real implementation, you would use a library like `golang.org/x/crypto/openpgp`
	// to sign the commit data with a user's GPG key. This involves finding the user's key,
	// creating a signature, and appending it to the commit message.
	fmt.Println("Signing commit (placeholder)...")
	// For this placeholder, we'll just add a fake signature block to the message.
	commit.Message += "\n\n-----BEGIN ZARK SIGNATURE-----\n"
	commit.Message += "gpcA/p1AyoA4oFATESgC5eU+A8vB8+A9vA8+A9vA8+A9vA8+A9vA8+A9vA8+A9vA8+A9\n"
	commit.Message += "-----END ZARK SIGNATURE-----\n"
	return commit, nil
}