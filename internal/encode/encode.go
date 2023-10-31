package encode

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/tag"
)

func MarshalStruct(buffer *bytes.Buffer, source any) error {
	if marshaller, ok := source.(encoding.BinaryMarshaler); ok {
		data, err := marshaller.MarshalBinary()

		if err != nil {
			return err
		}

		buffer.Write(data)
	}

	mesgBodyTypeRef := reflect.TypeOf(source)
	mesgBodyValueRef := reflect.ValueOf(source)

	if mesgBodyTypeRef.Kind() == reflect.Ptr {
		mesgBodyTypeRef = mesgBodyTypeRef.Elem()
		mesgBodyValueRef = mesgBodyValueRef.Elem()
	}

	for i := 0; i < mesgBodyValueRef.NumField(); i++ {
		fieldType := mesgBodyTypeRef.Field(i)
		fieldValue := mesgBodyValueRef.Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		rawTag := fieldType.Tag.Get(tag.Name)
		parsedTag, err := tag.NewMesgTag(rawTag)

		if err != nil {
			return fmt.Errorf("cannot parse tag of field %s.%s", mesgBodyTypeRef.Name(), fieldType.Name)
		}

		switch {
		case fieldType.Type.Kind() == reflect.Uint8:
			err = buffer.WriteUint8(uint8(fieldValue.Uint()))

		case fieldType.Type.Kind() == reflect.Uint16:
			err = buffer.WriteUint16(uint16(fieldValue.Uint()))

		case fieldType.Type.Kind() == reflect.Uint32:
			err = buffer.WriteUint32(uint32(fieldValue.Uint()))

		case fieldType.Type.Kind() == reflect.Slice && fieldType.Type.Elem().Kind() == reflect.Uint8:
			if parsedTag.Encoding != tag.EncodingRaw {
				return fmt.Errorf("unknown field %s.%s encoding: %s", mesgBodyTypeRef.Name(), fieldType.Name, parsedTag.Encoding)
			}

			_, err = buffer.Write(fieldValue.Bytes())

		case fieldType.Type.Kind() == reflect.String:
			if parsedTag.Encoding != tag.EncodingBCD {
				return fmt.Errorf("unknown field %s.%s encoding: %s", mesgBodyTypeRef.Name(), fieldType.Name, parsedTag.Encoding)
			}

			err = buffer.WriteBCD(fieldValue.String(), parsedTag.Length)

		case fieldType.Type.Kind() == reflect.Struct:
			err = MarshalStruct(buffer, fieldValue.Interface())
		}

		if err != nil {
			return err
		}
	}

	return nil
}
