package parser

import (
	"strings"
	"unicode"
)

// Token wraps a given token with context information
type Token struct {
	// The inferred type of the token
	Type TokenType

	// The actual token contents
	Value string

	// The line that this token appeared on, regardless of pipeline operators.
	// This is 1-based index.
	Line uint

	// The character within the line that this token started on, regardless of
	// pipeline operators. This is 1-based index.
	Character uint
}

type Tokens []Token

type tokenMode byte

const (
	tokenModeStart tokenMode = iota
	tokenModeKeyword
	tokenModeGeneric
	tokenModeComment
	tokenModeString
)

// tokenize reads an input string, and tokenize's it to PRQL standards.
//
// This function DOES perform normalization into the output tokens.
//
// Additionally, since the pipeline operator "|" and newline are synonymous they
// are normalized to the pipeline operator. This does not effect the tokens
// position within the source text.
func tokenize(source string) Tokens {
	tokens := make(Tokens, 0)

	mode := tokenModeStart
	var numLine uint = 1
	var numChar uint = 0
	var tknStartChar uint = 0
	hadPipeline := false
	hasContent := false
	inBlock := false
	lastChar := '\n'
	var tkn strings.Builder

	var typ TokenType = TokenTypeUnknown
	var lastTyp TokenType = TokenTypeUnknown

	var strChar rune
	strBlockLen := 1
	strBlockRem := 1

	// Local lambda to push the current token and reset the state
	pushToken := func() {
		if tkn.Len() > 0 {
			tokens = append(tokens, Token{
				Type:      typ,
				Value:     tkn.String(),
				Line:      numLine,
				Character: tknStartChar,
			})

			lastTyp = typ
			typ = TokenTypeUnknown
			tkn.Reset()
		}

		hasContent = false
	}

	pushPipeline := func() {
		if hadPipeline && lastTyp != TokenTypePipe {
			tokens = append(tokens, Token{
				Type:      TokenTypePipe,
				Value:     "|",
				Line:      numLine,
				Character: numChar,
			})

			hadPipeline = false
			lastTyp = TokenTypePipe
		}
	}

	for _, char := range source {
		numChar++
		switch mode {
		case tokenModeStart:
			if !unicode.IsSpace(char) {
				tknStartChar = numChar

				// Non-space character declares token has started
				// Ignore any whitespace before this
				if char == '|' {
					pushPipeline()
				} else if char == '#' {
					// Start comment
					mode = tokenModeComment
					typ = TokenTypeComment
				} else if char == '\'' || char == '"' {
					// Start string literal
					strChar = char
					mode = tokenModeString
					typ = TokenTypeString
				} else if strings.ContainsRune("[]()=", char) {
					// Insert operator for block start
					tokens = append(tokens, Token{
						Type:      TokenTypeOperator,
						Value:     string(char),
						Line:      numLine,
						Character: numChar,
					})
					lastTyp = TokenTypeOperator
				} else {
					if lastTyp == TokenTypeKeyword || lastTyp == TokenTypeGeneric {
						mode = tokenModeGeneric
						typ = TokenTypeUnknown
					} else {
						mode = tokenModeKeyword
						typ = TokenTypeKeyword
					}
					tkn.WriteRune(char)
				}
			}

		case tokenModeComment:
			// Comments run until the end of the line
			if char == '\n' {
				// End this comment line
				if hasContent {
					pushToken()
				}

				// Reset back to starting line
				mode = tokenModeStart
				hadPipeline = false
				hasContent = false
			} else {
				if hasContent {
					tkn.WriteRune(char)
				} else if !unicode.IsSpace(char) {
					tkn.WriteRune(char)
					hasContent = true
				}
			}

		case tokenModeString:

			if char == strChar {
				// Strings end when the starting character ends the same number of times
				if !hasContent {
					strBlockLen++
				} else {
					// Check if this is a "block" string
					if inBlock {
						// Deduct from chars remaining for end of string
						strBlockRem--
						if strBlockRem == 0 {
							// End met, push the token
							pushToken()
							mode = tokenModeStart
							strBlockLen = 1
							inBlock = false
						}
					} else {
						// Not a block, just push the token
						pushToken()
						hasContent = false
						mode = tokenModeStart
						strBlockLen = 1
					}
				}
			} else {
				if !hasContent {
					hasContent = true

					if strBlockLen >= 3 {
						inBlock = true
						strBlockRem = strBlockLen
					}
				} else if strBlockRem < strBlockLen {
					// We skipped thinking it might be end of block, but it wasn't.
					// So replace the missing chars we skipped.
					for i := (strBlockLen - strBlockRem); i > 0; i-- {
						tkn.WriteRune(strChar)
					}
					strBlockRem = strBlockLen
				}

				tkn.WriteRune(char)
			}

		case tokenModeKeyword:
			if unicode.IsSpace(char) {
				// keywords end with whitespace
				pushToken()
				mode = tokenModeStart
			} else {
				tkn.WriteRune(char)
				hasContent = true
				hadPipeline = true
			}

		case tokenModeGeneric:
			if unicode.IsSpace(char) {
				// keywords end with whitespace
				pushToken()
				mode = tokenModeStart
			} else {
				// Check for special cases
				if char == '\'' || char == '"' {
					// Could be F-STRING or S-STRING
					if tkn.Len() == 2 && (lastChar == 'f' || lastChar == 's') {
						// Start string literal
						tkn.Reset()
						strChar = char
						strBlockLen = 1
						mode = tokenModeString
						hasContent = false

						if lastChar == 'f' {
							typ = TokenTypeFString
						} else if lastChar == 's' {
							typ = TokenTypeSString
						}
					}
				} else {
					typ = TokenTypeGeneric

					tkn.WriteRune(char)
					hasContent = true
					hadPipeline = true
				}
			}
		}
		// END SWITCH

		if char == '\n' {
			if !inBlock {
				if hadPipeline {
					pushToken()
					hasContent = false
				}
				pushPipeline()
			}

			numLine++
			numChar = 0
		}
		lastChar = char
	}

	return tokens
}
