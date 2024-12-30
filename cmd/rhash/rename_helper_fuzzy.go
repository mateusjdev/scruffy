package rhash

import (
	"errors"
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

const (
	charset string = "abcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type FuzzyMachineOptions struct {
	uppercase bool
	truncate  uint8
}

func (fuzzyMachineOptions FuzzyMachineOptions) getChecksum() string {
	b := make([]byte, fuzzyMachineOptions.truncate)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	if fuzzyMachineOptions.uppercase {
		return strings.ToUpper(string(b))
	}
	return string(b)
}

func (fuzzyMachineOptions FuzzyMachineOptions) workOnFile(sourceFileInfo cfs.CustomFileInfo, destinationDirInfo cfs.CustomFileInfo) error {
	clog.Debugf("Working on file \"%s\"", sourceFileInfo.GetPath())

	extension := filepath.Ext(sourceFileInfo.GetPath())

	// If fails to rename, just generate a new name
	for {
		fileHash := fuzzyMachineOptions.getChecksum()
		destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

		// TODO(16): Check if has permission to move to destination
		err := cfs.SafeRename(sourceFileInfo.GetPath(), destination)
		if err == nil {
			clog.InfoIconf(clog.PrintIconSuccess, "\"%s\" -> %s", sourceFileInfo.GetPath(), destination)
			return nil
		} else if errors.Is(err, cfs.ErrSameFile) || errors.Is(err, cfs.ErrFileExists) {
			continue
		}

		return err
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
// TODO(23): Reuse enqueuePath for rename_helper_fuzzy and rename_helper_hash
func (fuzzyMachineOptions FuzzyMachineOptions) enqueuePath(inputPathInfo cfs.CustomFileInfo, outputPathInfo cfs.CustomFileInfo) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.GetPath())

	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		return fuzzyMachineOptions.workOnFile(inputPathInfo, outputPathInfo)
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
		return fuzzyMachineOptions.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
