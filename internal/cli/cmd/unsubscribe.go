package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/sverdejot/geemail/internal/cli/tui"
	"github.com/sverdejot/geemail/internal/core"
	"github.com/sverdejot/geemail/internal/core/auth"
)

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "unsubscribe from mailing list",
	RunE: func(cmd *cobra.Command, args []string) error {
		credentialsPath := "config/credentials.json"
		if envCredentialsPath := os.Getenv("CREDENTIALS_FILE"); envCredentialsPath != "" {
			credentialsPath = envCredentialsPath
		}
		b, err := os.ReadFile(credentialsPath)
		if err != nil {
			return fmt.Errorf("Unable to read client secret file: %v", err)
		}

		client, err := auth.NewHTTPClient(b)
		if err != nil {
			return fmt.Errorf("Unable to create default HTTP client: %v", err)
		}

		service, err := core.NewMessageService(client)
		if err != nil {
			return fmt.Errorf("Unable to create message service: %v", err)
		}
		if _, err := tea.NewProgram(tui.NewModel(service), tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("Error running program: %w", err)
		}
		return nil
	},
}
