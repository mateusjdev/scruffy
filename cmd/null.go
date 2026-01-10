package cmd

import (
	"mateusjdev/scruffy/cmd/check"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"

	"github.com/spf13/cobra"
)

var nullCmd = &cobra.Command{
	Use:   "null",
	Short: "Check if a file is only \\0 bytes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		clog.Debugf("Starting module::%s", cmd.Use)

		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		inputPathInfo, err := filesystem.StatPath(inputPath)
		if err != nil {
			return err
		}

		ignoreNonEmpty, err := cmd.Flags().GetBool("ignore-ok")
		if err != nil {
			return err
		}

		allowRename, err := cmd.Flags().GetBool("rename")
		if err != nil {
			return err
		}

		err = check.CheckPathForEmptyFiles(inputPathInfo, true, allowRename, ignoreNonEmpty)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(nullCmd)

	nullCmd.Flags().StringP("input", "i", "./", "Path to FILE which will be verified")
	nullCmd.Flags().BoolP("ignore-ok", "w", false, "Ignore files if has content.")
	nullCmd.Flags().BoolP("rename", "R", false, "Mark file extension as .empty")
}
