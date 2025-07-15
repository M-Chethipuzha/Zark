package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// SearchCmd creates the `zark search` command.
func SearchCmd() *cobra.Command {
	var author, message string
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for commits in the repository",
		Long:  "Search for commits by message, author, or content.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository")
			}

			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			results, err := core.SearchCommits(repo, query, author, message)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Println("No commits found.")
				return nil
			}

			for _, commit := range results {
				fmt.Printf("commit %s\n", commit.Hash())
				fmt.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
				fmt.Printf("Date:   %s\n", commit.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"))
				fmt.Printf("\n\t%s\n\n", commit.Message)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Search for commits by author")
	cmd.Flags().StringVar(&message, "message", "", "Search for commits by message content")

	return cmd
}