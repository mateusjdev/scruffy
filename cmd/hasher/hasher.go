package hasher

import (
	"crypto"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"mateusjdev/scruffy/cmd/filesystem"
	"os"

	"lukechampine.com/blake3"
)

const (
	HashAlgorithmBlake2b string = "blake2b"
	HashAlgorithmBlake3  string = "blake3"
	HashAlgorithmMD5     string = "md5"
	HashAlgorithmSHA1    string = "sha1"
	HashAlgorithmSHA256  string = "sha256"
	HashAlgorithmSHA512  string = "sha512"
)

var (
	ErrUnknownHashMethod = errors.New("hash method not valid")
)

type Hasher struct {
	hash.Hash
}

// TODO(8): Work on length/truncate flag
// Chosse between "--hash SHA224 ..." or "--hash SHA2 --length 224"
func NewHasher(algorithm string, length int) (*Hasher, error) {
	switch algorithm {
	case HashAlgorithmBlake2b:
		// length: fixed_256_bits (256, 384, 512)
		return &Hasher{crypto.BLAKE2b_256.New()}, nil
	case HashAlgorithmBlake3:
		return &Hasher{blake3.New(length/2, nil)}, nil
	case HashAlgorithmMD5:
		// length: fixed_128_bits
		return &Hasher{crypto.MD5.New()}, nil
	case HashAlgorithmSHA1:
		// length: fixed_160_bits
		return &Hasher{crypto.SHA1.New()}, nil
	case HashAlgorithmSHA256:
		// length: fixed_256_bits
		return &Hasher{crypto.SHA256.New()}, nil
	case HashAlgorithmSHA512:
		// length: fixed_512_bits
		return &Hasher{crypto.SHA512.New()}, nil
	}
	return nil, ErrUnknownHashMethod
}

func (hasher *Hasher) Checksum(fileInfo *filesystem.PathInfo) (string, error) {
	if fileInfo == nil || !fileInfo.IsRegularFile() {
		return "", errors.New("trying to hash a non file")
	}

	file, err := os.Open(fileInfo.Path())
	if err != nil {
		return "", err
	}

	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	hasher.Reset()

	return hashString, nil
}
