package renamer

import (
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
func RenameFromPath(
	renamer Renamer,
	recursive bool,
	inputPathInfo *filesystem.PathInfo,
	outputPathInfo *filesystem.PathInfo,
) error {
	return filepath.WalkDir(
		inputPathInfo.Path(),
		func(path string, di fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if di.IsDir() {
				if inputPathInfo.Path() == path {
					return nil
				}

				// if --recursive, reuse walking directory instead of walking through files
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

				return filepath.SkipDir
			}

			fileInfo, err := filesystem.StatPath(path)
			if err != nil {
				return err
			}

			clog.Debugf("Working on file \"%s\"", fileInfo.Path())
			return renamer.renameFile(fileInfo, outputPathInfo)
		},
	)
}
