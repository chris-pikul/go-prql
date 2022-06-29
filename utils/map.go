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
