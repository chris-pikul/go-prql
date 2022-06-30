package syntax

import (
	"strconv"
	"strings"

	"github.com/chris-pikul/go-prql/utils"
)

// Header represents the top-level "prql" dialect and version declaration.
//
// Syntax: "prql dialect:{string} version:{integer}"
type Header struct {
	Version utils.Optional[int]
	Dialect Dialect
}

// String returns the PRQL expression for defining this Header.
func (h Header) String() string {
	var str strings.Builder
	str.WriteString("prsql")

	if ver, ok := h.Version.Get(); ok {
		str.WriteString(" version:")
		str.WriteString(strconv.Itoa(*ver))
	}

	if h.Dialect != DialectGeneric {
		str.WriteString(" dialect:")
		str.WriteString(h.Dialect.String())
	}

	return str.String()
}
