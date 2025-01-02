/*
Copyright © 2024 Mateus J. <git@mateusj.dev>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scruffy",
	Short: "A digital maid with helpfull scripts",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Print debug logs")

	// INFO: Sets LogLevel to Warning
	rootCmd.PersistentFlags().BoolP("silent", "s", false, "SHHHHHHH! Doesn't print to stdout (some scripts will run way faster!)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show more information about what are being done")

	// TODO(9): Check need of setting log level via flags (Ex: --log INFO, DEBUG, WARNING, ...)
	rhashCmd.MarkFlagsMutuallyExclusive("debug", "silent")
	rhashCmd.MarkFlagsMutuallyExclusive("verbose", "silent")
}
