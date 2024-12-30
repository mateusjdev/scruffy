package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/clog"
	"os"

	"lukechampine.com/blake3"
)

type HashAlgorithm uint8

const (
	HashNotSet HashAlgorithm = iota
	HashSHA512
	HashSHA384
	HashSHA256
	HashSHA224
	HashSHA1
	HashMD5
	HashBlake3
	HashBlake2b
	HashFuzzy
)

// TODO(18): Merge checkHash and getHashMachine?
// checkHash makes a fast check if Hash is valid, getHashMachine create the HashMachine
// benchmark time taken by both
func checkHash(hash string) (HashAlgorithm, error) {
	switch hash {
	case "sha512":
		return HashSHA512, nil
	case "sha256":
		return HashSHA256, nil
	case "sha1":
		return HashSHA1, nil
	case "md5":
		return HashMD5, nil
	case "blake3":
		return HashBlake3, nil
	case "blake2b":
		return HashBlake2b, nil
	case "fuzzy":
		return HashFuzzy, nil
	default:
		return HashNotSet, errors.New("hash method not valid")
	}
}

// TODO(8): Work on lenght/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --lenght 224"

// TODO(13): Create a HashMachine interface, add Options
// Options: {algorithm, lenght/truncate, uppercase, ...}
func getHashMachine(algorithm HashAlgorithm, lenght int) (hash.Hash, error) {
	switch algorithm {
	case HashBlake2b:
		// lenght: fixed_256_bits (256, 384, 512)
		return crypto.BLAKE2b_256.New(), nil
	case HashBlake3:
		return blake3.New(lenght, nil), nil
	case HashFuzzy:
		panic("Not implemented")
	case HashMD5:
		// lenght: fixed_128_bits
		return crypto.MD5.New(), nil
	case HashSHA1:
		// lenght: fixed_160_bits
		return crypto.SHA1.New(), nil
	case HashSHA256:
		// lenght: fixed_256_bits
		return crypto.SHA256.New(), nil
	case HashSHA512:
		// lenght: fixed_512_bits
		return crypto.SHA512.New(), nil
	}
	return nil, errors.New("unable to create a hasher")
}

func hashFile(fileInfo CustomFileInfo, hasher hash.Hash) (string, error) {
	if fileInfo.PathType != PathIsFile {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.Path)
	clog.CheckIfError(err)
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	// TODO(17): Benchark reuse(.Reset()) vs recreate
	hasher.Reset()

	return hashString, nil
}
