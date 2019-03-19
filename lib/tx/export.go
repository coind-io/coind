package tx

import (
	"encoding/json"
	"errors"
	"io"
	"math/big"

	"github.com/coind-io/coind/lib/crypto"
)

type ExportType uint16

const (
	P2PKHType ExportType = 0x00
)

func (etype ExportType) String() string {
	switch etype {
	case P2PKHType:
		return "p2pkh"
	default:
		return "unknow"
	}
}

func (etype *ExportType) MarshalJSON() ([]byte, error) {
	return json.Marshal(etype.String())
}

type Voucher struct {
	Digest []byte
	Redeem []byte
	Signs  []*crypto.Signature
}

type Export interface {
	Type() ExportType
	GetAmount() *big.Int
	SetAmount(*big.Int)
	GetAddress() (*Address, error)
	SetAddress(*Address) error
	UnLock(voucher *Voucher) error
	Decode(r io.Reader) error
	Encode(w io.Writer) error
}

func NewExport(etype ExportType) (Export, error) {
	switch etype {
	case P2PKHType:
		return NewP2PKH(), nil
	default:
		return nil, errors.New("not exists export type")
	}
}
