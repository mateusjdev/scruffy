package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ApplicationName = "Scruffy"
	// TODO: Set Version through Makefile
	ApplicationVersion = "v0.2.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of scruffy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s", ApplicationName, ApplicationVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
