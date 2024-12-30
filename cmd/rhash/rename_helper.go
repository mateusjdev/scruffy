package rhash

import (
	"errors"
	"fmt"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"
)

type CustomFileInfo struct {
	Path string
	PathType
}

type CustomFutureFileInfo struct {
	DestinationDir string
	Filename       string
	Extension      string
}

// TODO(16): Check if has permission to move to destination
func checkIfExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return true, nil
}

func safeRename(source string, destination string) bool {
	// BUG: if source is lowercase and output is uppercase, they are reported as diferent files
	// This causes every file to be renamed as "uppercase_1.ext"
	if source == destination {
		clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", source)
		return true
	}

	exist, err := checkIfExists(destination)
	clog.CheckIfError(err)
	if !exist {
		os.Rename(source, destination)
		clog.InfoIconf(clog.PrintIconSuccess, "\"%s\" -> %s", source, destination)
		return true
	}

	return false
}

func workOnFile(sourceFileInfo CustomFileInfo, destinationDirInfo CustomFileInfo, hashMachine HashMachine) error {
	clog.Debugf("Working on file \"%s\"", sourceFileInfo.Path)

	// Get new filename (Hash)
	// TODO(15): Add --hash fuzzy
	fileHash, err := hashMachine.getChecksum(sourceFileInfo)
	clog.CheckIfError(err)

	extension := filepath.Ext(sourceFileInfo.Path)
	destination := filepath.Join(destinationDirInfo.Path, fileHash+extension)

	// TODO(16): Check if has permission to move to destination
	if safeRename(sourceFileInfo.Path, destination) {
		return nil
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.Path, newFileName)
		if safeRename(sourceFileInfo.Path, destination) {
			return nil
		}

		counter++
	}
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	clog.CheckIfError(err)
	return stat.IsDir()
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func enqueuePath(inputPathInfo CustomFileInfo, outputPathInfo CustomFileInfo, hashMachine HashMachine) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.Path)

	if inputPathInfo.PathType == PathIsFile {
		return workOnFile(inputPathInfo, outputPathInfo, hashMachine)
	}

	if inputPathInfo.PathType != PathIsDirectory {
		clog.Errorf("Not a valid file or directory")
		clog.ExitBecause(clog.ErrCodeGeneric)
	}

	filepath.WalkDir(inputPathInfo.Path, func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// TODO(7): Work on recursive flag
		// Recurse into directories
		if di.IsDir() {
			if inputPathInfo.Path == path {
				return nil
			}
			return filepath.SkipDir
		}

		return workOnFile(CustomFileInfo{path, PathIsFile}, outputPathInfo, hashMachine)
	})

	return nil
}
