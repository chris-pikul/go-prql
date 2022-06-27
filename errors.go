package prql

// PRQLError wraps a generic Go error with more context for PRQL parsing and
// generation.
//
// Implements the "error" interface for interoperability
type PRQLError struct {
	Err error
}

// Error implements the `error` interface. Returns the underlying error object
// PRQLError.Err.Error() string.
func (e PRQLError) Error() string {
	return e.Err.Error()
}

// String implements the `string` interface. Returns the underlying error object
// PRQLError.Err.Error() string.
func (e PRQLError) String() string {
	return e.Error()
}

// NewPRQLError creates a new wrapped error using the provided parent "error"
// object. The parent can be nil.
func NewPRQLError(parent error) PRQLError {
	return PRQLError{parent}
}
