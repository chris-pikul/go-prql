package utils

// KeyForValue returns the key (or nil if not found) within a given map that
// matches the given comparable value.
//
// The second returned value is a boolean "ok" if found at all.
func KeyForValue[K comparable, V comparable](tgt map[K]V, value V) (*K, bool) {
	for key, val := range tgt {
		if val == value {
			return &key, true
		}
	}

	return nil, false
}

// InvertMap returns a new map with all the given key-values on the target map,
// swapped. That is, the keys become values, and the values become the keys.
//
// Both the key and value must be comparable.
func InvertMap[K comparable, V comparable](tgt map[K]V) map[V]K {
	result := make(map[V]K, len(tgt))
	for key, val := range tgt {
		result[val] = key
	}
	return result
}
