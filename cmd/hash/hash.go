package hash

import (
	"errors"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"mateusjdev/scruffy/cmd/hasher"
	"path/filepath"
	"strings"
)

func hashFile(hasher *hasher.Hasher, fileInfo *filesystem.PathInfo, truncate uint8, uppercase bool) (string, error) {
	hash, err := hasher.Checksum(fileInfo)
	if err != nil {
		return "", err
	}

	if truncate != 0 {
		hash = hash[0:truncate]
	}

	if uppercase {
		hash = strings.ToUpper(hash)
	}

	return hash, nil
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func HashFromPath(hasher *hasher.Hasher, recursive bool, inputPathInfo *filesystem.PathInfo, truncate uint8, uppercase bool) error {
	if inputPathInfo.IsRegularFile() {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.Path())
		hash, err := hashFile(hasher, inputPathInfo, truncate, uppercase)
		if err != nil {
			return err
		}

		// Imprimir hash
		clog.Infof("%s - %s", hash, inputPathInfo.Path())
		return nil
	}

	if !inputPathInfo.IsDir() {
		return errors.New("input is not a valid file or directory")
	}

	// TODO(21): Check WalkDir error/return
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

				HashFromPath(
					hasher,
					recursive,
					recursePathInfo,
					truncate,
					uppercase,
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
		hash, err := hashFile(hasher, fileInfo, truncate, uppercase)
		if err != nil {
			return err
		}

		// Imprimir hash
		clog.Infof("%s - %s", hash, fileInfo.Path())
		return nil
	})
}
