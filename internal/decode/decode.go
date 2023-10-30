package decode

import (
	"encoding"
	"fmt"
	"io"
	"reflect"

	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/tag"
)

type RefField struct {
	StructName string
	FieldName  string
	Tag        string
	TypeRef    reflect.Type
	ValueRef   reflect.Value
}

func UnmarshalStruct(reader *bytes.Reader, target any) error {
	if unmarshaller, ok := target.(encoding.BinaryUnmarshaler); ok {
		data, err := io.ReadAll(reader)

		if err != nil {
			return err
		}

		return unmarshaller.UnmarshalBinary(data)
	}

	refFields := make([]*RefField, 0, 50)
	refFields = append(refFields, &RefField{
		TypeRef:  reflect.TypeOf(target),
		ValueRef: reflect.ValueOf(target),
	})

	for len(refFields) > 0 {
		head := refFields[0]
		refFields = refFields[1:]

		var (
			rawData   any
			err       error
			parsedTag *tag.MesgTag
		)

		switch {
		case head.TypeRef.Kind() == reflect.Uint8:
			rawData, err = reader.ReadUint8()

		case head.TypeRef.Kind() == reflect.Uint16:
			rawData, err = reader.ReadUint16()

		case head.TypeRef.Kind() == reflect.Uint32:
			rawData, err = reader.ReadUint32()

		case head.TypeRef.Kind() == reflect.Pointer:
			if head.ValueRef.IsNil() {
				head.ValueRef.Set(reflect.New(head.TypeRef.Elem()))
			}

			fields := make([]*RefField, 0, len(refFields)+1)
			fields = append(fields, &RefField{
				StructName: head.StructName,
				FieldName:  head.FieldName,
				Tag:        head.Tag,
				TypeRef:    head.TypeRef.Elem(),
				ValueRef:   head.ValueRef.Elem(),
			})
			fields = append(fields, refFields...)
			refFields = fields

		case head.TypeRef.Kind() == reflect.Slice && head.TypeRef.Elem().Kind() == reflect.Uint8:
			parsedTag, err = tag.NewMesgTag(head.Tag)

			if err != nil {
				return fmt.Errorf("cannot parse tag of field %s.%s", head.StructName, head.FieldName)
			}

			if parsedTag.Length == -1 {
				rawData, err = io.ReadAll(reader)
			} else if parsedTag.Length > 0 {
				rawData, err = reader.ReadFixedBytes(parsedTag.Length)
			} else {
				return fmt.Errorf("missing byte length for field %s.%s", head.StructName, head.FieldName)
			}

		case head.TypeRef.Kind() == reflect.String:
			parsedTag, err = tag.NewMesgTag(head.Tag)

			if err != nil {
				return fmt.Errorf("cannot parse tag of field %s.%s", head.StructName, head.FieldName)
			}

			if parsedTag.Length < 1 {
				return fmt.Errorf("missing byte length for field %s.%s", head.StructName, head.FieldName)
			}

			rawData, err = reader.ReadFixedString(parsedTag.Length)

		case head.TypeRef.Kind() == reflect.Struct:
			fields := refFields
			structFields := make([]*RefField, 0, head.TypeRef.NumField())

			for i := 0; i < head.TypeRef.NumField(); i++ {
				structField := head.TypeRef.Field(i)
				fieldValueRef := head.ValueRef.Field(i)

				if structField.PkgPath != "" {
					fmt.Printf("%s.%s pkgPath %s\n", head.TypeRef.Name(), structField.Name, structField.PkgPath)
					continue
				}

				structFields = append(structFields, &RefField{
					StructName: head.TypeRef.Name(),
					FieldName:  structField.Name,
					Tag:        structField.Tag.Get(tag.Name),
					TypeRef:    structField.Type,
					ValueRef:   fieldValueRef,
				})
			}

			refFields = make([]*RefField, 0, len(fields)+len(structFields))
			refFields = append(refFields, structFields...)
			refFields = append(refFields, fields...)

		default:
			return fmt.Errorf("unsupport field type %s", head.TypeRef)
		}

		if err != nil {
			return err
		}

		rawDataTypeRef := reflect.TypeOf(rawData)
		rawDataValueRef := reflect.ValueOf(rawData)

		if rawData != nil && rawDataTypeRef.Kind() != head.TypeRef.Kind() {
			return fmt.Errorf("can't set %s to field %s.%s, want %s", rawDataTypeRef.Kind(), head.StructName, head.FieldName, head.TypeRef.Kind())
		}

		if rawData != nil {
			head.ValueRef.Set(rawDataValueRef)
		}
	}

	return nil
}
