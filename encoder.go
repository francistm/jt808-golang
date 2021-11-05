package jt808

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// Marshal 编译一个消息体到字节数组
func Marshal(ptr *MessagePack) ([]byte, error) {
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

func marshalBody(writer io.Writer, packBody interface{}) error {
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

		var writerErr error

		if err != nil {
			return fmt.Errorf("cannot parse tag of field %s.%s", refMesgBodyType.Name(), fieldType.Name)
		}

		switch fieldValue.Kind() {
		case reflect.Uint8:
			writerErr = writeUint8(fieldValue.Interface().(uint8), writer)

		case reflect.Uint16:
			writerErr = writeUint16(fieldValue.Interface().(uint16), writer)

		case reflect.Uint32:
			writerErr = writeUint32(fieldValue.Interface().(uint32), writer)

		case reflect.Slice:
			if tag.fieldDataEncoding == tagEncodingNone {
				_, writerErr = writer.Write(fieldValue.Interface().([]byte))
			} else {
				return fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.fieldDataEncoding)
			}

		case reflect.String:
			if tag.fieldDataEncoding == tagEncodingBCD {
				writerErr = writeBCD(fieldValue.Interface().(string), writer)
			} else {
				return fmt.Errorf("unknown field %s.%s encoding: %s", refMesgBodyType.Name(), fieldType.Name, tag.fieldDataEncoding)
			}

		case reflect.Struct:
			writerErr = marshalBody(writer, fieldValue.Addr().Interface())

		case reflect.Ptr:
			writerErr = marshalBody(writer, fieldValue.Interface())
		}

		if writerErr != nil {
			return writerErr
		}
	}

	return nil
}

func (messagePack *MessagePack) marshalBody() ([]byte, error) {
	packBody := messagePack.PackBody
	bodyBytesWriter := new(bytes.Buffer)

	if err := marshalBody(bodyBytesWriter, packBody); err != nil {
		return nil, err
	}

	return bodyBytesWriter.Bytes(), nil
}
