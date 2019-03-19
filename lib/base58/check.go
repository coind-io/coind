package base58

import (
	"crypto/sha256"
	"errors"
)

func checksum(src []byte) [4]byte {
	first := sha256.Sum256(src)
	second := sha256.Sum256(first[:])
	cksum := [4]byte{}
	copy(cksum[:], second[:4])
	return cksum
}

func CheckEncode(src []byte, version byte) string {
	raw := make([]byte, 0, 1+len(src)+4)
	raw = append(raw, version)
	raw = append(raw, src[:]...)
	cksum := checksum(raw)
	raw = append(raw, cksum[:]...)
	return string(encode(raw))
}

func CheckDecode(value string) ([]byte, byte, error) {
	raw, err := decode([]byte(value))
	if err != nil {
		return nil, 0, err
	}
	if len(raw) < 5 {
		return nil, 0, errors.New("invalid format: version and/or checksum bytes missing")
	}
	version := raw[0]
	cksum := [4]byte{}
	copy(cksum[:], raw[len(raw)-4:])
	if checksum(raw[:len(raw)-4]) != cksum {
		return nil, 0, errors.New("checksum error")
	}
	return raw[1 : len(raw)-4], version, nil
}
