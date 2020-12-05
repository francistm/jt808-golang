package jt808

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/francistm/jt808-golang/jt808/message"
)

// Unmarshal 由二进制解析一个完整的消息包
func Unmarshal(buf []byte, messagePack *MessagePack) error {
	var checksum byte

	if buf[0] != 0x7e {
		return fmt.Errorf("invalid prefix byte 0x%.2X", buf[0])
	}

	if buf[len(buf)-1] != 0x7e {
		return fmt.Errorf("invalid suffix byte 0x%.2X", buf[0])
	}

	buf = buf[1 : len(buf)-1]

	buf, err := unescapeChars(buf)

	if err != nil {
		return err
	}

	c, err := computeChecksum(buf[0 : len(buf)-1])

	if err != nil {
		return err
	}

	checksum = c

	reader := bytes.NewReader(buf)

	// read header, ( 12 or 12 + 4 bytes depends on is multiple package message)
	headerBuf := make([]byte, 16)

	if _, err := reader.Read(headerBuf); err != nil {
		return err
	}

	if err := UnmarshalHeader(headerBuf, &messagePack.PackHeader); err != nil {
		return err
	}

	// is not a multiple package, reverse reader 4 bytes back because there's no package bytes
	if !messagePack.PackHeader.Property.IsMultiplePackage {
		for i := 0; i < 4; i++ {
			_ = reader.UnreadByte()
		}
	}

	// read bytes according header body data length
	bodyBuf := make([]byte, messagePack.PackHeader.Property.BodyByteLength)

	if _, err := reader.Read(bodyBuf); err != nil {
		return err
	}

	// update PackBody field from readed bytes to struct
	if err := messagePack.unmarshalBody(bodyBuf); err != nil {
		return err
	}

	// update checksum in message pack
	bs, err := reader.ReadByte()

	if err != nil {
		return err
	}

	messagePack.Checksum = bs
	messagePack.ChecksumValid = bs == checksum

	return nil
}

func unmarshalBody(reader io.Reader, packBody interface{}) error {
	refMesgBodyType := reflect.TypeOf(packBody).Elem()
	refMesgBodyValue := reflect.ValueOf(packBody).Elem()

	for i := 0; i < refMesgBodyValue.NumField(); i++ {
		fieldType := refMesgBodyType.Field(i)
		fieldValue := refMesgBodyValue.Field(i)

		rawTag, hasTag := fieldType.Tag.Lookup(tagName)

		// embed struct field is kind of struct
		if fieldValue.Kind() != reflect.Struct && !hasTag {
			continue
		}

		tag, err := parseTag(rawTag)

		if err != nil {
			return fmt.Errorf("cannot parse tag of field %s.%s", refMesgBodyType.Name(), fieldType.Name)
		}

		var readerErr error
		var readerValue interface{}

		switch fieldValue.Kind() {
		case reflect.Uint8:
			readerValue, readerErr = readUint8(reader)

		case reflect.Uint16:
			readerValue, readerErr = readUint16(reader)

		case reflect.Uint32:
			readerValue, readerErr = readUint32(reader)

		case reflect.Ptr:
			structType := fieldValue.Type().Elem()
			structPtr := reflect.New(structType).Interface()
			readerValue = structPtr
			readerErr = unmarshalBody(reader, structPtr)

		case reflect.Slice:
			if tag.fieldDataEncoding == tagEncodingNone {
				readerValue, readerErr = ioutil.ReadAll(reader)
			} else {
				return fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.fieldDataEncoding)
			}

		case reflect.String:
			if tag.fieldDataLength < 1 {
				return fmt.Errorf("field %s.%s with string must set byte length", refMesgBodyType.Name(), fieldType.Name)
			}

			if tag.fieldDataEncoding == tagEncodingBCD {
				readerValue, readerErr = readBCD(reader, tag.fieldDataLength)
			} else if tag.fieldDataEncoding == tagEncodingNone {
				readerValue, readerErr = readBytes(reader, tag.fieldDataLength)
			} else {
				return fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.fieldDataEncoding)
			}

		case reflect.Struct:
			structType := fieldValue.Type()
			structPtr := reflect.New(structType).Interface()
			readerValue = structPtr
			readerErr = unmarshalBody(reader, structPtr)
		}

		if readerErr != nil {
			return readerErr
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set %s.%s field value", refMesgBodyType.Name(), fieldType.Name)
		}

		refReaderValue := reflect.ValueOf(readerValue)

		switch refReaderValue.Kind() {
		case reflect.Ptr:
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue.Set(reflect.ValueOf(refReaderValue.Interface()))
			} else {
				fieldValue.Set(reflect.ValueOf(refReaderValue.Elem().Interface()))
			}

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.String, reflect.Slice:
			fieldValue.Set(reflect.ValueOf(readerValue))

		default:
		}
	}

	return nil
}

func (messagePack *MessagePack) unmarshalBody(buf []byte) error {
	reader := bytes.NewReader(buf)

	switch messagePack.PackHeader.MessageID {
	case 0x0001:
		messagePack.PackBody = new(message.Body0001)
	case 0x0200:
		messagePack.PackBody = new(message.Body0200)
	case 0x0801:
		messagePack.PackBody = new(message.Body0801)
	default:
		return fmt.Errorf("unsupported messageId: 0x%.4X", messagePack.PackHeader.MessageID)
	}

	return unmarshalBody(reader, messagePack.PackBody)
}
