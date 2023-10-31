package bytes

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (b *Buffer) WriteUint8(i uint8) error {
	return b.WriteByte(i)
}

func (b *Buffer) WriteUint16(i uint16) error {
	data := make([]byte, 2)

	binary.BigEndian.PutUint16(data, i)

	if _, err := b.Write(data); err != nil {
		return err
	}

	return nil
}

func (b *Buffer) WriteUint32(i uint32) error {
	data := make([]byte, 4)

	binary.BigEndian.PutUint32(data, i)

	if _, err := b.Write(data); err != nil {
		return err
	}

	return nil
}

func (b *Buffer) WriteBCD(s string, size int) error {
	if size%2 != 0 {
		return errors.New("size must be even")
	}

	if len(s)%2 != 0 && len(s) < size*2 {
		s = "0" + s
	} else if len(s)%2 != 0 && len(s) > size*2 {
		s = s[:size*2]
	}

	data, err := hex.DecodeString(s)

	if err != nil {
		return err
	}

	if len(data) < size {
		data = append(make([]byte, size-len(data)), data...)
	} else if len(data) > size {
		data = data[:size]
	}

	if _, err := b.Write(data); err != nil {
		return err
	}

	return nil
}

func (b *Buffer) WriteFixedString(s string, size int) {
	if len(s) < size {
		s = s + strings.Repeat("\x00", size-len(s))
	} else if len(s) > size {
		s = s[:size]
	}

	b.WriteString(s)
}
