package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/mimetype"

	"github.com/spf13/cobra"
)

var (
	mimeInputPath string
	ignoreOk      bool
)

var mimetypeCmd = &cobra.Command{
	Use:   "mimetype",
	Short: "Check the mimetype of a file",
	Run: func(cmd *cobra.Command, args []string) {
		tmpPath, err := cfs.ValidatePath(mimeInputPath, skipGitCheck, cfs.PathIsFile, cfs.PathIsDirectory)
		clog.PanicIf(err)
		mimetype.CheckMimeType(*tmpPath, ignoreOk)
	},
}

func init() {
	rootCmd.AddCommand(mimetypeCmd)

	mimetypeCmd.Flags().StringVarP(&mimeInputPath, "input", "i", "./", "Path to DIR/FILE which will be verified")

	mimetypeCmd.Flags().BoolVarP(&ignoreOk, "ignore-ok", "w", false, "Ignore files if extension match mimetype")
}
