package cfs

import (
	"errors"
	"fmt"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type PathType uint8

const (
	PathIsUnknown PathType = iota
	PathIsNonExistent
	PathIsFile
	PathIsDirectory
)

type CustomFileInfo struct {
	path     string
	pathType PathType
}

var (
	ErrSameFile     = errors.New("args are the same file")
	ErrRenameFailed = errors.New("couldn't rename file")
	ErrFileExists   = errors.New("file already exist")
)

func IsSameVolume(path1, path2 string) bool {
	return filepath.VolumeName(path1) == filepath.VolumeName(path2)
}

// TODO(16): Check if has permission to move to destination
// ?: return cfs.CustomFileInfo?
func SafeRename(source string, destination string) error {
	// BUG: if source is lowercase and output is uppercase, they are reported as diferent files
	// This causes every file to be renamed as "uppercase_1.ext"
	// Check need of (os.SameFile)

	if source == destination {
		return ErrSameFile
	}

	if _, err := os.Stat(destination); err == nil {
		return ErrFileExists
	} else if errors.Is(err, os.ErrNotExist) {
		return os.Rename(source, destination)
	}

	return ErrRenameFailed
}

func GetValidatedPath(path string) (CustomFileInfo, error) {
	path, err := filepath.Abs(path)
	clog.PanicIf(err)
	stat, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return GetUnvalidatedPath(path, PathIsNonExistent), nil
		}
		return GetUnvalidatedPath(path, PathIsUnknown), err
	}
	if stat.IsDir() {
		return GetUnvalidatedPath(path, PathIsDirectory), nil
	}
	return GetUnvalidatedPath(path, PathIsFile), nil
}

func GetUnvalidatedPath(path string, pathType PathType) CustomFileInfo {
	return CustomFileInfo{
		path:     path,
		pathType: pathType,
	}
}

func (fileInfo CustomFileInfo) GetPath() string {
	return fileInfo.path
}

func (fileInfo CustomFileInfo) GetPathType() PathType {
	return fileInfo.pathType
}

// Using go-git because it doesn't require git binary
func IsGitRepo(path string) bool {
	_, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil && err == git.ErrRepositoryNotExists {
		return false
	}
	return true
}

func ValidatePath(path string, skipGitCheck bool, validPathType ...PathType) (*CustomFileInfo, error) {
	if path == "" {
		clog.PanicReturning(fmt.Errorf("path is empty or invalid"), clog.ErrUserGeneric)
	}

	if len(validPathType) == 0 {
		clog.PanicReturning(fmt.Errorf("Couldn't validate path"), clog.ErrCodeGeneric)
	}

	tmpPath, err := GetValidatedPath(path)
	clog.PanicIf(err)

	isValid := false
	for _, pathType := range validPathType {
		if tmpPath.GetPathType() == pathType {
			isValid = true
		}
	}
	if !isValid {
		clog.PanicReturning(fmt.Errorf("Path %s is not a valid file or a directory\n", tmpPath), clog.ErrUserInput)
		return nil, nil
	}

	if IsGitRepo(tmpPath.GetPath()) {
		if !skipGitCheck {
			clog.PanicReturning(fmt.Errorf("%s is in a git repo", tmpPath.GetPath()), clog.ErrUserGeneric)
		}
		clog.Infof("%s is in a git repo", tmpPath.GetPath())
	}

	return &tmpPath, nil
}
