package do

// KeyValueBy slice to map, key value specified by iteratee
func KeyValueBy[K comparable, E, V any](collection []E, iteratee func(item E) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for _, v := range collection {
		k, r := iteratee(v)
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
