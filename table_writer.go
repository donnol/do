package do

import (
	"fmt"
	"reflect"
)

type TableWriter interface {
	SetHeader(headers []string)
	SetRows(rows [][]string)
}

// TableWrite use T's field name as column, T's field value as cell, which T must be a struct.
func TableWrite[T any](w TableWriter, list []T) {
	if len(list) == 0 {
		return
	}

	var headers []string
	rows := make([][]string, 0, len(list))
	for index, item := range list {
		val := reflect.ValueOf(item)

		if index == 0 {
			typ := val.Type()
			if typ.Kind() != reflect.Struct {
				panic(fmt.Errorf("the list element type is not struct"))
			}

			headers = make([]string, 0, typ.NumField())
			for i := 0; i < typ.NumField(); i++ {
				headers = append(headers, typ.Field(i).Name)
			}
		}

		row := make([]string, 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			row = append(row, fmt.Sprintf("%v", val.Field(i).Interface()))
		}
		rows = append(rows, row)
	}

	w.SetHeader(headers)
	w.SetRows(rows)
}
