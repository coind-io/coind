package block

import (
	"bytes"
	"io"

	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
	"github.com/coind-io/coind/lib/tx"
)

type Block struct {
	Header *Header        `json:"header"`
	TxHash []*hash.Hash   `json:"txhash"`
	TxList []*tx.Transfer `json:"txlist"`
}

func NewBlock(height uint64) *Block {
	b := new(Block)
	b.Header = NewHeader(height)
	b.TxHash = []*hash.Hash{}
	b.TxList = []*tx.Transfer{}
	return b
}

func (b *Block) AddTransfer(tx *tx.Transfer) error {
	txhash, err := tx.Hash()
	if err != nil {
		return err
	}
	for _, item := range b.TxHash {
		if item.IsEqual(txhash) {
			return nil
		}
	}
	b.TxHash = append(b.TxHash, txhash)
	b.TxList = append(b.TxList, tx)
	root, err := hash.ComputeMerkleRoot(b.TxHash)
	if err != nil {
		return err
	}
	b.Header.MerkleRoot = root
	return nil
}

func (b *Block) Prev() *hash.Hash {
	return b.Header.PrevBlock
}

func (b *Block) Height() uint64 {
	return b.Header.Height
}

func (b *Block) Timestamp() uint64 {
	return b.Header.Timestamp
}

func (b *Block) MerkleRoot() *hash.Hash {
	return b.Header.MerkleRoot
}

func (b *Block) Creator() *crypto.PubKey {
	return b.Header.Creator
}

func (b *Block) Hash() (*hash.Hash, error) {
	return b.Header.Hash()
}

func (b *Block) Sign(pk *crypto.PrivKey) error {
	return b.Header.Sign(pk)
}

func (b *Block) Fork() (*Block, error) {
	header, err := b.Header.Fork()
	if err != nil {
		return nil, err
	}
	dst := new(Block)
	dst.Header = header
	dst.TxHash = []*hash.Hash{}
	dst.TxList = []*tx.Transfer{}
	return dst, nil
}

func (b *Block) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.Encode(b.Header)
	enc.EncodeUint16(uint16(len(b.TxHash)))
	for i := 0; i < len(b.TxHash); i++ {
		enc.Encode(b.TxHash[i])
	}
	return enc.Error()
}

func (b *Block) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	b.Header = new(Header)
	dec.Decode(b.Header)
	b.TxHash = []*hash.Hash{}
	length := dec.DecodeUint16()
	for i := 0; i < int(length); i++ {
		txhash := hash.New(hash.Hash256Size)
		dec.Decode(txhash)
		b.TxHash = append(b.TxHash, txhash)
	}
	return dec.Error()
}

func (b *Block) Size() uint64 {
	buf := bytes.NewBuffer(nil)
	err := b.Encode(buf)
	if err != nil {
		return 0
	}
	return uint64(len(buf.Bytes()))
}
