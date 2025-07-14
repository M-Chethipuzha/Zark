package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// CheckoutCmd creates the `zark checkout` command.
func CheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout [branch-or-commit]",
		Short: "Switch branches or restore working tree files",
		Long:  "This command switches branches or restores working tree files. It updates the files in the working tree to match the version in the specified branch or commit.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			// Call the core logic for checkout
			return core.Checkout(repo, args[0])
		},
	}
}