package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// SaveCmd creates the `zark save` command.
func SaveCmd() *cobra.Command {
	var message string
	var sign bool
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save staged changes to the repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if message == "" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter commit message: ")
				msg, _ := reader.ReadString('\n')
				message = strings.TrimSpace(msg)
			}

			if message == "" {
				return fmt.Errorf("commit message cannot be empty")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository")
			}

			// Secret Scanning
			if err := core.ScanForSecrets(repo); err != nil {
				return err
			}

			commitHash, err := core.CreateCommit(repo, message, sign)
			if err != nil {
				return err
			}

			fmt.Printf("Saved changes in commit %s\n", commitHash[:8])
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	cmd.Flags().BoolVarP(&sign, "sign", "s", false, "Sign the commit with GPG")

	return cmd
}