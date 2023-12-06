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
