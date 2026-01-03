package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/hasher"
	"mateusjdev/scruffy/cmd/renamer"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ARGS_MIN_TRUNCATE = 8
	ARGS_MAX_TRUNCATE = 128
)

var (
	hashAlgorithm string
	random        bool

	displayAbsPath bool
	dryRun         bool
	force          bool
	recursive      bool
	truncate       uint8
	uppercase      bool
	skipGitCheck   bool
	inputPath      string
	outputPath     string
)

var renameCmd = &cobra.Command{
	Use:     "rename [source]",
	Aliases: []string{"rhash"},
	Short:   "Batch rename files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check LogLevel (global-flags)
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			clog.SetLogLevel(clog.LevelDebug)
		}

		silent, _ := cmd.Flags().GetBool("silent")
		if silent {
			clog.SetLogLevel(clog.LevelWarning)
		}
		// Check LogLevel (global-flags)

		if truncate < ARGS_MIN_TRUNCATE {
			clog.Panicf("--truncate is very low, choose >= %d", ARGS_MIN_TRUNCATE)
		}

		if truncate > ARGS_MAX_TRUNCATE {
			clog.Panicf("--truncate is very high, choose <= %d", ARGS_MAX_TRUNCATE)
		}

		skipGitCheck = force || debug || dryRun

		// TODO: validate hash here
		hashAlgorithm = strings.ToLower(hashAlgorithm)

		clog.Debugf(`Args:
	dry-run: %t
	silent: %t
	recursive: %t
	absolute-path: %t
	skipGitCheck: %t
	uppercase: %t
	truncate: %d
	inputPath: %s
	outputPath: %s
	hash: %s
	random: %t`, dryRun, silent, recursive, displayAbsPath, skipGitCheck, uppercase, truncate, inputPath, outputPath, hashAlgorithm, random)
	},
	Run: func(cmd *cobra.Command, args []string) {
		clog.Debugf("Starting module::%s", cmd.Use)

		if inputPath == "" {
			clog.Panicf("--input is empty or invalid")
		}

		inputPathInfo, err := cfs.StatPath(inputPath)
		clog.PanicIf(err)
		// TODO: Wrap if os.ErrNotExist
		/*
			if inputPathInfo.GetPathType() == cfs.PathIsNonExistent {
				clog.Panicf("Source path %s is not a valid file or a directory\n", inputPath)
			}
		*/

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
		outputPathInfo, err := cfs.StatPath(outputPath)
		clog.PanicIf(err)
		if !outputPathInfo.IsDir() {
			clog.Panicf("Destination folder \"%s\" is not a valid directory\n", outputPath)
		}

		isInGitRepo, err := cfs.IsPathInGitRepo(inputPathInfo)
		if isInGitRepo {
			if !skipGitCheck {
				clog.Panicf("%s is in a git repo", inputPathInfo.Path())
			}
			clog.Infof("%s is in a git repo", inputPathInfo.Path())
		}

		if inputPathInfo.Path() != outputPathInfo.Path() {
			isInGitRepo, err := cfs.IsPathInGitRepo(outputPathInfo)
			clog.PanicIf(err)
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
		clog.PanicIf(err)

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
			clog.PanicIf(err)
		} else {
			// HASH_MACHINE
			mHasher, err := hasher.NewHasher(hashAlgorithm, int(truncate))
			clog.PanicIf(err)

			renameMethod, err = renamer.NewHashRenamer(
				mHasher,
				uppercase,
				truncate,
				dryRun,
				displayAbsPath,
				currentWorkDir,
			)
			clog.PanicIf(err)
		}
		// PATH_WALK
		renamer.RenameFromPath(renameMethod, recursive, inputPathInfo, outputPathInfo)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	// Rename :: Methods
	renameCmd.Flags().StringVarP(&hashAlgorithm, "hash", "H", "blake3", "Use file hash [md5/blake3/blake2b/sha1/sha256/sha512]")
	renameCmd.Flags().BoolVarP(&random, "random", "R", false, "Use random characters.")
	renameCmd.MarkFlagsMutuallyExclusive("hash", "random")

	// Rename :: Path
	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3) || Drop -i and use 'scruffy rhash $i $2 $3'
	renameCmd.Flags().StringVarP(&inputPath, "input", "i", "./", "Path to DIR/FILE which will be hashed")

	// INFO: If --output/defaultOutputPath is not declared, it will be the same as --input/defaultInputPath
	renameCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Location were hashed files will be stored")
	renameCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recurse DIRs, when enabled, will not accept a target directory")

	// Rename :: Options
	renameCmd.Flags().BoolVarP(&uppercase, "uppercase", "U", false, "Convert characters to UPPERCASE")
	// recommended max filename is 256
	renameCmd.Flags().Uint8VarP(&truncate, "truncate", "t", 32, "Truncate filename (Beetween 8 and 128)")

	// Rename :: Logging
	renameCmd.Flags().BoolVarP(&displayAbsPath, "absolute-path", "A", false, "Print absolute paths relative when logging")
	renameCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Don't rename files")

	// Rename :: Other
	// Ignore git checks
	renameCmd.Flags().BoolVarP(&force, "force", "F", false, "Ignore git checks")

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
