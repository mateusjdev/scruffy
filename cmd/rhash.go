package cmd

import (
	"mateusjdev/scruffy/cmd/rhash"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rhashCmd)

	// TODO: add viper and XDG
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	// defaultHash := rename.BLAKE3
	// rootCmd.Flags().VarP(&defaultHash, "hash", "H", "hash that will be used: [md5/blake3/blake2b/sha1/sha256/sha512/fuzzy]")

	rhashCmd.Flags().StringP("hash", "H", "blake3", "hash that will be used: [md5/blake3/blake2b/sha1/sha256/sha512/fuzzy]")

	// TODO add multiple inputs: -i $1 -i $2 -i $3
	defaultInputPath := "./"
	rhashCmd.Flags().StringP("input", "i", defaultInputPath, "Files that will be hashed")

	// if not declared --output/defaultOutputPath will be the same as --input/defaultInputPath
	defaultOutputPath := ""
	rhashCmd.Flags().StringP("output", "o", defaultOutputPath, "Location were hashed files will be stored")

	rhashCmd.Flags().BoolP("dry-run", "d", false, "Doesn't rename or delete files'")
	// TODO: Imply --force
	rhashCmd.Flags().BoolP("debug", "D", false, "Print debug logs")
	rhashCmd.Flags().BoolP("silent", "s", false, "SHHHHHHH! Doesn't print to stdout (runs way faster!)")
	rhashCmd.Flags().BoolP("uppercase", "u", false, "Convert characters to UPPERCASE when possible")
	rhashCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept output folder")
	rhashCmd.Flags().BoolP("force", "F", false, "Ignore git checks")
	rhashCmd.Flags().BoolP("verbose", "v", false, "Show full path")

	rhashCmd.Flags().Int8P("lenght", "l", 16, "Truncate filename")
	// rhashCmd.Flags().Int8P("lenght", "l", 16, "Lenght used in filename for blake3 and fuzzy algorithms")

	// TODO: -v INFO, DEBUG, WARNING, ...
	rhashCmd.MarkFlagsMutuallyExclusive("debug", "silent")
	rhashCmd.MarkFlagsMutuallyExclusive("verbose", "silent")

	// For now recursive and output will be exclusive
	rhashCmd.MarkFlagsMutuallyExclusive("recursive", "output")
	rhashCmd.MarkFlagFilename("input")
	rhashCmd.MarkFlagDirname("output")
}

// TODO: dynamic
var rhashCmd = &cobra.Command{
	Use:   "rhash",
	Short: "Rename files to their hash sum",
	Run:   rhash.RenameFilesToHash,
}
