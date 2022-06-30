package syntax

import (
	"fmt"

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
