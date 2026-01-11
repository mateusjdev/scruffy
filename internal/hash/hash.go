package hash

import (
	"io/fs"
	"mateusjdev/scruffy/internal/clog"
	"mateusjdev/scruffy/internal/filesystem"
	"mateusjdev/scruffy/internal/hasher"
	"path/filepath"
	"strings"
)

func hashFile(
	hasher *hasher.Hasher,
	fileInfo *filesystem.PathInfo,
	truncate uint8,
	uppercase bool,
) (string, error) {
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
func HashFromPath(
	hasher *hasher.Hasher,
	recursive bool,
	inputPathInfo *filesystem.PathInfo,
	truncate uint8,
	uppercase bool,
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

			hash, err := hashFile(hasher, fileInfo, truncate, uppercase)
			if err != nil {
				return err
			}

			// Imprimir hash
			clog.Infof("%s - %s", hash, fileInfo.Path())
			return nil
		},
	)
}
