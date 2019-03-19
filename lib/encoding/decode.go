package encoding

import (
	"encoding/binary"
	"io"
	"math/big"
)

type Decoder struct {
	err error
	r   io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	dec := new(Decoder)
	dec.err = nil
	dec.r = r
	return dec
}

func (dec *Decoder) Decode(id IDecoder) {
	if dec.err != nil {
		return
	}
	dec.err = id.Decode(dec.r)
	return
}

func (dec *Decoder) DecodeUint8() uint8 {
	if dec.err != nil {
		return 0
	}
	buf := make([]byte, 1)
	n, err := dec.r.Read(buf)
	if n <= 0 || err != nil {
		dec.err = ErrEof
		return 0
	}
	return uint8(buf[0])
}

func (dec *Decoder) DecodeUint16() uint16 {
	if dec.err != nil {
		return 0
	}
	buf := make([]byte, 2)
	n, err := dec.r.Read(buf)
	if n <= 0 || err != nil {
		dec.err = ErrEof
		return 0
	}
	return binary.BigEndian.Uint16(buf)
}

func (dec *Decoder) DecodeUint32() uint32 {
	if dec.err != nil {
		return 0
	}
	buf := make([]byte, 4)
	n, err := dec.r.Read(buf)
	if n <= 0 || err != nil {
		dec.err = ErrEof
		return 0
	}
	return binary.BigEndian.Uint32(buf)
}

func (dec *Decoder) DecodeUint64() uint64 {
	if dec.err != nil {
		return 0
	}
	buf := make([]byte, 8)
	n, err := dec.r.Read(buf)
	if n <= 0 || err != nil {
		dec.err = ErrEof
		return 0
	}
	return binary.BigEndian.Uint64(buf)
}

func (dec *Decoder) DecodeVarint() uint64 {
	if dec.err != nil {
		return 0
	}
	first := dec.DecodeUint8()
	if dec.err != nil {
		return 0
	}
	if first == 0xFD {
		u16 := dec.DecodeUint16()
		if dec.err != nil {
			return 0
		}
		return uint64(u16)
	}
	if first == 0xFE {
		u32 := dec.DecodeUint32()
		if dec.err != nil {
			return 0
		}
		return uint64(u32)
	}
	if first == 0xFF {
		u64 := dec.DecodeUint64()
		if dec.err != nil {
			return 0
		}
		return u64
	}
	return uint64(first)
}

func (dec *Decoder) DecodeByte() byte {
	if dec.err != nil {
		return 0
	}
	number := dec.DecodeUint8()
	if dec.err != nil {
		return 0
	}
	return byte(number)
}

func (dec *Decoder) DecodeBool() bool {
	if dec.err != nil {
		return false
	}
	number := dec.DecodeUint8()
	if dec.err != nil {
		return false
	}
	if number == 0 {
		return false
	}
	return true
}

func (dec *Decoder) DecodeBytes() []byte {
	if dec.err != nil {
		return []byte{}
	}
	length := dec.DecodeVarint()
	if dec.err != nil {
		return []byte{}
	}
	if length == 0 {
		return []byte{}
	}
	buf := make([]byte, int(length))
	n, err := dec.r.Read(buf)
	if n <= 0 || err != nil {
		dec.err = ErrEof
		return []byte{}
	}
	return buf
}

func (dec *Decoder) DecodeBInt() *big.Int {
	if dec.err != nil {
		return big.NewInt(0)
	}
	raw := dec.DecodeBytes()
	if dec.err != nil {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(raw)
}

func (dec *Decoder) DecodeString() string {
	if dec.err != nil {
		return ""
	}
	raw := dec.DecodeBytes()
	if dec.err != nil {
		return ""
	}
	return string(raw)
}

func (dec *Decoder) Error() error {
	return dec.err
}
