package rhash

import (
	"fmt"
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
)

type MachineOptions struct {
	Uppercase      bool
	Truncate       uint8
	DryRun         bool
	AbsolutePath   bool
	CurrentWorkDir string
}

type RenameHelper interface {
	workOnFile(*cfs.PathInfo, *cfs.PathInfo) error
	getChecksum(*cfs.PathInfo) (string, error)
}

type RenameMachine interface {
	RenameHelper
}

func ReportOperation(options MachineOptions, operation Operation, source *cfs.PathInfo, destinationPath string) {
	var fSource string
	// TODO: parse isSameVolume on rhash/parse.go (before)
	if options.AbsolutePath {
		fSource = source.Path()
	} else {
		var err error
		if cfs.IsSameVolume(options.CurrentWorkDir, source.Path()) {
			fSource, err = filepath.Rel(options.CurrentWorkDir, source.Path())
			if err != nil {
				fSource = source.Path()
			}
		} else {
			fSource = source.Path()
		}

		if cfs.IsSameVolume(options.CurrentWorkDir, destinationPath) {
			destinationPath, err = filepath.Rel(options.CurrentWorkDir, destinationPath)
			if err != nil {
				destinationPath = destinationPath
			}
		} else {
			destinationPath = destinationPath
		}
	}

	switch operation {
	case OperationSameFile:
		clog.Infof("file \"%s\" already match its hash", fSource)
	case OperationRenamed:
		clog.InfoSuccessf("\"%s\" -> %s", fSource, destinationPath)
	case OperationDryRun:
		clog.Infof("\"%s\" -> \"%s\"", fSource, destinationPath)
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func EnqueuePath(renameMachine RenameMachine, recursive bool, inputPathInfo, outputPathInfo *cfs.PathInfo) error {

	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.Path())

	if inputPathInfo.IsRegularFile() {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.Path())
		return renameMachine.workOnFile(inputPathInfo, outputPathInfo)
	}

	if !inputPathInfo.IsDir() {
		clog.PanicReturning(fmt.Errorf("Not a valid file or directory"), clog.ErrCodeGeneric)
	}

	// TODO(21): Check WalkDir error/return
	filepath.WalkDir(inputPathInfo.Path(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.Path() == path {
				return nil
			}

			if recursive {
				recursePathInfo, err := cfs.StatPath(path)
				clog.PanicIf(err)

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

		fileInfo, err := cfs.StatPath(path)
		clog.PanicIf(err)

		clog.Debugf("Working on file \"%s\"", fileInfo.Path())
		return renameMachine.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
