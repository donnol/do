package do

// FieldWithAlias use alias if not empty, or use defaultAlias
func FieldWithAlias(field, alias, defaultAlias string) string {
	if alias == "" {
		alias = defaultAlias
	}
	if alias == "" {
		return field
	}
	return alias + "." + field
}

// WithWhere return condition if t is not a zero value of T
func WithWhere[T comparable](
	t T,
	cond func(field string, value interface{}) string,
	field string,
	value ...T,
) string {
	useCond := false
	switch v := any(t).(type) {
	case interface{ IsZero() bool }:
		useCond = !v.IsZero()
	default:
		useCond = !IsZero(t)
	}
	if !useCond {
		return ""
	}

	// replace t with value
	if len(value) != 0 {
		t = value[0]
	}

	return cond(field, t)
}
