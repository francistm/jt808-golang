package jt808

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

//go:generate go run github.com/francistm/jt808-golang/tools/generator/decoder

// Marshal 编译一个消息体到字节数组
func Marshal[T any](ptr *MessagePack[T]) ([]byte, error) {
	var buf bytes.Buffer

	bodyBytes, err := ptr.marshalBody()

	if err != nil {
		return nil, err
	}

	ptr.PackHeader.Property.BodyByteLength = uint16(len(bodyBytes))

	if b, err := marshalHeader(&ptr.PackHeader); err != nil {
		return nil, err
	} else if _, err := buf.Write(b); err != nil {
		return nil, err
	}

	if _, err := buf.Write(bodyBytes); err != nil {
		return nil, err
	}

	if checksum, err := calculateChecksum(buf.Bytes()); err != nil {
		return nil, err
	} else if err := buf.WriteByte(checksum); err != nil {
		return nil, err
	}

	escapedBytes, err := encodeBytes(buf.Bytes())

	if err != nil {
		return nil, err
	}

	var finalBuf bytes.Buffer

	finalBuf.WriteByte(identifyByte)

	if _, err := finalBuf.Write(escapedBytes); err != nil {
		return nil, err
	}

	finalBuf.WriteByte(identifyByte)

	return finalBuf.Bytes(), nil
}

func marshalBody[T any](writer io.Writer, packBody T) error {
	refMesgBodyType := reflect.TypeOf(packBody)
	refMesgBodyValue := reflect.ValueOf(packBody)

	if refMesgBodyType.Kind() == reflect.Ptr {
		refMesgBodyType = refMesgBodyType.Elem()
		refMesgBodyValue = refMesgBodyValue.Elem()
	}

	for i := 0; i < refMesgBodyValue.NumField(); i++ {
		fieldType := refMesgBodyType.Field(i)
		fieldValue := refMesgBodyValue.Field(i)

		rawTag, hasTag := fieldType.Tag.Lookup(tagName)

		// embed struct field is kind of struct
		if fieldValue.Kind() != reflect.Struct && !hasTag {
			continue
		}

		tag, err := parseMesgTag(rawTag)

		if err != nil {
			return fmt.Errorf("cannot parse tag of field %s.%s", refMesgBodyType.Name(), fieldType.Name)
		}

		switch {
		case fieldType.Type.Kind() == reflect.Uint8:
			err = writeUint8(uint8(fieldValue.Uint()), writer)

		case fieldType.Type.Kind() == reflect.Uint16:
			err = writeUint16(uint16(fieldValue.Uint()), writer)

		case fieldType.Type.Kind() == reflect.Uint32:
			err = writeUint32(uint32(fieldValue.Uint()), writer)

		case fieldType.Type.Kind() == reflect.Slice && fieldType.Type.Elem().Kind() == reflect.Uint8:
			if tag.dataEncoding == tagEncodingNone {
				_, err = writer.Write(fieldValue.Bytes())
			} else {
				err = fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.dataEncoding)
			}

		case fieldType.Type.Kind() == reflect.String:
			if tag.dataEncoding == tagEncodingBCD {
				err = writeBCD(fieldValue.String(), writer)
			} else {
				err = fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.dataEncoding)
			}

		case fieldType.Type.Kind() == reflect.Struct:
			err = marshalBody(writer, fieldValue.Interface())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (messagePack *MessagePack[T]) marshalBody() ([]byte, error) {
	bodyBytesWriter := new(bytes.Buffer)

	if err := marshalBody(bodyBytesWriter, messagePack.PackBody); err != nil {
		return nil, err
	}

	return bodyBytesWriter.Bytes(), nil
}
