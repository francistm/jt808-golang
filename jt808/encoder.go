package jt808

import (
	"bytes"
)

// Marshal 编译一个消息体到字节数组
func Marshal(ptr *MessagePack) ([]byte, error) {
	var buf bytes.Buffer
	var bodyBytes []byte

	b, err := ptr.PackBody.marshalBody()

	if err != nil {
		return nil, err
	}

	bodyBytes = b

	ptr.PackHeader.Property.BodyByteLength = uint16(len(bodyBytes))

	if b, err := marshalHeader(&ptr.PackHeader); err != nil {
		return nil, err
	} else if _, err := buf.Write(b); err != nil {
		return nil, err
	}

	if _, err := buf.Write(bodyBytes); err != nil {
		return nil, err
	}

	if checksum, err := computeChecksum(buf.Bytes()); err != nil {
		return nil, err
	} else if err := buf.WriteByte(checksum); err != nil {
		return nil, err
	}

	escapedBytes, err := escapeChars(buf.Bytes())

	if err != nil {
		return nil, err
	}

	var finalBuf bytes.Buffer

	finalBuf.WriteByte(0x7e)

	if _, err := finalBuf.Write(escapedBytes); err != nil {
		return nil, err
	}

	finalBuf.WriteByte(0x7e)

	return finalBuf.Bytes(), nil
}
