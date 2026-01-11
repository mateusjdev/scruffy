package mimetype

import (
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func checkFile(path *filesystem.PathInfo, ignoreOk, allowRename bool) error {
	clog.Debugf("Working on file \"%s\"", path.Path())
	mtype, err := mimetype.DetectFile(path.Path())
	if err != nil {
		return err
	}

	mime := mtype.String()
	extension := mtype.Extension()
	fullpath := path.Path()
	if extension == "" {
		clog.Warningf("%s Unknown - %s", mime, fullpath)
		return nil
	}

	extensionMatch := filepath.Ext(path.Path()) == extension
	if !extensionMatch {
		if allowRename {
			err := filesystem.RenameExtension(path, extension)
			if err != nil {
				return err
			}

			clog.Warningf("RENAMED %s %s - %s", mime, extension, fullpath)
			return nil
		}

		clog.Infof("%s %s - %s", mime, extension, fullpath)
		return nil
	}

	if ignoreOk {
		return nil
	}

	clog.InfoSuccessf("%s %s - %s", mime, extension, fullpath)
	return nil
}

func CheckPath(
	pathInfo *filesystem.PathInfo,
	ignoreExtMatch bool,
	allowRename bool,
	recursive bool,
) error {
	return filepath.WalkDir(
		pathInfo.Path(),
		func(path string, di fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if di.IsDir() {
				if recursive {
					return nil
				}

				if pathInfo.Path() == path {
					return nil
				}

				return filepath.SkipDir
			}

			fileInfo, err := filesystem.StatPath(path)
			if err != nil {
				return err
			}

			return checkFile(fileInfo, ignoreExtMatch, allowRename)
		},
	)
}
