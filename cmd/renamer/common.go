package renamer

import (
	"errors"
	"io/fs"
	"mateusjdev/scruffy/cmd/clog"
	"mateusjdev/scruffy/cmd/filesystem"
	"path/filepath"
)

type Operation int

const (
	OperationSameFile Operation = iota
	OperationRenamed
	OperationDryRun
)

type RenamerOptions struct {
	Truncate       uint8
	Uppercase      bool
	DryRun         bool
	DisplayAbsPath bool
	CurrentWorkDir string
}

type RenameMethod interface {
	renameFile(*filesystem.PathInfo, *filesystem.PathInfo) error
	createFileName(*filesystem.PathInfo) (string, error)
}

type Renamer interface {
	RenameMethod
}

func ReportOperation(options RenamerOptions, operation Operation, source *filesystem.PathInfo, destinationPath string) {
	var fSource string
	// TODO: parse isSameVolume on rhash/parse.go (before)
	if options.DisplayAbsPath {
		fSource = source.Path()
	} else {
		var err error
		if filesystem.IsSameVolume(options.CurrentWorkDir, source.Path()) {
			fSource, err = filepath.Rel(options.CurrentWorkDir, source.Path())
			if err != nil {
				fSource = source.Path()
			}
		} else {
			fSource = source.Path()
		}

		if filesystem.IsSameVolume(options.CurrentWorkDir, destinationPath) {
			relativePath, err := filepath.Rel(options.CurrentWorkDir, destinationPath)
			if err != nil {
				destinationPath = relativePath
			}
		}
	}

	switch operation {
	case OperationSameFile:
		clog.Infof("file \"%s\" already match its hash", fSource)
	case OperationRenamed:
		clog.InfoSuccessf("\"%s\" -> %s", fSource, destinationPath)
	case OperationDryRun:
		clog.Infof("\"%s\" -> \"%s\"", fSource, destinationPath)
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func RenameFromPath(renamer Renamer, recursive bool, inputPathInfo, outputPathInfo *filesystem.PathInfo) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.Path())

	if !outputPathInfo.IsDir() {
		return errors.New("output is not a valid file or directory")
	}

	if inputPathInfo.IsRegularFile() {
		clog.Debugf("Working on file \"%s\"", inputPathInfo.Path())
		return renamer.renameFile(inputPathInfo, outputPathInfo)
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

				RenameFromPath(
					renamer,
					recursive,
					recursePathInfo,
					recursePathInfo,
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
		return renamer.renameFile(fileInfo, outputPathInfo)
	})
}
