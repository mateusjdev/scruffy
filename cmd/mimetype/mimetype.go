package mimetype

import (
	"errors"
	"fmt"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func renameExtension(path *filesystem.PathInfo, ext string) error {
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
			err := renameExtension(path, extension)
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

func CheckPath(pathInfo *filesystem.PathInfo, ignoreOk, allowRename bool) error {
	if pathInfo.IsRegularFile() {
		return checkFile(pathInfo, ignoreOk, allowRename)
	}

	if !pathInfo.IsDir() {
		return errors.New("Not a valid file or directory")
	}

	return filepath.WalkDir(pathInfo.Path(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if pathInfo.Path() == path {
				return nil
			}

			// Recursive call of directories
			// Ignores WalkDir
			recursePathInfo, err := filesystem.StatPath(path)
			if err != nil {
				return err
			}

			CheckPath(recursePathInfo, ignoreOk, allowRename)
			return filepath.SkipDir
		}

		fileInfo, err := filesystem.StatPath(path)
		if err != nil {
			return err
		}

		clog.Debugf("Working on file \"%s\"", fileInfo.Path())
		return checkFile(fileInfo, ignoreOk, allowRename)
	})
}
