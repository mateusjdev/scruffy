package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/clog"
	"os"
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

// TODO(8): Work on lenght/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --lenght 224"
func getHashAlgorithm(hash string, lenght int) (hash.Hash, error) {

	hash = strings.ToLower(hash)

	switch hash {
	case "blake2b":
		// lenght: fixed_256_bits (256, 384, 512)
		return crypto.BLAKE2b_256.New(), nil
	case "blake3":
		return blake3.New(lenght/2, nil), nil
	case "fuzzy":
		panic("Not implemented")
	case "md5":
		// lenght: fixed_128_bits
		return crypto.MD5.New(), nil
	case "sha1":
		// lenght: fixed_160_bits
		return crypto.SHA1.New(), nil
	case "sha256":
		// lenght: fixed_256_bits
		return crypto.SHA256.New(), nil
	case "sha512":
		// lenght: fixed_512_bits
		return crypto.SHA512.New(), nil
	}
	return nil, errors.New("hash method not valid")
}

func (hashMachine HashMachine) getChecksum(fileInfo CustomFileInfo) (string, error) {
	if fileInfo.PathType != PathIsFile {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.Path)
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
