package syntax

import (
	"strconv"
	"strings"

	"github.com/chris-pikul/go-prql/utils"
)

// Query represents the top-level "prql" dialect and version declaration
type Query struct {
	Version utils.Optional[int]
	Dialect Dialect
}

// String returns the PRQL expression for defining this Query.
func (q Query) String() string {
	var str strings.Builder
	str.WriteString("prsql")

	if ver, ok := q.Version.Get(); ok {
		str.WriteString(" version:")
		str.WriteString(strconv.Itoa(*ver))
	}

	if q.Dialect != DialectGeneric {
		str.WriteString(" dialect:")
		str.WriteString(q.Dialect.String())
	}

	return str.String()
}
