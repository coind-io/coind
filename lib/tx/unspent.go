package tx

import (
	"bytes"
	"io"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type Unspent struct {
	TxHash  *hash.Hash `json:"txhash"`
	TxIndex uint16     `json:"txindex"`
	Export  Export     `json:"export"`
}

func NewUnspent() *Unspent {
	unspent := new(Unspent)
	unspent.TxHash = hash.NewHash256()
	unspent.TxIndex = 0
	unspent.Export = nil
	return unspent
}

func (unspent *Unspent) Key() []byte {
	buf := bytes.NewBuffer(nil)
	enc := encoding.NewEncoder(buf)
	enc.Encode(unspent.TxHash)
	enc.EncodeUint16(unspent.TxIndex)
	return buf.Bytes()
}

func (unspent *Unspent) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.Encode(unspent.TxHash)
	enc.EncodeUint16(unspent.TxIndex)
	enc.EncodeByte(byte(unspent.Export.Type()))
	enc.Encode(unspent.Export)
	return enc.Error()
}

func (unspent *Unspent) Decode(r io.Reader) (err error) {
	dec := encoding.NewDecoder(r)
	unspent.TxHash = hash.NewHash256()
	dec.Decode(unspent.TxHash)
	unspent.TxIndex = dec.DecodeUint16()
	etype := dec.DecodeByte()
	unspent.Export, err = NewExport(ExportType(etype))
	if err != nil {
		return err
	}
	dec.Decode(unspent.Export)
	return dec.Error()
}
