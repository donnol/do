package do

func Zero[T any]() T {
	var v T
	return v
}

// Pointer return a pointer of type T
func Pointer[T any]() *T { return new(T) }
