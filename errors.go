package prql

import (
	"fmt"
)

// Error wraps a generic Go error with more context for PRQL parsing and
// generation.
//
// Implements the "error" interface for interoperability
type Error struct {
	Type ErrorType
	Err  error
}

// String implements the `string` interface. Formats the error into a printable
// version followin the pattern: "PRQL {type} error: {message}". This for
// example would look like "PRQL SYNTAX error: unknown keyword 'test'" for a
// syntax error.
func (e Error) String() string {
	return fmt.Sprintf("PRQL %s error: %s", e.Type.String(), e.Err.Error())
}

// Error implements the `error` interface. Returns the underlying error object
// Error.Err.Error() string. Does not perform any formatting to the returned
// error, for that use prql.Error.String().
func (e Error) Error() string {
	return e.Err.Error()
}

// NewError creates a new Error.
func NewError(errType ErrorType, parent error) Error {
	return Error{errType, parent}
}

// NewSyntaxErrorf generates a new Error with a pre-defined ErrorType of
// ErrorTypeSyntax, and provides a string style formatting for generating the
// parernt error object within it.
func NewSyntaxErrorf(format string, args ...interface{}) Error {
	return Error{
		ErrorTypeSyntax,
		fmt.Errorf(format, args...),
	}
}
