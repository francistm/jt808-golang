package bytes

import (
	"bytes"
	"errors"
	"io"
)

func CalcChecksum(buf []byte) (byte, error) {
	if len(buf) < 2 {
		return 0, errors.New("buf is less than 2 bytes")
	}

	reader := bytes.NewReader(buf)
	checksum, err := reader.ReadByte()

	if err != nil {
		return 0, err
	}

	for {
		var b byte
		var err error

		b, err = reader.ReadByte()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return 0, err
		}

		checksum ^= b
	}

	return checksum, nil
}
