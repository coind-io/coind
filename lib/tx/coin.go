package tx

import (
	"bytes"
	"io"
	"math/big"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type Coin struct {
	TxHash  *hash.Hash `json:"txhash"`
	TxIndex uint16     `json:"txindex"`
	Amount  *big.Int   `json:"amount"`
	Redeem  []byte     `json:"redeem"`
}

func NewCoin() *Coin {
	coin := new(Coin)
	coin.TxHash = hash.NewHash256()
	coin.TxIndex = 0
	coin.Amount = big.NewInt(0)
	coin.Redeem = []byte{}
	return coin
}

func (c *Coin) Key() []byte {
	buf := bytes.NewBuffer(nil)
	enc := encoding.NewEncoder(buf)
	enc.Encode(c.TxHash)
	enc.EncodeUint16(c.TxIndex)
	return buf.Bytes()
}

func (c *Coin) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.Encode(c.TxHash)
	enc.EncodeUint16(c.TxIndex)
	enc.EncodeBInt(c.Amount)
	enc.EncodeBytes(c.Redeem)
	return enc.Error()
}

func (c *Coin) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	c.TxHash = hash.NewHash256()
	dec.Decode(c.TxHash)
	c.TxIndex = dec.DecodeUint16()
	c.Amount = dec.DecodeBInt()
	c.Redeem = dec.DecodeBytes()
	return dec.Error()
}
