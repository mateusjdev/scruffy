package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/cfs"
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

// TODO(8): Work on length/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --length 224"
func GetHashAlgorithm(hash string, length int) (hash.Hash, error) {
	switch hash {
	case HashAlgorithmBlake2b:
		// length: fixed_256_bits (256, 384, 512)
		return crypto.BLAKE2b_256.New(), nil
	case HashAlgorithmBlake3:
		return blake3.New(length/2, nil), nil
	case HashAlgorithmMD5:
		// length: fixed_128_bits
		return crypto.MD5.New(), nil
	case HashAlgorithmSHA1:
		// length: fixed_160_bits
		return crypto.SHA1.New(), nil
	case HashAlgorithmSHA256:
		// length: fixed_256_bits
		return crypto.SHA256.New(), nil
	case HashAlgorithmSHA512:
		// length: fixed_512_bits
		return crypto.SHA512.New(), nil
	}
	return nil, ErrUnknownHashMethod
}

func (hashMachine HashMachine) getChecksum(fileInfo *cfs.PathInfo) (string, error) {
	if !fileInfo.IsRegularFile() {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.Path())
	if err != nil {
		return "", err
	}

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

func (hashMachine HashMachine) workOnFile(sourceFileInfo, destinationDirInfo *cfs.PathInfo) error {
	fileHash, err := hashMachine.getChecksum(sourceFileInfo)
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
