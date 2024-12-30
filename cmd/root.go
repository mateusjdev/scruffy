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

func init() {}
