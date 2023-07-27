package do

// PtrTo [deprecate] return pointer of a new value copied from v
func PtrTo[T any](v T) *T { return &v }

func Zero[T any]() T {
	var v T
	return v
}
