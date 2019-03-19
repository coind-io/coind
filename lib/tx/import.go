package tx

import (
	"bytes"
	"io"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type Import struct {
	Genesis bool       `json:"genesis"`
	TxHash  *hash.Hash `json:"txhash"`
	TxIndex uint16     `json:"txindex"`
	Redeem  []byte     `json:"redeem"`
}

func NewImport() *Import {
	imp := new(Import)
	imp.Genesis = false
	imp.TxHash = hash.NewHash256()
	imp.TxIndex = 0
	imp.Redeem = []byte{}
	return imp
}

func (imp *Import) Key() []byte {
	buf := bytes.NewBuffer(nil)
	enc := encoding.NewEncoder(buf)
	enc.Encode(imp.TxHash)
	enc.EncodeUint16(imp.TxIndex)
	return buf.Bytes()
}

func (imp *Import) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeBool(imp.Genesis)
	enc.Encode(imp.TxHash)
	enc.EncodeUint16(imp.TxIndex)
	enc.EncodeBytes(imp.Redeem)
	return enc.Error()
}

func (imp *Import) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	imp.Genesis = dec.DecodeBool()
	imp.TxHash = hash.NewHash256()
	dec.Decode(imp.TxHash)
	imp.TxIndex = dec.DecodeUint16()
	imp.Redeem = dec.DecodeBytes()
	return dec.Error()
}
