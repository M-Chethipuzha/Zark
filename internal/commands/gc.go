package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"zark/internal/core"
)

// GCCmd creates the `zark gc` (garbage collect) command.
func GCCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gc",
		Short: "Cleanup unnecessary files and optimize the local repository",
		Long:  "This command runs a number of housekeeping tasks within the current repository, such as compressing file revisions (to save disk space and increase performance).",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			repo := core.NewRepository(cwd)
			if !repo.Exists() {
				return fmt.Errorf("not a zark repository (or any of the parent directories)")
			}

			packer, err := core.NewPacker(repo)
			if err != nil {
				return fmt.Errorf("failed to initialize packer: %w", err)
			}

			objectCount := packer.GetObjectCount()
			if objectCount == 0 {
				fmt.Println("Repository is already optimized. Nothing to do.")
				return nil
			}

			// --- Beginner-Friendly Additions ---
			fmt.Printf("This command will optimize your repository for better performance and less disk space.\n")
			fmt.Printf("It will pack %d loose objects into a single, more efficient file.\n\n", objectCount)
			fmt.Print("Do you want to continue? (y/N) ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(response)) != "y" {
				fmt.Println("Optimization cancelled.")
				return nil
			}
			// ---

			fmt.Println("\nPacking objects...")
			packHash, err := packer.PackObjects()
			if err != nil {
				// Provide a more specific error message if packing fails
				if strings.Contains(err.Error(), "no loose objects to pack") {
					fmt.Println("No loose objects found to pack. Your repository is already optimized.")
					return nil
				}
				return fmt.Errorf("failed to pack objects: %w", err)
			}

			fmt.Println("\nOptimization complete!")
			fmt.Printf("Created packfile: pack-%s.pack\n", packHash)
			return nil
		},
	}
}