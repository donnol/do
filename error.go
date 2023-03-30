package do

import "fmt"

type Error[T any] struct {
	inner T
}

func NewError[T any](inner T) error {
	return &Error[T]{
		inner: inner,
	}
}

func (e Error[T]) Error() string {
	return fmt.Sprintf("%v", e.inner)
}

func (e Error[T]) Inner() T {
	return e.inner
}

type e struct {
	raw error
}

func (e e) Error() string {
	return fmt.Sprintf("raw error: %v", e.raw)
}

func (e e) Raw() error {
	return e.raw
}

func efrom(err error) *e {
	return &e{
		raw: err,
	}
}

func MatchError(v any) bool {
	_, ok := ConvertError(v)
	return ok
}

func ConvertError(v any) (*Error[*e], bool) {
	r, ok := v.(*Error[*e])
	return r, ok
}
