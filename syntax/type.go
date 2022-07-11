package syntax

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/chris-pikul/go-prql/utils"
)

// Type is an enumeration (byte) which declares the type declaration related to
// a PRSQL value/parameter
type Type byte

const (
	// TypeUnknown represents an unknown erronous type declaration
	TypeUnknown Type = iota

	// TypeBoolean represents a boolean
	TypeBoolean

	// TypeInteger represents an integer (signed or unsigned)
	TypeInteger

	// TypeScalar represents a scalar
	TypeScalar

	// TypeFloat represents a floating point number
	TypeFloat

	// TypeString represents a string (variable length by default)
	TypeString

	// TypeDate represents a temporal date (no time included)
	TypeDate

	// TypeTime represents a temporal time (no date included)
	TypeTime

	// TypeTimestamp represents a temporal timestamp, date and time are included
	TypeTimestamp

	// TypeTable represents a table reference
	TypeTable

	// TypeColumn represents a column reference
	TypeColumn
)

// holds types -> string mapping
var typeStringMap = map[Type]string{
	TypeUnknown:   "unknown",
	TypeBoolean:   "boolean",
	TypeInteger:   "integer",
	TypeScalar:    "scalar",
	TypeFloat:     "float",
	TypeString:    "string",
	TypeDate:      "date",
	TypeTime:      "time",
	TypeTimestamp: "timestamp",
	TypeTable:     "table",
	TypeColumn:    "column",
}

// holds string -> type mapping
var typeTypeMap = utils.InvertMap(typeStringMap)

// String returns the string representation of Type
func (t Type) String() string {
	if str, ok := typeStringMap[t]; ok {
		return str
	}

	return "unknown"
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. Converts
// string representations of Type into the type iteself, or returns an error if
// the string is invalid.
func (t *Type) UnmarshalText(text []byte) error {
	str := string(text)
	if typ, ok := typeTypeMap[str]; ok {
		*t = typ
	}

	return fmt.Errorf("invalid Type '%s'", str)
}

var reNumeric = regexp.MustCompile(`(?i)^[+-]?(?:\d*\.?)\d+(?:e[+-]?\d+)?[df]?$`)

// inferNumeric checks if the given string is a numeric value or mathmatical
// expression.
func inferNumeric(literal string) bool {
	return reNumeric.MatchString(literal)
}

// InferType attempts to guess (infer) at the type of a given literal string.
// It returns the type which may be TypeUnknown, and optionally an error if the
// inferance cannot be performed.
//
// A point of confusion may be when inferring a reference to either a table or
// column since these are string names that require context within the scope of
// their expression.
func InferType(literal string) (Type, error) {
	if len(literal) == 0 {
		return TypeUnknown, errors.New("cannot infer type on empty string")
	}

	// Check if string, based on first character being a quote
	first := literal[0]
	if literal[0] == '\'' || literal[0] == '"' {
		// Ensure valid by checking last character matching first
		if literal[len(literal)-1] != first {
			return TypeString, fmt.Errorf(`string literal "%s" does not terminate with the same character %c`, literal, first)
		}

		return TypeString, nil
	}

	if inferNumeric(literal) {
		return TypeScalar, nil
	}

	return TypeUnknown, nil
}

// ReflectType returns the Type which is correct for the given value (interface{})
// using reflection.
func ReflectType(val any) Type {
	reflType := reflect.ValueOf(val)

	// If it's a pointer, dereference it for the proper type
	if reflType.Kind() == reflect.Ptr {
		reflType = reflType.Elem()
	}

	// Check for the primitive types first
	switch reflType.Type().Kind() {
	case reflect.Bool:
		return TypeBoolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return TypeInteger
	case reflect.Float32, reflect.Float64:
		return TypeFloat
	case reflect.String:
		return TypeString
	case reflect.Struct:
		// Special cases for structs, they may be times.
		reflInterf := reflType.Interface()
		switch reflInterf.(type) {
		case time.Time:
			// We have no way of knowing whether this is anything but a timestamp
			return TypeTimestamp
		}
	}

	return TypeUnknown
}
