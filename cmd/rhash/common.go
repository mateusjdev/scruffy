package rhash

import (
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"
)

type Operation int

const (
	OperationSameFile Operation = iota
	OperationRenamed
	OperationDryRun
)

const (
	HashAlgorithmBlake2b string = "blake2b"
	HashAlgorithmBlake3  string = "blake3"
	HashAlgorithmMD5     string = "md5"
	HashAlgorithmSHA1    string = "sha1"
	HashAlgorithmSHA256  string = "sha256"
	HashAlgorithmSHA512  string = "sha512"

	HashAlgorithmFuzzy string = "fuzzy"
)

type MachineOptions struct {
	Uppercase      bool
	Truncate       uint8
	DryRun         bool
	AbsolutePath   bool
	CurrentWorkDir string
}

type RenameHelper interface {
	workOnFile(cfs.CustomFileInfo, cfs.CustomFileInfo) error
	getChecksum(cfs.CustomFileInfo) (string, error)
}

type RenameMachine interface {
	RenameHelper
}

func ReportOperation(options MachineOptions, operation Operation, source, destination cfs.CustomFileInfo) {
	var fSource, fDestination string
	// TODO: parse isSameVolume on rhash/parse.go (before)
	if options.AbsolutePath {
		fSource = source.GetPath()
		fDestination = destination.GetPath()
	} else {
		var err error
		if cfs.IsSameVolume(options.CurrentWorkDir, source.GetPath()) {
			fSource, err = filepath.Rel(options.CurrentWorkDir, source.GetPath())
			if err != nil {
				fSource = source.GetPath()
			}
		} else {
			fSource = source.GetPath()
		}

		if cfs.IsSameVolume(options.CurrentWorkDir, destination.GetPath()) {
			fDestination, err = filepath.Rel(options.CurrentWorkDir, destination.GetPath())
			if err != nil {
				fDestination = destination.GetPath()
			}
		} else {
			fDestination = destination.GetPath()
		}
	}

	switch operation {
	case OperationSameFile:
		clog.Infof("file \"%s\" already match its hash", fSource)
	case OperationRenamed:
		clog.InfoSuccessf("\"%s\" -> %s", fSource, fDestination)
	case OperationDryRun:
		clog.Infof("\"%s\" -> \"%s\"", fSource, fDestination)
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func EnqueuePath(renameMachine RenameMachine, recursive bool, inputPathInfo, outputPathInfo cfs.CustomFileInfo) error {

	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.GetPath())

	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.GetPath())
		return renameMachine.workOnFile(inputPathInfo, outputPathInfo)
	}

	if inputPathInfo.GetPathType() != cfs.PathIsDirectory {
		clog.Errorf("Not a valid file or directory")
		clog.ExitBecause(clog.ErrCodeGeneric)
	}

	// TODO(21): Check WalkDir error/return
	filepath.WalkDir(inputPathInfo.GetPath(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.GetPath() == path {
				return nil
			}

			if recursive {
				recursePathInfo, err := cfs.GetValidatedPath(path)
				clog.CheckIfError(err)

				EnqueuePath(
					renameMachine,
					recursive,
					recursePathInfo,
					recursePathInfo,
				)
			}

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := cfs.GetValidatedPath(path)
		clog.CheckIfError(err)

		clog.Debugf("Working on file \"%s\"", fileInfo.GetPath())
		return renameMachine.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
