package rhash

import (
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"
	"strings"
)

type Operation int

const (
	OperationSameFile Operation = iota
	OperationRenamed
	OperationDryRun
)

type MachineOptions struct {
	uppercase bool
	truncate  uint8
	dryRun    bool
	verbose   bool
}

type RenameHelper interface {
	workOnFile(cfs.CustomFileInfo, cfs.CustomFileInfo) error
	getChecksum(cfs.CustomFileInfo) (string, error)
}

type RenameMachine interface {
	RenameHelper
}

func StripCommonPrefix(path1, path2 string) (string, string) {
	commonPrefix := filepath.Dir(path1)
	for {
		if strings.HasPrefix(path2, commonPrefix) {
			break
		}
		commonPrefix = filepath.Dir(commonPrefix)
	}
	commonPrefix = commonPrefix + string(os.PathSeparator)
	return strings.TrimPrefix(path1, commonPrefix), strings.TrimPrefix(path2, commonPrefix)
}

func ReportOperation(options MachineOptions, operation Operation, source cfs.CustomFileInfo, destination cfs.CustomFileInfo) {
	var fSource string
	var fDestination string
	if options.verbose {
		fSource = source.GetPath()
		fDestination = destination.GetPath()
	} else {
		fSource, fDestination = StripCommonPrefix(
			source.GetPath(),
			destination.GetPath(),
		)
	}

	switch operation {
	case OperationSameFile:
		clog.Infof("file \"%s\" already hashed", fSource)
	case OperationRenamed:
		clog.InfoSuccessf("\"%s\" -> %s", fSource, fDestination)
	case OperationDryRun:
		clog.Infof("\"%s\" -> %s", fSource, fDestination)
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func EnqueuePath(renameMachine RenameMachine, inputPathInfo cfs.CustomFileInfo, outputPathInfo cfs.CustomFileInfo) error {

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

		// TODO(7): Work on recursive flag
		// Recurse into directories

		if di.IsDir() {
			if inputPathInfo.GetPath() == path {
				return nil
			}
			return filepath.SkipDir
		}

		fileInfo := cfs.GetUnvalidatedPath(path, cfs.PathIsFile)

		clog.Debugf("Working on file \"%s\"", fileInfo.GetPath())
		return renameMachine.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
