package encoding

import (
	"errors"
	"io"
)

var ErrEof = errors.New("coind/lib/encoding: got EOF, can not get the next byte")

type IEncoder interface {
	Encode(w io.Writer) error
}

type IDecoder interface {
	Decode(r io.Reader) error
}
