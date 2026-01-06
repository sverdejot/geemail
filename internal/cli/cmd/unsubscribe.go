package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/sverdejot/geemail/internal/cli/tui"
	"github.com/sverdejot/geemail/internal/core"
	"github.com/sverdejot/geemail/internal/core/token"
)

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "unsubscribe from mailing list",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := token.NewHTTPClient(ctx)
		if err != nil {
			return fmt.Errorf("Unable to create default HTTP client: %v", err)
		}

		service, err := core.NewMessageService(ctx, client)
		if err != nil {
			return fmt.Errorf("Unable to create message service: %v", err)
		}
        m, err := tui.NewRoot(ctx, service)
        if err != nil {
            return err
        }
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("Error running program: %w", err)
		}
		return nil
	},
}
