package mimetype

import (
	"fmt"
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func checkFile(path *cfs.PathInfo, ignoreOk bool) error {
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

func CheckPath(inputPathInfo *cfs.PathInfo, ignoreOk bool) {
	if inputPathInfo.IsRegularFile() {
		err := checkFile(inputPathInfo, ignoreOk)
		clog.PanicIf(err)
	}

	if !inputPathInfo.IsDir() {
		clog.PanicReturning(fmt.Errorf("Not a valid file or directory"), clog.ErrCodeGeneric)
	}

	filepath.WalkDir(inputPathInfo.Path(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.Path() == path {
				return nil
			}

			recursePathInfo, err := cfs.StatPath(path)
			clog.PanicIf(err)
			CheckPath(recursePathInfo, ignoreOk)

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := cfs.StatPath(path)
		clog.PanicIf(err)

		clog.Debugf("Working on file \"%s\"", fileInfo.Path())
		err = checkFile(fileInfo, ignoreOk)
		clog.PanicIf(err)
		return nil
	})

}
