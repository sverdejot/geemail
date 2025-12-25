package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "geemail",
    Short: "Fast, bulk Gmail inbox cleanup",
    Long:  "A TUI tool to aggressively reduce unread Gmail messages",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(unsubscribeCmd)
    rootCmd.AddCommand(generateTokenCmd)
    rootCmd.AddCommand(summaryCmd)
}
