package block

import (
	"bytes"
	"io"
	"time"

	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type Header struct {
	Version    uint16            `json:"version"`
	ChainID    uint64            `json:"chainid"`
	Height     uint64            `json:"height"`
	Timestamp  uint64            `json:"timestamp"`
	PrevBlock  *hash.Hash        `json:"prevblock"`
	MerkleRoot *hash.Hash        `json:"merkleroot"`
	Creator    *crypto.PubKey    `json:"creator"`
	Signature  *crypto.Signature `json:"signature"`
}

func NewHeader(height uint64) *Header {
	header := new(Header)
	header.Version = BLOCK_VERSION
	header.ChainID = BLOCK_CHAINID
	header.Height = height
	header.Timestamp = uint64(time.Now().Unix())
	header.PrevBlock = hash.NewHash256()
	header.MerkleRoot = hash.NewHash256()
	header.Creator = crypto.NewPubKey()
	header.Signature = crypto.NewSignature()
	return header
}

func (h *Header) Encode(w io.Writer) error {
	enc := encoding.NewEncoder(w)
	enc.EncodeUint16(h.Version)
	enc.EncodeUint64(h.ChainID)
	enc.EncodeUint64(h.Height)
	enc.EncodeUint64(h.Timestamp)
	enc.Encode(h.PrevBlock)
	enc.Encode(h.MerkleRoot)
	enc.Encode(h.Creator)
	enc.Encode(h.Signature)
	return enc.Error()
}

func (h *Header) Decode(r io.Reader) error {
	dec := encoding.NewDecoder(r)
	h.Version = dec.DecodeUint16()
	h.ChainID = dec.DecodeUint64()
	h.Height = dec.DecodeUint64()
	h.Timestamp = dec.DecodeUint64()
	h.PrevBlock = hash.NewHash256()
	dec.Decode(h.PrevBlock)
	h.MerkleRoot = hash.NewHash256()
	dec.Decode(h.MerkleRoot)
	h.Creator = crypto.NewPubKey()
	dec.Decode(h.Creator)
	h.Signature = crypto.NewSignature()
	dec.Decode(h.Signature)
	return dec.Error()
}

func (h *Header) Hash() (*hash.Hash, error) {
	buf := bytes.NewBuffer(nil)
	enc := encoding.NewEncoder(buf)
	enc.EncodeUint16(h.Version)
	enc.EncodeUint64(h.ChainID)
	enc.EncodeUint64(h.Height)
	enc.EncodeUint64(h.Timestamp)
	enc.Encode(h.PrevBlock)
	enc.Encode(h.MerkleRoot)
	enc.Encode(h.Creator)
	err := enc.Error()
	if err != nil {
		return nil, err
	}
	return hash.SumDoubleHash256(buf.Bytes()), nil
}

func (h *Header) Sign(pk *crypto.PrivKey) error {
	digest, err := h.Hash()
	if err != nil {
		return err
	}
	h.Signature = pk.Sign(digest.Bytes())
	return nil
}

func (h *Header) Fork() (*Header, error) {
	prev, err := h.Hash()
	if err != nil {
		return nil, err
	}
	dst := NewHeader(h.Height + 1)
	dst.PrevBlock = prev
	return dst, nil
}
