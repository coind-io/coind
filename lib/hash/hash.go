package hash

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"

	"github.com/coind-io/coind/lib/encoding"
)

type Hash struct {
	bytes []byte
	size  int
}

func New(size int) *Hash {
	hash := new(Hash)
	if size == 0 {
		size = Hash256Size
	}
	hash.bytes = make([]byte, size)
	hash.size = size
	return hash
}

func NewHash160() *Hash {
	return New(Hash160Size)
}

func NewHash256() *Hash {
	return New(Hash256Size)
}

func NewHash160FromBytes(bytes []byte) (*Hash, error) {
	hash := NewHash160()
	err := hash.Update(bytes)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func NewHash256FromHex(hexstr string) (*Hash, error) {
	bytes, err := hex.DecodeString(hexstr)
	if err != nil {
		return nil, err
	}
	hash := NewHash256()
	err = hash.Update(bytes)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (hash *Hash) Update(bytes []byte) error {
	if len(bytes) != hash.size {
		return errors.New("hash length error")
	}
	for i := range bytes {
		hash.bytes[i] = bytes[i]
	}
	return nil
}

func (hash *Hash) IsEqual(other *Hash) bool {
	if other == nil {
		return false
	}
	if len(other.bytes) != len(hash.bytes) {
		return false
	}
	for i := len(hash.bytes) - 1; i >= 0; i-- {
		if hash.bytes[i] != other.bytes[i] {
			return false
		}
	}
	return true
}

func (hash *Hash) IsZero() bool {
	zero := New(hash.size)
	return zero.IsEqual(hash)
}

func (hash *Hash) Size() int {
	return hash.size
}

func (hash *Hash) Bytes() []byte {
	return hash.bytes
}

func (hash *Hash) String() string {
	return hex.EncodeToString(hash.bytes)
}

func (hash *Hash) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	bytes := dec.DecodeBytes()
	if dec.Error() != nil {
		return dec.Error()
	}
	hash.Update(bytes)
	return nil
}

func (hash *Hash) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeBytes(hash.bytes)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (hash *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(hash.String())
}
