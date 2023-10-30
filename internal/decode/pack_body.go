package decode

import (
	"encoding"
	"fmt"
	"io"
	"reflect"

	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/tag"
)

func UnmarshalPackBody(reader *bytes.Reader, target any) error {
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

		var readerValue any

		switch fieldValue.Kind() {
		case reflect.Uint8:
			readerValue, err = reader.ReadUint8()

		case reflect.Uint16:
			readerValue, err = reader.ReadUint16()

		case reflect.Uint32:
			readerValue, err = reader.ReadUint32()

		case reflect.Ptr:
			structType := fieldValue.Type().Elem()
			structPtr := reflect.New(structType).Interface()
			readerValue = structPtr
			err = UnmarshalPackBody(reader, structPtr)

		case reflect.Slice:
			if parsedTag.Encoding != tag.EncodingRaw {
				return fmt.Errorf("unknown field %s.%s encoding: %s", mesgBodyTypeRef.Name(), fieldType.Name, parsedTag.Encoding)
			}

			readerValue, err = io.ReadAll(reader)

		case reflect.String:
			if parsedTag.Length < 1 {
				return fmt.Errorf("field %s.%s with string must set byte length", mesgBodyTypeRef.Name(), fieldType.Name)
			}

			switch parsedTag.Encoding {
			case tag.EncodingBCD:
				readerValue, err = reader.ReadBCD(parsedTag.Length)

			case tag.EncodingRaw:
				readerValue, err = reader.ReadFixedBytes(parsedTag.Length)

			default:
				return fmt.Errorf("unknown field %s.%s encoding: %s", mesgBodyTypeRef.Name(), fieldType.Name, parsedTag.Encoding)
			}

		case reflect.Struct:
			structType := fieldValue.Type()
			structPtr := reflect.New(structType).Interface()
			readerValue = structPtr
			err = UnmarshalPackBody(reader, structPtr)
		}

		if err != nil {
			return err
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set %s.%s field value", mesgBodyTypeRef.Name(), fieldType.Name)
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
		}
	}

	return nil
}
