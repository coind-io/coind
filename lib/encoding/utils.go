package encoding

import (
	"bytes"
)

func Marshal(ie IEncoder) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := ie.Encode(buf)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(raw []byte, id IDecoder) error {
	buf := bytes.NewBuffer(raw)
	err := id.Decode(buf)
	if err != nil {
		return err
	}
	return nil
}
