package rhash

import (
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"
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

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func EnqueuePath(renameMachine RenameMachine, inputPathInfo cfs.CustomFileInfo, outputPathInfo cfs.CustomFileInfo) error {

	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.GetPath())

	if inputPathInfo.GetPathType() == cfs.PathIsFile {
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
		return renameMachine.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
