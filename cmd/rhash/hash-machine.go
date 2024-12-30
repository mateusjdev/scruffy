package rhash

import (
	"crypto"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/common"
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

func checkHash(hash string) (HashAlgorithm, error) {
	// TODO: improve? map?
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

// TODO: add lenght
func createHasherBlake2b() (hash.Hash, error) {
	return crypto.BLAKE2b_256.New(), nil
}
func createHasherBlake3(lenght int) (hash.Hash, error) {
	// lenght: length_1:1
	return blake3.New(lenght, nil), nil
}

// TODO: fuzzy hash
func createHasherFuzzy(_ int) (hash.Hash, error) {
	// TODO: length
	panic("Not implemented")
}
func createHasherMD5() (hash.Hash, error) {
	// lenght: fixed_128_bits
	return crypto.MD5.New(), nil
}
func createHasherSHA1() (hash.Hash, error) {
	// lenght: fixed_160_bits
	return crypto.SHA1.New(), nil
}

// TODO: chosse between "--hash SHA224,256, ..." or "--hash SHA2 --lenght 512"
func createHasherSHA256() (hash.Hash, error) {
	// lenght: fixed_256_bits
	return crypto.SHA256.New(), nil
}
func createHasherSHA512() (hash.Hash, error) {
	// lenght: fixed_512_bits
	return crypto.SHA512.New(), nil
}

// TODO: Create interface hashMachine (hash.Hash, Uppercase)
func getHashMachine(algorithm HashAlgorithm, lenght int) (hash.Hash, error) {
	switch algorithm {
	case HashBlake2b:
		return createHasherBlake2b()
	case HashBlake3:
		return createHasherBlake3(lenght)
	case HashFuzzy:
		return createHasherFuzzy(lenght)
	case HashMD5:
		return createHasherMD5()
	case HashSHA1:
		return createHasherSHA1()
	case HashSHA256:
		return createHasherSHA256()
	case HashSHA512:
		return createHasherSHA512()
	}
	return nil, errors.New("unable to create a hasher")
}

func hashFile(fileInfo CustomFileInfo, hasher hash.Hash) (string, error) {
	// TODO: benchmark reuse
	hasher.Reset()

	if fileInfo.PathType != PathIsFile {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.Path)
	common.CheckIfError(err)
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}
