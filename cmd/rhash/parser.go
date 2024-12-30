package rhash

import (
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ARGS_MIN_TRUNCATE = 8

	HashAlgorithmBlake2b string = "blake2b"
	HashAlgorithmBlake3  string = "blake3"
	HashAlgorithmMD5     string = "md5"
	HashAlgorithmSHA1    string = "sha1"
	HashAlgorithmSHA256  string = "sha256"
	HashAlgorithmSHA512  string = "sha512"

	HashAlgorithmFuzzy string = "fuzzy"
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
	hash = strings.ToLower(hash)

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

	if truncate < ARGS_MIN_TRUNCATE {
		clog.Errorf("--truncate is very low, choose >= %d", ARGS_MIN_TRUNCATE)
		clog.ExitBecause(clog.ErrUserGeneric)
	}

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

	if outputPath == "" {
		if cmd.Flags().Lookup("output").Changed {
			clog.Errorf("--output is empty or invalid")
			clog.ExitBecause(clog.ErrUserGeneric)
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
		clog.ExitBecause(clog.ErrUserGeneric)
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

	// FUZZY_MACHINE

	if hash == HashAlgorithmFuzzy {
		fuzzyMachineOptions := FuzzyMachineOptions{
			uppercase: uppercase,
			truncate:  truncate,
			// INFO: For file naming this will be random,
			// but at least will show the destination path
			dryRun: dryRun,
		}
		fuzzyMachineOptions.enqueuePath(inputPathInfo, outputPathInfo)
	} else {
		// HASH_MACHINE

		hashAlgorithm, err := getHashAlgorithm(hash, int(truncate))
		clog.CheckIfError(err)
		hashMachine := HashMachine{
			Machine: hashAlgorithm,
			Options: HashMachineOptions{
				uppercase: uppercase,
				truncate:  truncate,
				dryRun:    dryRun,
			},
		}

		// PATH_WALK
		hashMachine.enqueuePath(inputPathInfo, outputPathInfo)
	}
}
