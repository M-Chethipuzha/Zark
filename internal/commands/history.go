package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// HistoryCmd creates the `zark history` command.
func HistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Show commit history",
		Long:  "Display the commit history of the current branch, starting from the most recent commit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// Call the core logic for showing history
			return core.ShowHistory(repo)
		},
	}
}
