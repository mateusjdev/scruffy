package renamer

import (
	"errors"
	"mateusjdev/scruffy/cmd/cfs"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

const (
	charset string = "abcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	seed       *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	charsetLen int        = len(charset)
)

type fuzzyRenamerOptions RenamerOptions

func NewFuzzyRenamer(uppercase bool, truncate uint8, dryRun, displayAbsolutePath bool, currentWorkDir string) (Renamer, error) {
	return fuzzyRenamerOptions{
		Uppercase:      uppercase,
		Truncate:       truncate,
		DryRun:         dryRun,
		DisplayAbsPath: displayAbsolutePath,
		CurrentWorkDir: currentWorkDir,
	}, nil
}

func (renameOptions fuzzyRenamerOptions) createFileName(_ *cfs.PathInfo) (string, error) {
	b := make([]byte, renameOptions.Truncate)
	for i := range b {
		b[i] = charset[seed.Intn(charsetLen)]
	}
	if renameOptions.Uppercase {
		return strings.ToUpper(string(b)), nil
	}
	return string(b), nil
}

func (renameOptions fuzzyRenamerOptions) renameFile(sourceFileInfo, destinationDirInfo *cfs.PathInfo) error {
	extension := filepath.Ext(sourceFileInfo.Path())

	// If fails to rename, just generate a new name

	if renameOptions.DryRun {
		fileHash, _ := renameOptions.createFileName(nil)
		destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

		ReportOperation(
			RenamerOptions(renameOptions),
			OperationDryRun,
			sourceFileInfo,
			destination,
		)

		return nil
	}

	for {
		fileHash, _ := renameOptions.createFileName(sourceFileInfo)
		destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

		// TODO(16): Check if has permission to move to destination
		err := cfs.SafeRename(sourceFileInfo, destination)
		if err == nil {
			ReportOperation(
				RenamerOptions(renameOptions),
				OperationRenamed,
				sourceFileInfo,
				destination,
			)
			return nil
		} else if errors.Is(err, cfs.ErrSameFile) || errors.Is(err, cfs.ErrFileExists) {
			continue
		}

		return err
	}
}
