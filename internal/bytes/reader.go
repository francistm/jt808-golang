package bytes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Reader struct {
	*bytes.Reader
}

func NewReader(in []byte) *Reader {
	return &Reader{
		Reader: bytes.NewReader(in),
	}
}

func (r *Reader) ReadUint8() (uint8, error) {
	buf := make([]byte, 1)

	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

func (r *Reader) ReadUint16() (uint16, error) {
	buf := make([]byte, 2)

	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(buf), nil
}

func (r *Reader) ReadUint32() (uint32, error) {
	buf := make([]byte, 4)

	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(buf), nil
}

func (r *Reader) ReadFixedBytes(size int) ([]byte, error) {
	buf := make([]byte, size)

	n, err := r.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[0:n], nil
}

func (r *Reader) ReadBCD(size int) (string, error) {
	buf := make([]byte, size)

	if _, err := r.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

func (r *Reader) ReadFixedString(size int) (string, error) {
	buf := make([]byte, size)
	readSize, err := r.Read(buf)

	if err != nil {
		return "", err
	}

	return string(buf[0:readSize]), nil
}

func (r *Reader) UnreadFixedBytes(size int) error {
	for i := 0; i < size; i++ {
		if err := r.UnreadByte(); err != nil {
			return err
		}
	}

	return nil
}
