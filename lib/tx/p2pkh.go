package tx

import (
	"errors"
	"io"
	"math/big"

	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type P2PKH struct {
	Etype  ExportType `json:"etype"`
	Hash   *hash.Hash `json:"hash"`
	Amount *big.Int   `json:"amount"`
}

func NewP2PKH() *P2PKH {
	pkh := new(P2PKH)
	pkh.Etype = P2PKHType
	pkh.Hash = hash.NewHash160()
	pkh.Amount = big.NewInt(0)
	return pkh
}

func (pkh *P2PKH) Type() ExportType {
	return pkh.Etype
}

func (pkh *P2PKH) GetAddress() (*Address, error) {
	address := NewAddress()
	err := address.Update(byte(P2PKHType), pkh.Hash.Bytes())
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (pkh *P2PKH) SetAddress(address *Address) error {
	atype := address.atype
	etype := ExportType(atype)
	if pkh.Type() != etype {
		return errors.New("not exists address type")
	}
	h, err := hash.NewHash160FromBytes(address.Digest())
	if err != nil {
		return err
	}
	pkh.Hash = h
	return nil
}

func (pkh *P2PKH) GetAmount() *big.Int {
	return pkh.Amount
}

func (pkh *P2PKH) SetAmount(amount *big.Int) {
	pkh.Amount = amount
}

func (pkh *P2PKH) UnLock(voucher *Voucher) error {
	pubkey, err := crypto.NewPubKeyFromBytes(voucher.Redeem)
	if err != nil {
		return err
	}
	if pkh.Hash.IsEqual(pubkey.Hash()) == false {
		return errors.New("Inconsistent pkh unable to redeem")
	}
	for _, sign := range voucher.Signs {
		if pubkey.Verify(voucher.Digest, sign) == true {
			return nil
		}
	}
	return errors.New("no signature matching")
}

func (pkh *P2PKH) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.Encode(pkh.Hash)
	enc.EncodeBInt(pkh.Amount)
	return enc.Error()
}

func (pkh *P2PKH) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	pkh.Hash = hash.NewHash160()
	dec.Decode(pkh.Hash)
	pkh.Amount = dec.DecodeBInt()
	return dec.Error()
}
