package rhash

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
	seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type FuzzyMachineOptions MachineOptions

func (fuzzyMachineOptions FuzzyMachineOptions) getChecksum(_ cfs.CustomFileInfo) (string, error) {
	b := make([]byte, fuzzyMachineOptions.truncate)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	if fuzzyMachineOptions.uppercase {
		return strings.ToUpper(string(b)), nil
	}
	return string(b), nil
}

func (fuzzyMachineOptions FuzzyMachineOptions) workOnFile(sourceFileInfo cfs.CustomFileInfo, destinationDirInfo cfs.CustomFileInfo) error {
	extension := filepath.Ext(sourceFileInfo.GetPath())

	// If fails to rename, just generate a new name

	if fuzzyMachineOptions.dryRun {
		fileHash, _ := fuzzyMachineOptions.getChecksum(sourceFileInfo)
		destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

		destinationFileInfo := cfs.GetUnvalidatedPath(destination, cfs.PathIsFile)

		ReportOperation(
			MachineOptions(fuzzyMachineOptions),
			OperationDryRun,
			sourceFileInfo,
			destinationFileInfo,
		)

		return nil
	}

	for {
		fileHash, _ := fuzzyMachineOptions.getChecksum(sourceFileInfo)
		destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

		// TODO(16): Check if has permission to move to destination
		err := cfs.SafeRename(sourceFileInfo.GetPath(), destination)
		if err == nil {
			ReportOperation(
				MachineOptions(fuzzyMachineOptions),
				OperationRenamed,
				sourceFileInfo,
				cfs.GetUnvalidatedPath(destination, cfs.PathIsFile),
			)
			return nil
		} else if errors.Is(err, cfs.ErrSameFile) || errors.Is(err, cfs.ErrFileExists) {
			continue
		}

		return err
	}
}
