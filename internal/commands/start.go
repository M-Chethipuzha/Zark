package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// StartCmd creates the `zark start` command.
func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Initialize a new Zark repository",
		Long:  "Initialize a new Zark repository in the current directory, creating the necessary .zark directory structure.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if repo.Exists() {
				// Re-initializing is safe, but we'll inform the user.
				fmt.Printf("Reinitialized existing Zark repository in %s/.zark\n", cwd)
			}

			if err := repo.Init(); err != nil {
				return fmt.Errorf("failed to initialize repository: %w", err)
			}

			if !repo.Exists() {
				fmt.Printf("Initialized empty Zark repository in %s/.zark\n", cwd)
			}
			return nil
		},
	}
}