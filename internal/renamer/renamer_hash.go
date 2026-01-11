package renamer

import (
	"errors"
	"fmt"
	"mateusjdev/scruffy/internal/filesystem"
	"mateusjdev/scruffy/internal/hasher"
	"path/filepath"
	"strings"
)

type hashRenamerOptions RenamerOptions

type HashMachine struct {
	Hasher  *hasher.Hasher
	Options hashRenamerOptions
}

func NewHashRenamer(hasher *hasher.Hasher, uppercase bool, truncate uint8, dryRun, displayAbsolutePath bool, currentWorkDir string) (Renamer, error) {
	return HashMachine{
		Hasher: hasher,
		Options: hashRenamerOptions{
			Uppercase:      uppercase,
			Truncate:       truncate,
			DryRun:         dryRun,
			DisplayAbsPath: displayAbsolutePath,
			CurrentWorkDir: currentWorkDir,
		},
	}, nil
}

func (hashMachine HashMachine) createFileName(fileInfo *filesystem.PathInfo) (string, error) {
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

func (hashMachine HashMachine) renameFile(sourceFileInfo, destinationDirInfo *filesystem.PathInfo) error {
	fileHash, err := hashMachine.createFileName(sourceFileInfo)
	if err != nil {
		return err
	}

	extension := filepath.Ext(sourceFileInfo.Path())
	destination := filepath.Join(destinationDirInfo.Path(), fileHash+extension)

	if hashMachine.Options.DryRun {
		if sourceFileInfo.Path() == destination {
			ReportOperation(
				RenamerOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destination,
			)
		} else {
			ReportOperation(
				RenamerOptions(hashMachine.Options),
				OperationDryRun,
				sourceFileInfo,
				destination,
			)
		}
		return nil
	}

	// TODO(16): Check if has permission to move to destination
	err = filesystem.SafeRename(sourceFileInfo, destination)
	if err == nil {
		ReportOperation(
			RenamerOptions(hashMachine.Options),
			OperationRenamed,
			sourceFileInfo,
			destination,
		)
		return nil
	} else if errors.Is(err, filesystem.ErrSameFile) {
		ReportOperation(
			RenamerOptions(hashMachine.Options),
			OperationSameFile,
			sourceFileInfo,
			destination,
		)
		return nil
	} else if !errors.Is(err, filesystem.ErrFileExists) {
		return err
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.Path(), newFileName)

		err = filesystem.SafeRename(sourceFileInfo, destination)
		if err == nil {
			ReportOperation(
				RenamerOptions(hashMachine.Options),
				OperationRenamed,
				sourceFileInfo,
				destination,
			)
			return nil
		} else if errors.Is(err, filesystem.ErrSameFile) {
			ReportOperation(
				RenamerOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destination,
			)
			return nil
		} else if !errors.Is(err, filesystem.ErrFileExists) {
			return err
		}

		counter++
	}
}
