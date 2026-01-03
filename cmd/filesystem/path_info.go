package filesystem

import (
	"errors"
	"os"
	"path/filepath"
)

type PathInfo struct {
	path          string
	isDir         bool
	isRegularFile bool
}

func (fileInfo *PathInfo) Path() string {
	return fileInfo.path
}

func (fileinfo *PathInfo) IsDir() bool {
	return fileinfo.isDir
}

// HACK: Verify file.Modes()
func (fileinfo *PathInfo) IsRegularFile() bool {
	return fileinfo.isRegularFile
}

func StatPath(path string) (*PathInfo, error) {
	if path == "" {
		return nil, errors.New("path is empty or invalid")
	}

	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(abspath)
	if err != nil {
		return nil, err
	}

	var pathInfo PathInfo

	pathInfo.path = abspath
	pathInfo.isDir = stat.IsDir()
	pathInfo.isRegularFile = stat.Mode().IsRegular()

	if !pathInfo.isDir && !pathInfo.isRegularFile {
		return nil, errors.New("error stating path")
	}

	return &pathInfo, nil
}
