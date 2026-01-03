package filesystem

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrSameFile     = errors.New("args are the same file")
	ErrRenameFailed = errors.New("couldn't rename file")
	ErrFileExists   = errors.New("file already exist")
)

func IsSameVolume(path1, path2 string) bool {
	return filepath.VolumeName(path1) == filepath.VolumeName(path2)
}

// TODO(16): Check if has permission to move to destination
// ?: return filesystem.CustomFileInfo?
func SafeRename(source *PathInfo, destination string) error {
	// BUG: if source is lowercase and output is uppercase, they are reported as diferent files
	// This causes every file to be renamed as "uppercase_1.ext"
	// Check need of (os.SameFile)

	if source.Path() == destination {
		return ErrSameFile
	}

	if _, err := os.Stat(destination); err == nil {
		return ErrFileExists
	} else if errors.Is(err, os.ErrNotExist) {
		return os.Rename(source.Path(), destination)
	}

	return ErrRenameFailed
}
