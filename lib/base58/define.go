package base58

import (
	"strconv"
)

const ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const ENCODED_ZERO = '1'

var INDEXES [256]byte

func init() {
	for i := 0; i < len(INDEXES); i++ {
		INDEXES[i] = 0xFF
	}
	for i := 0; i < len(ALPHABET); i++ {
		INDEXES[ALPHABET[i]] = byte(i)
	}
}

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base58 data at input byte " + strconv.FormatInt(int64(e), 10)
}
