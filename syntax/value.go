package syntax

// Value holds an individual value within a PRQL expression. This includes a
// Type declaration, and the value itself which is maintained as a interface{}.
type Value[T any] struct {
	typ   Type
	value T
}

func (v Value[T]) Type() Type {
	return v.typ
}

func (v Value[T]) Get() T {
	return v.value
}

func NewValue[T any](val T) Value[T] {
	return Value[T]{
		typ:   TypeUnknown,
		value: val,
	}
}
