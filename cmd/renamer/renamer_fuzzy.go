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

type FuzzyMachineOptions MachineOptions

func (fuzzyMachineOptions FuzzyMachineOptions) createFileName(_ *cfs.PathInfo) (string, error) {
	b := make([]byte, fuzzyMachineOptions.Truncate)
	for i := range b {
		b[i] = charset[seed.Intn(charsetLen)]
	}
	if fuzzyMachineOptions.Uppercase {
		return strings.ToUpper(string(b)), nil
	}
	return string(b), nil
}

func (fuzzyMachineOptions FuzzyMachineOptions) renameFile(sourceFileInfo, destinationDirInfo *cfs.PathInfo) error {
	extension := filepath.Ext(sourceFileInfo.Path())

	// If fails to rename, just generate a new name

	if fuzzyMachineOptions.DryRun {
		fileHash, _ := fuzzyMachineOptions.createFileName(nil)
		destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

		ReportOperation(
			MachineOptions(fuzzyMachineOptions),
			OperationDryRun,
			sourceFileInfo,
			destination,
		)

		return nil
	}

	for {
		fileHash, _ := fuzzyMachineOptions.createFileName(sourceFileInfo)
		destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

		// TODO(16): Check if has permission to move to destination
		err := cfs.SafeRename(sourceFileInfo, destination)
		if err == nil {
			ReportOperation(
				MachineOptions(fuzzyMachineOptions),
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
