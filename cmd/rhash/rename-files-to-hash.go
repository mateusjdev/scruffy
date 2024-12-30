package rhash

import (
	"errors"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/common"
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

func isGitRepo(path string) bool {
	_, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil && err == git.ErrRepositoryNotExists {
		return false
	}
	return true
}

func RenameFilesToHash(cmd *cobra.Command, args []string) {

	// move to rootArgs

	// Move to rootArgs
	// clog.Debugf("Starting %s", filepath.Base(os.Args[0]))
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		clog.SetLogLevel(clog.LevelDebug)
	}
	silent, _ := cmd.Flags().GetBool("silent")
	if silent {
		clog.SetLogLevel(clog.LevelWarning)
	}

	clog.Debugf("Starting module::%s", cmd.Use)

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	recursive, _ := cmd.Flags().GetBool("recursive")
	verbose, _ := cmd.Flags().GetBool("verbose")
	mForce, _ := cmd.Flags().GetBool("force")
	force := mForce || debug
	uppercase, _ := cmd.Flags().GetBool("uppercase")
	// TODO: Check MAX_PATH on windows
	lenght, _ := cmd.Flags().GetInt8("lenght")
	inputPath, _ := cmd.Flags().GetString("input")
	outputPath, _ := cmd.Flags().GetString("output")
	hash, _ := cmd.Flags().GetString("hash")
	// hash := cmd.Flags().Lookup("hash")

	clog.Debugf(`Args:
	dry-run: %t
	silent: %t
	recursive: %t
	verbose: %t
	force: %t
	uppercase: %t
	lenght: %d
	inputPath: %s
	outputPath: %s
	hash: %s`, dryRun, silent, recursive, verbose, force, uppercase, lenght, inputPath, outputPath, hash)

	if inputPath == "" {
		clog.Errorf("--input is empty or invalid")
		clog.ExitBecause(clog.USER_ERROR)
	}

	inputPathType, err := isValidPath(inputPath)
	common.CheckIfError(err)
	if inputPathType == PathIsNonExistent {
		clog.Errorf("Source path %s is not a valid file or a directory\n", inputPath)
		clog.ExitBecause(clog.USER_ERROR)
	}

	inputPathAbs, err := filepath.Abs(inputPath)
	common.CheckIfError(err)

	if outputPath == "" {
		if cmd.Flags().Lookup("output").Changed {
			clog.Errorf("--output is empty or invalid")
			clog.ExitBecause(clog.USER_ERROR)
		}

		if inputPathType == PathIsFile {
			outputPath = filepath.Dir(inputPathAbs)
		} else {
			outputPath = inputPathAbs
		}
	}

	// TODO: create outputPath if doesnt exist
	outputPathType, err := isValidPath(outputPath)
	common.CheckIfError(err)
	if outputPathType != PathIsDirectory {
		clog.Errorf("Destination folder \"%s\" is not a valid directory\n", outputPath)
		clog.ExitBecause(clog.USER_ERROR)
	}

	outputPathAbs, err := filepath.Abs(outputPath)
	common.CheckIfError(err)

	if isGitRepo(inputPathAbs) {
		if !force && !dryRun {
			clog.Errorf("%s is in a git repo", inputPathAbs)
			clog.ExitBecause(clog.USER_ERROR)
		}
		clog.Infof("%s is in a git repo", inputPathAbs)
	}

	if isGitRepo(outputPathAbs) {
		if !force && !dryRun && inputPathAbs != outputPathAbs {
			clog.Errorf("%s is in a git repo", outputPathAbs)
			clog.ExitBecause(clog.USER_ERROR)
		}
		clog.Infof("%s is in a git repo", outputPathAbs)
	}

	clog.Debugf("inputPathAbs: %s", inputPathAbs)
	clog.Debugf("outputPathAbs: %s", outputPathAbs)

	// HASH_MACHINE

	algorithm, err := checkHash(hash)
	if err != nil {
		clog.Errorf("%s: (%s)", err, hash)
	}

	clog.Debugf("Creating HashMachine with algorithm ID: (%d)", int(algorithm))

	hashMachine, err := getHashMachine(algorithm, int(lenght))
	common.CheckIfError(err)

	// PATH_WALK

	inputPathInfo := CustomFileInfo{inputPathAbs, inputPathType}
	outputPathInfo := CustomFileInfo{outputPathAbs, outputPathType}

	clog.Debugf("Enqueue")
	enqueuePath(inputPathInfo, outputPathInfo, hashMachine)
}
