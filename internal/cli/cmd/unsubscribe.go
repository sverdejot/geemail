package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unsubscribeCmd = &cobra.Command{
    Use: "unsubscribe",
    Short: "unsubscribe from mailing list",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Launching unsubscribe TUI...")
		return nil
    },
}
