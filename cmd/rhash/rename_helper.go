package rhash

import (
	"errors"
	"fmt"
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"
)

func workOnFile(sourceFileInfo cfs.CustomFileInfo, destinationDirInfo cfs.CustomFileInfo, hashMachine HashMachine) error {
	clog.Debugf("Working on file \"%s\"", sourceFileInfo.GetPath())

	// Get new filename (Hash)
	// TODO(15): Add --hash fuzzy
	fileHash, err := hashMachine.getChecksum(sourceFileInfo)
	clog.CheckIfError(err)

	extension := filepath.Ext(sourceFileInfo.GetPath())
	destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

	// TODO(16): Check if has permission to move to destination
	err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
	if errors.Is(err, cfs.ErrSameFile) {
		clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", sourceFileInfo.GetPath())
		return nil
	} else if !errors.Is(err, cfs.ErrFileExists) {
		return err
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.GetPath(), newFileName)

		err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
		if errors.Is(err, cfs.ErrSameFile) {
			clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", sourceFileInfo.GetPath())
			return nil
		} else if !errors.Is(err, cfs.ErrFileExists) {
			return err
		}

		counter++
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func enqueuePath(inputPathInfo cfs.CustomFileInfo, outputPathInfo cfs.CustomFileInfo, hashMachine HashMachine) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.GetPath())

	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		return workOnFile(inputPathInfo, outputPathInfo, hashMachine)
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
		return workOnFile(fileInfo, outputPathInfo, hashMachine)
	})

	return nil
}
