package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// TODO: dynamic
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of scruffy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Scruffy v0.2.0")
	},
}
