package do

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type Finder[R any] interface {
	// return query sql and args
	Query() (query string, args []any)

	// new a result type object, not the same, to receive every row
	// fields must be pointer type of result object's field, and it is match with query sql's select column one by one
	NewScanObjAndFields(colTypes []*sql.ColumnType) (r *R, fields []any)
}

// FindFunc return query, args, and object generator; object generator should return new object and it's related fields when called; FindFunc is a Finder too.
type FindFunc[R any] func() (query string, args []any, genObj func(colTypes []*sql.ColumnType) (r *R, fields []any))

func (f FindFunc[R]) Query() (query string, args []any) {
	query, args, _ = f()
	return query, args
}

func (f FindFunc[R]) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *R, fields []any) {
	_, _, genObj := f()
	return genObj(colTypes)
}

var (
	_ Finder[int] = (FindFunc[int])(nil)
)

type Storer interface {
	*sql.DB | *sql.Tx | *sql.Conn
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func FindList[S Storer, F Finder[R], R any](db S, finder F, res *[]R) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...) // sql里select了n列
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		obj, fields := finder.NewScanObjAndFields(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		*res = append(*res, *obj)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindFirst[S Storer, F Finder[R], R any](db S, finder F, res *R) (err error) {
	var r []R
	err = FindList(db, finder, &r)
	if err != nil {
		return
	}
	if len(r) > 0 {
		*res = r[0]
	}
	return
}

func FindListByFunc[S Storer, F FindFunc[R], R any](db S, finder F, res *[]R) (err error) {
	query, args, genObj := finder()
	rows, err := db.QueryContext(context.TODO(), query, args...) // sql里select了n列
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		obj, fields := genObj(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		*res = append(*res, *obj)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindFirstByFunc[S Storer, F FindFunc[R], R any](db S, finder F, res *R) (err error) {
	var r []R
	err = FindListByFunc(db, finder, &r)
	if err != nil {
		return
	}
	if len(r) > 0 {
		*res = r[0]
	}
	return
}

func FindAll[S Storer, F Finder[R], R any](db S, finder F, initial R) (r []R, err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...) // sql里select了n列
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		obj, fields := finder.NewScanObjAndFields(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		r = append(r, *obj)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindOne[S Storer, F Finder[R], R any](db S, finder F, initial R) (r R, err error) {
	res, err := FindAll(db, finder, initial)
	if err != nil {
		return r, err
	}
	if len(res) > 0 {
		r = res[0]
	}
	return
}

func PrintFields(fields []any) {
	fmt.Println("=== begin print fields")
	for i := range fields {
		vall := reflect.ValueOf(fields[i])
		if vall.Kind() == reflect.Pointer {
			vall = vall.Elem()
		}
		fmt.Printf("field: %v\n", vall.Interface())
	}
	fmt.Println("=== end print fields")
}
