// Package parser is responsible for taking incoming PRSQL strings and parsing,
// tokenizing, and generating a workable AST (Abstract Syntax Tree).
package parser

import "github.com/chris-pikul/go-prql"

// AST Dummy type alias.
type AST = map[string]string

// Parse takes the incoming PRQL query as a string, and attempts to
// parse/tokenize it into a working AST.
//
// Returns the AST, and a PRQLError for any errors occuring during parsing.
func Parse(source string) (AST, *prql.Error) {
	return nil, nil
}
