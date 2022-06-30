package syntax

import (
	"fmt"

	"github.com/chris-pikul/go-prql/utils"
)

// Dialect is a byte enum representing the accepted dialects. This is declared
// by the top-level "prql" expression.
type Dialect byte

const (
	DialectGeneric Dialect = iota
	DialectANSI
	DialectBigQuery
	DialectClickHouse
	DialectHive
	DialectMSSQL
	DialectMYSQL
	DialectPostgres
	DialectSQLite
	DialectSnowflake
)

// holds Dialect -> string mapping
var dialectStringMap = map[Dialect]string{
	DialectGeneric:    "generic",
	DialectANSI:       "ansi",
	DialectBigQuery:   "bigquery",
	DialectClickHouse: "clickhouse",
	DialectHive:       "hive",
	DialectMSSQL:      "mssql",
	DialectMYSQL:      "mysql",
	DialectPostgres:   "postgres",
	DialectSQLite:     "sqlite",
	DialectSnowflake:  "snowflake",
}

// holds string -> Dialect mapping
var dialectDialectMap = utils.InvertMap(dialectStringMap)

// String returns the string representation of the underlying Dialect enum. If
// invalid, defaults to returning "generic".
func (d Dialect) String() string {
	if str, ok := dialectStringMap[d]; ok {
		return str
	}

	return "generic"
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. Allows for
// reading strings into their Dialect enum value. Returns an error if the text
// is invalid and does not match a known dialect.
//
// Important: this is CASE-SENSITIVE.
func (d *Dialect) UnmarshalText(text []byte) error {
	str := string(text)
	if dial, ok := dialectDialectMap[str]; ok {
		*d = dial
	}
	return fmt.Errorf("invalid Dialect '%s'", str)
}
