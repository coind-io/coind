package crypto

import (
	"errors"
	"io"

	"golang.org/x/crypto/ed25519"

	"github.com/coind-io/coind/lib/base58"
	"github.com/coind-io/coind/lib/encoding"
)

type PrivKey struct {
	bytes [PrivKeySize]byte
}

func NewPrivKey() *PrivKey {
	pk := new(PrivKey)
	return pk
}

func NewPrivKeyFromBytes(bytes []byte) (*PrivKey, error) {
	pk := NewPrivKey()
	err := pk.Update(bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func NewPrivKeyFromWIF(wif_pk string) (*PrivKey, error) {
	bytes, btype, err := base58.CheckDecode(wif_pk)
	if err != nil {
		return nil, err
	}
	if btype != PrivKeyType {
		return nil, errors.New("unknow wallet import format")
	}
	pk := NewPrivKey()
	err = pk.Update(bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func (pk *PrivKey) Update(bytes []byte) error {
	size := len(bytes)
	if size != PrivKeySize {
		return errors.New("Illegal PrivKey Length")
	}
	for i := range bytes {
		pk.bytes[i] = bytes[i]
	}
	return nil
}

func (pk *PrivKey) PubKey() *PubKey {
	edkey := ed25519.NewKeyFromSeed(pk.bytes[:])
	edpub := make([]byte, ed25519.PublicKeySize)
	copy(edpub, edkey[32:])
	pub := NewPubKey()
	pub.Update(edpub)
	return pub
}

func (pk *PrivKey) Sign(data []byte) *Signature {
	edkey := ed25519.NewKeyFromSeed(pk.bytes[:])
	edsign := ed25519.Sign(edkey, data)
	sign, _ := NewSignatureFromBytes(edsign)
	return sign
}

func (pk *PrivKey) ToBytes() []byte {
	return pk.bytes[:]
}

func (pk *PrivKey) ToWIF() string {
	return base58.CheckEncode(pk.bytes[:], PrivKeyType)
}

func (pk *PrivKey) String() string {
	return pk.ToWIF()
}

func (pk *PrivKey) Decode(r io.Reader) error {
	n, err := r.Read(pk.bytes[:])
	if n <= 0 || err != nil {
		return encoding.ErrEof
	}
	return nil
}

func (pk *PrivKey) Encode(w io.Writer) error {
	_, err := w.Write(pk.bytes[:])
	if err != nil {
		return err
	}
	return nil
}
