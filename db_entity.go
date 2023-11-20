package do

// EntityWithTotal used to scan col1, col2, col3 to T, and scan total to Total when columns is [col1, col2, col3, total] if *T is interface{ ValuePtrs() []any }
type EntityWithTotal[T any] struct {
	Inner T     // data
	Total int64 `db:"total"` // total
}

// ValuePtrs return T's fields and Total's pointer if *T is interface{ ValuePtrs() []any }, or return Inner and Total's pointer; use by rows.Scan(s.ValuePtrs()....)
func (s *EntityWithTotal[T]) ValuePtrs() []any {
	switch v := any(&s.Inner).(type) {
	case interface{ ValuePtrs() []any }:
		ptrs := v.ValuePtrs()
		ptrs = append(ptrs, &s.Total)
		return ptrs
	default:
		return []any{
			&s.Inner,
			&s.Total,
		}
	}
}
