package core

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestShowHistory(t *testing.T) {
	t.Run("History on a repo with commits", func(t *testing.T) {
		// FIX: Correctly handle 3 return values from the setup function.
		repo, initialCommitHash, cleanup := setupTestRepo(t)
		// FIX: Add the deferred call to the cleanup function.
		defer cleanup()

		// Capture stdout to test the output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		if err := ShowHistory(repo); err != nil {
			t.Fatalf("ShowHistory failed: %v", err)
		}

		w.Close()
		os.Stdout = oldStdout // Restore stdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "commit "+initialCommitHash) {
			t.Errorf("History output should contain the initial commit hash. Got:\n%s", output)
		}
		if !strings.Contains(output, "initial commit") {
			t.Errorf("History output should contain the initial commit message. Got:\n%s", output)
		}
	})

	t.Run("History on a repo with no commits", func(t *testing.T) {
		// Create a new, empty repo
		tmpDir := t.TempDir()
		repo := NewRepository(tmpDir)
		repo.Init()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		if err := ShowHistory(repo); err != nil {
			t.Fatalf("ShowHistory failed: %v", err)
		}

		w.Close()
		os.Stdout = oldStdout // Restore stdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := strings.TrimSpace(buf.String())

		if output != "No commits yet." {
			t.Errorf("Expected 'No commits yet.', got '%s'", output)
		}
	})
}