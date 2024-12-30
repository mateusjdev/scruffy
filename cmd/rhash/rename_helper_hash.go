package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"mateusjdev/scruffy/cmd/cfs"
	"mateusjdev/scruffy/cmd/clog"
	"os"
	"path/filepath"
	"strings"

	"lukechampine.com/blake3"
)

type HashMachineOptions struct {
	uppercase bool
	truncate  uint8
}

type HashMachine struct {
	Machine hash.Hash
	Options HashMachineOptions
}

var (
	ErrUnknownHashMethod = errors.New("hash method not valid")
)

// TODO(8): Work on lenght/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --lenght 224"
func getHashAlgorithm(hash string, lenght int) (hash.Hash, error) {
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

	// TODO(17): Benchark reuse(.Reset()) vs recreate
	hashMachine.Machine.Reset()

	if hashMachine.Options.truncate != 0 {
		hashString = hashString[0:hashMachine.Options.truncate]
	}

	if hashMachine.Options.uppercase {
		hashString = strings.ToUpper(hashString)
	}

	return hashString, nil
}

func (hashMachine HashMachine) workOnFile(sourceFileInfo cfs.CustomFileInfo, destinationDirInfo cfs.CustomFileInfo) error {
	clog.Debugf("Working on file \"%s\"", sourceFileInfo.GetPath())

	fileHash, err := hashMachine.getChecksum(sourceFileInfo)
	clog.CheckIfError(err)

	extension := filepath.Ext(sourceFileInfo.GetPath())
	destination := filepath.Join(destinationDirInfo.GetPath(), fileHash+extension)

	// TODO(16): Check if has permission to move to destination
	err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
	if err == nil {
		clog.InfoIconf(clog.PrintIconSuccess, "\"%s\" -> %s", sourceFileInfo.GetPath(), destination)
		return nil
	} else if errors.Is(err, cfs.ErrSameFile) {
		clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", sourceFileInfo.GetPath())
		return nil
	} else if !errors.Is(err, cfs.ErrFileExists) {
		return err
	}

	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d%s", fileHash, counter, extension)
		destination := filepath.Join(destinationDirInfo.GetPath(), newFileName)

		err = cfs.SafeRename(sourceFileInfo.GetPath(), destination)
		if err == nil {
			clog.InfoIconf(clog.PrintIconSuccess, "\"%s\" -> %s", sourceFileInfo.GetPath(), destination)
			return nil
		} else if errors.Is(err, cfs.ErrSameFile) {
			clog.InfoIconf(clog.PrintIconNothing, "file %s already hashed", sourceFileInfo.GetPath())
			return nil
		} else if !errors.Is(err, cfs.ErrFileExists) {
			return err
		}

		counter++
	}
}

// TODO(14): Check need of path validation or continue to use CustomFileInfo
func (hashMachine HashMachine) enqueuePath(inputPathInfo cfs.CustomFileInfo, outputPathInfo cfs.CustomFileInfo) error {
	clog.Debugf("Enqueued: \"%s\"", inputPathInfo.GetPath())

	if inputPathInfo.GetPathType() == cfs.PathIsFile {
		return hashMachine.workOnFile(inputPathInfo, outputPathInfo)
	}

	if inputPathInfo.GetPathType() != cfs.PathIsDirectory {
		clog.Errorf("Not a valid file or directory")
		clog.ExitBecause(clog.ErrCodeGeneric)
	}

	// TODO(21): Check WalkDir error/return
	filepath.WalkDir(inputPathInfo.GetPath(), func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// TODO(7): Work on recursive flag
		// Recurse into directories
		if di.IsDir() {
			if inputPathInfo.GetPath() == path {
				return nil
			}
			return filepath.SkipDir
		}

		fileInfo := cfs.GetUnvalidatedPath(path, cfs.PathIsFile)
		return hashMachine.workOnFile(fileInfo, outputPathInfo)
	})

	return nil
}
