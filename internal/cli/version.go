package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is injected at build time via:
//
//	-ldflags "-X github.com/mirakuta-dev/mirakuta/internal/cli.Version=<tag>"
//
// Local `go build` without ldflags falls back to "dev".
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Mirakuta",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Mirakuta %s\n", Version)
	},
}
