package cmd

import (
	"mateusjdev/scruffy/cmd/rhash"

	"github.com/spf13/cobra"
)

var rhashCmd = &cobra.Command{
	Use:   "rhash",
	Short: "Rename files to their hash sum",
	Run:   rhash.RenameFilesToHash,
}

func init() {
	rootCmd.AddCommand(rhashCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	rhashCmd.Flags().StringP("hash", "H", "blake3", "hash that will be used: [md5/blake3/blake2b/sha1/sha256/sha512/fuzzy]")

	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3)
	// TODO(2): Drop -i and use 'scruffy rhash $i $2 $3'
	rhashCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be hashed")

	// INFO: If --output/defaultOutputPath is not declared, it will be the same as --input/defaultInputPath
	rhashCmd.Flags().StringP("output", "o", "", "Location were hashed files will be stored")

	rhashCmd.Flags().BoolP("abbreviate-path", "a", false, "Abbreviate paths when logging")

	rhashCmd.Flags().BoolP("dry-run", "d", false, "Don't rename files'")

	rhashCmd.Flags().BoolP("uppercase", "U", false, "Convert characters to UPPERCASE")

	// Ignore git checks
	rhashCmd.Flags().BoolP("force", "F", false, "Ignore git checks")

	// Truncate filename, max value is 256(uint8)
	rhashCmd.Flags().Uint8P("truncate", "t", 32, "Truncate filename")

	rhashCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept a target directory")

	// TODO(10): Recreate folder structure on destination Dir
	// For now --recursive and --output will be mutually exclusive
	// If --output is declared, rhash will not recurse into --input folders
	// If --recursive is declared, rhash will not accept another folder as output
	rhashCmd.MarkFlagsMutuallyExclusive("recursive", "output")

	// TODO: rootCmd.silent x rhashCmd.*
	// Why dry-run if nothing will be shown on screen?
	// rhashCmd.MarkFlagsMutuallyExclusive("silent", "dry-run")
	// Why abreviate paths if nothing will be shown on screen?
	// rhashCmd.MarkFlagsMutuallyExclusive("silent", "abbreviate-path")

	rhashCmd.MarkFlagFilename("input")
	rhashCmd.MarkFlagDirname("output")
}
