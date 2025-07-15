package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// BranchCmd creates the `zark branch` command.
func BranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "List, create, or delete branches",
	}

	cmd.AddCommand(branchCreateCmd())
	cmd.AddCommand(branchListCmd())

	return cmd
}

func branchCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			var branchName string
			if len(args) > 0 {
				branchName = args[0]
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter branch name: ")
				name, _ := reader.ReadString('\n')
				branchName = strings.TrimSpace(name)
			}

			if branchName == "" {
				return fmt.Errorf("branch name cannot be empty")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository")
			}

			return core.CreateBranch(repo, branchName)
		},
	}
}

func branchListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository")
			}

			return core.ListBranches(repo)
		},
	}
}