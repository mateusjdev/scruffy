package mimetype

import (
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func checkMimeType(path cfs.CustomFileInfo) {
	clog.Debugf("Working on file \"%s\"", path.GetPath())
	mtype, err := mimetype.DetectFile(path.GetPath())
	if err != nil {
		panic("ERROR")
	}
	clog.InfoSuccessf("%s %s", mtype.String(), mtype.Extension())

}

func CheckMimeType(inputPathInfo cfs.CustomFileInfo) {
	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		checkMimeType(inputPathInfo)
	}

	if inputPathInfo.GetPathType() != cfs.PathIsDirectory {
		clog.Errorf("Not a valid file or directory")
		clog.ExitBecause(clog.ErrCodeGeneric)
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
			clog.CheckIfError(err)
			CheckMimeType(recursePathInfo)

			// Skip walk(dir) from --recuse anyway, this helps ensure destination folder will be respected
			return filepath.SkipDir
		}

		fileInfo, err := cfs.GetValidatedPath(path)
		clog.CheckIfError(err)

		clog.Debugf("Working on file \"%s\"", fileInfo.GetPath())
		checkMimeType(fileInfo)
		return nil
	})

}
