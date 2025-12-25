package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sverdejot/geemail/internal/core/auth"
)

var generateTokenCmd = &cobra.Command{
    Use: "generate-token",
    Short: "Login into gmail using the CLI",
    RunE: func(cmd *cobra.Command, args []string) error {
        b, err := os.ReadFile("config/credentials.json")
        if err != nil {
            return fmt.Errorf("Unable to read client secret file: %v", err)
        }

        if err := auth.GenerateToken(b); err != nil {
            return fmt.Errorf("Cannot generate token: %v", err)
        }
        return nil
    },
}
