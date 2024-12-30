package cmd

import (
	"mateusjdev/scruffy/cmd/rhash"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rhashCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	rhashCmd.Flags().StringP("hash", "H", "blake3", "hash that will be used: [md5/blake3/blake2b/sha1/sha256/sha512/fuzzy]")

	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3)
	rhashCmd.Flags().StringP("input", "i", "./", "Files that will be hashed")

	// INFO: If --output/defaultOutputPath is not declared, it will be the same as --input/defaultInputPath
	rhashCmd.Flags().StringP("output", "o", "", "Location were hashed files will be stored")

	// TODO(3): Work on dry-run flag
	rhashCmd.Flags().BoolP("dry-run", "d", false, "Doesn't rename or delete files'")

	// TODO(4): Move to rootCmd
	rhashCmd.Flags().BoolP("debug", "D", false, "Print debug logs")
	// INFO: Sets LogLevel to Warning
	rhashCmd.Flags().BoolP("silent", "s", false, "SHHHHHHH! Doesn't print to stdout (runs way faster!)")
	// \ TODO(4): Move to rootCmd

	// TODO(5): Work on verbose flag
	rhashCmd.Flags().BoolP("verbose", "v", false, "Show full path")

	// TODO(6): Work on uppercase flag
	rhashCmd.Flags().BoolP("uppercase", "u", false, "Convert characters to UPPERCASE when possible")

	// TODO(7): Work on recursive flag
	rhashCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept output folder")

	// INFO: Ignore git checks, maybe will do something more later
	rhashCmd.Flags().BoolP("force", "F", false, "Ignore git checks")

	// TODO(8): Work on lenght/truncate flag
	// Lenght used in filename for blake3 and fuzzy algorithms
	rhashCmd.Flags().Int8P("lenght", "l", 16, "Truncate filename")

	// TODO(9): Check need of setting log level via flags (Ex: --log INFO, DEBUG, WARNING, ...)
	rhashCmd.MarkFlagsMutuallyExclusive("debug", "silent")
	rhashCmd.MarkFlagsMutuallyExclusive("verbose", "silent")

	// TODO(10): Recreate folder structure on destination Dir
	// For now --recursive and --output will be mutually exclusive
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
