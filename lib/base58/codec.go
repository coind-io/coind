package base58

import (
	"math/big"
)

func decode(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return src, nil
	}
	// leading zero bytes
	zeros := int(0)
	for i := 0; i < len(src); i++ {
		if src[i] != ENCODED_ZERO {
			break
		}
		zeros = zeros + 1
	}
	// decode values
	number := new(big.Int)
	radix := big.NewInt(58)
	for i := 0; i < len(src); i++ {
		b := INDEXES[src[i]]
		if b == 0xFF {
			return nil, CorruptInputError(i)
		}
		number.Mul(number, radix)
		number.Add(number, big.NewInt(int64(b)))
	}
	// handel zeros
	value := number.Bytes()
	dest := make([]byte, zeros, zeros+len(value))
	dest = append(dest, value...)
	return dest, nil
}

func encode(src []byte) []byte {
	if len(src) == 0 {
		return src
	}
	// bas58 zeros prefix handle
	preifx := make([]byte, 0)
	for i := 0; i < len(src); i++ {
		if src[i] != 0 {
			break
		}
		preifx = append(preifx, ENCODED_ZERO)
	}
	// encode value
	number := big.NewInt(0).SetBytes(src)
	radix := big.NewInt(58)
	zero := big.NewInt(0)
	dest := make([]byte, 0)
	for number.Cmp(zero) > 0 {
		mod := big.NewInt(0)
		mod = mod.Mod(number, radix)
		number = number.Div(number, radix)
		dest = append(dest, ALPHABET[mod.Int64()])
	}
	dest = reverse(dest)
	dest = append(preifx, dest...)
	return dest
}

func reverse(src []byte) []byte {
	dest := make([]byte, 0, len(src))
	for i := len(src) - 1; i >= 0; i-- {
		dest = append(dest, src[i])
	}
	return dest
}
