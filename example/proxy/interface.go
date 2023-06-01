package main

import "context"

//go:generate go install github.com/donnol/tools/cmd/tbc@latest
//go:generate tbc mock .
// Run `go generate` in this dir will generate mock.go file contains Mock struct for below interface.

type (
	UserStore interface {
		ById(ctx context.Context, id uint64) User
	}

	User struct {
		Id   uint64
		Name string
	}

	BookStore interface {
		ById(ctx context.Context, id uint64) Book
		Hook()
	}

	Book struct {
		Id       uint64
		Name     string
		UserId   uint64
		UserName string
	}
)

type (
	BookService interface {
		WithUser(ctx context.Context, id uint64) Book
	}
)
