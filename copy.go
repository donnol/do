package do

// Copy return a copy of v
func Copy[T any](v *T) *T {
	if v == nil {
		return nil
	}

	t := *v
	return &t
}
