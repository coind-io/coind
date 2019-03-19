package encoding

import (
	"encoding/binary"
	"io"
	"math/big"
)

type Encoder struct {
	err error
	w   io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	enc := new(Encoder)
	enc.err = nil
	enc.w = w
	return enc
}

func (enc *Encoder) Encode(ie IEncoder) {
	if enc.err != nil {
		return
	}
	enc.err = ie.Encode(enc.w)
	return
}

func (enc *Encoder) EncodeUint8(value uint8) {
	if enc.err != nil {
		return
	}
	buf := []byte{byte(value)}
	_, err := enc.w.Write(buf)
	if err != nil {
		enc.err = err
		return
	}
	return
}

func (enc *Encoder) EncodeUint16(value uint16) {
	if enc.err != nil {
		return
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	_, err := enc.w.Write(buf)
	if err != nil {
		enc.err = err
		return
	}
	return
}

func (enc *Encoder) EncodeUint32(value uint32) {
	if enc.err != nil {
		return
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)
	_, err := enc.w.Write(buf)
	if err != nil {
		enc.err = err
		return
	}
	return
}

func (enc *Encoder) EncodeUint64(value uint64) {
	if enc.err != nil {
		return
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	_, err := enc.w.Write(buf)
	if err != nil {
		enc.err = err
		return
	}
	return
}

func (enc *Encoder) EncodeVarint(value uint64) {
	if enc.err != nil {
		return
	}
	if value <= 0xFC {
		enc.EncodeUint8(uint8(value))
		return
	}
	if value <= 0xFFFF {
		enc.EncodeUint8(0xFD)
		enc.EncodeUint16(uint16(value))
		return
	}
	if value <= 0xFFFFFFFF {
		enc.EncodeUint8(0xFE)
		enc.EncodeUint32(uint32(value))
		return
	}
	enc.EncodeUint8(0xFF)
	enc.EncodeUint64(value)
	return
}

func (enc *Encoder) EncodeByte(value byte) {
	if enc.err != nil {
		return
	}
	enc.EncodeUint8(uint8(value))
	return
}

func (enc *Encoder) EncodeBool(value bool) {
	if enc.err != nil {
		return
	}
	var number uint8
	if value == true {
		number = 1
	}
	enc.EncodeUint8(number)
	if enc.err != nil {
		return
	}
	return
}

func (enc *Encoder) EncodeBytes(value []byte) {
	if enc.err != nil {
		return
	}
	enc.EncodeVarint(uint64(len(value)))
	if enc.err != nil {
		return
	}
	_, err := enc.w.Write(value)
	if err != nil {
		enc.err = err
		return
	}
	return
}

func (enc *Encoder) EncodeBInt(bint *big.Int) {
	if enc.err != nil {
		return
	}
	enc.EncodeBytes(bint.Bytes())
	return
}

func (enc *Encoder) EncodeString(value string) {
	if enc.err != nil {
		return
	}
	enc.EncodeBytes([]byte(value))
	return
}

func (enc *Encoder) Error() error {
	return enc.err
}
