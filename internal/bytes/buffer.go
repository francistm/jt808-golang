package bytes

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

func (b *Buffer) WriteBCD(s string) error {
	data, err := hex.DecodeString(s)

	if err != nil {
		return err
	}

	if _, err := b.Write(data); err != nil {
		return err
	}

	return nil
}
