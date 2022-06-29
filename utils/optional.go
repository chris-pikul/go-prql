package utils

// Optional is a generic struct holding a pointer to the type T. Additionally it
// holds an "Ok" value specifying if the value was explicitly set.
//
// Unlike using single pointers to declare optional values, the addition of the
// Ok method can mean that while the value maybe nil, it was explicitly set to
// nil, and is not just a zero value.
type Optional[T any] struct {
	value *T
	ok    bool
}

// Value returns the underlying pointer to the value of this Optional.
func (opt Optional[T]) Value() *T {
	return opt.value
}

// Ok returns a boolean on whether the value was explicitly set or not. An
// Optional.Value() can still return nil even if this returns true, meaning that
// the nil value is intended, and not a zero-value.
func (opt Optional[T]) Ok() bool {
	return opt.ok
}

// Get returns both the value pointer, and the "ok" boolean signifying the
// the current value, and whether it was explicitly set.
func (opt Optional[T]) Get() (*T, bool) {
	return opt.value, opt.ok
}

// Set explicitly sets the value of this Optional and turns the "Ok" flag to
// true, regardless of if the value is nil.
func (opt *Optional[T]) Set(value *T) {
	opt.value = value
	opt.ok = true
}

// Clear sets the value to nil, and clears the "Ok" flag to false.
func (opt *Optional[T]) Clear() {
	opt.value = nil
	opt.ok = false
}

// NewOptional creates a new generic Optional value of given type T.
// It takes an incoming pointer value of the type T, as well as a boolean flag
// on whether this is explicitly set.
//
// If the value parameter is not nil, the explicit flag is ignored and will
// create an Optional with "Ok" set to true.
func NewOptional[T any](value *T, explicit bool) Optional[T] {
	isOk := explicit
	if value != nil {
		isOk = true
	}

	return Optional[T]{
		value: value,
		ok:    isOk,
	}
}
