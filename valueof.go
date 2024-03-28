package do

import (
	"reflect"
)

// ValueOf return v.field's value, if v is a struct or map
func ValueOf[F any](v any, field string) (f F, ok bool) {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()

	var fv reflect.Value
	switch rt.Kind() {
	case reflect.Struct:
		fv = rv.FieldByName(field)
	case reflect.Map:
		fv = rv.MapIndex(reflect.ValueOf(field))
	default:
		return
	}
	if !fv.IsValid() {
		return
	}

	var success bool
	f, success = fv.Interface().(F)
	if !success {
		return
	}
	ok = true

	return
}
