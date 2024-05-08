package do

func Unique[S comparable](s []S) []S {
	r := make([]S, 0, len(s))

	m := make(map[S]struct{}, len(s))
	for _, e := range s {
		_, ok := m[e]
		if ok {
			continue
		}

		m[e] = struct{}{}
		r = append(r, e)
	}

	return r
}

func In[S comparable](s []S, e S) bool {
	for _, item := range s {
		if e == item {
			return true
		}
	}
	return false
}

func First[T any](s []T) (t T, ok bool) {
	if len(s) == 0 {
		return
	}
	return s[0], true
}

func Last[T any](s []T) (t T, ok bool) {
	l := len(s)
	if l == 0 {
		return
	}
	return s[l-1], true
}

func Index[T any](s []T, i int) (t T, ok bool) {
	l := len(s)
	if l == 0 || i > l-1 {
		return
	}
	return s[i], true
}

func FirstOr[T any](s []T) (t T) {
	t, _ = First(s)
	return
}

func LastOr[T any](s []T) (t T) {
	t, _ = Last(s)
	return
}

func IndexOr[T any](s []T, i int) (t T) {
	t, _ = Index(s, i)
	return
}
