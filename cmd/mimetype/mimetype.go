package mimetype

import (
	"fmt"
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func checkFile(path cfs.CustomFileInfo, ignoreOk bool) error {
	clog.Debugf("Working on file \"%s\"", path.GetPath())
	mtype, err := mimetype.DetectFile(path.GetPath())
	if err != nil {
		return err
	}

	mime := mtype.String()
	extension := mtype.Extension()
	fullpath := path.GetPath()
	if extension == "" {
		clog.Warningf("%s Unknown - %s", mime, fullpath)
		return nil
	}

	extensionMatch := filepath.Ext(path.GetPath()) == mtype.Extension()
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

func CheckPath(inputPathInfo cfs.CustomFileInfo, ignoreOk bool) {
	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		err := checkFile(inputPathInfo, ignoreOk)
		clog.PanicIf(err)
	}

	if inputPathInfo.GetPathType() != cfs.PathIsDirectory {
		clog.PanicReturning(fmt.Errorf("Not a valid file or directory"), clog.ErrCodeGeneric)
	}

	filepath.WalkDir(inputPathInfo.GetPath(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if di.IsDir() {
			if inputPathInfo.GetPath() == path {
				return nil
			}

			recursePathInfo, err := cfs.GetValidatedPath(path)
			clog.PanicIf(err)
			CheckPath(recursePathInfo, ignoreOk)

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := cfs.GetValidatedPath(path)
		clog.PanicIf(err)

		clog.Debugf("Working on file \"%s\"", fileInfo.GetPath())
		err = checkFile(fileInfo, ignoreOk)
		clog.PanicIf(err)
		return nil
	})

}
