package rhash

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"

	"github.com/spf13/cobra"
)

func RenameFilesToHash(cmd *cobra.Command, args []string) {

	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		clog.SetLogLevel(clog.LevelDebug)
	}
	silent, _ := cmd.Flags().GetBool("silent")
	if silent {
		clog.SetLogLevel(clog.LevelWarning)
	}

	clog.Debugf("Starting module::%s", cmd.Use)

	// TODO(3): Work on dry-run flag
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// TODO(5): Work on verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// TODO(6): Work on uppercase flag
	uppercase, _ := cmd.Flags().GetBool("uppercase")

	// TODO(7): Work on recursive flag
	recursive, _ := cmd.Flags().GetBool("recursive")

	skipGitCheck, _ := cmd.Flags().GetBool("force")
	skipGitCheck = skipGitCheck || debug || dryRun

	truncate, _ := cmd.Flags().GetUint8("truncate")

	inputPath, _ := cmd.Flags().GetString("input")
	outputPath, _ := cmd.Flags().GetString("output")

	hash, _ := cmd.Flags().GetString("hash")

	clog.Debugf(`Args:
	dry-run: %t
	silent: %t
	recursive: %t
	verbose: %t
	skipGitCheck: %t
	uppercase: %t
	truncate: %d
	inputPath: %s
	outputPath: %s
	hash: %s`, dryRun, silent, recursive, verbose, skipGitCheck, uppercase, truncate, inputPath, outputPath, hash)

	if inputPath == "" {
		clog.Errorf("--input is empty or invalid")
		clog.ExitBecause(clog.ErrUserGeneric)
	}

	inputPathInfo, err := cfs.GetValidatedPath(inputPath)
	clog.CheckIfError(err)
	if inputPathInfo.GetPathType() == cfs.PathIsNonExistent {
		clog.Errorf("Source path %s is not a valid file or a directory\n", inputPath)
		clog.ExitBecause(clog.ErrUserGeneric)
	}

	inputPathAbs, err := filepath.Abs(inputPath)
	clog.CheckIfError(err)

	if outputPath == "" {
		if cmd.Flags().Lookup("output").Changed {
			clog.Errorf("--output is empty or invalid")
			clog.ExitBecause(clog.ErrUserGeneric)
		}

		if inputPathInfo.GetPathType() == cfs.PathIsFile {
			outputPath = filepath.Dir(inputPathAbs)
		} else {
			outputPath = inputPathAbs
		}
	}

	// TODO(11): Create destinationPath if doesn't exist (maybe add a flag? force?)
	outputPathInfo, err := cfs.GetValidatedPath(outputPath)
	clog.CheckIfError(err)
	if outputPathInfo.GetPathType() != cfs.PathIsDirectory {
		clog.Errorf("Destination folder \"%s\" is not a valid directory\n", outputPath)
		clog.ExitBecause(clog.ErrUserGeneric)
	}

	outputPathAbs, err := filepath.Abs(outputPath)
	clog.CheckIfError(err)

	if cfs.IsGitRepo(inputPathAbs) {
		if !skipGitCheck {
			clog.Errorf("%s is in a git repo", inputPathAbs)
			clog.ExitBecause(clog.ErrUserGeneric)
		}
		clog.Infof("%s is in a git repo", inputPathAbs)
	}

	if cfs.IsGitRepo(outputPathAbs) {
		if !skipGitCheck && inputPathAbs != outputPathAbs {
			clog.Errorf("%s is in a git repo", outputPathAbs)
			clog.ExitBecause(clog.ErrUserGeneric)
		}
		clog.Infof("%s is in a git repo", outputPathAbs)
	}

	clog.Debugf("inputPathAbs: %s", inputPathAbs)
	clog.Debugf("outputPathAbs: %s", outputPathAbs)

	// HASH_MACHINE

	hashAlgorithm, err := getHashAlgorithm(hash, int(truncate))
	clog.CheckIfError(err)
	hashMachine := HashMachine{
		Machine: hashAlgorithm,
		Options: HashMachineOptions{
			uppercase: uppercase,
			truncate:  truncate,
		},
	}

	// PATH_WALK

	enqueuePath(inputPathInfo, outputPathInfo, hashMachine)
}
