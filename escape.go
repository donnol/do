package do

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"reflect"
)

// EscapeStruct escape struct field value with escaper, v must be a struct pointer
func EscapeStruct(v any, escaper Escaper) (err error) {
	val := reflect.ValueOf(v)

	typ := val.Type()
	if typ.Kind() != reflect.Pointer {
		return fmt.Errorf("v is not a pointer")
	}

	rval := val.Elem()
	rtyp := typ.Elem()
	if rtyp.Kind() != reflect.Struct {
		return fmt.Errorf("v is not a struct")
	}

	return escapeStruct(rtyp, rval, escaper)
}

func escapeStruct(rtyp reflect.Type, rval reflect.Value, escaper Escaper) (err error) {
	for i := 0; i < rtyp.NumField(); i++ {
		fieldTyp := rtyp.Field(i)
		fieldVal := rval.Field(i)

		if fieldTyp.Anonymous {
			if err = escapeStruct(fieldTyp.Type, fieldVal, escaper); err != nil {
				return
			}
			continue
		}

		if fieldTyp.Type.Kind() == reflect.Slice || fieldTyp.Type.Kind() == reflect.Array {
			for j := 0; j < fieldVal.Len(); j++ {
				if err = escapeStruct(fieldTyp.Type.Elem(), fieldVal.Index(j), escaper); err != nil {
					return
				}
			}
			continue
		}

		newVal, err := escaper(fieldVal.Interface())
		if err != nil {
			return fmt.Errorf("escape %+v(%+v) failed", fieldVal, fieldTyp)
		}
		if newVal != nil {
			fieldVal.Set(reflect.ValueOf(newVal))
		}
	}

	return
}

type Escaper func(field any) (any, error)

func XMLEscaper(field any) (r any, err error) {
	switch v := field.(type) {
	case string:
		buf := new(bytes.Buffer)
		err = xml.EscapeText(buf, []byte(v))
		if err != nil {
			return
		}
		return buf.String(), nil
	case []byte:
		buf := new(bytes.Buffer)
		err = xml.EscapeText(buf, []byte(v))
		if err != nil {
			return
		}
		return buf.Bytes(), nil
	}
	if err != nil {
		return
	}

	return
}

func HTMLEscaper(field any) (r any, err error) {
	switch v := field.(type) {
	case string:
		return html.EscapeString(v), nil
	case []byte:
		return StringToBytes(html.EscapeString(BytesToString(v))), nil
	}

	return
}
