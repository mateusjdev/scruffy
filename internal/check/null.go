package check

import (
	"io/fs"
	"mateusjdev/scruffy/internal/clog"
	"mateusjdev/scruffy/internal/filesystem"
	"path/filepath"
)

func CheckPathForEmptyFiles(
	inputPathInfo *filesystem.PathInfo,
	recursive bool,
	allowRename bool,
	ignoreNonEmpty bool,
) error {
	return filepath.WalkDir(
		inputPathInfo.Path(),
		func(path string, di fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if di.IsDir() {
				if recursive {
					return nil
				}

				if inputPathInfo.Path() == path {
					return nil
				}

				return filepath.SkipDir
			}

			fileInfo, err := filesystem.StatPath(path)
			if err != nil {
				return err
			}

			isEmpty, err := filesystem.IsFileEmpty(fileInfo)
			if err != nil {
				return err
			}

			if isEmpty {
				if allowRename {
					err := filesystem.RenameExtension(fileInfo, ".empty")
					if err != nil {
						return err
					}

					clog.Warningf("RENAMED File %s is empty!", fileInfo.Path())
				} else {
					clog.Warningf("File %s is empty!", fileInfo.Path())
				}
			} else if !ignoreNonEmpty {
				clog.Warningf("File %s has content!", fileInfo.Path())
			}

			return nil
		},
	)
}
