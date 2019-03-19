package tx

import (
	"bytes"
	"io"
	"math/big"
	"time"

	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type Transfer struct {
	Version   uint16              `json:"version"`
	Imports   []*Import           `json:"imports"`
	Exports   []Export            `json:"exports"`
	Timestamp *big.Int            `json:"timestamp"`
	Signature []*crypto.Signature `json:"signatures"`
}

func NewTransfer() *Transfer {
	tx := new(Transfer)
	tx.Version = 0
	tx.Imports = []*Import{}
	tx.Exports = []Export{}
	tx.Timestamp = big.NewInt(time.Now().Unix())
	tx.Signature = []*crypto.Signature{}
	return tx
}

func (tx *Transfer) AddImport(imp *Import) {
	if imp == nil {
		return
	}
	if len(tx.Imports) == 32 {
		return
	}
	tx.Imports = append(tx.Imports, imp)
	return
}

func (tx *Transfer) AddExport(export Export) {
	if export == nil {
		return
	}
	if len(tx.Exports) == 32 {
		return
	}
	tx.Exports = append(tx.Exports, export)
	return
}

func (tx *Transfer) Hash() (*hash.Hash, error) {
	buf := bytes.NewBuffer(nil)
	enc := encoding.NewEncoder(buf)

	enc.EncodeUint16(tx.Version)
	enc.EncodeUint16(uint16(len(tx.Imports)))
	for i := 0; i < len(tx.Imports); i++ {
		enc.Encode(tx.Imports[i])
	}
	enc.EncodeUint16(uint16(len(tx.Exports)))
	for i := 0; i < len(tx.Exports); i++ {
		enc.EncodeUint16(uint16(tx.Exports[i].Type()))
		enc.Encode(tx.Exports[i])
	}
	enc.EncodeBInt(tx.Timestamp)

	if enc.Error() != nil {
		return nil, enc.Error()
	}

	return hash.SumDoubleHash256(buf.Bytes()), nil
}

func (tx *Transfer) Sign(privkey *crypto.PrivKey) error {
	digest, err := tx.Hash()
	if err != nil {
		return err
	}
	sign := privkey.Sign(digest.Bytes())
	tx.Signature = append(tx.Signature, sign)
	return nil
}

func (tx *Transfer) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeUint16(tx.Version)
	enc.EncodeUint16(uint16(len(tx.Imports)))
	for i := 0; i < len(tx.Imports); i++ {
		enc.Encode(tx.Imports[i])
	}
	enc.EncodeUint16(uint16(len(tx.Exports)))
	for i := 0; i < len(tx.Exports); i++ {
		enc.EncodeUint16(uint16(tx.Exports[i].Type()))
		enc.Encode(tx.Exports[i])
	}
	enc.EncodeBInt(tx.Timestamp)
	enc.EncodeUint16(uint16(len(tx.Signature)))
	for i := 0; i < len(tx.Signature); i++ {
		enc.Encode(tx.Signature[i])
	}
	return enc.Error()
}

func (tx *Transfer) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	tx.Version = dec.DecodeUint16()

	length := int(dec.DecodeUint16())
	tx.Imports = make([]*Import, length)
	for i := 0; i < length; i++ {
		imp := new(Import)
		dec.Decode(imp)
		tx.Imports[i] = imp
	}

	length = int(dec.DecodeUint16())
	tx.Exports = make([]Export, length)
	for i := 0; i < length; i++ {
		etype := dec.DecodeUint16()
		exp, err := NewExport(ExportType(etype))
		if err != nil {
			return err
		}
		dec.Decode(exp)
		tx.Exports[i] = exp
	}

	tx.Timestamp = dec.DecodeBInt()

	length = int(dec.DecodeUint16())
	tx.Signature = make([]*crypto.Signature, length)
	for i := 0; i < length; i++ {
		sig := crypto.NewSignature()
		dec.Decode(sig)
		tx.Signature[i] = sig
	}
	return dec.Error()
}
