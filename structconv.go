package do

import (
	"fmt"
	"reflect"
)

// ConvSliceByName fill 'to' struct slice with 'from' by name
func ConvSliceByName[F, T any](from []F, to []T) {
	for i, f := range from {
		ConvByName(f, to[i])
	}
}

// ConvByName fill 'to' struct with 'from' by name
// like to.Name = from.Name
// to must be a struct pointer
func ConvByName[F, T any](from F, to T) {
	toByFieldNameReflect(from, to)
}

var (
	emptyValue = reflect.Value{}
)

func toByFieldNameReflect[F, T any](from F, to T) {
	fromValue := reflect.ValueOf(from)
	if fromValue.Type().Kind() == reflect.Pointer {
		fromValue = fromValue.Elem()
	}

	toValue := reflect.ValueOf(to)
	if toValue.Type().Kind() != reflect.Pointer {
		panic(fmt.Errorf("to is not a pointer"))
	}
	toElemValue := toValue.Elem()
	if toElemValue.Type().Kind() != reflect.Struct {
		panic(fmt.Errorf("to is not a struct"))
	}

	tobytoByFieldNameReflectValue(fromValue, toElemValue)
}

func tobytoByFieldNameReflectValue(fromValue, toElemValue reflect.Value) {
	for i := 0; i < toElemValue.Type().NumField(); i++ {
		field := toElemValue.Type().Field(i)
		fieldValue := toElemValue.Field(i)

		// 匿名
		if field.Anonymous {
			tobytoByFieldNameReflectValue(fromValue, fieldValue)
			continue
		}

		fromFieldValue := fromValue.FieldByName(field.Name)
		if fromFieldValue == emptyValue {
			continue
		}
		if fromFieldValue.Type().Kind() != fieldValue.Type().Kind() {
			continue
		}

		fieldValue.Set(fromFieldValue)
	}
}

func MakeSlice[T any](l int) (to []*T) {
	for i := 0; i < l; i++ {
		t := new(T)
		to = append(to, t)
	}
	return
}
