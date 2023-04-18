package do

// As assert the value v to type T
func As[T any](v any) T {
	return v.(T)
}
