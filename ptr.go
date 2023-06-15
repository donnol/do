package do

// PtrTo return v's pointer
func PtrTo[T any](v T) *T { return &v }

func Zero[T any]() T {
	var v T
	return v
}
