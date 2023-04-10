# do

[![GoDoc](https://pkg.go.dev/badge/github.com/donnol/do)](https://pkg.go.dev/github.com/donnol/do) [![Go Report Card](https://goreportcard.com/badge/github.com/donnol/do)](https://goreportcard.com/report/github.com/donnol/do)

Do something interesting.

Base Go1.18 with generic.

## Must

Panic if error is not nil, otherwise return some result except the error.

```go
package main

import (
	"fmt"

	"github.com/donnol/do"
)

func main() {
	do.Must(retErr()) // without result

	// specify result type with type parameter
	_ = do.Must1(retErrAndOneResult()) // with one result

	_, _ = do.Must2(retErrAndTwoResult()) // with two result
}

func retErr() error {
	return fmt.Errorf("a new error")
}

func retErrAndOneResult() (int, error) {
	return 1, fmt.Errorf("a new error")
}

func retErrAndTwoResult() (int, int, error) {
	return 0, 1, fmt.Errorf("a new error")
}
```

## Slice to map by key

```go
r := KeyValueBy([]string{"a", "aa", "aaa"}, func(str string) (string, int) {
	return str, len(str)
})
want := map[string]int{"a": 1, "aa": 2, "aaa": 3}
// r is what we want.
```

## Join

```go
r := NestedJoin([]Book{
	{Id: 1, Title: "hello", Author: 1},
	{Id: 2, Title: "world", Author: 1},
	{Id: 3, Title: "good", Author: 2},
	{Id: 4, Title: "job", Author: 2},
}, []User{
	{Id: 1, Name: "jd"},
	{Id: 2, Name: "jc"},
}, UserBookMatcher, func(j Book, k User) BookWithUser {
	return BookWithUser{
		Book:     j,
		UserName: k.Name,
	}
})
want := []BookWithUser{
	{Book{1, "hello", 1}, "jd"},
	{Book{2, "world", 1}, "jd"},
	{Book{3, "good", 2}, "jc"},
	{Book{4, "job", 2}, "jc"},
}
// r is what we want.
```

```go
r := HashJoin([]Book{
	{Id: 1, Title: "hello", Author: 1},
	{Id: 2, Title: "world", Author: 1},
	{Id: 3, Title: "good", Author: 2},
	{Id: 4, Title: "job", Author: 2},
}, []User{
	{Id: 1, Name: "jd"},
	{Id: 2, Name: "jc"},
}, func(item Book) uint64 {
	return item.Author
}, func(item User) uint64 {
	return item.Id
}, func(j Book, k User) BookWithUser {
	return BookWithUser{
		Book:     j,
		UserName: k.Name,
	}
})
want := []BookWithUser{
	{Book{1, "hello", 1}, "jd"},
	{Book{2, "world", 1}, "jd"},
	{Book{3, "good", 2}, "jc"},
	{Book{4, "job", 2}, "jc"},
}
// r is what we want.
```

## Send HTTP request

Send a http request with a simple function.

## Worker

A worker pool process job with a limited number Goroutine.

## DB connect and find

```go
// 0. open a db
var tdb *sql.DB

// 1. define a finder
type finderOfUser struct {
	id uint64
}

func (f *finderOfUser) Query() (query string, args []any) {
	query = `select * from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUser) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *UserForDB, fields []any) {
	r = &UserForDB{}
	fields = append(fields,
		&r.Id,
		&r.Name,
	)
	return
}

// 2. find
finder := &finderOfUser{
	id: 1,
}
r, err := do.FindAll(tdb, finder, (UserForDB{}))
if err != nil {
	panic(err)
}

// 3. find per batch
finder := &finderOfUser{
	id: 1,
}
// batchNum is 10, if there are 20 records, it will be processed in two parts
err := do.Batch(tdb, finder, 10, func(r []UserForDB) error {
	// Process this batch of data

	return nil
}
```
