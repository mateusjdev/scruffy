package rhash

import (
	"errors"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type PathType uint8

const (
	PathIsUnknown PathType = iota
	PathIsNonExistent
	PathIsFile
	PathIsDirectory
)

func isValidPath(path string) (PathType, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return PathIsNonExistent, nil
		}
		return PathIsUnknown, err
	}
	if stat.IsDir() {
		return PathIsDirectory, nil
	}
	return PathIsFile, nil
}

// Using go-git because it doesn't require git binary
func isGitRepo(path string) bool {
	_, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil && err == git.ErrRepositoryNotExists {
		return false
	}
	return true
}

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

	inputPathType, err := isValidPath(inputPath)
	clog.CheckIfError(err)
	if inputPathType == PathIsNonExistent {
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

		if inputPathType == PathIsFile {
			outputPath = filepath.Dir(inputPathAbs)
		} else {
			outputPath = inputPathAbs
		}
	}

	// TODO(11): Create destinationPath if doesn't exist (maybe add a flag? force?)
	outputPathType, err := isValidPath(outputPath)
	clog.CheckIfError(err)
	if outputPathType != PathIsDirectory {
		clog.Errorf("Destination folder \"%s\" is not a valid directory\n", outputPath)
		clog.ExitBecause(clog.ErrUserGeneric)
	}

	outputPathAbs, err := filepath.Abs(outputPath)
	clog.CheckIfError(err)

	if isGitRepo(inputPathAbs) {
		if !skipGitCheck {
			clog.Errorf("%s is in a git repo", inputPathAbs)
			clog.ExitBecause(clog.ErrUserGeneric)
		}
		clog.Infof("%s is in a git repo", inputPathAbs)
	}

	if isGitRepo(outputPathAbs) {
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

	inputPathInfo := CustomFileInfo{inputPathAbs, inputPathType}
	outputPathInfo := CustomFileInfo{outputPathAbs, outputPathType}

	enqueuePath(inputPathInfo, outputPathInfo, hashMachine)
}
