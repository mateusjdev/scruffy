package mimetype

import (
	"errors"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func checkFile(path *filesystem.PathInfo, ignoreOk bool) error {
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

	extensionMatch := filepath.Ext(path.Path()) == mtype.Extension()
	if !extensionMatch {
		clog.Warningf("%s %s - %s", mime, extension, fullpath)
		return nil
	}

	if ignoreOk {
		return nil
	}

	clog.InfoSuccessf("%s %s - %s", mime, extension, fullpath)
	return nil
}

func CheckPath(inputPathInfo *filesystem.PathInfo, ignoreOk bool) error {
	if inputPathInfo.IsRegularFile() {
		return checkFile(inputPathInfo, ignoreOk)
	}

	if !inputPathInfo.IsDir() {
		return errors.New("Not a valid file or directory")
	}

	return filepath.WalkDir(inputPathInfo.Path(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.Path() == path {
				return nil
			}

			recursePathInfo, err := filesystem.StatPath(path)
			if err != nil {
				return err
			}
			CheckPath(recursePathInfo, ignoreOk)

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := filesystem.StatPath(path)
		if err != nil {
			return err
		}

		clog.Debugf("Working on file \"%s\"", fileInfo.Path())
		return checkFile(fileInfo, ignoreOk)
	})
}
