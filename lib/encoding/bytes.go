package encoding

import (
	"encoding/binary"
)

func BConnect(slices ...[]byte) []byte {
	total := int(0)
	for _, s := range slices {
		total = total + len(s)
	}
	dst := make([]byte, 0, total)
	for _, s := range slices {
		for _, b := range s {
			dst = append(dst, b)
		}
	}
	return dst
}

func I2b(value uint64) []byte {
	dst := make([]byte, 8)
	binary.BigEndian.PutUint64(dst, value)
	return dst
}

func B2i(raw []byte) uint64 {
	return binary.BigEndian.Uint64(raw)
}
