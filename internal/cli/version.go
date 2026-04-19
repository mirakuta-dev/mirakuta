package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Mirakuta",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Mirakuta v%s\n", version)
	},
}
