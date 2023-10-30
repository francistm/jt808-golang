package decode

import (
	"encoding"
	"fmt"
	"io"
	"reflect"

	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/tag"
)

func UnmarshalStruct(reader *bytes.Reader, target any) error {
	if unmarshaller, ok := target.(encoding.BinaryUnmarshaler); ok {
		data, err := io.ReadAll(reader)

		if err != nil {
			return err
		}

		return unmarshaller.UnmarshalBinary(data)
	}

	mesgBodyTypeRef := reflect.TypeOf(target).Elem()
	mesgBodyValueRef := reflect.ValueOf(target).Elem()

	for i := 0; i < mesgBodyValueRef.NumField(); i++ {
		structField := mesgBodyTypeRef.Field(i)
		structFieldValueRef := mesgBodyValueRef.Field(i)

		if structField.PkgPath != "" {
			continue
		}

		rawTag := structField.Tag.Get(tag.Name)
		parsedTag, err := tag.NewMesgTag(rawTag)

		if err != nil {
			return fmt.Errorf("cannot parse tag of field %s.%s", mesgBodyTypeRef.Name(), structField.Name)
		}

		var readerValue any

		switch {
		case structField.Type.Kind() == reflect.Uint8:
			readerValue, err = reader.ReadUint8()

		case structField.Type.Kind() == reflect.Uint16:
			readerValue, err = reader.ReadUint16()

		case structField.Type.Kind() == reflect.Uint32:
			readerValue, err = reader.ReadUint32()

		case structField.Type.Kind() == reflect.Ptr:
			structType := structFieldValueRef.Type().Elem()
			structPtr := reflect.New(structType).Interface()
			readerValue = structPtr
			err = UnmarshalStruct(reader, structPtr)

		case structField.Type.Kind() == reflect.Slice && structField.Type.Elem().Kind() == reflect.Uint8:
			if parsedTag.Encoding != tag.EncodingRaw {
				return fmt.Errorf("unsupport encoding %s for field %s.%s", parsedTag.Encoding, mesgBodyTypeRef.Name(), structField.Name)
			}

			readerValue, err = io.ReadAll(reader)

		case structField.Type.Kind() == reflect.String:
			if parsedTag.Length < 1 {
				return fmt.Errorf("missing byte length for field %s.%s", mesgBodyTypeRef.Name(), structField.Name)
			}

			switch parsedTag.Encoding {
			case tag.EncodingBCD:
				readerValue, err = reader.ReadBCD(parsedTag.Length)

			case tag.EncodingRaw:
				readerValue, err = reader.ReadFixedBytes(parsedTag.Length)

			default:
				return fmt.Errorf("unsupport encoding %s for field %s.%s", parsedTag.Encoding, mesgBodyTypeRef.Name(), structField.Name)
			}

		case structField.Type.Kind() == reflect.Struct:
			readerValue = reflect.New(structField.Type).Interface()
			err = UnmarshalStruct(reader, readerValue)
		}

		if err != nil {
			return err
		}

		if !structFieldValueRef.CanSet() {
			return fmt.Errorf("cannot set %s.%s field value", mesgBodyTypeRef.Name(), structField.Name)
		}

		refReaderValue := reflect.ValueOf(readerValue)

		switch refReaderValue.Kind() {
		case reflect.Ptr:
			if structFieldValueRef.Kind() == reflect.Ptr {
				structFieldValueRef.Set(reflect.ValueOf(refReaderValue.Interface()))
			} else {
				structFieldValueRef.Set(reflect.ValueOf(refReaderValue.Elem().Interface()))
			}

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.String, reflect.Slice:
			structFieldValueRef.Set(reflect.ValueOf(readerValue))
		}
	}

	return nil
}
