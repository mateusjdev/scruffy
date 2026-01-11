package cmd

import (
	"errors"
	"fmt"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"mateusjdev/scruffy/cmd/hash"
	"mateusjdev/scruffy/cmd/hasher"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Check hash of files.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			return err
		}

		if debug || quiet {
			clog.Warningf("Nothing to do!")
			os.Exit(0)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		clog.Debugf("Starting module::%s", cmd.Use)

		truncate, err := cmd.Flags().GetUint8("truncate")
		if err != nil {
			return err
		}

		if truncate < ARGS_MIN_TRUNCATE {
			return fmt.Errorf("--truncate is very low, choose >= %d", ARGS_MIN_TRUNCATE)
		}

		if truncate > ARGS_MAX_TRUNCATE {
			return fmt.Errorf("--truncate is very high, choose <= %d", ARGS_MAX_TRUNCATE)
		}

		// TODO: validate hash here
		algorithm, err := cmd.Flags().GetString("algorithm")
		if err != nil {
			return err
		}

		algorithm = strings.ToLower(algorithm)

		uppercase, err := cmd.Flags().GetBool("uppercase")
		if err != nil {
			return err
		}

		displayAbsPath, err := cmd.Flags().GetBool("absolute-path")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		if inputPath == "" {
			return errors.New("--input is empty or invalid")
		}

		inputPathInfo, err := filesystem.StatPath(inputPath)
		if err != nil {
			return err
		}

		clog.Debugf(
			"Args:\n%s\n%s\n%s\n%s\n%s\n%s\n",
			fmt.Sprintf("recursive: %t", recursive),
			fmt.Sprintf("absolute-path: %t", displayAbsPath),
			fmt.Sprintf("uppercase: %t", uppercase),
			fmt.Sprintf("truncate: %d", truncate),
			fmt.Sprintf("inputPath: %s", inputPath),
			fmt.Sprintf("hashAlgorithm: %s", algorithm),
		)

		// HASH_MACHINE
		mHasher, err := hasher.NewHasher(algorithm, int(truncate))
		if err != nil {
			return err
		}

		// PATH_WALK
		return hash.HashFromPath(mHasher, recursive, inputPathInfo, truncate, uppercase)
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	// Rename :: Methods
	hashCmd.Flags().StringP("algorithm", "a", "blake3", "Use file hash [md5/blake3/blake2b/sha1/sha256/sha512]")

	// Rename :: Path
	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3) || Drop -i and use 'scruffy rhash $i $2 $3'
	hashCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be hashed")

	hashCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept a target directory")

	// Rename :: Options
	hashCmd.Flags().BoolP("uppercase", "U", false, "Convert characters to UPPERCASE")
	// recommended max filename is 256
	hashCmd.Flags().Uint8P("truncate", "t", 32, "Truncate filename (Beetween 8 and 128)")

	// Rename :: Logging
	hashCmd.Flags().BoolP("absolute-path", "P", false, "Print absolute paths relative when logging")
}
