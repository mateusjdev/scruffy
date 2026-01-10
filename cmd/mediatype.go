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

		mimetype.CheckPath(inputPathInfo, ignoreExtMatch)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mediatypeCmd)

	mediatypeCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be verified")
	mediatypeCmd.Flags().BoolP("ignore-ok", "w", false, "Ignore files if extension match mimetype")
}
