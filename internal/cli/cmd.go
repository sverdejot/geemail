package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/sverdejot/geemail/internal/cli/tui"
	"github.com/sverdejot/geemail/internal/gmail"
	"github.com/sverdejot/geemail/internal/gmail/auth"
)

var rootCmd = &cobra.Command{
	Use:   "geemail",
	Short: "Fast, bulk Gmail inbox cleanup",
	Long:  "A TUI tool to aggressively reduce unread Gmail messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return fmt.Errorf("error reading flag: %w", err)
		}

		ctx := cmd.Context()
		client, err := auth.NewHTTPClient(ctx)
		if err != nil {
			return fmt.Errorf("unable to create default HTTP client: %v", err)
		}

		service, err := gmail.NewMessageService(ctx, client)
		if err != nil {
			return fmt.Errorf("unable to create message service: %v", err)
		}
		m, err := tui.NewRoot(ctx, service, dryRun)
		if err != nil {
			return err
		}
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

func Execute() {
	rootCmd.Flags().Bool("dry-run", false, "Simulate all the actions without performing any change")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
