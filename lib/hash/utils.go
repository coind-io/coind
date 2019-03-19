package hash

import (
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/ripemd160"
)

func BytesLimit(bytes []byte, limit int) []byte {
	result := make([]byte, limit)
	size := len(bytes)
	if size != limit {
		return result
	}
	for i := range bytes {
		result[i] = bytes[i]
	}
	return result[:]
}

func SumHash160(bytes []byte) *Hash {
	ripemd := ripemd160.New()
	ripemd.Write(bytes)
	hash := New(Hash160Size)
	hash.Update(ripemd.Sum(nil))
	return hash
}

func SumHash256(bytes []byte) *Hash {
	first := sha256.Sum256(bytes)
	hash := New(Hash256Size)
	hash.Update(first[:])
	return hash
}

func SumDoubleHash256(bytes []byte) *Hash {
	first := sha256.Sum256(bytes)
	second := sha256.Sum256(first[:])
	hash := New(Hash256Size)
	hash.Update(second[:])
	return hash
}

func ComputeMerkleRoot(hashes []*Hash) (*Hash, error) {
	if len(hashes) == 0 {
		return nil, errors.New("NewMerkleTree input no item error.")
	}
	if len(hashes) == 1 {
		return hashes[0], nil
	}
	tree, err := NewMerkleTree(hashes)
	if err != nil {
		return nil, err
	}
	return tree.Root.Hash, nil
}
