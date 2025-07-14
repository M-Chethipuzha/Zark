package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// BranchCmd creates the `zark branch` command.
func BranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branch [branch-name]",
		Short: "List, create, or delete branches",
		Long:  "If no branch name is provided, it lists all local branches. Otherwise, it creates a new branch pointing to the current HEAD.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// If no arguments, list branches. Otherwise, create one.
			if len(args) == 0 {
				return core.ListBranches(repo)
			} else {
				return core.CreateBranch(repo, args[0])
			}
		},
	}
}