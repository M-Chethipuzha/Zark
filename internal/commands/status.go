package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// StatusCmd creates the `zark status` command.
func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the working tree status",
		Long:  "Displays paths that have differences between the index file and the current HEAD commit, paths that have differences between the working tree and the index file, and paths in the working tree that are not tracked by Zark.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// Call the core logic for getting status
			return core.GetStatus(repo)
		},
	}
}