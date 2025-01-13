package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"
	"strings"

	"lukechampine.com/blake3"
)

type HashMachineOptions MachineOptions

type HashMachine struct {
	Machine hash.Hash
	Options HashMachineOptions
}

var (
	ErrUnknownHashMethod = errors.New("hash method not valid")
)

// TODO(8): Work on lenght/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --lenght 224"
func GetHashAlgorithm(hash string, lenght int) (hash.Hash, error) {
	switch hash {
	case HashAlgorithmBlake2b:
		// lenght: fixed_256_bits (256, 384, 512)
		return crypto.BLAKE2b_256.New(), nil
	case HashAlgorithmBlake3:
		return blake3.New(lenght/2, nil), nil
	case HashAlgorithmMD5:
		// lenght: fixed_128_bits
		return crypto.MD5.New(), nil
	case HashAlgorithmSHA1:
		// lenght: fixed_160_bits
		return crypto.SHA1.New(), nil
	case HashAlgorithmSHA256:
		// lenght: fixed_256_bits
		return crypto.SHA256.New(), nil
	case HashAlgorithmSHA512:
		// lenght: fixed_512_bits
		return crypto.SHA512.New(), nil
	}
	return nil, ErrUnknownHashMethod
}

func (hashMachine HashMachine) getChecksum(fileInfo cfs.CustomFileInfo) (string, error) {
	if fileInfo.GetPathType() != cfs.PathIsFile {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.GetPath())
	clog.CheckIfError(err)
	defer file.Close()
	if _, err := io.Copy(hashMachine.Machine, file); err != nil {
		return "", err
	}
	hashInBytes := hashMachine.Machine.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	hashMachine.Machine.Reset()

	if hashMachine.Options.Truncate != 0 {
		hashString = hashString[0:hashMachine.Options.Truncate]
	}

	if hashMachine.Options.Uppercase {
		hashString = strings.ToUpper(hashString)
	}

	return hashString, nil
}

func (hashMachine HashMachine) workOnFile(sourceFileInfo cfs.CustomFileInfo, destinationDirInfo cfs.CustomFileInfo) error {
	fileHash, err := hashMachine.getChecksum(sourceFileInfo)

	clog.CheckIfError(err)

	extension := filepath.Ext(sourceFileInfo.GetPath())
	destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

	destinationFileInfo := cfs.GetUnvalidatedPath(destination, cfs.PathIsFile)

	if hashMachine.Options.DryRun {
		if sourceFileInfo.GetPath() == destination {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destinationFileInfo,
			)
		} else {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationDryRun,
				sourceFileInfo,
				destinationFileInfo,
			)
		}
		return nil
	}

	// TODO(16): Check if has permission to move to destination
	err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
	if err == nil {
		ReportOperation(
			MachineOptions(hashMachine.Options),
			OperationRenamed,
			sourceFileInfo,
			destinationFileInfo,
		)
		return nil
	} else if errors.Is(err, cfs.ErrSameFile) {
		ReportOperation(
			MachineOptions(hashMachine.Options),
			OperationSameFile,
			sourceFileInfo,
			destinationFileInfo,
		)
		return nil
	} else if !errors.Is(err, cfs.ErrFileExists) {
		return err
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.GetPath(), newFileName)
		destinationFileInfo = cfs.GetUnvalidatedPath(destination, cfs.PathIsFile)

		err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
		if err == nil {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationRenamed,
				sourceFileInfo,
				destinationFileInfo,
			)

			return nil
		} else if errors.Is(err, cfs.ErrSameFile) {
			ReportOperation(
				MachineOptions(hashMachine.Options),
				OperationSameFile,
				sourceFileInfo,
				destinationFileInfo,
			)
			return nil
		} else if !errors.Is(err, cfs.ErrFileExists) {
			return err
		}

		counter++
	}
}
