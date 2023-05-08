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
