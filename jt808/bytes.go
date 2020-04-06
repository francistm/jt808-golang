package jt808

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func readBCD(reader io.Reader, byteLen int) (string, error) {
	buf := make([]byte, byteLen)

	if _, err := reader.Read(buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

func readUint8(reader io.Reader) (uint8, error) {
	buf := make([]byte, 1)

	if _, err := reader.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

func readUint16(reader io.Reader) (uint16, error) {
	buf := make([]byte, 2)

	if _, err := reader.Read(buf); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(buf), nil
}

// 消息发送前对消息中 0x7e, 0x7d 进行转义
// 需先将消息体进行转义，然后在首位增加 0x7e 的标识位字节
func escapeChars(buf []byte) ([]byte, error) {
	var writer bytes.Buffer
	var reader = bytes.NewReader(buf)

	for {
		b, err := reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		switch b {
		case 0x7d:
			_, err = writer.Write([]byte{0x7d, 0x01})
		case 0x7e:
			_, err = writer.Write([]byte{0x7d, 0x02})
		default:
			err = writer.WriteByte(b)
		}

		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}

// 消息收到后，对其中 0x7d01, 0x7d02 进行还原
// 需先去除首位的 0x7e 标识位字节后，再进行消息体转移
func unescapeChars(buf []byte) ([]byte, error) {
	var writer bytes.Buffer
	var reader = bytes.NewReader(buf)

	for {
		b, err := reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if b != 0x7d {
			err := writer.WriteByte(b)

			if err != nil {
				return nil, err
			}

			continue
		}

		if nextByte, err := reader.ReadByte(); err != nil {
			return nil, err
		} else {
			switch nextByte {
			case 0x01:
				writer.WriteByte(0x7d)
			case 0x02:
				writer.WriteByte(0x7e)
			default:
				return nil, errors.New("invalid char after 0x7e when unescape")
			}
		}
	}

	return writer.Bytes(), nil
}
