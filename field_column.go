package do

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// ObjectAndFieldsHelper if *T is interface{ ValuePtrs() []any } then use its results as fields, otherwise use reflect to get fields by ColumnType
func ObjectAndFieldsHelper[T any](fieldMappers ...func(string) string) func(colTypes []*sql.ColumnType) (r *T, fields []any) {
	var fieldMapper func(string) string
	if len(fieldMappers) != 0 {
		fieldMapper = fieldMappers[0]
	}

	return func(colTypes []*sql.ColumnType) (r *T, fields []any) {
		r = new(T)
		switch vi := any(r).(type) {
		case interface{ ValuePtrs() []any }:
			fields = append(fields, vi.ValuePtrs()...)

		// Scan 将从数据库读取的列转换为以下常见的 Go 类型和 sql 包提供的特殊类型：
		//
		// *string
		// *[]byte
		// *int, *int8, *int16, *int32, *int64
		// *uint, *uint8, *uint16, *uint32, *uint64
		// *bool
		// *float32, *float64
		// *interface{}
		// *RawBytes
		// *Rows (cursor value)
		// any type implementing Scanner (see Scanner docs)
		//
		case *string,
			*[]byte,
			*bool,
			*float32, *float64,
			*int, *int8, *int16, *int32, *int64,
			*uint, *uint8, *uint16, *uint32, *uint64:
			fields = append(fields, r)

		default:
			fields = FieldsByColumnType(r, colTypes, fieldMapper)
		}
		return
	}
}

// FieldsByColumnType t is a struct pointer, and use it's field match column name to receive scan value. It will use db tag to get column name first, or lower case field name. You can specify fieldMapper to control column name with field name
func FieldsByColumnType(t any, colTypes []*sql.ColumnType, fieldMapper func(string) string) (fields []any) {
	validName := make(map[string]struct{})
	for _, ct := range colTypes {
		validName[ct.Name()] = struct{}{}
	}

	nameValues := fieldsByColumnName(t, validName, fieldMapper)

	for _, ct := range colTypes {
		fields = append(fields, nameValues[ct.Name()])
	}

	return
}

func fieldsByColumnName(t any, validName map[string]struct{}, fieldMapper func(string) string) (nameValues map[string]any) {
	val := reflect.ValueOf(t)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr {
		panic(fmt.Errorf("t must be a struct pointer"))
	}
	val = val.Elem()
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("t must be a struct pointer"))
	}

	return fieldsByColumnNameInner(typ, val, validName, fieldMapper)
}

func fieldsByColumnNameInner(typ reflect.Type, val reflect.Value, validName map[string]struct{}, fieldMapper func(string) string) (nameValues map[string]any) {
	nameValues = make(map[string]any)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if field.Anonymous {
			nv := fieldsByColumnNameInner(value.Type(), value, validName, fieldMapper)
			nameValues = MergeKeyValue(nameValues, nv)
		} else {
			fieldName := ""
			if fieldMapper == nil {
				fieldName = field.Tag.Get("db")
				if fieldName == "" {
					fieldName = strings.ToLower(field.Name)
				}
			} else {
				fieldName = fieldMapper(field.Name)
			}

			if _, ok := validName[fieldName]; !ok {
				continue
			}

			nameValues[fieldName] = value.Addr().Interface()
		}
	}

	return
}
