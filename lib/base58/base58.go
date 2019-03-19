package base58

func Encode(src []byte) string {
	return string(encode(src))
}

func Decode(value string) ([]byte, error) {
	return decode([]byte(value))
}
