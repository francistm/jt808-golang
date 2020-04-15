package jt808

import (
	"bytes"
	"errors"
	"fmt"
)

// 由二进制解析一个完整的消息包
func Unmarshal(buf []byte, v *MessagePack) error {
	var checksum byte

	if buf[0] != 0x7e {
		return errors.New(fmt.Sprintf("invalid prefix byte 0x%.2X", buf[0]))
	}

	if buf[len(buf)-1] != 0x7e {
		return errors.New(fmt.Sprintf("invalid suffix byte 0x%.2X", buf[0]))
	}

	buf = buf[1 : len(buf)-1]

	if b, err := unescapeChars(buf); err != nil {
		return err
	} else {
		buf = b
	}

	if c, err := computeChecksum(buf[0 : len(buf)-1]); err != nil {
		return err
	} else {
		checksum = c
	}

	reader := bytes.NewReader(buf)

	// read header, ( 12 or 12 + 4 bytes depends on is multiple package message)
	headerBuf := make([]byte, 16)

	if _, err := reader.Read(headerBuf); err != nil {
		return err
	}

	if err := UnmarshalHeader(headerBuf, &v.PackHeader); err != nil {
		return err
	}

	// is not a multiple package, reverse reader 4 bytes back because there's no package bytes
	if !v.PackHeader.Property.IsMultiplePackage {
		for i := 0; i < 4; i++ {
			_ = reader.UnreadByte()
		}
	}

	// read bytes according header body data length
	bodyBuf := make([]byte, v.PackHeader.Property.BodyByteLength)

	if _, err := reader.Read(bodyBuf); err != nil {
		return err
	}

	// update PackBody field from readed bytes to struct
	if err := unmarshalBody(bodyBuf, v); err != nil {
		return err
	}

	// update checksum in message pack
	if b, err := reader.ReadByte(); err != nil {
		return err
	} else {
		v.Checksum = b
		v.ChecksumValid = b == checksum
	}

	return nil
}

func unmarshalBody(buf []byte, ptr *MessagePack) error {
	var unmarshalFunc func([]byte) (PackBody, error)

	switch ptr.PackHeader.MessageId {
	case 0x0001:
		unmarshalFunc = func(b []byte) (PackBody, error) {
			return unmarshalBody0001(b)
		}
	case 0x0200:
		unmarshalFunc = func(b []byte) (PackBody, error) {
			return unmarshalBody0200(b)
		}
	case 0x0801:
		unmarshalFunc = func(b []byte) (PackBody, error) {
			return unmarshalBody0801(b)
		}
	default:
		return errors.New(fmt.Sprintf("unsupported messageId: 0x%.4X", ptr.PackHeader.MessageId))
	}

	if unmarshalFunc == nil {
		return errors.New(fmt.Sprintf("missing unmarshal function for messageId: 0x%.4X", ptr.PackHeader.MessageId))
	}

	// if this's a multiple package, dont' unmarshal it at this moment.
	// store the body bytes, unmarshal function to messagePack struct,
	// and wait until `func (*MessagePack) ConcatAndUnmarshal` been called
	if ptr.PackHeader.Property.IsMultiplePackage {
		ptr.bodyBuf = buf
		ptr.unmarshalFunc = unmarshalFunc

		return nil
	}

	packBody, err := unmarshalFunc(buf)

	if err != nil {
		return err
	}

	ptr.PackBody = packBody

	return nil
}
