package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/sverdejot/geemail/internal/core/auth"
	"github.com/sverdejot/geemail/internal/core"
)

var summaryCmd = &cobra.Command{
    Use: "summary",
    Short: "Shows a summary of the latest unread messages",
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

        t0 := time.Now()
        contents, err := service.GetContent(cmd.Context())
        t1 := time.Now().Sub(t0)

        if err != nil {
            return fmt.Errorf("Unable to retrieve latest messages: %v", err)
        }

        fmt.Println(core.FormatCount(contents))
        fmt.Println("took: ", t1)

        return nil
    },
}
