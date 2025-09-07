package cmd

import (
	"errors"
	"fmt"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	ffprobeInputPath string
)

var ffprobeCmd = &cobra.Command{
	Use:   "ffprobe",
	Short: "Alias to ffprobe",
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, err := cmd.Flags().GetString("input")
		if err != nil || inputPath == "" {
			clog.Errorf("Error getting user input!")
			clog.ExitBecause(clog.ErrCodeGeneric)
		}
		if _, err := os.Stat(inputPath); err == nil {
			ffprobeCmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", inputPath)
			stdout, err := ffprobeCmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			// TODO: Use Json
			fmt.Println(string(stdout))
		} else if errors.Is(err, os.ErrNotExist) {
			clog.Errorf("File \"%s\" doesn't exist!", inputPath)
			clog.ExitBecause(clog.ErrUserInput)
		} else {
			clog.Errorf(err.Error())
			clog.ExitBecause(clog.ErrCodeGeneric)
		}
	},
}

func init() {
	rootCmd.AddCommand(ffprobeCmd)

	ffprobeCmd.Flags().StringVarP(&ffprobeInputPath, "input", "i", "", "Path to DIR/FILE which will be probed")
}
