package prql

import "fmt"

// ErrorType is a Go style enum (internally a byte) for the type of error
// that occured. Intented to be used within the Error object to
// differentiate between different types, or sources, of errors.
type ErrorType byte

const (
	// ErrorTypeUnknown signifies the type of error is unknown, and as such
	// should represent a fatal error. The most common explanation for the type
	// being unknown would be an improperly instantiated Error object, or
	// an unknown runtime panic that was caught be recover().
	//
	// Encoded as "UNKNOWN"
	ErrorTypeUnknown ErrorType = iota

	// ErrorTypeSyntax signifies the error was generated during parsing, and
	// is considered a syntax error (client error) relating to the input PRQL
	// query.
	//
	// Encoded as "SYNTAX"
	ErrorTypeSyntax
)

// String returns a string representation of the ErrorType enum. By default, Go
// will use this for encoding as well.
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

// UnmarshalText implements the encoding.TextUnmarshaler interface. Allows for
// deserialization to read string values and convert into the underlying
// ErrorType enum.
func (t *ErrorType) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "SYNTAX":
		*t = ErrorTypeSyntax
	default:
		*t = ErrorTypeUnknown
		return fmt.Errorf("ErrorType '%s' is invalid", str)
	}
	return nil
}
