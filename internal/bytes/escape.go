package bytes

import (
	"bytes"
	"errors"
	"io"
)

// 消息发送前对消息中 0x7e, 0x7d 进行转义
// 需先将消息体进行转义，然后在首位增加 0x7e 的标识位字节
func Escape(buf []byte) ([]byte, error) {
	var writer bytes.Buffer
	var reader = bytes.NewReader(buf)

	for {
		b, err := reader.ReadByte()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
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
func Unescape(buf []byte) ([]byte, error) {
	var writer bytes.Buffer
	var reader = bytes.NewReader(buf)

	for {
		b, err := reader.ReadByte()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		if b != 0x7d {
			err := writer.WriteByte(b)

			if err != nil {
				return nil, err
			}

			continue
		}

		nextByte, err := reader.ReadByte()

		if err != nil {
			return nil, err
		}

		switch nextByte {
		case 0x01:
			writer.WriteByte(0x7d)
		case 0x02:
			writer.WriteByte(0x7e)
		default:
			return nil, errors.New("invalid char after 0x7e when unescape")
		}
	}

	return writer.Bytes(), nil
}
