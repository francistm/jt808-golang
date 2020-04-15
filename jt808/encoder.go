package jt808

import (
	"bytes"
)

func Marshal(ptr *MessagePack) ([]byte, error) {
	var buf bytes.Buffer

	return buf.Bytes(), nil
}
