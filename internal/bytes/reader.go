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

func (r *Reader) ReadBytes(len int) ([]byte, error) {
	buf := make([]byte, len)

	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *Reader) ReadBCD(len int) (string, error) {
	buf := make([]byte, len)

	if _, err := r.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}
