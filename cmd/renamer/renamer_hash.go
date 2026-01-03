package renamer

import (
	"errors"
	"fmt"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/hasher"
	"path/filepath"
	"strings"
)

type HashMachineOptions MachineOptions

type HashMachine struct {
	Hasher  *hasher.Hasher
	Options HashMachineOptions
}

func (hashMachine HashMachine) createFileName(fileInfo *cfs.PathInfo) (string, error) {
	hash, err := hashMachine.Hasher.Checksum(fileInfo)
	if err != nil {
		return "", err
	}

	if hashMachine.Options.Truncate != 0 {
		hash = hash[0:hashMachine.Options.Truncate]
	}

	if hashMachine.Options.Uppercase {
		hash = strings.ToUpper(hash)
	}

	return hash, nil
}

func (hashMachine HashMachine) renameFile(sourceFileInfo, destinationDirInfo *cfs.PathInfo) error {
	fileHash, err := hashMachine.createFileName(sourceFileInfo)
	if err != nil {
		return err
	}

	extension := filepath.Ext(sourceFileInfo.Path())
	destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

	if hashMachine.Options.DryRun {
		if sourceFileInfo.Path() == destination {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destination,
			)
		} else {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationDryRun,
				sourceFileInfo,
				destination,
			)
		}
		return nil
	}

	// TODO(16): Check if has permission to move to destination
	err = cfs.SafeRename(sourceFileInfo, destination)
	if err == nil {
		ReportOperation(
			MachineOptions(hashMachine.Options),
			OperationRenamed,
			sourceFileInfo,
			destination,
		)
		return nil
	} else if errors.Is(err, cfs.ErrSameFile) {
		ReportOperation(
			MachineOptions(hashMachine.Options),
			OperationSameFile,
			sourceFileInfo,
			destination,
		)
		return nil
	} else if !errors.Is(err, cfs.ErrFileExists) {
		return err
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.Path(), newFileName)

		err = cfs.SafeRename(sourceFileInfo, destination)
		if err == nil {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationRenamed,
				sourceFileInfo,
				destination,
			)
			return nil
		} else if errors.Is(err, cfs.ErrSameFile) {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destination,
			)
			return nil
		} else if !errors.Is(err, cfs.ErrFileExists) {
			return err
		}

		counter++
	}
}
