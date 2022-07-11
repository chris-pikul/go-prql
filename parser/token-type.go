package parser

// TokenType declares the type, or context in which a token is held.
type TokenType byte

const (
	// TokenTypeUnknown is a zero-value default for unknown and possibly
	// erronous token.
	TokenTypeUnknown TokenType = iota

	// TokenTypePipe represents a pipeline-operator "|", either imlicit at the end
	// of a line, or explicitly present.
	TokenTypePipe

	// TokenTypeComment holds the contents of a comment, for preservation.
	TokenTypeComment

	// TokenTypeKeyword declares that the token should be a transform keyword
	TokenTypeKeyword

	// TokenTypeGeneric declares that the tokens content is considered some kind
	// of generic user-space content such as an alias, or reference, etc.
	TokenTypeGeneric

	// TokenTypeOperator represents an operator which matters within the context
	// of surrounding tokens. Such as an "=", or ":", or otherwise
	TokenTypeOperator

	// TokenTypeString represents the entire strings contents as parsed within
	// the start and end characters (or combination of them).
	TokenTypeString

	// TokenTypeFString represents the entire contents of a f-string.
	TokenTypeFString

	// TokenTypeSString represents the entire contents of a s-string.
	TokenTypeSString
)

func (t TokenType) String() string {
	switch t {
	case TokenTypePipe:
		return "PIPE"
	case TokenTypeComment:
		return "COMMENT"
	case TokenTypeKeyword:
		return "KEYWORD"
	case TokenTypeGeneric:
		return "GENERIC"
	case TokenTypeOperator:
		return "OPERATOR"
	case TokenTypeString:
		return "STRING"
	case TokenTypeFString:
		return "F-STRING"
	case TokenTypeSString:
		return "S-STRING"
	default:
		return "UNKNOWN"
	}
}
