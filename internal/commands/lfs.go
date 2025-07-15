package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// LFSCmd creates the `zark lfs` command.
func LFSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lfs",
		Short: "Manage large files with Large File Storage (LFS)",
		Long:  `LFS replaces large files such as audio samples, videos, datasets, and graphics with tiny text pointers inside Zark. The large files themselves are stored on a remote server.`,
	}

	cmd.AddCommand(lfsTrackCmd())

	return cmd
}

func lfsTrackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "track [pattern]",
		Short: "Start tracking a file pattern with LFS",
		Long:  `Configures Zark to store files matching a path pattern with LFS. For example, 'zark lfs track "*.zip"' will track all .zip files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository. Please run 'zark start' to initialize a new repository")
			}

			pattern := args[0]

			// --- Beginner-Friendly Additions ---
			fmt.Printf("This will configure Zark to track files matching '%s' using Large File Storage (LFS).\n", pattern)
			fmt.Println("Instead of storing the large file directly in the repository, Zark will store a lightweight pointer.")
			fmt.Println("This keeps your repository small and fast.")
			// ---

			if err := core.TrackLFS(repo, pattern); err != nil {
				return err
			}

			fmt.Printf("\nSuccessfully tracking '%s'.\n", pattern)
			fmt.Println("Next steps: Make sure '.zarkattributes' is added and committed to your repository.")
			fmt.Println("  1. Run 'zark add .zarkattributes'")
			fmt.Println("  2. Run 'zark add your_large_file.ext'")
			fmt.Println("  3. Run 'zark save -m \"Track large files\"'")
			return nil
		},
	}
}