package filesystem

import (
	"errors"
	"fmt"
	"io"
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

// TODO: Remove old extension
func RenameExtension(path *PathInfo, ext string) error {
	if path.IsDir() {
		return errors.New("won't rename a directory")
	}

	dir := filepath.Dir(path.Path())
	base := filepath.Base(path.Path())
	inc_n := 0
	inc := ""
	for {
		err := SafeRename(path, filepath.Join(dir, base+inc+ext))

		if err == nil {
			break
		}

		if err == ErrFileExists {
			inc_n++
			inc = fmt.Sprintf("_%d", inc_n)
			continue
		} else {
			return err
		}
	}

	return nil
}

func IsFileEmpty(path *PathInfo) (bool, error) {
	file, err := os.Open(path.Path())
	if err != nil {
		return false, err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	if stat.Size() == 0 {
		return true, nil
	}

	var isEmpty bool = true
	for isEmpty {
		b := make([]byte, 4096)
		n, err := file.Read(b)

		if err == io.EOF {
			break
		}

		for i, v := range b {
			if i >= n {
				break
			}

			if v != 0 {
				isEmpty = false
				break
			}
		}
	}

	return isEmpty, nil
}
