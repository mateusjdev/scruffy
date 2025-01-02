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

	rhashCmd.Flags().BoolP("dry-run", "d", false, "Doesn't rename or delete files'")

	rhashCmd.Flags().BoolP("uppercase", "U", false, "Convert characters to UPPERCASE when possible")

	rhashCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept output folder")

	// INFO: Ignore git checks, maybe will do something more later
	rhashCmd.Flags().BoolP("force", "F", false, "Ignore git checks")

	// INFO: Max value is 256(uint8)
	rhashCmd.Flags().Uint8P("truncate", "t", 32, "Truncate filename")

	// TODO(10): Recreate folder structure on destination Dir
	// For now --recursive and --output will be mutually exclusive
	// If --output is declared, rhash will not recurse into --input folders
	// If --recursive is declared, rhash will not accept another folder as output
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
