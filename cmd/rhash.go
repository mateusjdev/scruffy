package cmd

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/rhash"
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
	absolutePath   bool
	dryRun         bool
	force          bool
	hash           string
	recursive      bool
	truncate       uint8
	uppercase      bool
	skipGitCheck   bool
	inputPath      string
	outputPath     string
	currentWorkDir string
)

var rhashCmd = &cobra.Command{
	Use:   "rhash [source]",
	Short: "Rename files to their hash sum",
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
			clog.Errorf("--truncate is very low, choose >= %d", ARGS_MIN_TRUNCATE)
			clog.ExitBecause(clog.ErrUserInput)
		}

		if truncate > ARGS_MAX_TRUNCATE {
			clog.Errorf("--truncate is very high, choose <= %d", ARGS_MAX_TRUNCATE)
			clog.ExitBecause(clog.ErrUserInput)
		}

		skipGitCheck = force || debug || dryRun

		// TODO: validate hash here
		hash = strings.ToLower(hash)

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
	hash: %s`, dryRun, silent, recursive, absolutePath, skipGitCheck, uppercase, truncate, inputPath, outputPath, hash)
	},
	Run: func(cmd *cobra.Command, args []string) {
		clog.Debugf("Starting module::%s", cmd.Use)

		if inputPath == "" {
			clog.Errorf("--input is empty or invalid")
			clog.ExitBecause(clog.ErrUserGeneric)
		}

		inputPathInfo, err := cfs.GetValidatedPath(inputPath)
		clog.CheckIfError(err)
		if inputPathInfo.GetPathType() == cfs.PathIsNonExistent {
			clog.Errorf("Source path %s is not a valid file or a directory\n", inputPath)
			clog.ExitBecause(clog.ErrUserInput)
		}

		if outputPath == "" {
			if cmd.Flags().Lookup("output").Changed {
				clog.Errorf("--output is empty or invalid")
				clog.ExitBecause(clog.ErrUserInput)
			}

			if inputPathInfo.GetPathType() == cfs.PathIsFile {
				outputPath = filepath.Dir(inputPathInfo.GetPath())
			} else {
				outputPath = inputPathInfo.GetPath()
			}
		}

		// TODO(11): Create destinationPath if doesn't exist (maybe add a flag? force?)
		outputPathInfo, err := cfs.GetValidatedPath(outputPath)
		clog.CheckIfError(err)
		if outputPathInfo.GetPathType() != cfs.PathIsDirectory {
			clog.Errorf("Destination folder \"%s\" is not a valid directory\n", outputPath)
			clog.ExitBecause(clog.ErrUserInput)
		}

		if cfs.IsGitRepo(inputPathInfo.GetPath()) {
			if !skipGitCheck {
				clog.Errorf("%s is in a git repo", inputPathInfo.GetPath())
				clog.ExitBecause(clog.ErrUserGeneric)
			}
			clog.Infof("%s is in a git repo", inputPathInfo.GetPath())
		}

		if inputPathInfo.GetPath() != outputPathInfo.GetPath() && cfs.IsGitRepo(outputPathInfo.GetPath()) {
			if !skipGitCheck {
				clog.Errorf("%s is in a git repo", outputPathInfo.GetPath())
				clog.ExitBecause(clog.ErrUserGeneric)
			}
			clog.Infof("%s is in a git repo", outputPathInfo.GetPath())
		}

		clog.Debugf("inputPathInfo.GetPath(): %s", inputPathInfo.GetPath())
		clog.Debugf("outputPathInfo.GetPath(): %s", outputPathInfo.GetPath())

		cwd, err := os.Getwd()
		clog.CheckIfError(err)
		currentWorkDir = cwd

		if hash == rhash.HashAlgorithmFuzzy {
			// FUZZY_MACHINE
			fuzzyMachineOptions := rhash.FuzzyMachineOptions{
				Uppercase: uppercase,
				Truncate:  truncate,
				// INFO: For file naming this (dryRun) will be random,
				// but at least it will show the destination path
				DryRun:         dryRun,
				AbsolutePath:   absolutePath,
				CurrentWorkDir: currentWorkDir,
			}

			// PATH_WALK
			rhash.EnqueuePath(fuzzyMachineOptions, recursive, inputPathInfo, outputPathInfo)
		} else {
			// HASH_MACHINE
			hashAlgorithm, err := rhash.GetHashAlgorithm(hash, int(truncate))
			clog.CheckIfError(err)
			hashMachine := rhash.HashMachine{
				Machine: hashAlgorithm,
				Options: rhash.HashMachineOptions{
					Uppercase:      uppercase,
					Truncate:       truncate,
					DryRun:         dryRun,
					AbsolutePath:   absolutePath,
					CurrentWorkDir: currentWorkDir,
				},
			}

			// PATH_WALK
			rhash.EnqueuePath(hashMachine, recursive, inputPathInfo, outputPathInfo)
		}
	},
}

func init() {
	rootCmd.AddCommand(rhashCmd)

	// TODO(1): Add https://github.com/spf13/viper for configuration
	// TODO(1a): Use XDG Base Directory Specification
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scruffy.yaml)")

	rhashCmd.Flags().StringVarP(&hash, "hash", "H", "blake3", "hash that will be used: [md5/blake3/blake2b/sha1/sha256/sha512/fuzzy]")

	// TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3)
	// TODO(2): Drop -i and use 'scruffy rhash $i $2 $3'
	rhashCmd.Flags().StringVarP(&inputPath, "input", "i", "./", "Path to DIR/FILE which will be hashed")

	// INFO: If --output/defaultOutputPath is not declared, it will be the same as --input/defaultInputPath
	rhashCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Location were hashed files will be stored")

	rhashCmd.Flags().BoolVarP(&absolutePath, "absolute-path", "A", false, "Print absolute paths relative when logging")

	rhashCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Don't rename files")

	rhashCmd.Flags().BoolVarP(&uppercase, "uppercase", "U", false, "Convert characters to UPPERCASE")

	// Ignore git checks
	rhashCmd.Flags().BoolVarP(&force, "force", "F", false, "Ignore git checks")

	// Truncate filename, max value is 256(uint8)
	rhashCmd.Flags().Uint8VarP(&truncate, "truncate", "t", 32, "Truncate filename (Beetween 8 and 128)")

	rhashCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recurse DIRs, when enabled, will not accept a target directory")

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
