package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zark/internal/commands"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "zark",
		Short: "Next-generation version control system",
		Long:  "Zark is a user-friendly, high-performance version control system",
	}

	rootCmd.AddCommand(commands.StartCmd())
	rootCmd.AddCommand(commands.SaveCmd())
	rootCmd.AddCommand(commands.HistoryCmd())
	rootCmd.AddCommand(commands.AddCmd())
	rootCmd.AddCommand(commands.StatusCmd())
	rootCmd.AddCommand(commands.CheckoutCmd())
	rootCmd.AddCommand(commands.BranchCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}