package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/mimetype"

	"github.com/spf13/cobra"
)

var (
	mimeInputPath string
)

var mimetypeCmd = &cobra.Command{
	Use:   "mimetype",
	Short: "Check the mimetype of a file",
	Run: func(cmd *cobra.Command, args []string) {
		tmpPath, err := cfs.ValidatePath(mimeInputPath, true, cfs.PathIsFile, cfs.PathIsDirectory)
		clog.CheckIfError(err)
		clog.InfoSuccessf(tmpPath.GetPath())
		mimetype.CheckMimeType(*tmpPath)
	},
}

func init() {
	rootCmd.AddCommand(mimetypeCmd)

	mimetypeCmd.Flags().StringVarP(&mimeInputPath, "input", "i", "./", "Path to DIR/FILE which will be verified")
}
