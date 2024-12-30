package rhash

import (
	"errors"
	"fmt"
	"hash"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/common"
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

// TODO: Check for destination permissions, if not granted, this will report ok for creating a new file
func checkIfExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return true, nil
}

func safeRename(source string, destination string) bool {
	if source == destination {
		clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", source)
		return true
	}

	exist, _ := checkIfExists(destination)
	if !exist {
		os.Rename(source, destination)
		clog.InfoIconf(clog.PrintIconSuccess, "\"%s\" -> %s", source, destination)
		return true
	}

	return false
}

func workOnFile(sourceFileInfo CustomFileInfo, destinationDirInfo CustomFileInfo, hashMachine hash.Hash) error {
	clog.Debugf("Working on file \"%s\"", sourceFileInfo.Path)
	fileHash, err := hashFile(sourceFileInfo, hashMachine)
	common.CheckIfError(err)
	extension := filepath.Ext(sourceFileInfo.Path)
	destination := filepath.Join(destinationDirInfo.Path, fileHash+extension)

	// Check if has permission to move
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
	common.CheckIfError(err)
	return stat.IsDir()
}

// TODO: Verificar a necessidade de validar os caminhos
// func enqueuePath(inputPathAbs string, outputPathAbs string, hashMachine hash.Hash) error {
func enqueuePath(inputPathInfo CustomFileInfo, outputPathInfo CustomFileInfo, hashMachine hash.Hash) error {

	if inputPathInfo.PathType == PathIsFile {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.Path)
		return workOnFile(inputPathInfo, outputPathInfo, hashMachine)
	}

	filepath.WalkDir(inputPathInfo.Path, func(path string, di fs.DirEntry, err error) error {
		clog.Debugf("Walking on %s", path)
		if err != nil {
			return err
		}

		// TODO: recurse into directories
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
