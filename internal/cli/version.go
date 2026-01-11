package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"  //nolint:unusued

var versionCmd = &cobra.Command{
    Use: "version",
    Short: "prints the current geemail version",
    Run: func(_ *cobra.Command, _ []string) {
        fmt.Printf("geemail version: %s\n", version)  
    },
}
