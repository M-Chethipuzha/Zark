package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

func SaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "Save changes to the repository",
		Long:  "Save all changes with an interactive commit message",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// Get commit message
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter commit message: ")
			message, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read commit message: %w", err)
			}
			message = strings.TrimSpace(message)

			if message == "" {
				return fmt.Errorf("commit message cannot be empty")
			}

			// Load config
			config, err := repo.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Create storage
			storage := core.NewStorage(repo.ObjectsDir)

			// Load or create index
			index, err := core.LoadIndex(repo.IndexPath)
			if err != nil {
				index = core.NewIndex()
			}

			// Add all files to index
			err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip .zark directory
				if strings.Contains(path, ".zark") {
					return nil
				}

				// Skip directories
				if info.IsDir() {
					return nil
				}

				// Read file content
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", path, err)
				}

				// Create blob
				blob := core.NewBlob(content)
				if err := storage.Store(blob); err != nil {
					return fmt.Errorf("failed to store blob: %w", err)
				}

				// Add to index
				relPath, _ := filepath.Rel(cwd, path)
				index.Add(relPath, blob.Hash(), "100644", info.Size(), info.ModTime())

				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to add files to index: %w", err)
			}

			// Create tree from index
			var treeEntries []core.TreeEntry
			for _, entry := range index.Entries {
				treeEntries = append(treeEntries, core.TreeEntry{
					Mode: entry.Mode,
					Name: entry.Path,
					Hash: entry.Hash,
					Type: "blob",
				})
			}

			tree := core.NewTree(treeEntries)
			if err := storage.Store(tree); err != nil {
				return fmt.Errorf("failed to store tree: %w", err)
			}

			// Get parent commit (if exists)
			var parent string
			if headData, err := os.ReadFile(repo.HeadPath); err == nil {
				headRef := strings.TrimSpace(string(headData))
				if strings.HasPrefix(headRef, "ref: ") {
					refPath := filepath.Join(repo.ZarkDir, headRef[5:])
					if refData, err := os.ReadFile(refPath); err == nil {
						parent = strings.TrimSpace(string(refData))
					}
				}
			}

			// Create commit
			commit := core.NewCommit(tree.Hash(), parent, config.User.Name, config.User.Email, message)
			if err := storage.Store(commit); err != nil {
				return fmt.Errorf("failed to store commit: %w", err)
			}

			// Update HEAD
			headRefPath := filepath.Join(repo.RefsDir, "heads", "main")
			if err := os.MkdirAll(filepath.Dir(headRefPath), 0755); err != nil {
				return fmt.Errorf("failed to create refs directory: %w", err)
			}

			if err := os.WriteFile(headRefPath, []byte(commit.Hash()+"\n"), 0644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}

			// Save index
			if err := index.Save(repo.IndexPath); err != nil {
				return fmt.Errorf("failed to save index: %w", err)
			}

			fmt.Printf("Saved changes in commit %s\n", commit.Hash()[:8])
			return nil
		},
	}
}