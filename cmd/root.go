package cmd

import (
	"mateusjdev/scruffy/cmd/clog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scruffy",
	Short: "A digital janitor with helpfull scripts",
	RunE: func(cmd *cobra.Command, args []string) error {
		clog.Debugf("Starting module::%s", cmd.Use)

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		if debug {
			clog.SetLogLevel(clog.LevelDebug)
		}

		return nil
	},
}

func init() {
	// Sets LogLevel to Debug
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Set log-level to debug")

	// Sets LogLevel to Warning
	rootCmd.PersistentFlags().BoolP("silent", "s", false, "Set log-level to warning (some scripts will run way faster!)")

	// TODO(9): Check the need to configure LogLevel via flag (Ex: --log-level DEBUG, INFO, WARNING, ERROR)
	rootCmd.MarkFlagsMutuallyExclusive("debug", "silent")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		clog.Errorf("%s\n", err)
		os.Exit(1)
	}
}
