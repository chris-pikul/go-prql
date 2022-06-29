package prql

// ErrorType is a Go style enum (internally a byte) for the type of error
// that occured. Intented to be used within the Error object to
// differentiate between different types, or sources, of errors.
type ErrorType byte

const (
	// ErrorTypeUnknown signifies the type of error is unknown, and as such
	// should represent a fatal error. The most common explanation for the type
	// being unknown would be an improperly instantiated Error object, or
	// an unknown runtime panic that was caught be recover().
	ErrorTypeUnknown ErrorType = iota

	// ErrorTypeSyntax signifies the error was generated during parsing, and
	// is considered a syntax error (client error) relating to the input PRQL
	// query.
	ErrorTypeSyntax
)

func (t ErrorType) String() string {
	switch t {
	case ErrorTypeSyntax:
		return "SYNTAX"
	}

	return "UNKNOWN"
}

// Valid returns true if the ErrorType is within the acceptable range of
// known errors. This excludes the ErrorTypeUnknown constant, as those are
// reserved for truely uknown or zero-value errors.
func (t ErrorType) Valid() bool {
	return t > ErrorTypeUnknown && t <= ErrorTypeSyntax
}

// Error wraps a generic Go error with more context for PRQL parsing and
// generation.
//
// Implements the "error" interface for interoperability
type Error struct {
	Type ErrorType
	Err  error
}

// Error implements the `error` interface. Returns the underlying error object
// Error.Err.Error() string.
func (e Error) Error() string {
	return e.Err.Error()
}

// String implements the `string` interface. Returns the underlying error object
// Error.Err.Error() string.
func (e Error) String() string {
	return e.Error()
}

// NewError creates a new Error.
func NewError(errType ErrorType, parent error) Error {
	return Error{errType, parent}
}
