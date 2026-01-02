package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/mimetype"

	"github.com/spf13/cobra"
)

var mediatypeCmd = &cobra.Command{
	Use:     "mediatype",
	Aliases: []string{"mimetype"},
	Short:   "Check the mimetype of a file",
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, err := cmd.Flags().GetString("input")
		clog.PanicIf(err)

		inputPathInfo, err := cfs.StatPath(inputPath)
		clog.PanicIf(err)

		ignoreExtMatch, err := cmd.Flags().GetBool("ignore-ok")
		clog.PanicIf(err)

		mimetype.CheckPath(inputPathInfo, ignoreExtMatch)
	},
}

func init() {
	rootCmd.AddCommand(mediatypeCmd)

	mediatypeCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be verified")
	mediatypeCmd.Flags().BoolP("ignore-ok", "w", false, "Ignore files if extension match mimetype")
}
