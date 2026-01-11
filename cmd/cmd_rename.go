package cmd

import (
	"fmt"
	"mateusjdev/scruffy/internal/clog"
	"mateusjdev/scruffy/internal/filesystem"
	"mateusjdev/scruffy/internal/hasher"
	"mateusjdev/scruffy/internal/renamer"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"rhash"},
	Short:   "Batch rename files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		clog.Debugf("Starting module::%s", cmd.Use)

		silent, _ := cmd.Flags().GetBool("silent")
		if silent {
			clog.SetLogLevel(clog.LevelWarning)
		}

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
		hashAlgorithm, err := cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}

		hashAlgorithm = strings.ToLower(hashAlgorithm)

		uppercase, err := cmd.Flags().GetBool("uppercase")
		if err != nil {
			return err
		}

		displayAbsPath, err := cmd.Flags().GetBool("absolute-path")
		if err != nil {
			return err
		}

		random, err := cmd.Flags().GetBool("random")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}

		skipGitCheck, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		clog.Debugf("Args:\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			fmt.Sprintf("dry-run: %t", dryRun),
			fmt.Sprintf("silent: %t", silent),
			fmt.Sprintf("recursive: %t", recursive),
			fmt.Sprintf("absolute-path: %t", displayAbsPath),
			fmt.Sprintf("skipGitCheck: %t", skipGitCheck),
			fmt.Sprintf("uppercase: %t", uppercase),
			fmt.Sprintf("truncate: %d", truncate),
			fmt.Sprintf("inputPath: %s", inputPath),
			fmt.Sprintf("outputPath: %s", outputPath),
			fmt.Sprintf("hash: %s", hashAlgorithm),
			fmt.Sprintf("random: %t", random),
		)

		if inputPath == "" {
			clog.Panicf("--input is empty or invalid")
		}

		inputPathInfo, err := filesystem.StatPath(inputPath)
		if err != nil {
			return err
		}

		if outputPath == "" {
			if cmd.Flags().Lookup("output").Changed {
				clog.Panicf("--output is empty or invalid")
			}

			if inputPathInfo.IsRegularFile() {
				outputPath = filepath.Dir(inputPathInfo.Path())
			} else {
				outputPath = inputPathInfo.Path()
			}
		}

		// TODO(11): Create destinationPath if doesn't exist (maybe add a flag? force?)
		outputPathInfo, err := filesystem.StatPath(outputPath)
		if err != nil {
			return err
		}
		if !outputPathInfo.IsDir() {
			clog.Panicf("Destination folder \"%s\" is not a valid directory\n", outputPath)
		}

		isInGitRepo, err := filesystem.IsPathInGitRepo(inputPathInfo)
		if isInGitRepo {
			if !skipGitCheck {
				clog.Panicf("%s is in a git repo", inputPathInfo.Path())
			}
			clog.Infof("%s is in a git repo", inputPathInfo.Path())
		}

		if inputPathInfo.Path() != outputPathInfo.Path() {
			isInGitRepo, err := filesystem.IsPathInGitRepo(outputPathInfo)
			if err != nil {
				return err
			}
			if isInGitRepo {
				if !skipGitCheck {
					clog.Panicf("%s is in a git repo", outputPathInfo.Path())
				}
				clog.Infof("%s is in a git repo", outputPathInfo.Path())
			}
		}

		clog.Debugf("inputPathInfo.GetPath(): %s", inputPathInfo.Path())
		clog.Debugf("outputPathInfo.GetPath(): %s", outputPathInfo.Path())

		currentWorkDir, err := os.Getwd()
		if err != nil {
			return err
		}

		var renameMethod renamer.Renamer
		if random {
			// FUZZY_MACHINE
			renameMethod, err = renamer.NewFuzzyRenamer(
				uppercase,
				truncate,
				// INFO: For file naming this (dryRun) will be random,
				// but at least it will show the destination path
				dryRun,
				displayAbsPath,
				currentWorkDir,
			)
			if err != nil {
				return err
			}
		} else {
			// HASH_MACHINE
			mHasher, err := hasher.NewHasher(hashAlgorithm, int(truncate))
			if err != nil {
				return err
			}

			renameMethod, err = renamer.NewHashRenamer(
				mHasher,
				uppercase,
				truncate,
				dryRun,
				displayAbsPath,
				currentWorkDir,
			)
			if err != nil {
				return err
			}
		}
		// PATH_WALK
		renamer.RenameFromPath(renameMethod, recursive, inputPathInfo, outputPathInfo)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	// Rename :: Methods
	renameCmd.Flags().StringP("hash", "H", "blake3", "Use file hash [md5/blake3/blake2b/sha1/sha256/sha512]")
	renameCmd.Flags().BoolP("random", "R", false, "Use random characters.")
	renameCmd.MarkFlagsMutuallyExclusive("hash", "random")

	// Rename :: Path
	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3) || Drop -i and use 'scruffy rhash $i $2 $3'
	renameCmd.Flags().StringP("input", "i", "./", "Path to DIR/FILE which will be hashed")

	// INFO: If --output/defaultOutputPath is not declared, it will be the same as --input/defaultInputPath
	renameCmd.Flags().StringP("output", "o", "", "Location were hashed files will be stored")
	renameCmd.Flags().BoolP("recursive", "r", false, "Recurse DIRs, when enabled, will not accept a target directory")

	// Rename :: Options
	renameCmd.Flags().BoolP("uppercase", "U", false, "Convert characters to UPPERCASE")
	// recommended max filename is 256
	renameCmd.Flags().Uint8P("truncate", "t", 32, "Truncate filename (Beetween 8 and 128)")

	// Rename :: Logging
	renameCmd.Flags().BoolP("absolute-path", "A", false, "Print absolute paths relative when logging")
	renameCmd.Flags().BoolP("dry-run", "d", false, "Don't rename files")

	// Rename :: Other
	// Ignore git checks
	renameCmd.Flags().BoolP("force", "F", false, "Ignore git checks")

	// TODO(10): Recreate folder structure on destination Dir
	// For now --recursive and --output will be mutually exclusive
	// If --output is declared, rhash will not recurse into --input folders
	// If --recursive is declared, rhash will not accept another folder as output
	renameCmd.MarkFlagsMutuallyExclusive("recursive", "output")

	// TODO: rootCmd.silent x rhashCmd.*
	// Why dry-run if nothing will be shown on screen?
	// rhashCmd.MarkFlagsMutuallyExclusive("silent", "dry-run")
	// Why abreviate paths if nothing will be shown on screen?
	// rhashCmd.MarkFlagsMutuallyExclusive("silent", "abbreviate-path")

	renameCmd.MarkFlagDirname("output")
}
