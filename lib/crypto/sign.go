package crypto

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"

	"github.com/coind-io/coind/lib/encoding"
)

type Signature struct {
	bytes [SignSize]byte
}

func NewSignature() *Signature {
	sign := new(Signature)
	return sign
}

func NewSignatureFromBytes(bytes []byte) (*Signature, error) {
	sign := new(Signature)
	err := sign.Update(bytes)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func (s *Signature) Update(bytes []byte) error {
	size := len(bytes)
	if size != SignSize {
		return errors.New("Illegal Signature Length")
	}
	for i := range bytes {
		s.bytes[i] = bytes[i]
	}
	return nil
}

func (s *Signature) String() string {
	return hex.EncodeToString(s.bytes[:])
}

func (s *Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Signature) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	bytes := dec.DecodeBytes()
	if dec.Error() != nil {
		return dec.Error()
	}
	return s.Update(bytes)
}

func (s *Signature) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeBytes(s.bytes[:])
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}
