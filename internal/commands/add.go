package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// AddCmd creates the `zark add` command.
func AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [file...]",
		Short: "Add file contents to the index",
		Long:  "This command updates the index using the current content found in the working tree, preparing the content for the next commit.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// Call the core logic for adding files
			return core.AddFiles(repo, args)
		},
	}
}