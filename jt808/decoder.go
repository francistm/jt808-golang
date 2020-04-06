package jt808

import (
	"bytes"
	"errors"
	"reflect"
)

// 由二进制解析一个完整的消息包
func Unmarshal(buf []byte, v *MessagePack) error {
	reader := bytes.NewReader(buf)

	// read first byte as prefix
	if b, err := reader.ReadByte(); err != nil {
		return err
	} else if b != 0x7e {
		return errors.New("invalid prefix byte")
	}

	// read header, ( 12 or 14 byte depends on is multiple package message)
	headerBuf := make([]byte, 14)

	if _, err := reader.Read(headerBuf); err != nil {
		return err
	}

	if err := UnmarshalHeader(headerBuf, &v.PackHeader); err != nil {
		return err
	}

	// is not a multiple package, reverse reader 2 bytes back because there's no package bytes
	if !v.PackHeader.Property.IsMultiplePackage {
		_ = reader.UnreadByte()
		_ = reader.UnreadByte()
		// move read pointer 2 bytes back if there's no package property
	}

	// read bytes according header body data length
	bodyBuf := make([]byte, v.PackHeader.Property.BodyByteLength)

	if _, err := reader.Read(bodyBuf); err != nil {
		return err
	}

	// update PackBody field from readed bytes to struct
	if err := unmarshalPackBody(v.PackHeader.MessageId, bodyBuf, &v); err != nil {
		return err
	}

	// update checksum in message pack, but didn't valid it
	if b, err := reader.ReadByte(); err != nil {
		return err
	} else {
		v.Checksum = b
	}

	// read the last byte as suffix in whole message
	if b, err := reader.ReadByte(); err != nil {
		return err
	} else if b != 0x7e {
		return errors.New("invalid suffix byte")
	}

	return nil
}

func unmarshalPackBody(messageId uint16, buf []byte, ptr interface{}) (err error) {
	var unmarshalPackBody interface{}

	switch messageId {
	case 0x0001:
		unmarshalPackBody, err = unmarshalBody0001(buf)
	default:
		err = errors.New("unsupported messageId")
	}

	if err == nil {
		v := reflect.ValueOf(ptr)

		if !v.IsValid() {
			return errors.New("invalid target when unmarshal PackBody")
		}

		for v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}

		field := v.FieldByName("PackBody")

		if !field.CanSet() {
			return errors.New("target doesn't have PackBody field")
		}

		field.Set(reflect.ValueOf(unmarshalPackBody))
	}

	return
}
