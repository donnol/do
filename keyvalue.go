package do

// KeyValueBy slice to map, key value specified by iteratee
func KeyValueBy[K comparable, E, V any](collection []E, iteratee func(item E) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for i := range collection {
		k, r := iteratee(collection[i])
		result[k] = r
	}

	return result
}

// KeyBy slice to map, key specified by iteratee, value is slice element
func KeyBy[K comparable, E any](collection []E, iteratee func(item E) K) map[K]E {
	return KeyValueBy(collection, func(item E) (K, E) {
		return iteratee(item), item
	})
}

func Keys[K comparable, E any](collection map[K]E) []K {
	result := make([]K, len(collection))

	var i int
	for k := range collection {
		result[i] = k
		i++
	}

	return result
}

func Values[K comparable, E any](collection map[K]E) []E {
	result := make([]E, len(collection))

	var i int
	for _, v := range collection {
		result[i] = v
		i++
	}

	return result
}

// MergeKeyValue merge m2 into m1, will override m1 key value if both exists
func MergeKeyValue[K comparable, V any](m1, m2 map[K]V) map[K]V {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func MapFrom[T any, R interface{ From(T) }](s []T, initial func() R) []R {
	r := make([]R, len(s))
	for i, item := range s {
		// v必须是一个指针，如果直接传入R值会导致值一直一样；因此改为传入initial函数，用它来生成生成R对象值
		v := initial()
		v.From(item)
		r[i] = v
	}
	return r
}

func MapFrom2[T any, R interface{ From(T) R }](s []T) []R {
	var f R

	r := make([]R, len(s))
	for i, item := range s {
		// 有了返回值，因此不需要使用指针
		v := f.From(item)
		r[i] = v
	}
	return r
}

func MapTo[T interface{ To() R }, R any](s []T) []R {
	r := make([]R, len(s))
	for i, item := range s {
		r[i] = item.To()
	}
	return r
}

func MapTo2[T interface{ To(T) R }, R any](s []T) []R {
	r := make([]R, len(s))
	for i, item := range s {
		r[i] = item.To(item)
	}
	return r
}
