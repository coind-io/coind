package tx

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/coind-io/coind/lib/base58"
	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/encoding"
)

type Address struct {
	atype  byte
	digest [20]byte
}

func NewAddress() *Address {
	address := new(Address)
	address.atype = 0x01
	return address
}

func NewAddressFromBytes(bytes []byte) (*Address, error) {
	if len(bytes) != 21 {
		return nil, fmt.Errorf("invalid address from bytes length of %v, want %v", len(bytes), 21)
	}
	address := NewAddress()
	err := address.Update(bytes[0], bytes[1:])
	if err != nil {
		return nil, err
	}
	return address, nil
}

func NewAddressFromWIF(wif_addr string) (*Address, error) {
	digest, atype, err := base58.CheckDecode(wif_addr)
	if err != nil {
		return nil, err
	}
	address := NewAddress()
	err = address.Update(atype, digest)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func NewAddressFromPubKey(pub *crypto.PubKey) (*Address, error) {
	pkh := pub.Hash()
	address := NewAddress()
	err := address.Update(byte(P2PKHType), pkh.Bytes())
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (addr *Address) Update(atype byte, digest []byte) error {
	size := len(digest)
	if size != 20 {
		return fmt.Errorf("invalid address digest length of %v, want %v", size, 20)
	}
	addr.atype = atype
	copy(addr.digest[:], digest)
	return nil
}

func (addr *Address) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	atype := dec.DecodeByte()
	digest := dec.DecodeBytes()
	if dec.Error() != nil {
		return dec.Error()
	}
	err := addr.Update(atype, digest)
	if err != nil {
		return err
	}
	return nil
}

func (addr *Address) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeByte(addr.atype)
	enc.EncodeBytes(addr.digest[:])
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (addr *Address) Type() byte {
	return addr.atype
}

func (addr *Address) Digest() []byte {
	return addr.digest[:]
}

func (addr *Address) ToBytes() []byte {
	bytes := make([]byte, 0, 21)
	bytes = append(bytes, addr.atype)
	bytes = append(bytes, addr.digest[:]...)
	return bytes
}

func (addr *Address) ToWIF() string {
	return base58.CheckEncode(addr.digest[:], addr.atype)
}

func (addr *Address) String() string {
	return addr.ToWIF()
}

func (addr *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(addr.ToWIF())
}
