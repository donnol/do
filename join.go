package do

// NestedJoin like nested loop join
func NestedJoin[J, K, R any](
	left []J,
	right []K,
	match func(J, K) bool,
	mapper func(J, K) R,
) []R {
	var r = make([]R, 0, len(left))

	for _, j := range left {
		for _, k := range right {
			if !match(j, k) {
				continue
			}
			r = append(r, mapper(j, k))
		}
	}

	return r
}

// HashJoin like hash join
func HashJoin[K comparable, LE, RE, R any](
	left []LE,
	right []RE,
	lk func(item LE) K,
	rk func(item RE) K,
	mapper func(LE, RE) R,
) []R {
	var r = make([]R, 0, len(left))

	rm := KeyBy(right, rk)

	for _, le := range left {
		k := lk(le)
		re := rm[k]
		r = append(r, mapper(le, re))
	}

	return r
}
