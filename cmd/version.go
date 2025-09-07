package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const ApplicationName = "Scruffy"

var ApplicationVersion = "build from source"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of scruffy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s - %s", ApplicationName, ApplicationVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
