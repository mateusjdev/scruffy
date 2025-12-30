package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

var (
	mimeInputPath string
)

var mimetypeCmd = &cobra.Command{
	Use:   "mimetype",
	Short: "Check the mimetype of a file",
	Run: func(cmd *cobra.Command, args []string) {
		tmpPath, err := cfs.ValidatePath(mimeInputPath, skipGitCheck, cfs.PathIsFile)
		clog.CheckIfError(err)
		mtype, err := mimetype.DetectFile(tmpPath.GetPath())
		if err != nil {
			panic("ERROR")
		}
		clog.InfoSuccessf("%s %s", mtype.String(), mtype.Extension())
	},
}

func init() {
	rootCmd.AddCommand(mimetypeCmd)

	mimetypeCmd.Flags().StringVarP(&mimeInputPath, "input", "i", "./", "Path to DIR/FILE which will be verified")
}
