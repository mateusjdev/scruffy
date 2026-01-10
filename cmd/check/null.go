package check

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"os"
	"path/filepath"
)

func renameExtension(path *filesystem.PathInfo, ext string) error {
	if path.IsDir() {
		return errors.New("won't rename a directory")
	}

	// clog.Errorf(path.Path())
	// clog.Errorf(filepath.Dir(path.Path()))
	// return errors.ErrUnsupported

	dir := filepath.Dir(path.Path())
	base := filepath.Base(path.Path())
	inc_n := 0
	inc := ""
	for {
		err := filesystem.SafeRename(path, filepath.Join(dir, base+inc+ext))

		if err == nil {
			break
		}

		if err == filesystem.ErrFileExists {
			inc_n++
			inc = fmt.Sprintf("_%d", inc_n)
			continue
		} else {
			return err
		}
	}

	return nil
}

func isFileEmpty(path *filesystem.PathInfo) (bool, int, error) {
	file, err := os.Open(path.Path())
	if err != nil {
		return false, 0, err
	}

	defer file.Close()

	isEmpty := true
	byteFound := 0

	for isEmpty {
		b := make([]byte, 4096)
		n, err := file.Read(b)

		if err != nil && err != io.EOF {
			return false, 0, err
		}

		for i, v := range b {
			if i > n {
				break
			}

			if v != 0 {
				byteFound += i
				isEmpty = false
				break
			}
		}

		if isEmpty {
			byteFound += n
		}

		if err == io.EOF {
			break
		}
	}

	return isEmpty, byteFound, nil
}

func CheckPathForEmptyFiles(inputPathInfo *filesystem.PathInfo, recursive, allowRename, ignoreNonEmpty bool) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.Path())

	if inputPathInfo.IsRegularFile() {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.Path())

		isEmpty, foundAt, err := isFileEmpty(inputPathInfo)
		if err != nil {
			return err
		}

		if isEmpty {
			if allowRename {
				err := renameExtension(inputPathInfo, ".empty")
				if err != nil {
					return err
				}

				clog.Warningf("RENAMED File %s is empty!", inputPathInfo.Path())
			} else {
				clog.Warningf("File %s is empty!", inputPathInfo.Path())
			}
		} else if !ignoreNonEmpty {
			clog.Warningf("File %s has content at %d!", inputPathInfo.Path(), foundAt)
		}

		return nil
	}

	if !inputPathInfo.IsDir() {
		return errors.New("input is not a valid file or directory")
	}

	return filepath.WalkDir(inputPathInfo.Path(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.Path() == path {
				return nil
			}

			if recursive {
				recursePathInfo, err := filesystem.StatPath(path)
				if err != nil {
					return err
				}

				CheckPathForEmptyFiles(
					recursePathInfo,
					recursive,
					allowRename,
					ignoreNonEmpty,
				)
			}

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := filesystem.StatPath(path)
		if err != nil {
			return err
		}

		clog.Debugf("Working on file \"%s\"", fileInfo.Path())
		isEmpty, foundAt, err := isFileEmpty(fileInfo)
		if err != nil {
			return err
		}

		if isEmpty {
			if allowRename {
				err := renameExtension(fileInfo, ".empty")
				if err != nil {
					return err
				}

				clog.Warningf("RENAMED File %s is empty!", fileInfo.Path())
			} else {
				clog.Warningf("File %s is empty!", fileInfo.Path())
			}
		} else if !ignoreNonEmpty {
			clog.Warningf("File %s has content at %d!", fileInfo.Path(), foundAt)
		}

		return nil
	})
}
