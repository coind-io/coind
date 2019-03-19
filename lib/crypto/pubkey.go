package crypto

import (
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/crypto/ed25519"

	"github.com/coind-io/coind/lib/base58"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type PubKey struct {
	bytes [PubKeySize]byte
}

func NewPubKey() *PubKey {
	pub := new(PubKey)
	return pub
}

func NewPubKeyFromBytes(bytes []byte) (*PubKey, error) {
	pub := new(PubKey)
	err := pub.Update(bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func (pub *PubKey) Update(bytes []byte) error {
	size := len(bytes)
	if size != PubKeySize {
		return errors.New("Illegal PubKey Length")
	}
	for i := range bytes {
		pub.bytes[i] = bytes[i]
	}
	return nil
}

func (pub *PubKey) Verify(digest []byte, sign *Signature) bool {
	edpub := ed25519.PublicKey(pub.bytes[:])
	return ed25519.Verify(edpub, digest, sign.bytes[:])
}

func (pub *PubKey) IsEqual(other *PubKey) bool {
	if other == nil {
		return false
	}
	for i := len(pub.bytes) - 1; i >= 0; i-- {
		if pub.bytes[i] != other.bytes[i] {
			return false
		}
	}
	return true
}

func (pub *PubKey) IsZero() bool {
	zero := NewPubKey()
	return zero.IsEqual(pub)
}

func (pub *PubKey) Hash() *hash.Hash {
	first := hash.SumHash256(pub.bytes[:])
	second := hash.SumHash160(first.Bytes())
	return second
}

func (pub *PubKey) Bytes() []byte {
	return pub.bytes[:]
}

func (pub *PubKey) ToWIF() string {
	return base58.CheckEncode(pub.bytes[:], PubKeyType)
}

func (pub *PubKey) String() string {
	return pub.ToWIF()
}

func (pub *PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.String())
}

func (pub *PubKey) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	bytes := dec.DecodeBytes()
	if dec.Error() != nil {
		return dec.Error()
	}
	return pub.Update(bytes)
}

func (pub *PubKey) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeBytes(pub.bytes[:])
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}
