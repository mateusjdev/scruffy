package cmd

import (
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"mateusjdev/scruffy/cmd/mimetype"

	"github.com/spf13/cobra"
)

var mediatypeCmd = &cobra.Command{
	Use:     "mediatype",
	Aliases: []string{"mimetype", "mime"},
	Short:   "Check the mimetype of a file",
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

		ignoreExtMatch, err := cmd.Flags().GetBool("ignore-ok")
		if err != nil {
			return err
		}

		allowRename, err := cmd.Flags().GetBool("rename")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return mimetype.CheckPath(
			inputPathInfo,
			ignoreExtMatch,
			allowRename,
			recursive,
		)
	},
}

func init() {
	rootCmd.AddCommand(mediatypeCmd)

	mediatypeCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be verified.")
	mediatypeCmd.Flags().BoolP("ignore-ok", "w", false, "Ignore files if extension match mimetype.")
	mediatypeCmd.Flags().BoolP("rename", "R", false, "Rename file extensions.")
	mediatypeCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs.")
}
